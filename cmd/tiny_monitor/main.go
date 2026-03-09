package main

import (
	"image/color"
	"machine"
	"time"

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
					"cpu:"+line+"%",
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
