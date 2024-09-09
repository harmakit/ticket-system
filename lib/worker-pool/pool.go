package worker_pool

import (
	"context"
	"sync"
)

type Task[In any, Out any] struct {
	Callback func(In) (Out, error)
	Args     In
}

func (t Task[In, Out]) Run() TaskResult[Out] {
	result, err := t.Callback(t.Args)
	return TaskResult[Out]{
		Result: result,
		Err:    err,
	}
}

type TaskResult[Out any] struct {
	Result Out
	Err    error
}

type WorkerPool[In any, Out any] interface {
	Submit(ctx context.Context, tasks []Task[In, Out]) error
	Close() error
}

type pool[In any, Out any] struct {
	wg           *sync.WaitGroup
	workersCount int
	tasks        chan Task[In, Out]
	results      chan TaskResult[Out]
	closed       bool
}

func New[In any, Out any](ctx context.Context, workersCount int) (WorkerPool[In, Out], <-chan TaskResult[Out]) {
	tasks := make(chan Task[In, Out], workersCount)
	results := make(chan TaskResult[Out], workersCount)
	wg := &sync.WaitGroup{}

	p := &pool[In, Out]{
		wg:           wg,
		workersCount: workersCount,
		tasks:        tasks,
		results:      results,
		closed:       false,
	}

	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go p.worker(ctx, tasks, results)
	}

	go func() {
		select {
		case <-ctx.Done():
			_ = p.Close()
			return
		}
	}()

	go func() {
		p.wg.Wait()
		_ = p.Close()
	}()

	return p, results
}

func (p *pool[In, Out]) Close() error {
	if p.closed {
		return ErrWorkerPoolClosed
	}
	p.closed = true

	close(p.tasks)
	p.wg.Wait()
	close(p.results)

	return nil
}

func (p *pool[In, Out]) Submit(ctx context.Context, tasks []Task[In, Out]) error {
	if p.closed {
		return ErrWorkerPoolClosed
	}

	go func() {
		for _, task := range tasks {
			select {
			case <-ctx.Done():
				return
			case p.tasks <- task:
			}
		}
	}()

	return nil
}

func (p *pool[In, Out]) worker(ctx context.Context, tasks <-chan Task[In, Out], results chan<- TaskResult[Out]) {
	defer p.wg.Done()

	for task := range tasks {
		select {
		case <-ctx.Done():
			return
		default:
			result := task.Run()
			results <- result
		}
	}
}
