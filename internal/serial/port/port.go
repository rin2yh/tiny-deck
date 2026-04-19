package port

import (
	"errors"
	"log"
	"strings"
	"time"

	"go.bug.st/serial"
)

const (
	retryInterval = 200 * time.Millisecond
	retryCount    = 25
)

// Open は usbmodem デバイスを開く。
// launchd が USB attach イベントで起動するため、/dev/cu.usbmodem* が
// 列挙されるようになるまで最大 5 秒 (200ms * 25) リトライする。
func Open(mode *serial.Mode) (serial.Port, error) {
	var lastErr error
	for range retryCount {
		port, err := open(mode)
		if err == nil {
			return port, nil
		}
		lastErr = err
		time.Sleep(retryInterval)
	}
	return nil, lastErr
}

func open(mode *serial.Mode) (serial.Port, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return nil, err
	}
	for _, p := range ports {
		if !strings.Contains(p, "usbmodem") {
			continue
		}
		port, err := serial.Open(p, mode)
		if err != nil {
			return nil, err
		}
		log.Printf("connected: %s", p)
		return port, nil
	}
	return nil, errors.New("serial port not found")
}
