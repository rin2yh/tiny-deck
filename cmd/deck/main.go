package main

import (
	"machine"
	"time"

	"github.com/rin2yh/tiny-deck/internal/display"
	"github.com/rin2yh/tiny-deck/internal/driver"
	"github.com/rin2yh/tiny-deck/internal/joystick"
	"github.com/rin2yh/tiny-deck/internal/keyboard"
)

func main() {
	machine.InitADC()
	disp := display.NewDevice()

	ws := driver.NewWS2812B(machine.GPIO1)
	scanner := keyboard.Setup()
	keyboard.StartupAnimation(ws, scanner)

	layers := keyboard.NewLayerState()
	layerKey := keyboard.NewLongPressDetector(int(keyboard.KeyC3R2), 1000*time.Millisecond)

	serial := machine.Serial
	serial.Configure(machine.UARTConfig{})

	enc := keyboard.NewEncoder(machine.GPIO3, machine.GPIO4, machine.GPIO2)
	js := joystick.NewJoystick(false, true)
	dp := display.New(disp, func() string { return layers.Current().String() })

	for {
		if layerKey.Update(scanner.Scan()) {
			layers.Toggle()
		}
		handleDisplayCommand(dp.Update(serial), ws, scanner, &dp)
		keyboard.DispatchVolume(enc.Update())
		joystick.DispatchMouse(js.Update())
		time.Sleep(10 * time.Millisecond)
	}
}

func handleDisplayCommand(cmd display.Command, ws *driver.WS2812B, scanner *keyboard.Scanner, dp *display.Display) {
	switch cmd {
	case display.CommandNotify:
		keyboard.NotificationAnimation(ws, scanner.KeyCount())
	case display.CommandMetricsChanged:
		keyboard.RenderMetrics(ws, scanner.KeyCount(), dp.CPUPercent(), dp.MemPercent())
	}
}
