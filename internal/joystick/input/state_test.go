package input

import "testing"

func TestStationaryJoystickOffsetDoesNotKeepActivityAlive(t *testing.T) {
	var state State
	first := state.Update(7, -3, false)
	if !first.HasActivity() {
		t.Fatal("first position change should count as activity")
	}

	for i := 0; i < 3; i++ {
		if event := state.Update(7, -3, false); event.HasActivity() {
			t.Fatalf("unchanged position reported activity: %+v", event)
		}
	}
}

func TestAxisMovementAndClickCountAsActivity(t *testing.T) {
	var state State
	state.Update(4, 2, false)

	if event := state.Update(5, 2, false); !event.HasActivity() {
		t.Fatal("axis movement should count as activity")
	}
	if event := state.Update(5, 2, true); !event.HasActivity() || !event.Click {
		t.Fatalf("button press should count as activity: %+v", event)
	}
	if event := state.Update(5, 2, false); event.HasActivity() || !event.Release {
		t.Fatalf("button release should not count as activity: %+v", event)
	}
}
