package main

import (
	"machine"
	keyboardpkg "machine/usb/hid/keyboard"
	"machine/usb/hid/mouse"
	"time"

	"github.com/rin2yh/tiny-deck/internal/display"
	"github.com/rin2yh/tiny-deck/internal/driver"
	"github.com/rin2yh/tiny-deck/internal/joystick"
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
	machine.InitADC()

	disp := ssd1306.NewI2C(machine.I2C0)
	disp.Configure(ssd1306.Config{
		Address:  0x3C,
		Width:    128,
		Height:   64,
		Rotation: ssd1306.ROTATION_180,
	})
	disp.ClearDisplay()

	ws := driver.NewWS2812B(machine.GPIO1)
	scanner := keyboard.Setup()
	keyboard.StartupAnimation(ws, scanner)

	layers := keyboard.NewLayerState()
	layerKey := keyboard.NewLongPressDetector(int(keyboard.KeyC3R2), 1000*time.Millisecond)

	// USB Serial
	serial := machine.Serial
	serial.Configure(machine.UARTConfig{})

	enc := keyboard.NewEncoder(machine.GPIO3, machine.GPIO4, machine.GPIO2)
	ms := mouse.Port()
	kb := keyboardpkg.Port()
	js := joystick.NewJoystick(false, true)

	dp := display.New(disp, func() string { return layers.Current().String() })
	for {
		if layerKey.Update(scanner.Scan()) {
			layers.Toggle()
		}
		cmd := dp.Update(serial)
		switch cmd {
		case display.CommandNotify:
			keyboard.NotificationAnimation(ws, scanner.KeyCount())
		case display.CommandMetricsChanged:
			keyboard.RenderMetrics(ws, scanner.KeyCount(), dp.CPUPercent(), dp.MemPercent())
		}
		switch enc.Update() {
		case keyboard.EncoderVolumeUp:
			kb.Press(keyboardpkg.KeyMediaVolumeInc)
		case keyboard.EncoderVolumeDown:
			kb.Press(keyboardpkg.KeyMediaVolumeDec)
		case keyboard.EncoderMute:
			kb.Press(keyboardpkg.KeyMediaMute)
		}
		ev := js.Update()
		if ev.DX != 0 || ev.DY != 0 {
			ms.Move(ev.DX, ev.DY)
		}
		if ev.Click {
			ms.Press(mouse.Left)
		}
		if ev.Release {
			ms.Release(mouse.Left)
		}
		time.Sleep(10 * time.Millisecond)
	}
}
