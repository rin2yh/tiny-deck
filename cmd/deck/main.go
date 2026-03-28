package main

import (
	"machine"
	"time"

	"github.com/rin2yh/tiny-deck/internal/display"
	"github.com/rin2yh/tiny-deck/internal/driver"
	"github.com/rin2yh/tiny-deck/internal/keyboard"
	"tinygo.org/x/drivers/ssd1306"
)

func main() {
	// I2C (OLED)
	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: 400 * machine.KHz,
		SDA:       machine.GPIO12,
		SCL:       machine.GPIO13,
	})

	disp := ssd1306.NewI2C(machine.I2C0)
	disp.Configure(ssd1306.Config{
		Address:  0x3C,
		Width:    128,
		Height:   64,
		Rotation: ssd1306.ROTATION_180,
	})
	disp.ClearDisplay()

	ws := driver.NewWS2812B(machine.GPIO1)
	matrix := keyboard.Setup()
	keyboard.StartupAnimation(ws, matrix)

	// USB Serial
	serial := machine.Serial
	serial.Configure(machine.UARTConfig{})

	m := display.NewMetrics(disp)
	for {
		m.Update(serial)
		time.Sleep(10 * time.Millisecond)
	}
}
