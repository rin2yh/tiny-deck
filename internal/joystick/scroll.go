package joystick

import "machine/usb/hid/mouse"

const (
	scrollTiltThreshold   = 100
	autoScrollPeriodTicks = 30
)

type Scroller struct {
	accumY     int
	autoScroll bool
	autoTicks  int
}

func NewScroller() *Scroller { return &Scroller{} }

func (s *Scroller) Dispatch(ev Event) {
	if ev.Click {
		s.autoScroll = !s.autoScroll
		s.autoTicks = 0
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

	if s.autoScroll {
		s.autoTicks++
		if s.autoTicks >= autoScrollPeriodTicks {
			s.autoTicks = 0
			ms.WheelDown()
		}
	}
}
