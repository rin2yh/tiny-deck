package keyboard

import "time"

type LongPressDetector struct {
	keyIndex    int
	holdTime    time.Duration
	pressStart  time.Time
	pressActive bool
	triggered   bool
}

func NewLongPressDetector(keyIndex int, holdTime time.Duration) *LongPressDetector {
	return &LongPressDetector{
		keyIndex: keyIndex,
		holdTime: holdTime,
	}
}

// Update は長押し閾値に達した瞬間だけ true を返す（その後はキーを離すまで false）
func (d *LongPressDetector) Update(keys []bool) bool {
	pressed := d.keyIndex < len(keys) && keys[d.keyIndex]

	if !pressed {
		d.pressActive = false
		d.triggered = false
		return false
	}

	if !d.pressActive {
		d.pressActive = true
		d.pressStart = time.Now()
		d.triggered = false
		return false
	}

	if d.triggered {
		return false
	}

	if time.Since(d.pressStart) >= d.holdTime {
		d.triggered = true
		return true
	}

	return false
}
