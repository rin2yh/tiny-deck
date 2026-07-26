package encoder

import (
	"machine"

	"tinygo.org/x/drivers/encoders"
)

// Event は 1 回のポーリングで観測した操作。用途を決めるのは呼び出し側で、
// 通常レイヤーでは音量、ゲームレイヤーではメニュー操作に割り当てる。
type Event struct {
	Delta   int  // ノッチ単位の回転量（時計回りが正）
	Pressed bool // 押し込みの押下エッジ
}

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

func (e *Encoder) Update() Event {
	var ev Event

	// PullUp接続のため反転
	pressed := !e.swPin.Get()
	ev.Pressed = pressed && !e.swLast
	e.swLast = pressed

	// Position()はrawValue/Precisionを返すため1ノッチ=delta1
	pos := e.enc.Position()
	if pos != e.lastPos {
		ev.Delta = pos - e.lastPos
		e.lastPos = pos
	}
	return ev
}
