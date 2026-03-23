package goroutine

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
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

func TestRangeChannel(t *testing.T) {
	channel := make(chan string)

	go func() {
		for i := 0; i < 5; i++ {
			channel <- "Data ke " + strconv.Itoa(i)
		}
		close(channel)
	}()

	for data := range channel {
		fmt.Println("Menerima data dari channel:", data)
	}
}

func TestSelectChannel(t *testing.T) {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go GiveMeResponse(ch1)
	go GiveMeResponse(ch2)

	counter := 0
	for {
		select {
		case data := <-ch1:
			fmt.Println("Menerima data dari channel 1:", data)
			counter++
		case data := <-ch2:
			fmt.Println("Menerima data dari channel 2:", data)
			counter++
		}
		if counter == 2 {
			break
		}
	}
}

func TestDefaultChannel(t *testing.T) {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go GiveMeResponse(ch1)
	go GiveMeResponse(ch2)

	counter := 0
	for {
		select {
		case data := <-ch1:
			fmt.Println("Menerima data dari channel 1:", data)
			counter++
		case data := <-ch2:
			fmt.Println("Menerima data dari channel 2:", data)
			counter++
		default:
			fmt.Println("Tidak ada data yang diterima")
			counter++
		}
		if counter == 3 {
			break
		}
	}
}

func TestRaceCondition(t *testing.T) {
	var x = 0
	for i := 0; i < 100; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				x = x + 1
			}
		}()
	}

	time.Sleep(5 * time.Second)

	fmt.Println("Counter:", x)
}

// MUTEX

func TestMutex(t *testing.T) {
	var x = 0
	var mutex sync.Mutex
	for i := 0; i < 100; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				mutex.Lock()
				x = x + 1
				mutex.Unlock()
			}
		}()
	}

	time.Sleep(5 * time.Second)

	fmt.Println("Counter:", x)

}

// RW MUTEX

type BankAccount struct {
	RWMutex sync.RWMutex
	Balance int
}

func (account *BankAccount) AddBalance(amount int) {
	account.RWMutex.Lock()
	account.Balance = account.Balance + amount
	account.RWMutex.Unlock()
}

func (account *BankAccount) GetBalance() int {
	account.RWMutex.RLock()
	balance := account.Balance
	account.RWMutex.RUnlock()
	return balance
}

func TestRWMutex(t *testing.T) {
	account := BankAccount{}

	for i := 0; i < 100; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				fmt.Println("Previous Balance:", account.GetBalance())
				account.AddBalance(1)
				fmt.Println("Balance After Add:", account.GetBalance())
			}
		}()
	}

	time.Sleep(15 * time.Second)
	fmt.Println("Balance:", account.GetBalance())
}

// DeADLOCK

type UserBalance struct {
	sync.Mutex
	Name    string
	Balance int
}

func (user *UserBalance) Transfer(target *UserBalance, amount int) {
	user.Lock()
	fmt.Println("Lock user1", user.Name)
	target.Lock()
	fmt.Println("Lock user2", target.Name)
	user.Balance = user.Balance - amount
	target.Balance = target.Balance + amount
	user.Unlock()
	target.Unlock()

}

func TestDeadLock(t *testing.T) {
	user1 := UserBalance{Name: "User 1", Balance: 1000000}
	user2 := UserBalance{Name: "User 2", Balance: 1000000}

	go user1.Transfer(&user2, 100000)
	go user2.Transfer(&user1, 200000)

	fmt.Println("User 1 Balance:", user1.Balance)
	fmt.Println("User 2 Balance:", user2.Balance)

}

// Wait Group

func RunAsinkronus(group *sync.WaitGroup) {
	defer group.Done()
	fmt.Println("Run Asinkronus")
	group.Add(1)
	fmt.Println("Menunggu selesai")
	time.Sleep(2 * time.Second)
}

func TestWaitGroup(t *testing.T) {
	group := &sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		go RunAsinkronus(group)
	}
	group.Wait()
	fmt.Println("Selesai semua proses")
}

// ONCE
var counter = 0

func OnlyOnce() {
	counter++
}

func TestOnce(t *testing.T) {
	var once sync.Once
	group := &sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		go func() {
			defer group.Done()
			group.Add(1)
			once.Do(OnlyOnce)
		}()
	}

	group.Wait()
	fmt.Println("Counter:", counter)
}

// POOL
func TestPool(t *testing.T) {
	pool := sync.Pool{
		New: func() interface{} {
			return "New Data"
		},
	}

	pool.Put("Data 1")
	pool.Put("Data 2")
	pool.Put("Data 3")

	for i := 0; i < 10; i++ {
		go func() {
			data := pool.Get()
			fmt.Println(data)
			pool.Put(data)
		}()
	}
	time.Sleep(35 * time.Second)
}

// SYNC MAP

func AddToMap(syncMap *sync.Map, waitGroup *sync.WaitGroup, value int) {
	defer waitGroup.Done()
	waitGroup.Add(1)
	syncMap.Store("key"+strconv.Itoa(value), "value"+strconv.Itoa(value))

}
func TestSyncMap(t *testing.T) {
	data := &sync.Map{}
	group := &sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		go AddToMap(data, group, i)
	}

	group.Wait()

	data.Range(func(key, val interface{}) bool {
		fmt.Println(key, ":", val)
		return true
	})

}


//GOMAXPROCS

func TestGOMAXPROCS(t *testing.T) {
	fmt.Println("Jumlah CPU:", runtime.NumCPU())
	fmt.Println("GOMAXPROCS:", runtime.GOMAXPROCS(-1))
	fmt.Println("GOMAXPROCS:", runtime.NumGoroutine())

}