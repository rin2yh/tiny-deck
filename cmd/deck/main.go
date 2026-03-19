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
	colPins := []machine.Pin{
		machine.GPIO5,
		machine.GPIO6,
		machine.GPIO7,
		machine.GPIO8,
	}

	rowPins := []machine.Pin{
		machine.GPIO9,
		machine.GPIO10,
		machine.GPIO11,
	}

	for _, c := range colPins {
		c.Configure(machine.PinConfig{Mode: machine.PinOutput})
		c.Low()
	}

	for _, c := range rowPins {
		c.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	}

	keyboard.StartupAnimation(ws, len(rowPins)*len(colPins))

	// USB Serial
	serial := machine.Serial
	serial.Configure(machine.UARTConfig{})

	line := ""
	white := color.RGBA{255, 255, 255, 255}

	for {
		if serial.Buffered() > 0 {
			b, _ := serial.ReadByte()

			if b == '\n' {
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
				line = ""
			} else {
				line += string(b)
			}
		}

		time.Sleep(10 * time.Millisecond)
	}
}
