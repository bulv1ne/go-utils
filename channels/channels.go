package channels

import (
	"iter"
	"sync"
)

// PipeSeqToChannel pipes elements from a given sequence to a channel.
//
// This generic function takes a sequence of elements of any type (T),
// processes the sequence in a sequential manner, and pipes the elements
// into a buffered channel of the specified size. This function allows
// for decoupling of producers (the sequence) and consumers (the channel).
//
// Parameters:
//   - input: A sequence of type iter.Seq[T] that will be iterated over and
//     its elements sent into the channel.
//   - size: An integer representing the buffer size of the channel to be created.
//
// Returns:
//
//	A buffered channel of type chan T, where T is the same type as the elements
//	in the input sequence.
//
// Example usage:
//
//	seq := slices.Values([]int{1, 2, 3, 4, 5})
//	ch := PipeSeqToChannel(seq, 10)
//	for val := range ch {
//	    fmt.Println(val)
//	}
func PipeSeqToChannel[T any](input iter.Seq[T], size int) <-chan T {
	output := make(chan T, size)
	go func() {
		for row := range input {
			output <- row
		}
		close(output)
	}()
	return output
}

// MergeChannels merges multiple read-only channels into a single output channel.
//
// This function takes a variadic number of input channels (of any type) and merges their output
// into a single channel. The returned channel will receive values from all the input channels
// as they become available. The merged output channel is closed once all input channels are closed.
//
// Parameters:
// - chs: A variadic number of read-only channels from which the function will read values.
//
// Returns:
//   - A single read-only channel of the same type as the input channels. This channel will receive
//     all values from the input channels until all of them are closed.
//
// Example usage:
//
//	ch1 := make(chan int)
//	ch2 := make(chan int)
//	ch3 := MergeChannels(ch1, ch2)
//
//	// Write to ch1 and ch2 in separate goroutines, read from ch3
//	go func() {
//	    ch1 <- 1
//	    close(ch1)
//	}()
//	go func() {
//	    ch2 <- 2
//	    close(ch2)
//	}()
//	for val := range ch3 {
//	    fmt.Println(val) // Outputs: 1 or 2, depending on channel read order
//	}
func MergeChannels[T any](chs ...<-chan T) <-chan T {
	output := make(chan T, 1)
	wg := sync.WaitGroup{}
	wg.Add(len(chs))
	for _, ch := range chs {
		go func() {
			defer wg.Done()
			for v := range ch {
				output <- v
			}
		}()
	}
	go func() {
		wg.Wait()
		close(output)
	}()
	return output
}
