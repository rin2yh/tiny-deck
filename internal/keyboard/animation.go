package keyboard

import (
	"time"

	"github.com/rin2yh/tiny-deck/internal/color"
	"github.com/rin2yh/tiny-deck/internal/driver"
)

var outerKeys = []int{0, 1, 2, 3, 4, 7, 8, 9, 10, 11}

func NotificationAnimation(ws *driver.WS2812B, matrix int) {
	colors := make([]uint32, matrix)
	baseHue := 0
	n := len(outerKeys)
	start := time.Now()

	for time.Since(start) < 4*time.Second {
		for j, idx := range outerKeys {
			offset := (j * 256) / n
			pos := uint8((baseHue + offset) % 256)
			r, g, b := color.Wheel(pos)
			colors[idx] = ws.PackGRB(r, g, b)
		}
		ws.WriteRaw(colors)
		baseHue = (baseHue + 2) % 256
		time.Sleep(30 * time.Millisecond)
	}

	ws.WriteRaw(color.FillAll(matrix, 0))
}
