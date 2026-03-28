package display

import (
	"image/color"
	"machine"

	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
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

type Metrics struct {
	display *ssd1306.Device
	reader  serialReader
}

func NewMetrics(d *ssd1306.Device) Metrics {
	return Metrics{display: d}
}

func (m *Metrics) Update(serial machine.Serialer) {
	if m.reader.read(serial) {
		m.display.ClearBuffer()
		tinyfont.WriteLine(m.display, &freemono.Regular9pt7b, 10, 40, m.reader.line(), white)
		m.display.Display()
	}
}
