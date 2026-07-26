package arcade

import (
	"time"

	"tinygo.org/x/drivers/ssd1306"
)

const (
	// ゲームは 60 TPS 前提だが、I2C 400kHz では 128x64 の全画面転送だけで
	// 20ms 強かかる。32ms ごとに 2 フレーム進めて実時間の速度を合わせる。
	frameInterval = 32 * time.Millisecond
	stepsPerFrame = 2
)

type mode uint8

const (
	modeMenu mode = iota
	modePlaying
)

// Launcher は選択画面とゲーム本体を切り替える状態機械。
type Launcher struct {
	oled    *ssd1306.Device
	mode    mode
	cursor  int
	current Game

	lastFrame time.Time
	// フレーム間隔（32ms）より短い押下を取りこぼさないよう保持する
	jumpLatch bool
	keyHeld   bool
	night     bool
	menuDirty bool

	// ボードに時計が無く毎回同じ状態で起動するため、乱数シードは
	// 選択画面に留まったループ回数（＝プレイヤーの操作タイミング）で作る
	ticks uint32
}

func NewLauncher(oled *ssd1306.Device) *Launcher {
	return &Launcher{oled: oled, menuDirty: true}
}

// Enter は選択画面から始める。ゲームレイヤーに入ったときに呼ぶ。
func (l *Launcher) Enter() {
	l.mode = modeMenu
	l.current = nil
	l.menuDirty = true
	l.jumpLatch = false
	l.keyHeld = false
	l.setNight(false)
}

// Leave は遊んでいたゲームを破棄し、反転表示を元に戻す。
// ゲームレイヤーから抜けるときに呼ぶ。
func (l *Launcher) Leave() {
	l.current = nil
	l.setNight(false)
}

// Update は 1 ループ分の入力を処理し、必要なら画面を描き直す。
// keyDown は操作キーのいずれかが押されているか、rotate はエンコーダーの
// 回転ノッチ数、decide は押し込みの押下エッジ（選択画面では決定、
// プレイ中は終了して選択画面へ戻る）。
func (l *Launcher) Update(keyDown bool, rotate int, decide bool) {
	l.ticks++
	if keyDown && !l.keyHeld {
		l.jumpLatch = true
	}
	l.keyHeld = keyDown

	if l.mode == modeMenu {
		l.updateMenu(rotate, decide)
		return
	}
	l.updatePlaying(decide)
}

func (l *Launcher) updateMenu(rotate int, decide bool) {
	if rotate != 0 {
		l.cursor = (l.cursor + rotate) % len(entries)
		if l.cursor < 0 {
			l.cursor += len(entries)
		}
		l.menuDirty = true
	}

	if decide {
		// シードは 0 だとゲーム側の xorshift が回らないので必ず 1 以上にする
		l.current = entries[l.cursor].New(l.ticks | 1)
		l.mode = modePlaying
		l.jumpLatch = false
		l.lastFrame = time.Now()
		l.drawGame()
		return
	}

	if l.menuDirty {
		l.menuDirty = false
		l.drawMenu()
	}
}

func (l *Launcher) updatePlaying(decide bool) {
	if decide {
		l.mode = modeMenu
		l.current = nil
		l.menuDirty = true
		l.setNight(false)
		return
	}

	now := time.Now()
	if now.Sub(l.lastFrame) < frameInterval {
		return
	}
	l.lastFrame = now

	jump := l.jumpLatch
	l.jumpLatch = false
	for i := 0; i < stepsPerFrame; i++ {
		l.current.Update(jump)
		jump = false // 押下エッジは 1 フレームだけ
	}
	l.drawGame()
}
