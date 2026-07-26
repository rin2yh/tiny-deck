package arcade

import (
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

	// ゲームオーバー画面は y=20 の GAME OVER と y=32 の PRESS JUMP、
	// y=56 の地面が埋まっているので、その隙間に戻る導線を置く。
	exitHintY    = 48
	exitHintText = "PUSH:MENU"
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
	l.drawCentered("ENC:SELECT/START", menuHintY)
	l.oled.Display()
}

func (l *Launcher) drawGame() {
	l.oled.ClearBuffer()
	l.current.Draw(oledDisplay{l.oled})
	if l.current.Over() {
		l.drawCentered(exitHintText, exitHintY)
	}
	l.setNight(l.current.Night())
	l.oled.Display()
}

func (l *Launcher) drawCentered(s string, y int16) {
	_, w := tinyfont.LineWidth(&proggy.TinySZ8pt7b, s)
	width, _ := l.oled.Size()
	tinyfont.WriteLine(l.oled, &proggy.TinySZ8pt7b, (width-int16(w))/2, y, s, white)
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
