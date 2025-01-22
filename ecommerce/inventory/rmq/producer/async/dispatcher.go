package async

import (
    "context"
    "sync"
)

type Dispatcher struct {
    workers chan struct{}
    wg      sync.WaitGroup
}

func NewDispatcher(maxWorkers int) *Dispatcher {
    return &Dispatcher{
        workers: make(chan struct{}, maxWorkers),
    }
}

func (d *Dispatcher) Dispatch(task func() error) error {
    d.workers <- struct{}{}
    d.wg.Add(1)

    go func() {
        defer func() {
            <-d.workers
            d.wg.Done()
        }()
        
        _ = task()
    }()

    return nil
}

func (d *Dispatcher) DispatchSync(task func() error) error {
    d.workers <- struct{}{}
    defer func() { <-d.workers }()

    return task()
}

func (d *Dispatcher) Wait() {
    d.wg.Wait()
}

func (d *Dispatcher) WaitWithContext(ctx context.Context) error {
    done := make(chan struct{})
    go func() {
        d.wg.Wait()
        close(done)
    }()

    select {
    case <-done:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}