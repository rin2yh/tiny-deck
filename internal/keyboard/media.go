package keyboard

import "machine/usb/hid/keyboard"

func DispatchVolume(cmd EncoderCommand) {
	if cmd == EncoderNone {
		return
	}
	kb := keyboard.Port()
	switch cmd {
	case EncoderVolumeUp:
		kb.Press(keyboard.KeyMediaVolumeInc)
	case EncoderVolumeDown:
		kb.Press(keyboard.KeyMediaVolumeDec)
	case EncoderMute:
		kb.Press(keyboard.KeyMediaMute)
	}
}
