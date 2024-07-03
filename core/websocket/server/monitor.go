package server

import (
	"context"
	"github.com/bpcoder16/Water/logit"
	"runtime"
	"time"
)

func WebSocketMonitor(m *Manager) (err error) {
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-ticker.C:
			var mem runtime.MemStats
			runtime.ReadMemStats(&mem)
			logit.Context(context.Background()).InfoW(
				logit.DefaultMessageKey, "WebSocketMonitor",
				"WebSocketClientCnt", m.Len(),
				"runtime.NumGoroutine", runtime.NumGoroutine(),
				"runtime.NumCPU", runtime.NumCPU(),
				"MemoryAllocation", mem.Alloc/1024,
				"MemorySys", mem.Sys/1024,
			)
		}
	}
}
