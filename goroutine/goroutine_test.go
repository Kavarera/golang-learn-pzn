package goroutine

import (
	"fmt"
	"testing"
	"time"
)

func CreateGoroutine() {
	fmt.Println("Hello, 	Goroutine!")
}

func TestCreateGoroutine(t *testing.T) {
	go CreateGoroutine()

	fmt.Println("Ups..")

	time.Sleep(1 * time.Second)
}

func DisplayNumber(number int) {
	fmt.Println(number)
}

func TestManyGoroutine(t *testing.T) {
	for i := 0; i < 100000; i++ {
		go DisplayNumber(i)
	}
	time.Sleep(15 * time.Second)
}

func TestCreateChannel(t *testing.T) {
	channel := make(chan string)

	defer close(channel)

	go func() {
		time.Sleep(2 * time.Second)
		channel <- "Hello Channel"
		fmt.Println("Selesai mengirim data ke channel")
	}()

	data := <-channel
	fmt.Println("Menerima data dari channel:", data)

	time.Sleep(5 * time.Second)

}

func GiveMeResponse(channel chan string) {
	time.Sleep(2 * time.Second)
	channel <- "Hello Channel"
	fmt.Println("Selesai mengirim data ke channel")
}

func TestChannelAsParameter(t *testing.T) {
	channel := make(chan string)

	defer close(channel)

	go GiveMeResponse(channel)

	data := <-channel
	fmt.Println("Menerima data dari channel:", data)

	time.Sleep(5 * time.Second)
}

func OnlyIn(channel chan<- string) {
	time.Sleep(2 * time.Second)
	channel <- "Hello Channel"
	fmt.Println("Selesai mengirim data ke channel")
}

func OnlyOut(channel <-chan string) {
	data := <-channel
	fmt.Println("Menerima data dari channel:", data)
}

func TestInOutChannel(t *testing.T) {
	channel := make(chan string)

	defer close(channel)

	go OnlyIn(channel)

	go OnlyOut(channel)

	time.Sleep(5 * time.Second)
}

func TestBufferedChannel(t *testing.T) {
	channel := make(chan string, 3)

	defer close(channel)

	channel <- "Hello"

	fmt.Println("Satu data sudah masuk ke channel")
}
