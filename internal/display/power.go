package display

import (
	"time"

	"tinygo.org/x/drivers/ssd1306"
)

const (
	dimAfter = time.Minute
	offAfter = 3 * time.Minute
	// These are SSD1306 contrast-register values, not percentages of panel luminance.
	dimContrast    = 0x20
	normalContrast = 0xCF
)

type powerState uint8

const (
	powerNormal powerState = iota
	powerDim
	powerOff
)

// updatePower applies idle and host-heartbeat transitions. Host metrics stop
// while macOS sleeps or the serial connection is unavailable.
func (d *Display) updatePower(now time.Time, hostUnavailable bool) bool {
	if hostUnavailable && !d.hostUnavailable {
		d.hostUnavailable = true
		d.setPowerState(powerOff)
		return false
	}
	if !hostUnavailable && d.hostUnavailable {
		d.hostUnavailable = false
		d.lastActivity = now
		return d.setPowerState(powerNormal)
	}

	switch d.powerState {
	case powerNormal:
		if now.Sub(d.lastActivity) >= dimAfter {
			d.setPowerState(powerDim)
		}
	case powerDim:
		if now.Sub(d.lastActivity) >= offAfter {
			d.setPowerState(powerOff)
		}
	}
	return false
}

// setPowerState returns true when a previously off OLED has been turned on.
func (d *Display) setPowerState(next powerState) bool {
	if d.powerState == next {
		return false
	}
	wasOff := d.powerState == powerOff
	d.powerState = next

	switch next {
	case powerNormal:
		d.setContrast(normalContrast)
		d.oled.Command(ssd1306.DISPLAYON)
	case powerDim:
		d.setContrast(dimContrast)
	case powerOff:
		d.oled.Command(ssd1306.DISPLAYOFF)
	}

	return wasOff && next == powerNormal
}

func (d *Display) setContrast(value uint8) {
	d.oled.Command(ssd1306.SETCONTRAST)
	d.oled.Command(value)
}
