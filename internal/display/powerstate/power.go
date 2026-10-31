package powerstate

import "time"

const (
	dimAfter = time.Minute
	offAfter = 3 * time.Minute
)

type State uint8

const (
	Normal State = iota
	Dim
	Off
)

type Transition struct {
	State   State
	Changed bool
	Woke    bool
}

// Controller tracks OLED idle and host availability without hardware access.
type Controller struct {
	lastActivity    time.Time
	state           State
	hostUnavailable bool
}

func New(now time.Time) Controller {
	return Controller{lastActivity: now}
}

func (c *Controller) State() State { return c.state }

func (c *Controller) Activity(now time.Time) Transition {
	c.lastActivity = now
	return c.setState(Normal)
}

func (c *Controller) Update(now time.Time, hostUnavailable bool) Transition {
	if hostUnavailable && !c.hostUnavailable {
		c.hostUnavailable = true
		return c.setState(Off)
	}
	if !hostUnavailable && c.hostUnavailable {
		c.hostUnavailable = false
		c.lastActivity = now
		return c.setState(Normal)
	}

	switch c.state {
	case Normal:
		if now.Sub(c.lastActivity) >= dimAfter {
			return c.setState(Dim)
		}
	case Dim:
		if now.Sub(c.lastActivity) >= offAfter {
			return c.setState(Off)
		}
	}
	return Transition{State: c.state}
}

func (c *Controller) setState(next State) Transition {
	if c.state == next {
		return Transition{State: c.state}
	}
	wasOff := c.state == Off
	c.state = next
	return Transition{State: next, Changed: true, Woke: wasOff && next == Normal}
}
