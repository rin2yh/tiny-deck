package input

import "testing"

func TestStationaryJoystickOffsetDoesNotKeepActivityAlive(t *testing.T) {
	tests := []struct {
		name         string
		dx, dy       int
		wantActivity bool
	}{
		{name: "first nonzero position", dx: 7, dy: -3, wantActivity: true},
		{name: "same position sample one", dx: 7, dy: -3},
		{name: "same position sample two", dx: 7, dy: -3},
		{name: "same position sample three", dx: 7, dy: -3},
	}

	var state State
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if event := state.Update(tt.dx, tt.dy, false); event.HasActivity() != tt.wantActivity {
				t.Fatalf("HasActivity() = %t, want %t; event = %+v", event.HasActivity(), tt.wantActivity, event)
			}
		})
	}
}

func TestAxisMovementAndButtonEdges(t *testing.T) {
	tests := []struct {
		name         string
		dx, dy       int
		pressed      bool
		wantActivity bool
		wantClick    bool
		wantRelease  bool
	}{
		{name: "initial position", dx: 4, dy: 2, wantActivity: true},
		{name: "axis movement", dx: 5, dy: 2, wantActivity: true},
		{name: "button press", dx: 5, dy: 2, pressed: true, wantActivity: true, wantClick: true},
		{name: "button release", dx: 5, dy: 2, wantRelease: true},
	}

	var state State
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := state.Update(tt.dx, tt.dy, tt.pressed)
			if event.HasActivity() != tt.wantActivity || event.Click != tt.wantClick || event.Release != tt.wantRelease {
				t.Fatalf("event = %+v, want activity=%t click=%t release=%t", event, tt.wantActivity, tt.wantClick, tt.wantRelease)
			}
		})
	}
}
