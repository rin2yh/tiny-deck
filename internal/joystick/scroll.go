package joystick

import "machine/usb/hid/mouse"

const (
	scrollTiltThreshold   = 100
	autoScrollPeriodTicks = 30
)

type Scroller struct {
	accumY    int
	pressed   bool
	holdTicks int
}

func NewScroller() *Scroller { return &Scroller{} }

func (s *Scroller) Dispatch(ev Event) {
	if ev.Click {
		s.pressed = true
		s.holdTicks = 0
	}
	if ev.Release {
		s.pressed = false
	}

	ms := mouse.Port()

	if ev.DY != 0 {
		s.accumY += ev.DY
		for s.accumY >= scrollTiltThreshold {
			ms.WheelDown()
			s.accumY -= scrollTiltThreshold
		}
		for s.accumY <= -scrollTiltThreshold {
			ms.WheelUp()
			s.accumY += scrollTiltThreshold
		}
		return
	}
	s.accumY = 0

	if s.pressed {
		s.holdTicks++
		if s.holdTicks >= autoScrollPeriodTicks {
			s.holdTicks = 0
			ms.WheelDown()
		}
	}
}
