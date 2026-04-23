package led

import (
	"github.com/rin2yh/tiny-deck/internal/color"
	"github.com/rin2yh/tiny-deck/internal/driver"
	"github.com/rin2yh/tiny-deck/internal/keyboard"
)

var (
	// Left 2 columns, bottom-to-top
	cpuBarOrder = []int{
		int(keyboard.KeyC0R2), int(keyboard.KeyC1R2),
		int(keyboard.KeyC0R1), int(keyboard.KeyC1R1),
		int(keyboard.KeyC0R0), int(keyboard.KeyC1R0),
	}
	// Right 2 columns, bottom-to-top
	memBarOrder = []int{
		int(keyboard.KeyC2R2), int(keyboard.KeyC3R2),
		int(keyboard.KeyC2R1), int(keyboard.KeyC3R1),
		int(keyboard.KeyC2R0), int(keyboard.KeyC3R0),
	}
)

func RenderMetrics(ws *driver.WS2812B, cpuPct, memPct float32) {
	colors := buf[:]
	for i := range colors {
		colors[i] = 0
	}
	paintBar(ws, colors, cpuBarOrder, cpuPct)
	paintBar(ws, colors, memBarOrder, memPct)
	ws.WriteRaw(colors)
}

func ClearMetrics(ws *driver.WS2812B) {
	colors := buf[:]
	for i := range colors {
		colors[i] = 0
	}
	ws.WriteRaw(colors)
}

func paintBar(ws *driver.WS2812B, colors []uint32, order []int, pct float32) {
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	lit := int((pct*float32(len(order)))/100 + 0.5)
	r, g, b := color.Gradient(pct)
	c := ws.PackGRB(r, g, b)
	for i := 0; i < lit; i++ {
		colors[order[i]] = c
	}
}
