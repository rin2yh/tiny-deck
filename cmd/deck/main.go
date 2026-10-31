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
		jsEvent := js.Update()
		drawDisplay := layers.Current() != layer.Game
		cmd, woke := dp.Update(serial, drawDisplay)
		if drawDisplay {
			handleDisplayCommand(cmd, ws, &dp)
		} else if woke {
			games.Refresh()
		}

		if anyKeyPressed(keys) || ev.Delta != 0 || ev.Pressed || joystickMoved(jsEvent) {
			if dp.Activity(time.Now()) {
				if layers.Current() == layer.Game {
					games.Refresh()
				} else {
					dp.Redraw()
				}
			}
		}

		if layers.Current() == layer.Game {
			// ゲームが OLED を占有する間は音量・スクロール・LED を止める。
			// 指標だけは取り込んで、戻ったときに最新の値が出るようにする。
			games.Update(jumpPressed(keys), ev.Delta, ev.Pressed)
		} else {
			encoder.DispatchVolume(ev)
			scroller.Dispatch(jsEvent)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func anyKeyPressed(keys []bool) bool {
	for _, pressed := range keys {
		if pressed {
			return true
		}
	}
	return false
}

func joystickMoved(ev joystick.Event) bool {
	return ev.DX != 0 || ev.DY != 0 || ev.Click
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
