package display

import (
	"image/color"
	"machine"

	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinyfont"
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
	if serial.Buffered() == 0 {
		return false
	}
	b, _ := serial.ReadByte()
	return r.buf.add(b)
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
)

type Metrics struct {
	display        *ssd1306.Device
	reader         serialReader
	layerName      func() string
	lastLayer      string
	lastLayerLabel string
	lastLine       string
}

func NewMetrics(d *ssd1306.Device, layerName func() string) Metrics {
	return Metrics{display: d, layerName: layerName}
}

func (m *Metrics) Update(serial machine.Serialer) Command {
	lineChanged := m.reader.read(serial)
	currentLayer := m.layerName()
	layerChanged := currentLayer != m.lastLayer

	if lineChanged {
		line := m.reader.line()
		if line == "notify" {
			return CommandNotify
		}
		m.lastLine = line
	}
	if layerChanged {
		m.lastLayer = currentLayer
		m.lastLayerLabel = "[" + currentLayer + "]"
	}

	if !lineChanged && !layerChanged {
		return CommandNone
	}

	m.display.ClearBuffer()
	tinyfont.WriteLine(m.display, &proggy.TinySZ8pt7b, 10, 8, m.lastLayerLabel, white)
	tinyfont.WriteLine(m.display, &proggy.TinySZ8pt7b, 10, 40, m.lastLine, white)
	m.display.Display()
	return CommandNone
}
