package display

import (
	"machine"

	"tinygo.org/x/drivers/ssd1306"
)

const (
	oledI2CFrequency = 400 * machine.KHz
	oledAddress      = 0x3C
	oledWidth        = 128
	oledHeight       = 64
)

func NewDevice() *ssd1306.Device {
	machine.I2C0.Configure(machine.I2CConfig{
		Frequency: oledI2CFrequency,
		SDA:       machine.GPIO12,
		SCL:       machine.GPIO13,
	})
	dev := ssd1306.NewI2C(machine.I2C0)
	dev.Configure(ssd1306.Config{
		Address:  oledAddress,
		Width:    oledWidth,
		Height:   oledHeight,
		Rotation: ssd1306.ROTATION_180,
	})
	dev.ClearDisplay()
	return dev
}
