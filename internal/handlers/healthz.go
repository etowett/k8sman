package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime"
	"time"

	"k8sman/internal/buildinfo"

	"github.com/gin-gonic/gin"
)

func (h Handler) Healthz() func(c *gin.Context) {
	return func(c *gin.Context) {
		hostName, err := os.Hostname()
		if err != nil {
			slog.Error("could not get hostname", "err", err)
		}

		memStats := &runtime.MemStats{}
		runtime.ReadMemStats(memStats)

		currentTime := time.Now()
		tZone, offset := currentTime.Zone()

		c.JSON(http.StatusOK, map[string]interface{}{
			"success":    true,
			"env":        h.cfg.General.Stage,
			"build_time": "1",
			"build_id":   buildinfo.BuildID,
			"build_tag":  buildinfo.BuildTag,
			"time": map[string]interface{}{
				"now":      currentTime,
				"timezone": tZone,
				"offset":   offset,
			},
			"server": map[string]interface{}{
				"hostname":   hostName,
				"cpu":        runtime.NumCPU(),
				"goroutines": runtime.NumGoroutine(),
				"goarch":     runtime.GOARCH,
				"goos":       runtime.GOOS,
				"compiler":   runtime.Compiler,
				"memory": map[string]interface{}{
					"alloc":       fmt.Sprintf("%v MB", bytesToMb(memStats.Alloc)),
					"total_alloc": fmt.Sprintf("%v MB", bytesToMb(memStats.TotalAlloc)),
					"sys":         fmt.Sprintf("%v MB", bytesToMb(memStats.Sys)),
					"num_gc":      memStats.NumGC,
				},
			},
		})
	}
}

func bytesToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
