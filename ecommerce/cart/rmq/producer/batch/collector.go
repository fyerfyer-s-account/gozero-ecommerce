package batch

import (
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rmq/types"
    "sync"
    "time"
)

type Collector struct {
    batch     []*types.CartEvent
    callback  func([]*types.CartEvent) error
    size      int
    interval  time.Duration
    mu        sync.Mutex
    timer     *time.Timer
}

func NewCollector(batchSize int, flushInterval time.Duration) *Collector {
    c := &Collector{
        batch:    make([]*types.CartEvent, 0, batchSize),
        size:     batchSize,
        interval: flushInterval,
    }
    
    c.timer = time.NewTimer(flushInterval)
    go c.timeoutFlusher()
    
    return c
}

func (c *Collector) Add(event *types.CartEvent, callback func([]*types.CartEvent) error) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.batch = append(c.batch, event)
    c.callback = callback
    
    if len(c.batch) >= c.size {
        return c.flush()
    }
    
    return nil
}

func (c *Collector) flush() error {
    if len(c.batch) == 0 {
        return nil
    }
    
    events := make([]*types.CartEvent, len(c.batch))
    copy(events, c.batch)
    c.batch = c.batch[:0]
    
    if c.callback != nil {
        return c.callback(events)
    }
    return nil
}

func (c *Collector) timeoutFlusher() {
    for range c.timer.C {
        c.mu.Lock()
        _ = c.flush()
        c.mu.Unlock()
        c.timer.Reset(c.interval)
    }
}   