package main

import (
	"machine"
	"machine/usb/hid/mouse"
	"time"

	"github.com/rin2yh/tiny-deck/internal/display"
	"github.com/rin2yh/tiny-deck/internal/driver"
	"github.com/rin2yh/tiny-deck/internal/joystick"
	"github.com/rin2yh/tiny-deck/internal/keyboard"
	"github.com/rin2yh/tiny-deck/internal/volume"
	"tinygo.org/x/drivers/ssd1306"
)

// TinyGoでヒープ再確保を避けるためパッケージレベルで確保
var (
	serialVolUp   = []byte(volume.CmdUp + "\n")
	serialVolDown = []byte(volume.CmdDown + "\n")
	serialMute    = []byte(volume.CmdMute + "\n")
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
	layerKey := keyboard.NewLongPressDetector(11, 1000*time.Millisecond)

	// USB Serial
	serial := machine.Serial
	serial.Configure(machine.UARTConfig{})

	enc := keyboard.NewEncoder(machine.GPIO3, machine.GPIO4, machine.GPIO2)
	ms := mouse.Port()
	js := joystick.NewJoystick(false, true)

	m := display.NewMetrics(disp, func() string { return layers.Current().String() })
	for {
		if layerKey.Update(scanner.Scan()) {
			layers.Toggle()
		}
		cmd := m.Update(serial)
		if cmd == display.CommandNotify {
			keyboard.NotificationAnimation(ws, scanner.KeyCount())
		}
		switch enc.Update() {
		case keyboard.EncoderVolumeUp:
			serial.Write(serialVolUp)
		case keyboard.EncoderVolumeDown:
			serial.Write(serialVolDown)
		case keyboard.EncoderMute:
			serial.Write(serialMute)
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
