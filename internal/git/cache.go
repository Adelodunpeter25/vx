package git

import (
	"context"
	"sync"
)

type Cache struct {
	mu      sync.RWMutex
	snap    *Snapshot
	loading bool
}

func NewCache() *Cache {
	return &Cache{}
}

func (c *Cache) Snapshot() *Snapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.snap
}

func (c *Cache) Update(ctx context.Context, collector *Collector) error {
	snap, err := collector.Collect(ctx)
	if err != nil {
		return err
	}
	c.mu.Lock()
	c.snap = snap
	c.mu.Unlock()
	return nil
}

func (c *Cache) Status(path string) (Status, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.snap == nil {
		return StatusUnknown, false
	}
	status, ok := c.snap.Statuses[path]
	if ok {
		return status, true
	}
	status, ok = c.snap.DirStatuses[path]
	if ok {
		return status, true
	}
	return StatusClean, false
}

