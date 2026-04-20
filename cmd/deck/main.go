package main

import (
	"machine"
	"time"

	"github.com/rin2yh/tiny-deck/internal/display"
	"github.com/rin2yh/tiny-deck/internal/driver"
	"github.com/rin2yh/tiny-deck/internal/joystick"
	"github.com/rin2yh/tiny-deck/internal/keyboard"
	"github.com/rin2yh/tiny-deck/internal/keyboard/encoder"
	"github.com/rin2yh/tiny-deck/internal/keyboard/layer"
	"github.com/rin2yh/tiny-deck/internal/keyboard/led"
)

func main() {
	machine.InitADC()
	disp := display.NewDevice()

	ws := driver.NewWS2812B(machine.GPIO1)
	scanner := keyboard.Setup()
	led.StartupAnimation(ws)

	layers := layer.New()
	layerKey := keyboard.NewLongPressDetector(int(keyboard.KeyC3R2), 1000*time.Millisecond)

	serial := machine.Serial
	serial.Configure(machine.UARTConfig{})

	enc := encoder.New(machine.GPIO3, machine.GPIO4, machine.GPIO2)
	js := joystick.NewJoystick(false, true)
	dp := display.New(disp, func() string { return layers.Current().String() })

	for {
		if layerKey.Update(scanner.Scan()) {
			layers.Toggle()
		}
		handleDisplayCommand(dp.Update(serial), ws, &dp)
		encoder.DispatchVolume(enc.Update())
		joystick.DispatchMouse(js.Update())
		time.Sleep(10 * time.Millisecond)
	}
}

func handleDisplayCommand(cmd display.Command, ws *driver.WS2812B, dp *display.Display) {
	switch cmd {
	case display.CommandNotify:
		led.NotificationAnimation(ws)
	case display.CommandMetricsChanged:
		led.RenderMetrics(ws, dp.CPUPercent(), dp.MemPercent())
	}
}
