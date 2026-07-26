// Package arcade は OLED 上で遊べるゲームと、その選択画面を提供する。
// ゲーム本体は外部モジュール（github.com/rin2yh/dinosaur-game 等）が持ち、
// このパッケージは選択・起動・終了と ssd1306 への描画だけを受け持つ。
package arcade

import "github.com/rin2yh/dinosaur-game/game"

// Display は 128x64 1bit 画面の点灯要求。ゲーム側の描画インターフェースと
// 同一の型でなければゲームが Game を満たさなくなるため、独自定義ではなく
// 型エイリアスにしている。
type Display = game.Display

// Game は選択画面から起動できる 1 本のゲーム。
type Game interface {
	// Update は 1 フレーム進める。pressed は操作キーの押下エッジ。
	Update(pressed bool)
	Draw(d Display)
	// Night は画面を反転表示する状態かどうか。
	Night() bool
	// Over はゲームオーバー表示中かどうか。手が空くこのタイミングだけ
	// 選択画面へ戻る導線を画面に重ねる。
	Over() bool
}

// Entry は選択画面に並ぶ 1 項目。ゲームを増やすときは entries に足す。
type Entry struct {
	Name string
	New  func(seed uint32) Game
}

var entries = []Entry{
	{Name: "DINOSAUR", New: func(seed uint32) Game { return dinosaur{game.New(seed)} }},
}

// dinosaur は game.Game に Over を足すだけのアダプタ。ゲーム側の Mode 型を
// arcade のインターフェースに持ち込まないよう、ここで吸収する。
type dinosaur struct{ *game.Game }

func (d dinosaur) Over() bool { return d.Mode() == game.ModeGameOver }
