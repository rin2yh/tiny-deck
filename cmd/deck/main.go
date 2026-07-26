package main

import (
	"machine"
	"time"

	"github.com/rin2yh/tiny-deck/internal/arcade"
	"github.com/rin2yh/tiny-deck/internal/display"
	"github.com/rin2yh/tiny-deck/internal/driver"
	"github.com/rin2yh/tiny-deck/internal/joystick"
	"github.com/rin2yh/tiny-deck/internal/keyboard"
	"github.com/rin2yh/tiny-deck/internal/keyboard/encoder"
	"github.com/rin2yh/tiny-deck/internal/keyboard/layer"
	"github.com/rin2yh/tiny-deck/internal/keyboard/led"
)

// layerKeyIndex はレイヤー切り替え用の長押しキー。ゲーム中もここから抜けられる
// 必要があるので、ジャンプ入力からは除外する。
const layerKeyIndex = int(keyboard.KeyC3R2)

func main() {
	machine.InitADC()
	disp := display.NewDevice()

	ws := driver.NewWS2812B(machine.GPIO1)
	scanner := keyboard.Setup()
	led.StartupAnimation(ws)

	layers := layer.New()
	layerKey := keyboard.NewLongPressDetector(layerKeyIndex, 1000*time.Millisecond)

	serial := machine.Serial
	serial.Configure(machine.UARTConfig{})

	enc := encoder.New(machine.GPIO3, machine.GPIO4, machine.GPIO2)
	js := joystick.NewJoystick(false, true)
	scroller := joystick.NewScroller()
	dp := display.New(disp, func() string { return layers.Current().String() })
	games := arcade.NewLauncher(disp)

	for {
		keys := scanner.Scan()
		if layerKey.Update(keys) {
			layers.Next()
			if layers.Current() == layer.Game {
				games.Enter()
				led.ClearMetrics(ws)
			} else {
				games.Leave()
			}
		}

		ev := enc.Update()
		if layers.Current() == layer.Game {
			// ゲームが OLED を占有する間は音量・スクロール・LED を止める。
			// 指標だけは取り込んで、戻ったときに最新の値が出るようにする。
			dp.Update(serial, false)
			games.Update(jumpPressed(keys), ev.Delta, ev.Pressed)
		} else {
			handleDisplayCommand(dp.Update(serial, true), ws, &dp)
			encoder.DispatchVolume(ev)
			scroller.Dispatch(js.Update())
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// jumpPressed はレイヤーキー以外のいずれかが押されているかを返す。
func jumpPressed(keys []bool) bool {
	for i, pressed := range keys {
		if pressed && i != layerKeyIndex {
			return true
		}
	}
	return false
}

func handleDisplayCommand(cmd display.Command, ws *driver.WS2812B, dp *display.Display) {
	switch cmd {
	case display.CommandNotify:
		led.NotificationAnimation(ws)
	case display.CommandMetricsChanged:
		led.RenderMetrics(ws, dp.CPUPercent(), dp.MemPercent())
	case display.CommandMetricsStale:
		led.ClearMetrics(ws)
	}
}
