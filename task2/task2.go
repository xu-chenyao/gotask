package main

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

// ### ✅指针

// **题目 1：** 编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，然后在主函数中调用该函数并输出修改后的值。
// **考察点：** 指针的使用、值传递与引用传递的区别。

// addTenToPointer 函数接收一个整数指针 (*int)
func addTenToPointer(numPtr *int) {
	// 解引用指针，修改指针指向的值
	*numPtr += 10
	fmt.Printf("Inside function: Value pointed to is %d\n", *numPtr)
}

func task1() {
	myNumber := 5
	fmt.Printf("Before function call: myNumber = %d\n", myNumber)

	// 传递 myNumber 的内存地址给函数
	addTenToPointer(&myNumber)

	fmt.Printf("After function call: myNumber = %d\n", myNumber)
}

// **题目 2：** 实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
// **考察点：** 指针运算、切片操作。

// doubleSliceElements 接收一个整数切片的指针 (*[]int)
func doubleSliceElements(slicePtr *[]int) {
	// 解引用指针，获取底层的切片
	// 这里的 *slicePtr 会得到 []int 类型的值
	// 然后我们遍历这个切片，对每个元素进行修改
	for i := range *slicePtr {
		(*slicePtr)[i] *= 2 // 注意这里解引用和索引的优先级
	}
	fmt.Printf("Inside function: Modified slice = %v\n", *slicePtr)
}

func task2() {
	mySlice := []int{1, 2, 3, 4, 5}
	fmt.Printf("Before function call: mySlice = %v\n", mySlice)

	// 传递 mySlice 的内存地址（即指向切片头部的指针）
	doubleSliceElements(&mySlice)

	fmt.Printf("After function call: mySlice = %v\n", mySlice)
}

// ### ✅Goroutine

// **题目 1：** 编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
// **考察点：** `go` 关键字的使用、协程的并发执行。

// printOddNumbers 协程函数，打印奇数
func printOddNumbers(wg *sync.WaitGroup) {
	defer wg.Done() // 函数结束时通知 WaitGroup
	for i := 1; i <= 10; i += 2 {
		fmt.Printf("Odd: %d\n", i)
	}
}

// printEvenNumbers 协程函数，打印偶数
func printEvenNumbers(wg *sync.WaitGroup) {
	defer wg.Done() // 函数结束时通知 WaitGroup
	for i := 2; i <= 10; i += 2 {
		fmt.Printf("Even: %d\n", i)
	}
}

func task3() {
	var wg sync.WaitGroup // 声明一个 WaitGroup 来等待所有协程完成

	// 启动第一个协程
	wg.Add(1) // 增加一个计数器
	go printOddNumbers(&wg)

	// 启动第二个协程
	wg.Add(1) // 增加一个计数器
	go printEvenNumbers(&wg)

	// 等待所有协程完成
	wg.Wait()
	fmt.Println("Done printing numbers.")
}

// **题目 2：** 设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
// **考察点：** 协程原理、并发任务调度。

// Task 定义一个任务类型，这里简化为接收一个名称的函数
type Task func(name string)

// scheduleTask 调度并执行单个任务，并测量执行时间
func scheduleTask(wg *sync.WaitGroup, task Task, taskName string) {
	defer wg.Done() // 任务完成时通知 WaitGroup

	startTime := time.Now() // 记录开始时间
	task(taskName)          // 执行任务
	endTime := time.Now()   // 记录结束时间

	duration := endTime.Sub(startTime) // 计算持续时间
	fmt.Printf("Task '%s' finished in %v\n", taskName, duration)
}

func task4() {
	var wg sync.WaitGroup

	// 定义一些模拟任务
	task1 := func(name string) {
		fmt.Printf("Running %s...\n", name)
		time.Sleep(100 * time.Millisecond) // 模拟任务执行
	}

	task2 := func(name string) {
		fmt.Printf("Running %s...\n", name)
		time.Sleep(150 * time.Millisecond)
	}

	task3 := func(name string) {
		fmt.Printf("Running %s...\n", name)
		time.Sleep(80 * time.Millisecond)
	}

	// 调度任务并发执行
	tasks := map[string]Task{
		"Task Alpha": task1,
		"Task Beta":  task2,
		"Task Gamma": task3,
	}

	fmt.Println("Starting task scheduler...")
	for name, task := range tasks {
		wg.Add(1) // 为每个任务增加 WaitGroup 计数
		go scheduleTask(&wg, task, name)
	}

	// 等待所有任务完成
	wg.Wait()
	fmt.Println("All tasks completed.")
}

// ### ✅面向对象

// **题目 1：** 定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
// **考察点：** 接口的定义与实现、面向对象编程风格。

// Shape 接口定义了图形的通用行为
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Rectangle 结构体表示一个矩形
type Rectangle struct {
	Width  float64
	Height float64
}

// Area 方法实现了 Shape 接口的 Area()
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

// Perimeter 方法实现了 Shape 接口的 Perimeter()
func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// Circle 结构体表示一个圆形
type Circle struct {
	Radius float64
}

// Area 方法实现了 Shape 接口的 Area()
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

// Perimeter 方法实现了 Shape 接口的 Perimeter() (圆周长)
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

func task5() {
	// 创建 Rectangle 实例
	rect := Rectangle{Width: 10, Height: 5}
	// 创建 Circle 实例
	circ := Circle{Radius: 7}

	// 定义一个 Shape 类型的切片，可以存储任何实现了 Shape 接口的类型
	shapes := []Shape{rect, circ}

	// 遍历 shapes 切片，调用每个 Shape 的方法
	for _, s := range shapes {
		fmt.Printf("Shape Type: %T\n", s) // 输出具体类型
		fmt.Printf("  Area: %.2f\n", s.Area())
		fmt.Printf("  Perimeter: %.2f\n", s.Perimeter())
		fmt.Println("---")
	}
}

// **题目 2：** 使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
// **考察点：** 组合的使用、方法接收者。

// Person 结构体定义了人的基本信息
type Person struct {
	Name string
	Age  int
}

// Employee 结构体通过组合 Person 来复用其字段
type Employee struct {
	Person     // 嵌入 Person 结构体，实现了组合
	EmployeeID string
}

// PrintInfo 方法为 Employee 结构体定义
func (e Employee) PrintInfo() {
	// 直接访问嵌入结构体的字段
	fmt.Printf("Employee Name: %s\n", e.Name)
	fmt.Printf("Employee Age: %d\n", e.Age)
	fmt.Printf("Employee ID: %s\n", e.EmployeeID)
}

func task6() {
	// 创建 Employee 实例
	emp := Employee{
		Person: Person{ // 初始化嵌入的 Person 结构体
			Name: "Alice Smith",
			Age:  30,
		},
		EmployeeID: "EMP001",
	}

	// 调用 Employee 的 PrintInfo 方法
	emp.PrintInfo()

	fmt.Println("\n--- Another Employee ---")
	emp2 := Employee{
		Person:     Person{Name: "Bob Johnson", Age: 25},
		EmployeeID: "EMP002",
	}
	emp2.PrintInfo()
}

// ### ✅Channel

// **题目 1：** 编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
// **考察点：** 通道的基本使用、协程间通信。

// producer 协程：生成整数并发送到通道
func producer(ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done() // 确保协程完成时通知 WaitGroup
	for i := 1; i <= 10; i++ {
		ch <- i // 将整数发送到通道
		fmt.Printf("Producer: Sent %d\n", i)
	}
	close(ch) // 发送完所有数据后关闭通道，通知消费者没有更多数据了
	fmt.Println("Producer: Channel closed.")
}

// consumer 协程：从通道接收整数并打印
func consumer(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()       // 确保协程完成时通知 WaitGroup
	for num := range ch { // 循环从通道接收数据，直到通道被关闭且所有数据都被接收
		fmt.Printf("Consumer: Received %d\n", num)
	}
	fmt.Println("Consumer: Channel exhausted.")
}

func task7() {
	// 创建一个无缓冲通道
	dataChannel := make(chan int)
	var wg sync.WaitGroup

	// 启动生产者协程
	wg.Add(1)
	go producer(dataChannel, &wg)

	// 启动消费者协程
	wg.Add(1)
	go consumer(dataChannel, &wg)

	// 等待所有协程完成
	wg.Wait()
	fmt.Println("Program finished.")
}

// **题目 2：** 实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
// **考察点：** 通道的缓冲机制。

const (
	bufferSize  = 5   // 缓冲区大小
	numMessages = 100 // 要发送的消息数量
)

// bufferedProducer 协程：生成整数并发送到带缓冲的通道
func bufferedProducer(ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 1; i <= numMessages; i++ {
		ch <- i // 尝试发送，如果通道满会阻塞
		fmt.Printf("Producer: Sent %d (Channel size: %d/%d)\n", i, len(ch), cap(ch))
		time.Sleep(10 * time.Millisecond) // 模拟生产延迟
	}
	close(ch) // 发送完所有数据后关闭通道
	fmt.Println("Producer: Channel closed.")
}

// bufferedConsumer 协程：从通道接收整数并打印
func bufferedConsumer(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		// 从通道接收数据，使用 ok 检查通道是否已关闭且无更多数据
		num, ok := <-ch
		if !ok {
			fmt.Println("Consumer: Channel exhausted.")
			return // 通道已关闭且无数据，退出循环
		}
		fmt.Printf("Consumer: Received %d (Channel size: %d/%d)\n", num, len(ch), cap(ch))
		time.Sleep(50 * time.Millisecond) // 模拟消费延迟
	}
}

func task8() {
	// 创建一个带缓冲的通道
	bufferedChannel := make(chan int, bufferSize)
	var wg sync.WaitGroup

	// 启动生产者协程
	wg.Add(1)
	go bufferedProducer(bufferedChannel, &wg)

	// 启动消费者协程
	wg.Add(1)
	go bufferedConsumer(bufferedChannel, &wg)

	// 等待所有协程完成
	wg.Wait()
	fmt.Println("Program finished.")
}

// ### ✅锁机制

// **题目 1：** 编写一个程序，使用 `sync.Mutex` 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
// **考察点：** `sync.Mutex` 的使用、并发数据安全。

var (
	sharedCounter int        // 共享的计数器
	mutex         sync.Mutex // 互斥锁，用于保护 sharedCounter
)

const (
	numGoroutines          = 10
	incrementsPerGoroutine = 1000
)

// incrementCounter safely increments the sharedCounter
func incrementCounter(wg *sync.WaitGroup) {
	defer wg.Done() // 协程结束时通知 WaitGroup

	for i := 0; i < incrementsPerGoroutine; i++ {
		mutex.Lock()    // 获取锁，阻塞其他协程访问 sharedCounter
		sharedCounter++ // 临界区：安全地修改共享变量
		mutex.Unlock()  // 释放锁
	}
}

func task9() {
	var wg sync.WaitGroup

	fmt.Println("Starting counter program with Mutex...")

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1) // 增加 WaitGroup 计数
		go incrementCounter(&wg)
	}

	wg.Wait() // 等待所有协程完成

	// 理论上最终计数器值应为 numGoroutines * incrementsPerGoroutine
	expectedValue := numGoroutines * incrementsPerGoroutine
	fmt.Printf("Final Counter Value: %d (Expected: %d)\n", sharedCounter, expectedValue)

	if sharedCounter == expectedValue {
		fmt.Println("Counter is correct, Mutex worked!")
	} else {
		fmt.Println("Error: Counter is incorrect!")
	}
}

// **题目 2：** 使用原子操作（`sync/atomic` 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
// **考察点：** 原子操作、并发数据安全。

var (
	// 使用 atomic.Int64 类型来存储计数器，它提供了原子操作
	// 或者使用 uint64，取决于计数器是否可能为负数。这里使用 uint64 也可以。
	// 这里用 int64 类型作为例子，因为 atomic.AddInt64 可以直接操作它。
	atomicCounter int64
)

const (
	numGoroutinesAtomic          = 10
	incrementsPerGoroutineAtomic = 1000
)

// incrementCounterAtomic safely increments the atomicCounter using atomic operations
func incrementCounterAtomic(wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < incrementsPerGoroutineAtomic; i++ {
		// 使用 atomic.AddInt64 进行原子递增操作
		// 它保证了在多协程并发访问时，递增操作的完整性
		atomic.AddInt64(&atomicCounter, 1)
	}
}

func task10() {
	var wg sync.WaitGroup

	fmt.Println("Starting counter program with Atomic operations...")

	for i := 0; i < numGoroutinesAtomic; i++ {
		wg.Add(1)
		go incrementCounterAtomic(&wg)
	}

	wg.Wait()

	// 使用 atomic.LoadInt64 来原子地读取计数器的最终值
	finalValue := atomic.LoadInt64(&atomicCounter)
	expectedValue := int64(numGoroutinesAtomic * incrementsPerGoroutineAtomic)
	fmt.Printf("Final Counter Value: %d (Expected: %d)\n", finalValue, expectedValue)

	if finalValue == expectedValue {
		fmt.Println("Counter is correct, Atomic operations worked!")
	} else {
		fmt.Println("Error: Counter is incorrect!")
	}
}

func main() {
	task1()
	task2()
	task3()
	task4()
	task5()
	task6()
	task7()
	task8()
	task9()
	task10()
}
