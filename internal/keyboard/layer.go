package keyboard

type Layer int

const (
	LayerNumPad Layer = iota // 初期状態（デフォルト）
	LayerTemp
)

func (l Layer) String() string {
	switch l {
	case LayerNumPad:
		return "NumPad"
	case LayerTemp:
		return "Temp"
	default:
		return "Unknown"
	}
}

type LayerState struct {
	current Layer
}

func NewLayerState() *LayerState {
	return &LayerState{current: LayerNumPad}
}

func (s *LayerState) Current() Layer {
	return s.current
}

func (s *LayerState) Switch(l Layer) {
	s.current = l
}

func (s *LayerState) Toggle() {
	if s.current == LayerNumPad {
		s.current = LayerTemp
	} else {
		s.current = LayerNumPad
	}
}
