package logger

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	logDir   = "log"
	initOnce sync.Once
	writeMu  sync.Mutex
)

// Init ensures the logging directory exists (idempotent).
func Init(dir string) error {
	var err error
	initOnce.Do(func() {
		if dir != "" {
			logDir = dir
		}
		err = os.MkdirAll(logDir, 0o755)
	})
	return err
}

// Entry describes a single JSON log line.
type Entry struct {
	Timestamp string         `json:"time"`
	Action    string         `json:"action"`
	Email     string         `json:"email,omitempty"`
	IP        string         `json:"ip,omitempty"`
	Method    string         `json:"method"`
	Path      string         `json:"path"`
	Extra     map[string]any `json:"extra,omitempty"`
}

// LogRequest writes a JSON log for the given context/action.
func LogRequest(c *gin.Context, action, email string, extra map[string]any) {
	if c == nil || c.Request == nil {
		return
	}

	ip := c.Request.Header.Get("CF-Connecting-IP")
	if ip == "" {
		ip = c.ClientIP()
	}

	entry := Entry{
		Timestamp: time.Now().Format(time.RFC3339Nano),
		Action:    action,
		Email:     email,
		IP:        ip,
		Method:    c.Request.Method,
		Path:      c.FullPath(),
	}
	if entry.Path == "" && c.Request.URL != nil {
		entry.Path = c.Request.URL.Path
	}
	if len(extra) > 0 {
		entry.Extra = extra
	}

	write(entry)
}

func write(e Entry) {
	if err := Init(logDir); err != nil {
		return
	}
	data, err := json.Marshal(e)
	if err != nil {
		return
	}
	filename := filepath.Join(logDir, time.Now().Format("02_01_2006")+".json")

	writeMu.Lock()
	defer writeMu.Unlock()

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return
	}
	defer f.Close()

	_, _ = f.Write(append(data, '\n'))
}
