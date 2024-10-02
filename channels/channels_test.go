package channels

import (
	"fmt"
	"slices"
	"testing"
)

func TestPipeSeqToChannel(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	seq := slices.Values(slice)

	channel := PipeSeqToChannel(seq, 1)
	value, ok := <-channel
	if !ok {
		t.Fatal("Expected channel to receive value, but it did not")
	}
	if value != 1 {
		t.Fatalf("Expected 1, got %d", value)
	}
	value, ok = <-channel
	if !ok {
		t.Fatal("Expected channel to receive value, but it did not")
	}
	if value != 2 {
		t.Fatalf("Expected 2, got %d", value)
	}
	value, ok = <-channel
	if !ok {
		t.Fatal("Expected channel to receive value, but it did not")
	}
	if value != 3 {
		t.Fatalf("Expected 3, got %d", value)
	}

	remainder := []int{}
	for i := range channel {
		remainder = append(remainder, i)
	}
	if !slices.Equal(remainder, []int{4, 5, 6, 7, 8, 9}) {
		t.Fatalf("Expected remainder, got %v", remainder)
	}
}

func DownloadS3File(value int) int {
	return value + 1
}

func ProcessS3File(value int) int {
	return value + 1
}

func UploadS3File(value int) int {
	return value + 1
}

func SplitS3Keys(keys <-chan int) (<-chan int, <-chan int) {
	smallFiles := make(chan int, 10)
	bigFiles := make(chan int, 10)

	go func() {
		for key := range keys {
			if key%2 == 0 {
				smallFiles <- key
			} else {
				bigFiles <- key
			}
		}
		close(smallFiles)
		close(bigFiles)
	}()

	return smallFiles, bigFiles
}

func TestPipe(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	seq := slices.Values(slice)
	s3Keys := PipeSeqToChannel(seq, 5)

	small, big := SplitS3Keys(s3Keys)

	smallFiles := Workers(small, 1, 20, DownloadS3File)
	bigFiles := Workers(big, 1, 2, DownloadS3File)

	files := MergeChannels(smallFiles, bigFiles)

	processedFiles := Workers(files, 1, 3, ProcessS3File)
	uploadedFiles := Workers(processedFiles, 1, 10, UploadS3File)

	result := []int{}
	for i := range uploadedFiles {
		result = append(result, i)
	}
	slices.Sort(result)
	expected := []int{4, 5, 6, 7, 8, 9, 10, 11, 12}
	if !slices.Equal(expected, result) {
		t.Fatalf("Expected %v, got %v", expected, result)
	}
}

func TestPipeNoThread(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}

	for _, i := range slice {
		file := DownloadS3File(i)
		processedFile := ProcessS3File(file)
		uploadedFile := UploadS3File(processedFile)
		fmt.Println("Uploaded file", uploadedFile)
	}
}
