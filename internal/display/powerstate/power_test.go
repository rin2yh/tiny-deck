package powerstate

import (
	"testing"
	"time"
)

func TestUpdateDimsAndTurnsOffAtIdleThresholds(t *testing.T) {
	started := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	controller := New(started)

	assertState(t, controller.Update(started.Add(time.Minute-time.Second), false), Normal)
	assertTransition(t, controller.Update(started.Add(time.Minute), false), Dim, true, false)
	assertState(t, controller.Update(started.Add(3*time.Minute-time.Second), false), Dim)
	assertTransition(t, controller.Update(started.Add(3*time.Minute), false), Off, true, false)
}

func TestActivityWakesOLEDAndRestartsIdleTimer(t *testing.T) {
	started := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	controller := New(started)
	controller.Update(started.Add(time.Minute), false)
	controller.Update(started.Add(3*time.Minute), false)

	activityAt := started.Add(4 * time.Minute)
	assertTransition(t, controller.Activity(activityAt), Normal, true, true)
	assertState(t, controller.Update(activityAt.Add(time.Minute-time.Second), false), Normal)
	assertTransition(t, controller.Update(activityAt.Add(time.Minute), false), Dim, true, false)
}

func TestHostUnavailableTurnsOffAndReconnectRestartsIdleTimer(t *testing.T) {
	started := time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC)
	controller := New(started)

	assertTransition(t, controller.Update(started.Add(5*time.Second), true), Off, true, false)
	assertState(t, controller.Update(started.Add(10*time.Second), true), Off)

	reconnectedAt := started.Add(15 * time.Second)
	assertTransition(t, controller.Update(reconnectedAt, false), Normal, true, true)
	assertState(t, controller.Update(reconnectedAt.Add(time.Minute-time.Second), false), Normal)
	assertTransition(t, controller.Update(reconnectedAt.Add(time.Minute), false), Dim, true, false)
}

func assertState(t *testing.T, got Transition, want State) {
	t.Helper()
	if got.State != want {
		t.Fatalf("state = %v, want %v", got.State, want)
	}
	if got.Changed || got.Woke {
		t.Fatalf("transition = %+v, want no change", got)
	}
}

func assertTransition(t *testing.T, got Transition, state State, changed, woke bool) {
	t.Helper()
	if got.State != state || got.Changed != changed || got.Woke != woke {
		t.Fatalf("transition = %+v, want state=%v changed=%t woke=%t", got, state, changed, woke)
	}
}
