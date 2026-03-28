package main

import (
	"image/color"
	"machine"
	"time"

	"github.com/rin2yh/tiny-deck/internal/driver"
	"github.com/rin2yh/tiny-deck/internal/keyboard"
	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

type LineBuffer struct {
	data [128]byte
	n    int
}

func (b *LineBuffer) Add(c byte) (completed bool) {
	if c == '\n' {
		return true
	}
	if b.n < len(b.data) {
		b.data[b.n] = c
		b.n++
	}
	return false
}

func (b *LineBuffer) String() string {
	return string(b.data[:b.n])
}

func (b *LineBuffer) Reset() {
	b.n = 0
}

type SerialReader struct {
	buf LineBuffer
}

func (r *SerialReader) Read(serial machine.Serialer) (completed bool) {
	if serial.Buffered() == 0 {
		return false
	}

	b, _ := serial.ReadByte()
	return r.buf.Add(b)
}

func (r *SerialReader) Line() string {
	s := r.buf.String()
	r.buf.Reset()
	return s
}

func main() {
	// I2C (OLED)
	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: 400 * machine.KHz,
		SDA:       machine.GPIO12,
		SCL:       machine.GPIO13,
	})

	display := ssd1306.NewI2C(machine.I2C0)
	display.Configure(ssd1306.Config{
		Address:  0x3C,
		Width:    128,
		Height:   64,
		Rotation: ssd1306.ROTATION_180,
	})
	display.ClearDisplay()

	ws := driver.NewWS2812B(machine.GPIO1)
	matrix := keyboard.Setup()
	keyboard.StartupAnimation(ws, matrix)

	// USB Serial
	serial := machine.Serial
	serial.Configure(machine.UARTConfig{})

	var line string
	white := color.RGBA{255, 255, 255, 255}

	serialReader := SerialReader{}
	for {
		if serialReader.Read(serial) {
			line = serialReader.Line()
			display.ClearBuffer()
			tinyfont.WriteLine(
				display,
				&freemono.Regular9pt7b,
				10,
				40,
				line,
				white,
			)

			display.Display()
		}

		time.Sleep(10 * time.Millisecond)
	}
}
