package channels

import "fmt"

func ExampleWorkers() {
	input := make(chan int, 10)
	for i := 1; i <= 10; i++ {
		input <- i
	}
	close(input)

	squareFn := func(n int) int {
		return n * n
	}

	output := Workers(input, 5, 3, squareFn)

	for result := range output {
		fmt.Println(result) // Outputs squares of input numbers
	}
	//In this example, three worker goroutines are processing the input channel concurrently,
	// applying the `squareFn` function to each item, and sending the results to the output channel.
}
