package led

import (
	"time"

	"github.com/rin2yh/tiny-deck/internal/driver"
)

func StartupAnimation(ws *driver.WS2812B) {
	rainbowLoop(ws, time.Second, nil)
}
