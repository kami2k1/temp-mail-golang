package utils

import (
    "encoding/json"
    "os"
    "sync"
)

var (
    domains   []string
    domainsMu sync.RWMutex
)

func LoadDomains(path string) error {
    b, err := os.ReadFile(path)
    if err != nil {
        return err
    }
    var d []string
    if err := json.Unmarshal(b, &d); err != nil {
        return err
    }
    domainsMu.Lock()
    defer domainsMu.Unlock()
    domains = d
    return nil
}

func GetDomains() []string {
    domainsMu.RLock()
    defer domainsMu.RUnlock()
    out := make([]string, len(domains))
    copy(out, domains)
    return out
}

