package layer

type Layer int

const (
	NumPad Layer = iota // 初期状態（デフォルト）
	Temp
)

func (l Layer) String() string {
	switch l {
	case NumPad:
		return "NumPad"
	case Temp:
		return "Temp"
	default:
		return "Unknown"
	}
}

type State struct {
	current Layer
}

func New() *State {
	return &State{current: NumPad}
}

func (s *State) Current() Layer {
	return s.current
}

func (s *State) Toggle() {
	if s.current == NumPad {
		s.current = Temp
	} else {
		s.current = NumPad
	}
}
