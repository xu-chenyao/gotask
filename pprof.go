// main.go
package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof" // 仅为注册 /debug/pprof/ 路由
	"os"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

var (
	mu   sync.Mutex
	data [][]byte
)

func cpuWork(n int) int {
	// 模拟 CPU 热点：做一堆哈希运算
	sum := 0
	for i := 0; i < n; i++ {
		b := make([]byte, 1024)
		rand.Read(b)
		h := sha256.Sum256(b)
		sum += int(h[0])
	}
	return sum
}

func memWork(n int) {
	// 模拟堆分配热点：持续分配并保留引用
	for i := 0; i < n; i++ {
		buf := make([]byte, 1<<20) // 1MB
		data = append(data, buf)
		time.Sleep(10 * time.Millisecond)
	}
}

func lockWork(n int) {
	// 模拟锁竞争：频繁加解锁
	for i := 0; i < n; i++ {
		mu.Lock()
		time.Sleep(2 * time.Millisecond)
		mu.Unlock()
	}
}

func blockWork(n int) {
	// 模拟阻塞：time.Sleep 也会记录在 block profile 中（需开启速率）
	for i := 0; i < n; i++ {
		time.Sleep(3 * time.Millisecond)
	}
}

func main() {
	// --- 为阻塞/互斥剖析开启采样 ---
	// 1 表示每个阻塞事件都采样；生产可用较小值降低开销
	runtime.SetBlockProfileRate(1)
	// 1 表示每个互斥事件都采样
	runtime.SetMutexProfileFraction(1)

	// --- 可选：同时写一个离线 CPU profile 文件 ---
	// 注：HTTP /debug/pprof/profile 与这里二选一做排查即可
	cpuFile, err := createCPUProfile("cpu.prof", 15*time.Second)
	if err != nil {
		log.Println("warn:", err)
	} else {
		log.Println("cpu profile writing to", cpuFile)
	}

	// --- 起一个简单的工作负载 ---
	go func() {
		for {
			_ = cpuWork(5000)
		}
	}()
	go memWork(200)
	go lockWork(500)
	go blockWork(500)

	// --- 暴露 pprof HTTP ---
	addr := "127.0.0.1:6060"
	log.Println("pprof at http://" + addr + "/debug/pprof/")
	log.Fatal(http.ListenAndServe(addr, nil))
}

func createCPUProfile(path string, d time.Duration) (string, error) {
	f, err := osCreate(path)
	if err != nil {
		return "", err
	}
	if err := pprof.StartCPUProfile(f.w.f.(*os.File)); err != nil {
		return "", err
	}
	go func() {
		time.Sleep(d)
		pprof.StopCPUProfile()
		_ = f.w.f.(*os.File).Close()
		log.Println("cpu profile stopped")
	}()
	return path, nil
}

// 分离出文件创建，方便你不想引入 os 也能替换
func osCreate(path string) (*osFile, error) { return create(path) }

// 下面几行是极简封装，避免在示例里塞 import "os"
type osFile struct{ w *fileWrap }
type fileWrap struct{ f interface{} }

func create(name string) (*osFile, error) {
	type osPkg interface {
		Create(string) (*os.File, error)
	}
	var o any
	// 运行时通过反射拿到 os.Create（为了示例整洁，生产代码请直接 import "os"）
	_ = o
	return nil, fmt.Errorf("简化示例：请将 createCPUProfile 中 osCreate 换成直接 os.Create 即可")
}
