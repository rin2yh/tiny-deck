package display

import (
	"time"

	"github.com/rin2yh/tiny-deck/internal/display/powerstate"
	"tinygo.org/x/drivers/ssd1306"
)

const (
	// These are SSD1306 contrast-register values, not percentages of panel luminance.
	dimContrast    = 0x20
	normalContrast = 0xCF
)

// updatePower applies idle and host-heartbeat transitions. Host metrics stop
// while macOS sleeps or the serial connection is unavailable.
func (d *Display) updatePower(now time.Time, hostUnavailable bool) bool {
	return d.applyPowerTransition(d.power.Update(now, hostUnavailable))
}

func (d *Display) applyPowerTransition(transition powerstate.Transition) bool {
	if !transition.Changed {
		return false
	}

	switch transition.State {
	case powerstate.Normal:
		d.setContrast(normalContrast)
		d.oled.Command(ssd1306.DISPLAYON)
	case powerstate.Dim:
		d.setContrast(dimContrast)
	case powerstate.Off:
		d.oled.Command(ssd1306.DISPLAYOFF)
	}

	return transition.Woke
}

func (d *Display) setContrast(value uint8) {
	d.oled.Command(ssd1306.SETCONTRAST)
	d.oled.Command(value)
}
