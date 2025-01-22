package async

import (
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
    d.workers <- struct{}{} // Acquire worker
    d.wg.Add(1)
    
    go func() {
        defer func() {
            <-d.workers // Release worker
            d.wg.Done()
        }()
        
        _ = task()
    }()
    
    return nil
}

func (d *Dispatcher) DispatchSync(task func() error) error {
    d.workers <- struct{}{} // Acquire worker
    defer func() {
        <-d.workers // Release worker
    }()
    
    return task()
}

func (d *Dispatcher) Wait() {
    d.wg.Wait()
}