package display

import (
	"image/color"
	"machine"
	"strings"
	"time"

	"github.com/rin2yh/tiny-deck/internal/serial"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/proggy"
)

var white = color.RGBA{R: 255, G: 255, B: 255, A: 255}

type Command int

const (
	CommandNone Command = iota
	CommandNotify
	CommandMetricsChanged
	CommandMetricsStale
)

type Display struct {
	oled           *ssd1306.Device
	reader         serial.Reader
	metrics        metricsState
	gopher         gopherAnim
	layerName      func() string
	lastLayer      string
	lastLayerLabel string
	lastStale      bool
}

func New(oled *ssd1306.Device, layerName func() string) Display {
	return Display{oled: oled, layerName: layerName, gopher: newGopherAnim()}
}

func (d *Display) CPUPercent() float32 { return d.metrics.cpu }
func (d *Display) MemPercent() float32 { return d.metrics.mem }

func (d *Display) Update(s machine.Serialer) Command {
	now := time.Now()
	gopherTicked := d.gopher.tick(now)
	currentLayer := d.layerName()
	layerChanged := currentLayer != d.lastLayer
	metricsChanged := false

	for d.reader.Read(s) {
		l := d.reader.Line()
		if l == "notify" {
			return CommandNotify
		}
		if strings.HasPrefix(l, prefixCPU) || strings.HasPrefix(l, prefixMem) {
			if d.metrics.tryUpdate(l, now) {
				metricsChanged = true
			}
		}
	}

	if layerChanged {
		d.lastLayer = currentLayer
		d.lastLayerLabel = "[" + currentLayer + "]"
	}

	stale := d.metrics.isStale(now)
	staleChanged := stale != d.lastStale

	if !metricsChanged && !layerChanged && !gopherTicked && !staleChanged {
		return CommandNone
	}

	d.lastStale = stale
	d.draw(stale)
	if staleChanged && stale {
		return CommandMetricsStale
	}
	if metricsChanged {
		return CommandMetricsChanged
	}
	return CommandNone
}

func (d *Display) draw(stale bool) {
	d.oled.ClearBuffer()
	tinyfont.WriteLine(d.oled, &proggy.TinySZ8pt7b, 4, 8, d.lastLayerLabel, white)
	d.metrics.draw(d.oled, stale)
	d.gopher.draw(d.oled)
	d.oled.Display()
}
