# Go 语言最佳实践

## 1. 命名规范

### 变量命名

- 使用 **驼峰命名** (camelCase)：`userName`, `maxCount`
- 缩写词全大写：`userID`, `APIKey` (前两个字母大写)
- 避免使用单个字母，除了循环计数器：`i`, `j`, `k`

### 常量命名

- 全大写下划线：`MAX_SIZE`, `DEFAULT_TIMEOUT`
- 枚举类型使用特定前缀：`StatusOK`, `StatusPending`

### 函数命名

- 公开函数使用 **帕斯卡命名** (PascalCase)：`ReadFile`, `NewClient`
- 私有函数使用 **驼峰命名**：`readFile`, `newClient`

## 2. 错误处理

```go
// 简单错误处理
if err != nil {
    return err
}

// 带上下文的错误
if err != nil {
    return fmt.Errorf("failed to read config: %w", err)
}

// 错误检查后提前返回
if err := doSomething(); err != nil {
    // 处理错误
}
```

### 自定义错误

```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}
```

## 3. 并发编程

### Goroutine

```go
// 启动 goroutine
go func() {
    // 工作
}()

// 使用 channel 通信
ch := make(chan int)
go func() {
    ch <- 42
}()
result := <-ch
```

### Select 语句

```go
select {
case msg := <-ch1:
    fmt.Println("received", msg)
case msg := <-ch2:
    fmt.Println("received", msg)
case <-time.After(time.Second):
    fmt.Println("timeout")
}
```

### 同步

```go
var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(i int) {
        defer wg.Done()
        // 工作
    }(i)
}
wg.Wait()
```

## 4. 性能优化

### 内存分配

```go
// 预分配切片容量
s := make([]int, 0, 100)

// 复用对象
var buf bytes.Buffer
buf.Reset()
```

### 避免不必要的分配

```go
// 使用 strings.Builder
var b strings.Builder
b.WriteString("Hello")
b.WriteString(" ")
// ...

// 使用 sync.Pool
var pool sync.Pool
pool.Put(obj)
obj = pool.Get()
```

## 5. 测试

### 单元测试

```go
func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("Add(2, 3) = %d, want %d", result, 5)
    }
}
```

### Table-Driven 测试

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        a, b, want int
    }{
        {1, 2, 3},
        {0, 0, 0},
        {-1, 1, 0},
    }
    for _, tt := range tests {
        got := Add(tt.a, tt.b)
        if got != tt.want {
            t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
        }
    }
}
```

## 6. 代码组织

### 包结构

```
myproject/
├── cmd/
│   └── myapp/
│       └── main.go
├── pkg/
│   ├── utils/
│   │   └── helper.go
│   └── models/
│       └── user.go
├── internal/
│   └── server/
│       └── server.go
└── go.mod
```

### 接口设计

```go
// 定义小而专注的接口
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

// 组合接口
type ReadWriter interface {
    Reader
    Writer
}
```

## 7. 常用工具

### linter

```bash
# golangci-lint
golangci-lint run

# 具体规则
golangci-lint run --enable=govet,staticcheck
```

### 格式化

```bash
# gofmt
gofmt -w .

# goimports
goimports -w .
```

## 8. 调试技巧

```go
// 使用 %+v 打印结构体详情
fmt.Printf("%+v\n", myStruct)

// 使用 log 包
log.Printf("debug: %+v", data)

// 使用 runtime 信息
func logMemStats() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    log.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
}
```