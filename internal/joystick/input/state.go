package input

type Event struct {
	DX, DY      int
	Click       bool
	Release     bool
	axisChanged bool
}

// HasActivity reports whether the joystick position changed or it was clicked.
func (e Event) HasActivity() bool {
	return e.axisChanged || e.Click
}

type State struct {
	lastDX, lastDY int
	buttonPressed  bool
}

func (s *State) Update(dx, dy int, pressed bool) Event {
	event := Event{
		DX:          dx,
		DY:          dy,
		Click:       pressed && !s.buttonPressed,
		Release:     !pressed && s.buttonPressed,
		axisChanged: dx != s.lastDX || dy != s.lastDY,
	}
	s.buttonPressed = pressed
	s.lastDX = dx
	s.lastDY = dy
	return event
}
