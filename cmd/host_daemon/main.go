package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/rin2yh/tiny-deck/internal/serial/port"
	"github.com/rin2yh/tiny-deck/internal/volume"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"go.bug.st/serial"
	_ "modernc.org/sqlite"
)

const (
	interval       = time.Second * 2
	pollInterval   = 500 * time.Millisecond
	notifyMsg      = "notify\n"
	notificationDB = "Library/Group Containers/group.com.apple.usernoted/db2/db"
)

func main() {
	mode := &serial.Mode{BaudRate: 115200}
	sp, err := port.Open(mode)
	if err != nil {
		log.Fatalf("failed to open port: %v", err)
	}
	defer sp.Close()

	var mu sync.Mutex
	go watchNotifications(sp, &mu)
	go watchVolumeCommands(sp)

	for {
		v, err := cpu.Percent(interval, false)
		if err != nil {
			continue
		}
		vm, err := mem.VirtualMemory()
		if err != nil {
			continue
		}

		line := fmt.Sprintf("cpu:%.2f%%\nmem:%.2f%%\n", v[0], vm.UsedPercent)
		mu.Lock()
		_, err = sp.Write([]byte(line))
		mu.Unlock()
		if err != nil {
			log.Fatalf("serial write error: %v", err)
		}
	}
}

func watchNotifications(sp serial.Port, mu *sync.Mutex) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("failed to get home dir: %v", err)
		return
	}
	dbPath := filepath.Join(home, notificationDB)

	db, err := sql.Open("sqlite", "file:"+dbPath+"?mode=ro")
	if err != nil {
		log.Printf("failed to open notification DB: %v", err)
		return
	}
	defer db.Close()

	var lastMax float64 // 0 = not yet initialized; delivered_date is always > 0

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), pollInterval/2)
		var val float64
		err := db.QueryRowContext(ctx, "SELECT MAX(delivered_date) FROM record").Scan(&val)
		cancel()
		if err != nil {
			log.Printf("notification DB query error: %v", err)
			continue
		}
		if lastMax > 0 && val > lastMax {
			mu.Lock()
			_, err := sp.Write([]byte(notifyMsg))
			mu.Unlock()
			if err != nil {
				log.Printf("notify write error: %v", err)
			}
		}
		lastMax = val
	}
}

func watchVolumeCommands(sp serial.Port) {
	scanner := bufio.NewScanner(sp)
	for scanner.Scan() {
		line := scanner.Text()
		var script string
		switch line {
		case volume.CmdUp:
			script = "set volume output volume (output volume of (get volume settings) + 5)"
		case volume.CmdDown:
			script = "set volume output volume (output volume of (get volume settings) - 5)"
		case volume.CmdMute:
			script = "set volume output muted not (output muted of (get volume settings))"
		}
		if script != "" {
			if out, err := exec.Command("osascript", "-e", script).CombinedOutput(); err != nil {
				log.Printf("[vol] osascript error: %v, output: %s", err, out)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("volume command scanner error: %v", err)
	}
}
