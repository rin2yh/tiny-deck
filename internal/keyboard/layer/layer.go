package layer

type Layer int

const (
	Deck Layer = iota // 初期状態（デフォルト）
	Game              // ゲーム選択画面

	layerCount
)

func (l Layer) String() string {
	switch l {
	case Deck:
		return "Deck"
	case Game:
		return "Game"
	default:
		return "Unknown"
	}
}

type State struct {
	current Layer
}

func New() *State {
	return &State{current: Deck}
}

func (s *State) Current() Layer {
	return s.current
}

// Next は Deck → Game → Deck と巡回する。
func (s *State) Next() {
	s.current = (s.current + 1) % layerCount
}
