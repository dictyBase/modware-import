package concurrent

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

// Job represents a single unit of work to be processed
type Job[I any] struct {
	ID      string
	Payload I
	Meta    map[string]interface{} // Optional metadata
}

// Result represents the output of processing a job
type Result[O any] struct {
	JobID    string
	Output   O
	Error    error
	Duration time.Duration
}

// WorkerFunc defines the function that processes jobs
type WorkerFunc[I, O any] func(context.Context, Job[I]) (O, error)

// PoolStats tracks worker pool execution metrics
type PoolStats struct {
	JobsProcessed       int64
	JobsSucceeded       int64
	JobsFailed          int64
	TotalProcessingTime time.Duration
	mu                  sync.Mutex // For concurrent updates
}

// Pool manages a collection of workers for parallel task processing
type Pool[I, O any] struct {
	// Configuration
	Workers    int
	BufferSize int

	// Channels
	jobChan    chan Job[I]
	resultChan chan Result[O]
	errorChan  chan error
	doneChan   chan struct{}

	// Control
	ctx        context.Context
	cancelFunc context.CancelFunc
	wg         sync.WaitGroup

	// Worker factory function
	workerFunc WorkerFunc[I, O]

	// Metrics
	stats *PoolStats
}

// PoolOption defines functional options for configuring a Pool
type PoolOption[I, O any] func(*Pool[I, O])

// worker represents a single worker in the pool
type worker[I, O any] struct {
	id         int
	jobChan    <-chan Job[I]
	resultChan chan<- Result[O]
	errorChan  chan<- error
	workerFunc WorkerFunc[I, O]
	ctx        context.Context
}

// WithWorkers sets the number of workers in the pool
func WithWorkers[I, O any](n int) PoolOption[I, O] {
	return func(p *Pool[I, O]) {
		if n > 0 {
			p.Workers = n
		}
	}
}

// WithBufferSize sets the buffer size for the job and result channels
func WithBufferSize[I, O any](size int) PoolOption[I, O] {
	return func(p *Pool[I, O]) {
		if size > 0 {
			p.BufferSize = size
		}
	}
}

// WithContext sets a custom context for the pool
func WithContext[I, O any](ctx context.Context) PoolOption[I, O] {
	return func(p *Pool[I, O]) {
		if ctx != nil {
			// Cancel any existing context
			if p.cancelFunc != nil {
				p.cancelFunc()
			}

			// Create new cancellable context
			p.ctx, p.cancelFunc = context.WithCancel(ctx)
		}
	}
}

// NewPool creates a new worker pool with the given options
func NewPool[I, O any](workerFunc WorkerFunc[I, O], options ...PoolOption[I, O]) *Pool[I, O] {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &Pool[I, O]{
		Workers:    4,   // Default number of workers
		BufferSize: 100, // Default buffer size
		workerFunc: workerFunc,
		ctx:        ctx,
		cancelFunc: cancel,
		stats:      &PoolStats{},
	}

	// Apply options
	for _, option := range options {
		option(pool)
	}

	return pool
}

// Start initializes channels and starts worker goroutines
func (p *Pool[I, O]) Start() {
	p.jobChan = make(chan Job[I], p.BufferSize)
	p.resultChan = make(chan Result[O], p.BufferSize)
	p.errorChan = make(chan error, p.BufferSize)
	p.doneChan = make(chan struct{})

	p.wg.Add(p.Workers)
	for i := 0; i < p.Workers; i++ {
		w := &worker[I, O]{
			id:         i,
			jobChan:    p.jobChan,
			resultChan: p.resultChan,
			errorChan:  p.errorChan,
			workerFunc: p.workerFunc,
			ctx:        p.ctx,
		}
		go func(worker *worker[I, O]) {
			defer p.wg.Done()
			worker.start()
		}(w)
	}

	// Start background goroutine to signal when all workers have completed
	go func() {
		p.wg.Wait()
		close(p.doneChan)
	}()
}

// Submit adds a job to the pool
func (p *Pool[I, O]) Submit(payload I, meta ...map[string]interface{}) {
	job := Job[I]{
		ID:      uuid.NewString(),
		Payload: payload,
	}
	if len(meta) > 0 {
		job.Meta = meta[0]
	}

	select {
	case <-p.ctx.Done():
		return
	case p.jobChan <- job:
		// Job submitted
	}
}

// Results returns a channel that provides access to job results
func (p *Pool[I, O]) Results() <-chan Result[O] {
	return p.resultChan
}

// Errors returns a channel that provides access to errors
func (p *Pool[I, O]) Errors() <-chan error {
	return p.errorChan
}

// Close terminates the pool, waiting for all workers to finish
func (p *Pool[I, O]) Close() {
	// Signal no more jobs will be submitted
	close(p.jobChan)

	// Wait for all workers to finish
	<-p.doneChan

	// Close result channel
	close(p.resultChan)
	close(p.errorChan)
}

// CloseAndWait closes the pool and waits for all results to be collected
func (p *Pool[I, O]) CloseAndWait() []Result[O] {
	close(p.jobChan)

	// Collect all results
	results := make([]Result[O], 0)
	for result := range p.resultChan {
		results = append(results, result)
	}

	// Wait for all workers to finish
	<-p.doneChan

	close(p.errorChan)
	return results
}

// Cancel immediately stops all workers
func (p *Pool[I, O]) Cancel() {
	p.cancelFunc()
	<-p.doneChan

	// Clean up channels
	close(p.resultChan)
	close(p.errorChan)
}

// GetStats returns a copy of the current pool statistics
func (p *Pool[I, O]) GetStats() PoolStats {
	p.stats.mu.Lock()
	defer p.stats.mu.Unlock()
	return *p.stats
}

// updateStats atomically updates pool statistics
func (p *Pool[I, O]) updateStats(success bool, duration time.Duration) {
	atomic.AddInt64(&p.stats.JobsProcessed, 1)
	if success {
		atomic.AddInt64(&p.stats.JobsSucceeded, 1)
	} else {
		atomic.AddInt64(&p.stats.JobsFailed, 1)
	}

	p.stats.mu.Lock()
	p.stats.TotalProcessingTime += duration
	p.stats.mu.Unlock()
}

// ProcessBatch submits a batch of jobs and waits for completion
func (p *Pool[I, O]) ProcessBatch(items []I) []Result[O] {
	// Submit all items
	for _, item := range items {
		p.Submit(item)
	}

	// No more items will be submitted
	close(p.jobChan)

	// Collect all results
	results := make([]Result[O], 0, len(items))
	for result := range p.resultChan {
		results = append(results, result)
	}

	// Wait for all workers to finish
	<-p.doneChan

	close(p.errorChan)
	return results
}

// start begins the worker processing loop
func (w *worker[I, O]) start() {
	for {
		select {
		case <-w.ctx.Done():
			return
		case job, ok := <-w.jobChan:
			if !ok {
				return // Channel closed
			}
			startTime := time.Now()
			output, err := w.workerFunc(w.ctx, job)
			duration := time.Since(startTime)

			result := Result[O]{
				JobID:    job.ID,
				Output:   output,
				Error:    err,
				Duration: duration,
			}

			select {
			case <-w.ctx.Done():
				return
			case w.resultChan <- result:
				// Result delivered
			}

			if err != nil {
				select {
				case <-w.ctx.Done():
					return
				case w.errorChan <- err:
					// Error reported
				default:
					// Non-blocking error reporting
				}
			}
		}
	}
}