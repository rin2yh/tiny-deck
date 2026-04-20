package led

import (
	"time"

	"github.com/rin2yh/tiny-deck/internal/color"
	"github.com/rin2yh/tiny-deck/internal/driver"
	"github.com/rin2yh/tiny-deck/internal/keyboard"
)

// TinyGoでヒープ再確保を避けるためパッケージレベルで確保
var buf [keyboard.KeyCount]uint32

var outerKeys = []int{
	int(keyboard.KeyC0R0), int(keyboard.KeyC0R1), int(keyboard.KeyC0R2),
	int(keyboard.KeyC1R0), int(keyboard.KeyC1R1),
	int(keyboard.KeyC2R1), int(keyboard.KeyC2R2),
	int(keyboard.KeyC3R0), int(keyboard.KeyC3R1), int(keyboard.KeyC3R2),
}

// rainbowLoop は duration の間、targets のインデックスにレインボー色を流す。
// targets が nil の場合は全キーを対象とする。
func rainbowLoop(ws *driver.WS2812B, duration time.Duration, targets []int) {
	colors := buf[:]
	n := len(targets)
	if targets == nil {
		n = keyboard.KeyCount
	}
	baseHue := 0
	start := time.Now()
	for time.Since(start) < duration {
		for i := range colors {
			colors[i] = 0
		}
		for j := 0; j < n; j++ {
			idx := j
			if targets != nil {
				idx = targets[j]
			}
			offset := (j * 256) / n
			pos := uint8((baseHue + offset) % 256)
			r, g, b := color.Wheel(pos)
			colors[idx] = ws.PackGRB(r, g, b)
		}
		ws.WriteRaw(colors)
		baseHue = (baseHue + 2) % 256
		time.Sleep(30 * time.Millisecond)
	}
	for i := range colors {
		colors[i] = 0
	}
	ws.WriteRaw(colors)
}

func NotificationAnimation(ws *driver.WS2812B) {
	rainbowLoop(ws, 4*time.Second, outerKeys)
}
