package batch

import (
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/types"
    "sync"
    "time"
)

type Collector struct {
    batch    []*types.OrderEvent
    callback func([]*types.OrderEvent) error
    size     int
    interval time.Duration
    mu       sync.Mutex
    timer    *time.Timer
}

func NewCollector(batchSize int, flushInterval time.Duration) *Collector {
    c := &Collector{
        batch:    make([]*types.OrderEvent, 0, batchSize),
        size:     batchSize,
        interval: flushInterval,
    }
    c.timer = time.NewTimer(flushInterval)
    
    go c.timeoutFlusher()
    return c
}

func (c *Collector) Add(event *types.OrderEvent, callback func([]*types.OrderEvent) error) error {
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

    events := make([]*types.OrderEvent, len(c.batch))
    copy(events, c.batch)
    c.batch = c.batch[:0]

    c.timer.Reset(c.interval)

    if c.callback != nil {
        return c.callback(events)
    }
    return nil
}

func (c *Collector) timeoutFlusher() {
    for range c.timer.C {
        c.mu.Lock()
        c.flush()
        c.mu.Unlock()
    }
}