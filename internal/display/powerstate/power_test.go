package powerstate

import (
	"testing"
	"time"
)

func TestUpdateDimsAndTurnsOffAtIdleThresholds(t *testing.T) {
	started := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	runSteps(t, started, []powerStep{
		{name: "before dim threshold", after: time.Minute - time.Second, want: Transition{State: Normal}},
		{name: "dim at one minute", after: time.Minute, want: Transition{State: Dim, Changed: true}},
		{name: "before off threshold", after: 3*time.Minute - time.Second, want: Transition{State: Dim}},
		{name: "off at three minutes", after: 3 * time.Minute, want: Transition{State: Off, Changed: true}},
	})
}

func TestActivityWakesOLEDAndRestartsIdleTimer(t *testing.T) {
	started := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	runSteps(t, started, []powerStep{
		{name: "dim", after: time.Minute, want: Transition{State: Dim, Changed: true}},
		{name: "off", after: 3 * time.Minute, want: Transition{State: Off, Changed: true}},
		{name: "activity wakes OLED", after: 4 * time.Minute, activity: true, want: Transition{State: Normal, Changed: true, Woke: true}},
		{name: "before dim after activity", after: 5*time.Minute - time.Second, want: Transition{State: Normal}},
		{name: "dim after activity timeout", after: 5 * time.Minute, want: Transition{State: Dim, Changed: true}},
	})
}

func TestHostUnavailableTurnsOffAndReconnectRestartsIdleTimer(t *testing.T) {
	started := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	runSteps(t, started, []powerStep{
		{name: "host becomes unavailable", after: 5 * time.Second, hostUnavailable: true, want: Transition{State: Off, Changed: true}},
		{name: "still unavailable", after: 10 * time.Second, hostUnavailable: true, want: Transition{State: Off}},
		{name: "host reconnects", after: 15 * time.Second, want: Transition{State: Normal, Changed: true, Woke: true}},
		{name: "before dim after reconnect", after: 74 * time.Second, want: Transition{State: Normal}},
		{name: "dim after reconnect timeout", after: 75 * time.Second, want: Transition{State: Dim, Changed: true}},
	})
}

type powerStep struct {
	name            string
	after           time.Duration
	hostUnavailable bool
	activity        bool
	want            Transition
}

func runSteps(t *testing.T, started time.Time, steps []powerStep) {
	t.Helper()
	controller := New(started)
	for _, step := range steps {
		t.Run(step.name, func(t *testing.T) {
			now := started.Add(step.after)
			var got Transition
			if step.activity {
				got = controller.Activity(now)
			} else {
				got = controller.Update(now, step.hostUnavailable)
			}
			if got != step.want {
				t.Fatalf("transition = %+v, want %+v", got, step.want)
			}
		})
	}
}
