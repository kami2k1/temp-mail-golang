package utils

import (
    "encoding/json"
    "os"
    "sync"
    "time"
)

type blacklistStore struct {
    Items map[string]int64 `json:"items"`
}

var (
    blPath string
    blMu   sync.RWMutex
    bl     = blacklistStore{Items: map[string]int64{}}
)

const blacklistTTL = 24 * time.Hour

func InitBlacklist(path string) error {
    blMu.Lock()
    defer blMu.Unlock()
    blPath = path
   
    if b, err := os.ReadFile(path); err == nil {
        _ = json.Unmarshal(b, &bl)
    }
   
    cutoff := time.Now().Add(-blacklistTTL).Unix()
    for k, ts := range bl.Items {
        if ts < cutoff {
            delete(bl.Items, k)
        }
    }
    return saveLocked()
}

func saveLocked() error {
    if blPath == "" {
        return nil
    }
    b, _ := json.MarshalIndent(bl, "", "  ")
    return os.WriteFile(blPath, b, 0644)
}

func IsBlacklisted(email string) bool {
    blMu.RLock()
    defer blMu.RUnlock()
    ts, ok := bl.Items[email]
    if !ok {
        return false
    }
    return time.Unix(ts, 0).After(time.Now().Add(-blacklistTTL))
}

func AddToBlacklist(email string) {
    blMu.Lock()
    defer blMu.Unlock()
    bl.Items[email] = time.Now().Unix()
    _ = saveLocked()
}

