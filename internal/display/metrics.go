package display

import (
	"image/color"
	"machine"
	"strconv"
	"strings"
	"time"

	"github.com/rin2yh/tiny-deck/internal/volume"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
	"tinygo.org/x/tinyfont/proggy"
)

var white = color.RGBA{R: 255, G: 255, B: 255, A: 255}

type lineBuffer struct {
	data [128]byte
	n    int
}

func (b *lineBuffer) add(c byte) bool {
	if c == '\n' {
		return true
	}
	if b.n < len(b.data) {
		b.data[b.n] = c
		b.n++
	}
	return false
}

func (b *lineBuffer) string() string {
	return string(b.data[:b.n])
}

func (b *lineBuffer) reset() {
	b.n = 0
}

type serialReader struct {
	buf lineBuffer
}

func (r *serialReader) read(serial machine.Serialer) bool {
	for serial.Buffered() > 0 {
		b, _ := serial.ReadByte()
		if r.buf.add(b) {
			return true
		}
	}
	return false
}

func (r *serialReader) line() string {
	s := r.buf.string()
	r.buf.reset()
	return s
}

type Command int

const (
	CommandNone Command = iota
	CommandNotify
	CommandMetricsChanged
)

const (
	prefixCPU = "cpu:"
	prefixMem = "mem:"
)

const (
	volOverlayDuration = 1500 * time.Millisecond
	volumeStep         = 5
	volumeMax          = 100

	volBarOuterX = 3
	volBarOuterY = 27
	volBarOuterW = 122
	volBarOuterH = 12
	volBarInnerX = 5
	volBarInnerY = 29
	volBarInnerW = 118
	volBarInnerH = 8
)

type VolPendingKind uint8

const (
	VolPendingUp VolPendingKind = iota
	VolPendingDown
	VolPendingMute
)

type Metrics struct {
	display         *ssd1306.Device
	reader          serialReader
	layerName       func() string
	lastLayer       string
	lastLayerLabel  string
	lastCPU         string
	lastMem         string
	lastCPUVal      float32
	lastMemVal      float32
	volOverlayUntil time.Time
	volValue        int8
	volMuted        bool
	volDirty        bool
}

func (m *Metrics) CPUPercent() float32 { return m.lastCPUVal }
func (m *Metrics) MemPercent() float32 { return m.lastMemVal }

func parsePercent(line, prefix string) (float32, bool) {
	s := strings.TrimPrefix(line, prefix)
	s = strings.TrimSuffix(s, "%")
	s = strings.TrimSpace(s)
	v, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, false
	}
	return float32(v), true
}

func NewMetrics(d *ssd1306.Device, layerName func() string) Metrics {
	return Metrics{display: d, layerName: layerName, volValue: -1}
}

func (m *Metrics) ShowVolumePending(kind VolPendingKind) {
	m.volOverlayUntil = time.Now().Add(volOverlayDuration)
	switch kind {
	case VolPendingUp, VolPendingDown:
		m.volMuted = false
		if m.volValue >= 0 {
			delta := int16(volumeStep)
			if kind == VolPendingDown {
				delta = -delta
			}
			v := int16(m.volValue) + delta
			if v < 0 {
				v = 0
			} else if v > volumeMax {
				v = volumeMax
			}
			m.volValue = int8(v)
		}
	case VolPendingMute:
		m.volMuted = !m.volMuted
	}
	m.volDirty = true
}

func (m *Metrics) extendOverlayIfActive() {
	if m.volOverlayUntil.IsZero() {
		return
	}
	m.volOverlayUntil = time.Now().Add(volOverlayDuration)
	m.volDirty = true
}

func (m *Metrics) drawMetrics() {
	m.display.ClearBuffer()
	tinyfont.WriteLine(m.display, &proggy.TinySZ8pt7b, 4, 8, m.lastLayerLabel, white)
	tinyfont.WriteLine(m.display, &freemono.Regular9pt7b, 4, 40, m.lastCPU, white)
	tinyfont.WriteLine(m.display, &freemono.Regular9pt7b, 4, 58, m.lastMem, white)
	m.display.Display()
}

func (m *Metrics) drawVolumeOverlay() {
	m.display.ClearBuffer()
	label := "VOL"
	if m.volMuted {
		label = "MUTE"
	}
	tinyfont.WriteLine(m.display, &proggy.TinySZ8pt7b, 4, 10, label, white)
	if !m.volMuted {
		m.display.FillRectangle(volBarOuterX, volBarOuterY, volBarOuterW, 1, white)
		m.display.FillRectangle(volBarOuterX, volBarOuterY+volBarOuterH-1, volBarOuterW, 1, white)
		m.display.FillRectangle(volBarOuterX, volBarOuterY, 1, volBarOuterH, white)
		m.display.FillRectangle(volBarOuterX+volBarOuterW-1, volBarOuterY, 1, volBarOuterH, white)
		if m.volValue > 0 {
			w := int16(volBarInnerW) * int16(m.volValue) / volumeMax
			m.display.FillRectangle(volBarInnerX, volBarInnerY, w, volBarInnerH, white)
		}
	}
	m.display.Display()
}

func (m *Metrics) Update(serial machine.Serialer) Command {
	currentLayer := m.layerName()
	layerChanged := currentLayer != m.lastLayer
	metricsChanged := false

	for m.reader.read(serial) {
		line := m.reader.line()
		if line == "notify" {
			return CommandNotify
		}
		switch {
		case strings.HasPrefix(line, volume.PrefixCurrent):
			s := strings.TrimSpace(strings.TrimPrefix(line, volume.PrefixCurrent))
			if v, err := strconv.Atoi(s); err == nil {
				if v < 0 {
					v = 0
				} else if v > volumeMax {
					v = volumeMax
				}
				m.volValue = int8(v)
				m.extendOverlayIfActive()
			}
		case strings.HasPrefix(line, volume.PrefixMuted):
			s := strings.TrimSpace(strings.TrimPrefix(line, volume.PrefixMuted))
			m.volMuted = s == "1"
			m.extendOverlayIfActive()
		case strings.HasPrefix(line, prefixCPU):
			if line != m.lastCPU {
				m.lastCPU = line
				if v, ok := parsePercent(line, prefixCPU); ok {
					m.lastCPUVal = v
					metricsChanged = true
				}
			}
		case strings.HasPrefix(line, prefixMem):
			if line != m.lastMem {
				m.lastMem = line
				if v, ok := parsePercent(line, prefixMem); ok {
					m.lastMemVal = v
					metricsChanged = true
				}
			}
		}
	}

	if layerChanged {
		m.lastLayer = currentLayer
		m.lastLayerLabel = "[" + currentLayer + "]"
	}

	if !m.volOverlayUntil.IsZero() {
		if time.Now().After(m.volOverlayUntil) {
			m.volOverlayUntil = time.Time{}
			m.volDirty = false
			m.drawMetrics()
			return CommandNone
		}
		if m.volDirty {
			m.drawVolumeOverlay()
			m.volDirty = false
		}
		return CommandNone
	}

	if !metricsChanged && !layerChanged {
		return CommandNone
	}

	m.drawMetrics()
	if metricsChanged {
		return CommandMetricsChanged
	}
	return CommandNone
}
