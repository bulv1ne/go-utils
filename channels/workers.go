package channels

import "sync"

// Workers creates a pool of worker goroutines to process items from an input channel and returns
// a channel that collects the results.
//
// This function launches a specified number of worker goroutines (`workerCount`), each reading
// from the input channel (`input`). Each worker applies the provided function (`fn`) to the
// received values and sends the result to the output channel. The size of the output channel
// is defined by `channelSize`. The output channel will be closed when all workers have completed
// their tasks and the input channel is fully processed.
//
// Parameters:
// - input: A read-only channel from which values of type `T1` are received.
// - channelSize: The size of the output channel (buffered channel), which controls how many results can be stored before being read.
// - workerCount: The number of worker goroutines to process items concurrently.
// - fn: A function that takes a value of type `T1` and transforms it to a value of type `T2`.
//
// Returns:
//   - A read-only channel of type `T2` that will receive the results of applying the function `fn`
//     to each item from the input channel. This channel will be closed once all input values have been
//     processed.
func Workers[T1, T2 any](input <-chan T1, channelSize int, workerCount int, fn func(T1) T2) <-chan T2 {
	if channelSize < 1 {
		panic("channel size < 1")
	}
	if workerCount < 1 {
		panic("worker count < 1")
	}

	output := make(chan T2, channelSize)
	wg := sync.WaitGroup{}
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go func() {
			defer wg.Done()
			for v := range input {
				result := fn(v)
				output <- result
			}
		}()
	}
	go func() {
		wg.Wait()
		close(output)
	}()
	return output
}
