package model

import (
	"runtime"
	"time"
)

type Stats struct {
	Time         int64  `json:"time"`
	GoVersion    string `json:"go_version"`
	GoMaxProcs   int    `json:"go_max_procs"`
	NumCPU       int    `json:"num_cpu"`
	NumGoroutine int    `json:"num_goroutine"`
	MemSys       uint64 `json:"mem_sys"`
}

func NewStats() *Stats {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	return &Stats{
		Time:         time.Now().UnixNano(),
		GoVersion:    runtime.Version(),
		GoMaxProcs:   runtime.GOMAXPROCS(0),
		NumCPU:       runtime.NumCPU(),
		NumGoroutine: runtime.NumGoroutine(),
		MemSys:       mem.Sys,
	}
}
