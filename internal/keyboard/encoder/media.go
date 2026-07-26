package encoder

import "machine/usb/hid/keyboard"

func DispatchVolume(ev Event) {
	if ev == (Event{}) {
		return
	}
	kb := keyboard.Port()
	if ev.Pressed {
		kb.Press(keyboard.KeyMediaMute)
	}
	// 1 回のポーリングで複数ノッチ回っていても 1 段だけ動かす
	switch {
	case ev.Delta > 0:
		kb.Press(keyboard.KeyMediaVolumeInc)
	case ev.Delta < 0:
		kb.Press(keyboard.KeyMediaVolumeDec)
	}
}
