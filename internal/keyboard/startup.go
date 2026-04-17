package keyboard

import (
	"machine"
	"time"

	"github.com/rin2yh/tiny-deck/internal/color"
	"github.com/rin2yh/tiny-deck/internal/driver"
)

func Setup() *Scanner {
	colPins := []machine.Pin{
		machine.GPIO5,
		machine.GPIO6,
		machine.GPIO7,
		machine.GPIO8,
	}

	rowPins := []machine.Pin{
		machine.GPIO9,
		machine.GPIO10,
		machine.GPIO11,
	}

	for _, c := range colPins {
		c.Configure(machine.PinConfig{Mode: machine.PinOutput})
		c.Low()
	}

	for _, c := range rowPins {
		c.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	}

	return &Scanner{
		cols: colPins,
		rows: rowPins,
		keys: make([]bool, len(rowPins)*len(colPins)),
	}
}

func StartupAnimation(ws *driver.WS2812B, sc *Scanner) {
	matrix := sc.KeyCount()
	colors := make([]uint32, matrix)
	baseHue := 0
	start := time.Now()

	for time.Since(start) < time.Second {
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
