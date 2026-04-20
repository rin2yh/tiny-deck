package encoder

import (
	"machine"

	"tinygo.org/x/drivers/encoders"
)

type Command uint8

const (
	None Command = iota
	VolumeUp
	VolumeDown
	Mute
)

type Encoder struct {
	enc     *encoders.QuadratureDevice
	swPin   machine.Pin
	lastPos int
	swLast  bool
}

func New(pinA, pinB, pinSW machine.Pin) *Encoder {
	enc := encoders.NewQuadratureViaInterrupt(pinA, pinB)
	// Precision=4: 4パルス/ノッチ。Position()は内部でPrecisionで割るため1ノッチ=1
	enc.Configure(encoders.QuadratureConfig{Precision: 4})
	pinSW.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	return &Encoder{enc: enc, swPin: pinSW}
}

func (e *Encoder) Update() Command {
	// PullUp接続のため反転
	pressed := !e.swPin.Get()
	if pressed && !e.swLast {
		e.swLast = true
		return Mute
	}
	if !pressed {
		e.swLast = false
	}

	// Position()はrawValue/Precisionを返すため1ノッチ=delta1
	pos := e.enc.Position()
	delta := pos - e.lastPos
	if delta >= 1 {
		e.lastPos = pos
		return VolumeUp
	}
	if delta <= -1 {
		e.lastPos = pos
		return VolumeDown
	}
	return None
}
