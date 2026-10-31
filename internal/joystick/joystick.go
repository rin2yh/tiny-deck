package joystick

import "machine"

const (
	deadzoneRaw = 2000
	divisorRaw  = 1500
	calSamples  = 8
)

type Event struct {
	DX, DY      int
	Click       bool
	Release     bool
	axisChanged bool
}

// HasActivity reports whether the joystick position changed or a click occurred.
func (e Event) HasActivity() bool {
	return e.axisChanged || e.Click
}

type Joystick struct {
	adcX, adcY       machine.ADC
	btn              machine.Pin
	centerX, centerY int
	invertX, invertY bool
	btnLast          bool
	lastDX, lastDY   int
}

func NewJoystick(invertX, invertY bool) *Joystick {
	j := &Joystick{
		adcX:    machine.ADC{Pin: machine.GPIO29},
		adcY:    machine.ADC{Pin: machine.GPIO28},
		btn:     machine.GPIO0,
		invertX: invertX,
		invertY: invertY,
	}
	j.adcX.Configure(machine.ADCConfig{})
	j.adcY.Configure(machine.ADCConfig{})
	j.btn.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	var sx, sy int
	for i := 0; i < calSamples; i++ {
		sx += int(j.adcX.Get())
		sy += int(j.adcY.Get())
	}
	j.centerX = sx / calSamples
	j.centerY = sy / calSamples
	return j
}

func (j *Joystick) Update() Event {
	dx := scaleAxis(int(j.adcX.Get()), j.centerX, deadzoneRaw, divisorRaw, j.invertX)
	dy := scaleAxis(int(j.adcY.Get()), j.centerY, deadzoneRaw, divisorRaw, j.invertY)

	pressed := !j.btn.Get()
	ev := Event{
		DX:          dx,
		DY:          dy,
		Click:       pressed && !j.btnLast,
		Release:     !pressed && j.btnLast,
		axisChanged: dx != j.lastDX || dy != j.lastDY,
	}
	j.btnLast = pressed
	j.lastDX = dx
	j.lastDY = dy
	return ev
}
