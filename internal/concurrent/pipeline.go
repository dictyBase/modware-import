package concurrent

import (
	"context"
	"sync"
)

// BatchProcessor runs a worker pool for batch processing
type BatchProcessor[I, O any] struct {
	Pool         *Pool[I, O]
	BatchSize    int
	batchCount   int
	currentBatch []I
	batchMutex   sync.Mutex // For thread-safe batch operations
}

// NewBatchProcessor creates a new BatchProcessor
func NewBatchProcessor[I, O any](
	workerFunc WorkerFunc[I, O],
	batchSize int,
	options ...PoolOption[I, O],
) *BatchProcessor[I, O] {
	return &BatchProcessor[I, O]{
		Pool:         NewPool(workerFunc, options...),
		BatchSize:    batchSize,
		currentBatch: make([]I, 0, batchSize),
	}
}

// Add adds an item to the current batch, automatically submitting
// the batch when it reaches the configured size
func (bp *BatchProcessor[I, O]) Add(item I) bool {
	bp.batchMutex.Lock()
	defer bp.batchMutex.Unlock()

	bp.currentBatch = append(bp.currentBatch, item)

	if len(bp.currentBatch) >= bp.BatchSize {
		bp.submitCurrentBatchLocked()
		return true
	}
	return false
}

// AddBatch adds multiple items to the current batch, submitting batches
// as they reach the configured size
func (bp *BatchProcessor[I, O]) AddBatch(items []I) int {
	if len(items) == 0 {
		return 0
	}

	bp.batchMutex.Lock()
	defer bp.batchMutex.Unlock()

	batchesSubmitted := 0
	for _, item := range items {
		bp.currentBatch = append(bp.currentBatch, item)

		if len(bp.currentBatch) >= bp.BatchSize {
			bp.submitCurrentBatchLocked()
			batchesSubmitted++
		}
	}

	return batchesSubmitted
}

// submitCurrentBatchLocked submits the current batch of items
// Assumes the caller holds the batchMutex
func (bp *BatchProcessor[I, O]) submitCurrentBatchLocked() {
	if len(bp.currentBatch) == 0 {
		return
	}

	// Create a copy of the current batch
	batch := make([]I, len(bp.currentBatch))
	copy(batch, bp.currentBatch)

	// Submit each item in the batch individually
	for _, item := range batch {
		bp.Pool.Submit(item, map[string]interface{}{
			"batch_number": bp.batchCount,
			"batch_size":   len(batch),
		})
	}

	// Reset the current batch
	bp.currentBatch = make([]I, 0, bp.BatchSize)
	bp.batchCount++
}

// submitCurrentBatch submits the current batch of items
func (bp *BatchProcessor[I, O]) submitCurrentBatch() {
	bp.batchMutex.Lock()
	defer bp.batchMutex.Unlock()
	bp.submitCurrentBatchLocked()
}

// Flush submits the current batch even if it's not full
func (bp *BatchProcessor[I, O]) Flush() {
	bp.submitCurrentBatch()
}

// Start starts the worker pool
func (bp *BatchProcessor[I, O]) Start() {
	bp.Pool.Start()
}

// Results returns the results channel from the pool
func (bp *BatchProcessor[I, O]) Results() <-chan Result[O] {
	return bp.Pool.Results()
}

// Errors returns the error channel from the pool
func (bp *BatchProcessor[I, O]) Errors() <-chan error {
	return bp.Pool.Errors()
}

// Close flushes the current batch and closes the pool
func (bp *BatchProcessor[I, O]) Close() {
	bp.Flush()
	bp.Pool.Close()
}

// Cancel immediately stops all processing
func (bp *BatchProcessor[I, O]) Cancel() {
	bp.Pool.Cancel()
}

// GetBatchCount returns the current batch count
func (bp *BatchProcessor[I, O]) GetBatchCount() int {
	bp.batchMutex.Lock()
	defer bp.batchMutex.Unlock()
	return bp.batchCount
}

// GetCurrentBatchSize returns the current batch size
func (bp *BatchProcessor[I, O]) GetCurrentBatchSize() int {
	bp.batchMutex.Lock()
	defer bp.batchMutex.Unlock()
	return len(bp.currentBatch)
}

// Pipeline represents a multi-stage processing pipeline
type Pipeline[I, O, R any] struct {
	Processor     *BatchProcessor[I, O]
	ResultHandler func(context.Context, <-chan Result[O]) R
	ctx           context.Context
	cancelFunc    context.CancelFunc
}

// NewPipeline creates a new processing pipeline
func NewPipeline[I, O, R any](
	processor *BatchProcessor[I, O],
	resultHandler func(context.Context, <-chan Result[O]) R,
) *Pipeline[I, O, R] {
	ctx, cancel := context.WithCancel(context.Background())
	return &Pipeline[I, O, R]{
		Processor:     processor,
		ResultHandler: resultHandler,
		ctx:           ctx,
		cancelFunc:    cancel,
	}
}

// Start initializes the pipeline and starts processing
func (p *Pipeline[I, O, R]) Start() {
	p.Processor.Start()
}

// Process processes all items and returns the final result
func (p *Pipeline[I, O, R]) Process(items []I) R {
	// Start the processor if not already started
	p.Start()
	// Submit all items in batches
	p.Processor.AddBatch(items)
	// Signal no more items
	p.Processor.Flush()
	// Process results
	result := p.ResultHandler(p.ctx, p.Processor.Results())
	// Clean up
	p.Processor.Close()

	return result
}

// Cancel stops all processing
func (p *Pipeline[I, O, R]) Cancel() {
	p.cancelFunc()
	p.Processor.Cancel()
}

// WithContext sets a custom context for the pipeline
func (p *Pipeline[I, O, R]) WithContext(
	ctx context.Context,
) *Pipeline[I, O, R] {
	if ctx == nil {
		return p
	}

	p.cancelFunc()
	p.ctx, p.cancelFunc = context.WithCancel(ctx)
	return p
}
