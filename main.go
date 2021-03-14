package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/gizak/termui/v3"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func toGB(value uint64) float64 {
	return float64(value) / 1024 / 1024 / 1024
}

func manual() termui.Drawable {
	commands := []string{
		"q TO QUIT DEMO",
		"h PRINT HOST INFO",
		"m PRINT MEMORY USAGE",
	}

	p := widgets.NewParagraph()
	p.Title = "Manual"
	p.Text = strings.Join(commands, "\n")

	p.SetRect(0, 0, 50, 7)
	p.TextStyle.Fg = ui.ColorWhite

	return p
}

func memory() termui.Drawable {
	v, _ := mem.VirtualMemory()

	g3 := widgets.NewGauge()
	g3.Title = fmt.Sprintf(" Memory Total: %vGB, Used: %.2fGB ", toGB(v.Total), toGB(v.Used))
	g3.SetRect(0, 20, 50, 7)
	g3.Percent = int(v.UsedPercent)
	g3.BarColor = ui.ColorGreen
	g3.Label = fmt.Sprintf("%v%% (%.2fGBs free)", g3.Percent, toGB(v.Available))

	return g3
}

func hostInfo() termui.Drawable {

	v, _ := host.Info()

	p := widgets.NewParagraph()
	p.Title = "Host Info"

	p.Text = fmt.Sprintf(
		`Hostname       : %v
		OS             : %v
		Platform       : %v 
		PlatformFamily : %v 
		PlatformVersion: %v`,
		v.Hostname, v.OS, v.Platform, v.PlatformFamily, v.PlatformVersion)

	p.SetRect(0, 20, 50, 7)
	p.TextStyle.Fg = ui.ColorWhite

	return p
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	manual := manual()

	draw := func(el func() termui.Drawable) {
		if el == nil {
			ui.Render(manual)
		} else {
			ui.Render(manual, el())
		}

	}

	var current func() termui.Drawable

	tickerCount := 1
	draw(nil)
	tickerCount++
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(500 * time.Millisecond).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "m":
				current = memory
			case "h":
				current = hostInfo
			}

		case <-ticker:
			draw(current)
			tickerCount++
		}
	}
}
