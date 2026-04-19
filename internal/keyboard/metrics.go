package keyboard

import (
	"github.com/rin2yh/tiny-deck/internal/color"
	"github.com/rin2yh/tiny-deck/internal/driver"
)

var (
	cpuBarOrder = []int{2, 5, 1, 4, 0, 3}
	memBarOrder = []int{8, 11, 7, 10, 6, 9}
	// TinyGoでヒープ再確保を避けるためパッケージレベルで確保
	metricsBuf [12]uint32
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
