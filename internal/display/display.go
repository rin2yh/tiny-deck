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
	oled            *ssd1306.Device
	reader          serial.Reader
	metrics         metricsState
	gopher          gopherAnim
	layerName       func() string
	startedAt       time.Time
	lastActivity    time.Time
	powerState      powerState
	hostUnavailable bool
	lastLayer       string
	lastLayerLabel  string
	lastStale       bool
}

func New(oled *ssd1306.Device, layerName func() string) Display {
	now := time.Now()
	return Display{
		oled:         oled,
		layerName:    layerName,
		gopher:       newGopherAnim(),
		startedAt:    now,
		lastActivity: now,
	}
}

func (d *Display) CPUPercent() float32 { return d.metrics.cpu }
func (d *Display) MemPercent() float32 { return d.metrics.mem }

// Update はホストからの行を取り込み、必要なら OLED を描き直す。
// draw が false のときは描画だけを止める。OLED をゲームに明け渡している間も
// シリアルを詰まらせず、指標を最新に保つために使う。
func (d *Display) Update(s machine.Serialer, draw bool) (Command, bool) {
	now := time.Now()
	gopherTicked := d.gopher.tick(now)
	currentLayer := d.layerName()
	layerChanged := currentLayer != d.lastLayer
	metricsChanged := false

	for d.reader.Read(s) {
		l := d.reader.Line()
		if l == "notify" {
			return CommandNotify, false
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
	unavailable := stale && (d.metrics.hasData || now.Sub(d.startedAt) >= metricsTimeout)
	woke := d.updatePower(now, unavailable)

	if !metricsChanged && !layerChanged && !gopherTicked && !staleChanged && !woke {
		return CommandNone, false
	}

	d.lastStale = stale
	if draw && d.powerState != powerOff {
		d.draw(stale)
	}
	if staleChanged && stale {
		return CommandMetricsStale, woke
	}
	if metricsChanged {
		return CommandMetricsChanged, woke
	}
	return CommandNone, woke
}

// Activity records a physical input and restores normal brightness. The return
// value reports whether the OLED was off and needs its current screen redrawn.
func (d *Display) Activity(now time.Time) bool {
	d.lastActivity = now
	return d.setPowerState(powerNormal)
}

// Redraw refreshes the regular status screen after the OLED has been turned on.
func (d *Display) Redraw() { d.draw(d.metrics.isStale(time.Now())) }

func (d *Display) draw(stale bool) {
	d.oled.ClearBuffer()
	tinyfont.WriteLine(d.oled, &proggy.TinySZ8pt7b, 4, 8, d.lastLayerLabel, white)
	d.metrics.draw(d.oled, stale)
	d.gopher.draw(d.oled)
	d.oled.Display()
}
