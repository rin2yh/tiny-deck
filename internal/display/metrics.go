package display

import (
	"image/color"
	"machine"
	"strconv"
	"strings"
	"time"

	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
	"tinygo.org/x/tinyfont/gophers"
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
	gopherBaseX    = 104
	gopherBaseY    = 28
	gopherInterval = 220 * time.Millisecond
)

// 全身がちゃんと映る 32pt glyph のみ採用（目だけ／部分表示になる低背 glyph は除外）
var gopherPoses = [...]string{"A", "B", "W", "Y"}

var gopherOffsets = [4]struct{ dx, dy int16 }{
	{dx: +2, dy: 0},
	{dx: -2, dy: 0},
	{dx: +2, dy: -1},
	{dx: -2, dy: -1},
}

type Metrics struct {
	display        *ssd1306.Device
	reader         serialReader
	layerName      func() string
	lastLayer      string
	lastLayerLabel string
	lastCPU        string
	lastMem        string
	lastCPUVal     float32
	lastMemVal     float32
	gopherPhase    uint8
	lastGopherAt   time.Time
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
	return Metrics{display: d, layerName: layerName, lastGopherAt: time.Now()}
}

func (m *Metrics) tickGopher(now time.Time) bool {
	if now.Sub(m.lastGopherAt) < gopherInterval {
		return false
	}
	m.gopherPhase++
	m.lastGopherAt = now
	return true
}

func (m *Metrics) drawMetrics() {
	m.display.ClearBuffer()
	tinyfont.WriteLine(m.display, &proggy.TinySZ8pt7b, 4, 8, m.lastLayerLabel, white)
	tinyfont.WriteLine(m.display, &freemono.Regular9pt7b, 4, 40, m.lastCPU, white)
	tinyfont.WriteLine(m.display, &freemono.Regular9pt7b, 4, 58, m.lastMem, white)
	phase := m.gopherPhase & 0x03
	off := gopherOffsets[phase]
	tinyfont.WriteLine(m.display, &gophers.Regular32pt, gopherBaseX+off.dx, gopherBaseY+off.dy, gopherPoses[phase], white)
	m.display.Display()
}

func (m *Metrics) Update(serial machine.Serialer) Command {
	gopherTicked := m.tickGopher(time.Now())
	currentLayer := m.layerName()
	layerChanged := currentLayer != m.lastLayer
	metricsChanged := false

	for m.reader.read(serial) {
		line := m.reader.line()
		if line == "notify" {
			return CommandNotify
		}
		switch {
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

	if !metricsChanged && !layerChanged && !gopherTicked {
		return CommandNone
	}

	m.drawMetrics()
	if metricsChanged {
		return CommandMetricsChanged
	}
	return CommandNone
}
