package display

import (
	"strconv"
	"strings"
	"time"

	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freemono"
)

const (
	prefixCPU      = "cpu:"
	prefixMem      = "mem:"
	metricsTimeout = 5 * time.Second
)

type metricsState struct {
	lastCPULine string
	lastMemLine string
	cpu         float32
	mem         float32
	lastAt      time.Time
	hasData     bool
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

func (m *metricsState) tryUpdate(line string, now time.Time) bool {
	switch {
	case strings.HasPrefix(line, prefixCPU):
		if line == m.lastCPULine {
			return false
		}
		m.lastCPULine = line
		if v, ok := parsePercent(line, prefixCPU); ok {
			m.cpu = v
			m.lastAt = now
			m.hasData = true
			return true
		}
	case strings.HasPrefix(line, prefixMem):
		if line == m.lastMemLine {
			return false
		}
		m.lastMemLine = line
		if v, ok := parsePercent(line, prefixMem); ok {
			m.mem = v
			m.lastAt = now
			m.hasData = true
			return true
		}
	}
	return false
}

func (m *metricsState) isStale(now time.Time) bool {
	if !m.hasData {
		return true
	}
	return now.Sub(m.lastAt) >= metricsTimeout
}

func (m *metricsState) draw(d *ssd1306.Device, stale bool) {
	if stale {
		return
	}
	tinyfont.WriteLine(d, &freemono.Regular9pt7b, 4, 40, m.lastCPULine, white)
	tinyfont.WriteLine(d, &freemono.Regular9pt7b, 4, 58, m.lastMemLine, white)
}
