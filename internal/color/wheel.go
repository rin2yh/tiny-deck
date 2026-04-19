package color

func Wheel(pos uint8) (uint8, uint8, uint8) {
	pos = 255 - pos
	switch {
	case pos < 85:
		return 255 - pos*3, 0, pos * 3
	case pos < 170:
		pos -= 85
		return 0, pos * 3, 255 - pos*3
	default:
		pos -= 170
		return pos * 3, 255 - pos*3, 0
	}
}

func FillAll(n int, c uint32) []uint32 {
	colors := make([]uint32, n)
	for i := range colors {
		colors[i] = c
	}
	return colors
}

// Gradient は 0..100 に正規化済みの pct を前提とする。
func Gradient(pct float32) (uint8, uint8, uint8) {
	if pct <= 50 {
		t := pct / 50
		return 255, uint8(255 * t), 0
	}
	t := (pct - 50) / 50
	return uint8(255 * (1 - t)), 255, 0
}
