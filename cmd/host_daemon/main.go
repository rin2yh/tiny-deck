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
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/rin2yh/tiny-deck/internal/serial/port"
	"github.com/rin2yh/tiny-deck/internal/volume"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"go.bug.st/serial"
	_ "modernc.org/sqlite"
)

const (
	interval        = time.Second * 2
	pollInterval    = 500 * time.Millisecond
	notifyMsg       = "notify\n"
	notificationDB  = "Library/Group Containers/group.com.apple.usernoted/db2/db"
	lockFileName    = "tiny-deck-host.lock"
	maxDeliveredSQL = "SELECT MAX(delivered_date) FROM record"
)

func main() {
	lockFile, ok := acquireSingleInstanceLock()
	if !ok {
		os.Exit(0)
	}
	defer lockFile.Close()

	mode := &serial.Mode{BaudRate: 115200}
	sp, err := port.Open(mode)
	if err != nil {
		log.Fatalf("failed to open port: %v", err)
	}
	defer sp.Close()

	var mu sync.Mutex
	go watchNotifications(sp, &mu)
	go watchVolumeCommands(sp, &mu)

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
			log.Printf("serial write failed (device likely detached): %v", err)
			return
		}
	}
}

func acquireSingleInstanceLock() (*os.File, bool) {
	lockPath := filepath.Join(os.TempDir(), lockFileName)
	f, err := os.OpenFile(lockPath, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		log.Printf("failed to open lock file %s: %v", lockPath, err)
		return nil, false
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		_ = f.Close()
		return nil, false
	}
	return f, true
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

	ctx, cancel := context.WithTimeout(context.Background(), pollInterval/2)
	var lastMax float64
	err = db.QueryRowContext(ctx, maxDeliveredSQL).Scan(&lastMax)
	cancel()
	if err != nil {
		log.Printf("notification DB unreadable (likely needs Full Disk Access), disabling notification forwarding: %v", err)
		return
	}

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for range ticker.C {
		ctx, cancel := context.WithTimeout(context.Background(), pollInterval/2)
		var val float64
		err := db.QueryRowContext(ctx, maxDeliveredSQL).Scan(&val)
		cancel()
		if err != nil {
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

func watchVolumeCommands(sp serial.Port, mu *sync.Mutex) {
	scanner := bufio.NewScanner(sp)
	for scanner.Scan() {
		line := scanner.Text()
		var setExpr string
		switch line {
		case volume.CmdUp:
			setExpr = "set volume output volume ((output volume of (get volume settings)) + 5)"
		case volume.CmdDown:
			setExpr = "set volume output volume ((output volume of (get volume settings)) - 5)"
		case volume.CmdMute:
			setExpr = "set volume output muted (not (output muted of (get volume settings)))"
		default:
			continue
		}
		out, err := exec.Command(
			"osascript",
			"-e", setExpr,
			"-e", "set s to (get volume settings)",
			"-e", "((output volume of s) as text) & \",\" & ((output muted of s) as text)",
		).Output()
		if err != nil {
			log.Printf("[vol] osascript error: %v (out=%q)", err, string(out))
			continue
		}
		vol, muted, ok := parseVolumeReply(string(out))
		if !ok {
			log.Printf("[vol] unexpected osascript output: %q", string(out))
			continue
		}
		mutedInt := 0
		if muted {
			mutedInt = 1
		}
		reply := fmt.Sprintf("%s%d\n%s%d\n", volume.PrefixCurrent, vol, volume.PrefixMuted, mutedInt)
		mu.Lock()
		_, werr := sp.Write([]byte(reply))
		mu.Unlock()
		if werr != nil {
			log.Printf("[vol] serial write failed: %v", werr)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("volume command scanner error: %v", err)
	}
}

func parseVolumeReply(s string) (int, bool, bool) {
	s = strings.TrimSpace(s)
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		return 0, false, false
	}
	v, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, false, false
	}
	if v < 0 {
		v = 0
	} else if v > 100 {
		v = 100
	}
	muted := strings.TrimSpace(parts[1]) == "true"
	return v, muted, true
}
