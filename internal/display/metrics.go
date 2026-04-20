package display

import (
	"strconv"
	"strings"

	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

const (
	prefixCPU = "cpu:"
	prefixMem = "mem:"
)

type metricsState struct {
	lastCPULine string
	lastMemLine string
	cpu         float32
	mem         float32
}

func parsePercent(line, prefix string) (float32, bool) {
	s := strings.TrimPrefix(line, prefix)
	s = strings.TrimSuffix(s, "%")
	s = strings.TrimSpace(s)
	v, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, false
	}
	return float32(v), true
}

func (m *metricsState) tryUpdate(line string) bool {
	switch {
	case strings.HasPrefix(line, prefixCPU):
		if line == m.lastCPULine {
			return false
		}
		m.lastCPULine = line
		if v, ok := parsePercent(line, prefixCPU); ok {
			m.cpu = v
			return true
		}
	case strings.HasPrefix(line, prefixMem):
		if line == m.lastMemLine {
			return false
		}
		m.lastMemLine = line
		if v, ok := parsePercent(line, prefixMem); ok {
			m.mem = v
			return true
		}
	}
	return false
}

func (m *metricsState) draw(d *ssd1306.Device) {
	tinyfont.WriteLine(d, &freemono.Regular9pt7b, 4, 40, m.lastCPULine, white)
	tinyfont.WriteLine(d, &freemono.Regular9pt7b, 4, 58, m.lastMemLine, white)
}
