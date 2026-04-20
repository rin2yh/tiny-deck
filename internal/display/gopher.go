package display

import (
	"time"

	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/gophers"
)

const (
	gopherBaseX    = 104
	gopherBaseY    = 28
	gopherInterval = 220 * time.Millisecond
)

// 全身がちゃんと映る 32pt glyph のみ採用（目だけ／部分表示になる低背 glyph は除外）
var gopherPoses = [...]string{"A", "B", "W", "Y"}

var gopherOffsets = [4]struct{ dx, dy int16 }{
	{dx: +2, dy: 0},
	{dx: -2, dy: 0},
	{dx: +2, dy: -1},
	{dx: -2, dy: -1},
}

type gopherAnim struct {
	phase  uint8
	lastAt time.Time
}

func newGopherAnim() gopherAnim {
	return gopherAnim{lastAt: time.Now()}
}

func (g *gopherAnim) tick(now time.Time) bool {
	if now.Sub(g.lastAt) < gopherInterval {
		return false
	}
	g.phase++
	g.lastAt = now
	return true
}

func (g *gopherAnim) draw(d *ssd1306.Device) {
	phase := g.phase & 0x03
	off := gopherOffsets[phase]
	tinyfont.WriteLine(d, &gophers.Regular32pt, gopherBaseX+off.dx, gopherBaseY+off.dy, gopherPoses[phase], white)
}
