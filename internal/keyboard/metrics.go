package keyboard

import (
	"github.com/rin2yh/tiny-deck/internal/color"
	"github.com/rin2yh/tiny-deck/internal/driver"
)

var (
	// Left 2 columns, bottom-to-top
	cpuBarOrder = []int{
		int(KeyC0R2), int(KeyC1R2),
		int(KeyC0R1), int(KeyC1R1),
		int(KeyC0R0), int(KeyC1R0),
	}
	// Right 2 columns, bottom-to-top
	memBarOrder = []int{
		int(KeyC2R2), int(KeyC3R2),
		int(KeyC2R1), int(KeyC3R1),
		int(KeyC2R0), int(KeyC3R0),
	}
	// TinyGoでヒープ再確保を避けるためパッケージレベルで確保
	metricsBuf [KeyCount]uint32
)

func RenderMetrics(ws *driver.WS2812B, matrix int, cpuPct, memPct float32) {
	colors := metricsBuf[:matrix]
	for i := range colors {
		colors[i] = 0
	}
	paintBar(ws, colors, cpuBarOrder, cpuPct)
	paintBar(ws, colors, memBarOrder, memPct)
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
