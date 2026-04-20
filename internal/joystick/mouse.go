package joystick

import "machine/usb/hid/mouse"

func DispatchMouse(ev Event) {
	if ev.DX == 0 && ev.DY == 0 && !ev.Click && !ev.Release {
		return
	}
	ms := mouse.Port()
	if ev.DX != 0 || ev.DY != 0 {
		ms.Move(ev.DX, ev.DY)
	}
	if ev.Click {
		ms.Press(mouse.Left)
	}
	if ev.Release {
		ms.Release(mouse.Left)
	}
}
