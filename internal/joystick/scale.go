package joystick

import "math"

func scaleAxis(raw, center, deadzone, divisor int, invert bool) int {
	d := raw - center
	switch {
	case d > deadzone:
		d -= deadzone
	case d < -deadzone:
		d += deadzone
	default:
		return 0
	}
	v := d / divisor
	if invert {
		v = -v
	}
	if v > math.MaxInt8 {
		v = math.MaxInt8
	} else if v < math.MinInt8 {
		v = math.MinInt8
	}
	return v
}
