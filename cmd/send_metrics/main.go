package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"go.bug.st/serial"
)

const (
	interval  = time.Second * 2
	notifyMsg = "notify\n"
)

func main() {
	mode := &serial.Mode{BaudRate: 115200}
	port, err := openPort(mode)
	if err != nil {
		log.Fatalf("failed to open port: %w", err)
	}
	defer port.Close()

	var mu sync.Mutex
	go watchNotifications(port, &mu)

	for {
		v, err := cpu.Percent(interval, false)
		if err != nil {
			continue
		}

		line := fmt.Sprintf("cpu:%.2f%%\n", v[0])
		mu.Lock()
		_, err = port.Write([]byte(line))
		mu.Unlock()
		if err != nil {
			log.Println("error occurred: %w", err)
		}
	}
}

func watchNotifications(port serial.Port, mu *sync.Mutex) {
	cmd := exec.Command("log", "stream",
		"--predicate", `process == "UserNotificationsServer"`,
		"--style", "compact",
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("failed to get stdout pipe: %v", err)
		return
	}
	if err := cmd.Start(); err != nil {
		log.Printf("failed to start log stream: %v", err)
		return
	}
	defer cmd.Wait()

	var lastNotify time.Time
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		if time.Since(lastNotify) > time.Second {
			mu.Lock()
			port.Write([]byte(notifyMsg))
			mu.Unlock()
			lastNotify = time.Now()
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("log stream scanner error: %v", err)
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
