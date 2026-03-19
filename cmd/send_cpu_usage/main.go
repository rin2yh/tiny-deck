package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"go.bug.st/serial"
)

const (
	interval = time.Second * 2
)

func main() {
	mode := &serial.Mode{BaudRate: 115200}
	port, err := openPort(mode)
	if err != nil {
		log.Fatalf("failed to open port: %w", err)
	}
	defer port.Close()

	for {
		v, err := cpu.Percent(interval, false)
		if err != nil {
			continue
		}

		line := fmt.Sprintf("%.2f\n", v[0])
		_, err = port.Write([]byte(line))
		if err != nil {
			log.Println("error occurred: %w", err)
		}
	}
}

func openPort(mode *serial.Mode) (serial.Port, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return nil, err
	}

	for _, p := range ports {
		if strings.Contains(p, "usbmodem") {
			port, err := serial.Open(p, mode)
			if err == nil {
				fmt.Println("connected:", p)
				return port, nil
			}
		}
	}

	return nil, fmt.Errorf("serial port not found")
}
