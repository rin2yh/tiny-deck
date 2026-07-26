package arcade

import (
	"bytes"
	"image/color"

	"tinygo.org/x/drivers/ssd1306"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/proggy"
)

var white = color.RGBA{R: 255, G: 255, B: 255, A: 255}

const (
	menuTitleY  = 12
	menuFirstY  = 30
	menuLineH   = 12
	menuCursorX = 6
	menuLabelX  = 18
	menuHintY   = 62

	// どのボタンかは選択画面の ENC:SELECT/START と同じ書き方で示す。
	// 重ねる位置はゲームごとに違うので Entry.HintY が持つ。
	exitHintText = "ENC:MENU"

	// SSD1306 の 128x64 は 1 ページ 8 行なのでフレームバッファは 1KB。
	frameBufferSize = 128 * 64 / 8
)

// oledDisplay はゲームの点灯要求を ssd1306 のバッファへ転送する。
type oledDisplay struct{ d *ssd1306.Device }

func (o oledDisplay) SetPixel(x, y int) { o.d.SetPixel(int16(x), int16(y), white) }

func (l *Launcher) drawMenu() {
	l.oled.ClearBuffer()
	tinyfont.WriteLine(l.oled, &proggy.TinySZ8pt7b, 4, menuTitleY, "[Game]", white)
	for i, e := range entries {
		y := int16(menuFirstY + i*menuLineH)
		if i == l.cursor {
			tinyfont.WriteLine(l.oled, &proggy.TinySZ8pt7b, menuCursorX, y, ">", white)
		}
		tinyfont.WriteLine(l.oled, &proggy.TinySZ8pt7b, menuLabelX, y, e.Name, white)
	}
	l.drawCentered(&proggy.TinySZ8pt7b, "ENC:SELECT/START", menuHintY)
	l.flush()
}

func (l *Launcher) drawGame() {
	l.oled.ClearBuffer()
	l.current.Draw(oledDisplay{l.oled})
	if l.current.Over() {
		// ゲーム画面の文字（3x5）に合わせて同じ背丈の TomThumb で描く
		l.drawCentered(&tinyfont.TomThumb, exitHintText, entries[l.cursor].HintY)
	}
	l.setNight(l.current.Night())
	l.flush()
}

// flush は前回転送した内容と変わっていなければ I2C 転送を省く。
// 全画面転送は 20ms 強かかるのに対し 1KB の比較は桁違いに安く、
// タイトルとゲームオーバーは点滅以外静止しているので効きが大きい。
func (l *Launcher) flush() {
	buf := l.oled.GetBuffer()
	if l.sent && bytes.Equal(buf, l.sentBuf[:]) {
		return
	}
	copy(l.sentBuf[:], buf)
	l.sent = true
	l.oled.Display()
}

func (l *Launcher) drawCentered(f tinyfont.Fonter, s string, y int16) {
	_, w := tinyfont.LineWidth(f, s)
	width, _ := l.oled.Size()
	tinyfont.WriteLine(l.oled, f, (width-int16(w))/2, y, s, white)
}

// setNight はナイトモードを SSD1306 の表示反転コマンドで表現する。
// 全画素を塗り潰して描き分けるより I2C 転送が 1 コマンドで済む。
func (l *Launcher) setNight(on bool) {
	if on == l.night {
		return
	}
	l.night = on
	if on {
		l.oled.Command(ssd1306.INVERTDISPLAY)
	} else {
		l.oled.Command(ssd1306.NORMALDISPLAY)
	}
}
