package joystick

import "machine"

const (
	deadzoneRaw = 2000
	divisorRaw  = 1500
	calSamples  = 8
)

type Event struct {
	DX, DY  int
	Click   bool
	Release bool
}

type Joystick struct {
	adcX, adcY       machine.ADC
	btn              machine.Pin
	centerX, centerY int
	invertX, invertY bool
	btnLast          bool
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
		DX:      dx,
		DY:      dy,
		Click:   pressed && !j.btnLast,
		Release: !pressed && j.btnLast,
	}
	j.btnLast = pressed
	return ev
}
