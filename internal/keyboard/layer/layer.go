package layer

type Layer int

const (
	NumPad Layer = iota // 初期状態（デフォルト）
	Temp
	Game // ゲーム選択画面

	layerCount
)

func (l Layer) String() string {
	switch l {
	case NumPad:
		return "NumPad"
	case Temp:
		return "Temp"
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
	return &State{current: NumPad}
}

func (s *State) Current() Layer {
	return s.current
}

// Next は NumPad → Temp → Game → NumPad と巡回する。
func (s *State) Next() {
	s.current = (s.current + 1) % layerCount
}
