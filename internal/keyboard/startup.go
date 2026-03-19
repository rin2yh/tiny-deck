package keyboard

import (
	"time"

	"github.com/rin2yh/tiny-deck/internal/color"
	"github.com/rin2yh/tiny-deck/internal/driver"
)

func StartupAnimation(ws *driver.WS2812B, matrix int) {
	colors := make([]uint32, matrix)
	baseHue := 0
	start := time.Now()
	duration := time.Second

	for time.Since(start) < duration {
		for i := 0; i < matrix; i++ {
			offset := (i * 256) / matrix
			pos := uint8((baseHue + offset) % 256)

			r, g, b := color.Wheel(pos)
			colors[i] = ws.PackGRB(r, g, b)
		}

		ws.WriteRaw(colors)
		baseHue = (baseHue + 2) % 256
		time.Sleep(30 * time.Millisecond)
	}

	ws.WriteRaw(color.FillAll(matrix, 0))
}
