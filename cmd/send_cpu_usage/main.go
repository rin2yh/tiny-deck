package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

const (
	interval = time.Second * 2
)

func main() {
	for {
		v, err := cpu.Percent(interval, false)
		if err != nil {
			continue
		}

		fmt.Printf("%.2f\n", v[0])
	}
}
