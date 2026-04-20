package encoder

import "machine/usb/hid/keyboard"

func DispatchVolume(cmd Command) {
	if cmd == None {
		return
	}
	kb := keyboard.Port()
	switch cmd {
	case VolumeUp:
		kb.Press(keyboard.KeyMediaVolumeInc)
	case VolumeDown:
		kb.Press(keyboard.KeyMediaVolumeDec)
	case Mute:
		kb.Press(keyboard.KeyMediaMute)
	}
}
