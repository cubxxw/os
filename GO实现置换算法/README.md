# 替换算法实验

**计算并输出下述各种算法在内存容量为3块、4块下的缺页率。**

![操作系统_页面_2](https://s2.loli.net/2022/05/27/vTR7BSYbIe4iNk8.jpg)

![操作系统_页面_3](https://s2.loli.net/2022/05/27/2BSZLiCa9J5E4mt.jpg)

## ① 先进先出的算法（FIFO）  要求用数组或链表方法实现

#### FIFO缓存淘汰算法的实现(Go语言实现)

![image-20220526223420487](https://s2.loli.net/2022/05/26/8bOAXrRNzIU1emv.png)

**如上图示，实现 fifo算法 的缓存架构图：**

fifo 算法是淘汰缓存中最早添加的记录,即一个数据最先进入缓存，那么也应该最先被删除掉(队列先进先出)。
算法的实现比较简单：创建一个队列（一般通过双链表实现），新增记录添加到队首，淘汰队尾记录

1. map 用来存储键值对。这是实现缓存最简单直接的数据结构，因为它的查找记录和增加记录时间复杂度都是 O(1)

2. list.List 是go标准库提供的双链表。

   > 通过这个双链表存放具体的值，移动任意记录到队首的时间复杂度都是 O(1)，
   > 在队首增加记录的时间复杂度是 O(1)，删除任意一条记录的时间复杂度是 O(1)

**FIFO代码实现如下：**

**需要用到go_package中的list包，介绍：**

**LISt概述[¶](https://pkg.go.dev/container/list?utm_source=gopls#pkg-overview)**

包列表实现了一个双向链表。

遍历一个列表（其中 l 是一个 *List）：

```go
if  e := l.Front(); e != nil e = e.Next() { 
	// 用 e.Value 做一些事情
}
```

**还用到了runtime包，包的介绍：**

> Package runtime contains operations that interact with Go's runtime system, such as functions to control goroutines. It also includes the low-level type information used by the reflect package; see reflect's documentation for the programmable interface to the run-time type system.
>
> 包运行时包含与Go运行时系统交互的操作，比如控制goroutines的函数。它还包括反射包使用的低级类型信息;有关运行时类型系统的可编程接口，请参阅reflect的文档。

**runtime 调度器是个非常有用的东西，关于 runtime 包几个方法:**

- **NumCPU：**返回当前系统的 CPU 核数量
- **GOMAXPROCS：**设置最大的可同时使用的 CPU 核数
  通过runtime.GOMAXPROCS函数，应用程序何以在运行期间设置运行时系统中得P最大数量。但这会引起“Stop the World”。所以，应在应用程序最早的调用。并且最好是在运行Go程序之前设置好操作程序的环境变量GOMAXPROCS，而不是在程序中调用runtime.GOMAXPROCS函数。
  无论我们传递给函数的整数值是什么值，运行时系统的P最大值总会在1~256之间。

> go1.8后，默认让程序运行在多个核上,可以不用设置了
> go1.8前，还是要设置一下，可以更高效的利益cpu

- **Gosched：**让当前线程让出 cpu 以让其它线程运行,它不会挂起当前线程，因此当前线程未来会继续执行
  这个函数的作用是让当前 goroutine 让出 CPU，当一个 goroutine 发生阻塞，Go 会自动地把与该 goroutine 处于同一系统线程的其他 goroutine 转移到另一个系统线程上去，以使这些 goroutine 不阻塞。
- **Goexit：**退出当前 goroutine(但是defer语句会照常执行)
- **NumGoroutine：**返回正在执行和排队的任务总数
  runtime.NumGoroutine函数在被调用后，会返回系统中的处于特定状态的Goroutine的数量。这里的特指是指Grunnable\Gruning\Gsyscall\Gwaition。处于这些状态的Groutine即被看做是活跃的或者说正在被调度。
  注意：垃圾回收所在Groutine的状态也处于这个范围内的话，也会被纳入该计数器。
- **GOOS：**目标操作系统
- **runtime.GC：**会让运行时系统进行一次强制性的垃圾收集
  1.强制的垃圾回收：不管怎样，都要进行的垃圾回收。2.非强制的垃圾回收：只会在一定条件下进行的垃圾回收（即运行时，系统自上次垃圾回收之后新申请的堆内存的单元（也成为单元增量）达到指定的数值）。
- **GOROOT：**获取goroot目录
- **GOOS :** 查看目标操作系统 很多时候，我们会根据平台的不同实现不同的操作，就而已用GOOS了

**为我们创建一个FIFO的队列有一个很好的帮助**

```go
// TODO: 定义cache接口
type Cache interface {
	// 设置/添加一个缓存，如果key存在，则用新值覆盖旧值
	Set(key string, value interface{})
	// 通过key获取一个缓存值
	Get(key string) interface{}
	// 通过key删除一个缓存值
	Del(key string)
	// 删除 '最无用' 的一个缓存值
	DelOldest()
	// 获取缓存已存在的元素个数
	Len() int
	// 缓存中 元素 已经所占用内存的大小
	UseBytes() int
}

// TODO: 结构体，数组，切片，map,要求实现 Value 接口，该接口只有1个 Len 方法，返回占用内存的字节数
type Value interface {
	Len() int
}

// TODO: 定义fifo结构体
type fifo struct {
	// 缓存最大容量，单位字节
	// groupCache 使用的是最大存放 entry个数
	maxBytes int

	// 已使用的字节数，只包括值， key不算
	usedBytes int

	// 双链表
	ll *list.List
	// map的key是字符串，value是双链表中对应节点的指针
	cache map[string]*list.Element
}

// TODO: 定义key,value 结构
type entry struct {
	key   string
	value interface{}
}

// TODO: 计算出元素占用内存字节数
func (e *entry) Len() int {
	return CalcLen(e.value)
}

// TODO: 计算value占用内存大小
func CalcLen(value interface{}) int {
	var n int
	switch v := value.(type) {
	case Value: // 结构体，数组，切片，map,要求实现 Value 接口，该接口只有1个 Len 方法，返回占用的内存字节数，如果没有实现该接口，则panic
		n = v.Len()
	case string:
		if runtime.GOARCH == "amd64" {
			n = 16 + len(v)
		} else {
			n = 8 + len(v)
		}
	case bool, int8, uint8:
		n = 1
	case int16, uint16:
		n = 2
	case int32, uint32, float32:
		n = 4
	case int64, uint64, float64:
		n = 8
	case int, uint:
		if runtime.GOARCH == "amd64" {
			n = 8
		} else {
			n = 4
		}
	case complex64:
		n = 8
	case complex128:
		n = 16
	default:
		panic(fmt.Sprintf("%T is not implement cache.value", value))
	}

	return n
}

// TODO: 构造函数，创建一个新 Cache，如果 maxBytes 是0，则表示没有容量限制
func NewFifoCache(maxBytes int) Cache {
	return &fifo{
		maxBytes: maxBytes,
		ll:       list.New(),
		cache:    make(map[string]*list.Element),
	}
}

// TODO: 通过 Set 方法往 Cache 头部增加一个元素（如果已经存在，则移到头部，并修改值）
func (f *fifo) Set(key string, value interface{}) {
	if element, ok := f.cache[key]; ok {
		f.ll.MoveToFront(element)
		eVal := element.Value.(*entry)
		f.usedBytes = f.usedBytes - CalcLen(eVal.value) + CalcLen(value) // 更新占用内存大小
		element.Value = value
	} else {
		element := &entry{key, value}
		e := f.ll.PushFront(element) // 头部插入一个元素并返回该元素
		f.cache[key] = e

		f.usedBytes += element.Len()
	}

	// 如果超出内存长度，则删除队首的节点
	for f.maxBytes > 0 && f.maxBytes < f.usedBytes {
		f.DelOldest()
	}
}

// TODO: 获取指定元素
func (f *fifo) Get(key string) interface{} {
	if e, ok := f.cache[key]; ok {
		return e.Value.(*entry).value
	}

	return nil
}

// TODO: 删除指定元素
func (f *fifo) Del(key string) {
	if e, ok := f.cache[key]; ok {
		f.removeElement(e)
	}
}

// TODO: 删除最 '无用' 元素
func (f *fifo) DelOldest() {
	f.removeElement(f.ll.Back())
}

// TODO: 删除元素并更新内存占用大小
func (f *fifo) removeElement(e *list.Element) {
	if e == nil {
		return
	}

	f.ll.Remove(e)
	en := e.Value.(*entry)
	f.usedBytes -= en.Len()
	delete(f.cache, en.key)
}

// TODO: 缓存中元素个数
func (f *fifo) Len() int {
	return f.ll.Len()
}

// TODO: 缓存池占用内存大小
func (f *fifo) UseBytes() int {
	return f.usedBytes
}
```

**测试**

```go
测试：
func TestFifoCache(t *testing.T) {
	cache := NewFifoCache(512)

	key := "k1"
	cache.Set(key, 1)
	fmt.Printf("cache 元素个数：%d, 占用内存 %d 字节\n\n", cache.Len(), cache.UseBytes())

	val := cache.Get(key)
	fmt.Println(cmp.Equal(val, 1))
	cache.DelOldest()
	fmt.Printf("cache 元素个数：%d, 占用内存 %d 字节\n\n", cache.Len(), cache.UseBytes())
}
----------------------------------------------------------------------------------------------------------
结果：
=== RUN   TestFifoCache
cache 元素个数：1, 占用内存 8 字节

true
cache 元素个数：0, 占用内存 0 字节

--- PASS: TestFifoCache (0.00s)
PASS
```



## ② 最近最少使用算法（LRU） 要求用计数器或堆栈方法实现