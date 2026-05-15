# Bubble Tea 最佳实践指南

本文档总结了使用 Bubble Tea 构建 TUI 应用时的最佳实践，结合本项目实际经验编写。

## 1. 核心架构原则

### 1.1 单向数据流

Bubble Tea 遵循 Elm 架构的单向数据流模式：

```
Model → View → Msg → Update → Model → View → ...
```

- **Model** 持有全部状态
- **View** 是纯函数，只读取 Model 并返回渲染字符串
- **Update** 接收消息，返回新 Model 和可选的 Cmd
- **Cmd** 执行副作用（I/O、HTTP 请求、定时器等）

### 1.2 三大接口

```go
type Model interface {
    Init() Cmd                    // 初始化命令
    Update(msg Msg) (Model, Cmd)  // 处理消息，返回新状态和命令
    View() string                 // 渲染视图（纯函数）
}
```

## 2. View() 必须是纯函数

### 2.1 禁止在 View() 中修改状态

**错误示例：**
```go
func (m *Model) View() string {
    m.textarea.Prompt = "❯ "  // 错误！View 中修改了状态
    return m.textarea.View()
}
```

**正确示例：**
```go
// 在初始化时设置
func NewModel() *Model {
    ta := textarea.New()
    ta.Prompt = "❯ "  // 正确：在初始化时设置
    return &Model{textarea: ta}
}

func (m *Model) View() string {
    return m.textarea.View()  // 纯函数，只读
}
```

### 2.2 View() 可能被多次调用

Bubble Tea 运行时可能在每帧多次调用 `View()`，因此：
- 不要在 `View()` 中执行昂贵计算
- 不要在 `View()` 中分配大量内存
- 不要在 `View()` 中修改任何状态

## 3. Update() 中的阻塞操作

### 3.1 禁止在 Update() 中执行阻塞操作

**错误示例：**
```go
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "enter" {
            m.inputCh <- data  // 错误！如果 channel 满了会阻塞整个 UI
        }
    }
    return m, nil
}
```

**正确示例：**
```go
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "enter" {
            return m, sendInput(m.inputCh, data)
        }
    }
    return m, nil
}

func sendInput(ch chan Data, data Data) tea.Cmd {
    return func() tea.Msg {
        ch <- data  // 在 Cmd 中执行，不阻塞 Update
        return inputSentMsg{}
    }
}
```

### 3.2 使用 tea.Cmd 包装所有副作用

所有 I/O 操作、网络请求、文件读写都应该包装在 `tea.Cmd` 中：

```go
func fetchAPI(url string) tea.Cmd {
    return func() tea.Msg {
        resp, err := http.Get(url)
        if err != nil {
            return errorMsg{err}
        }
        return dataMsg{resp}
    }
}
```

## 4. Channel 的正确使用

### 4.1 Channel 关闭时必须返回有效消息

**错误示例：**
```go
func waitForOutput(ch chan Event) tea.Cmd {
    return func() tea.Msg {
        ev, ok := <-ch
        if !ok {
            return nil  // 错误！返回 nil 会静默终止监听链
        }
        return outputMsg{ev}
    }
}
```

**正确示例：**
```go
type channelClosedMsg struct{}

func waitForOutput(ch chan Event) tea.Cmd {
    return func() tea.Msg {
        ev, ok := <-ch
        if !ok {
            return channelClosedMsg{}  // 返回专用消息类型
        }
        return outputMsg{ev}
    }
}
```

### 4.2 持续监听 Channel

需要从 channel 持续接收消息时，必须在每次处理后重新订阅：

```go
func (m *Model) handleOutput(ev Event) (tea.Model, tea.Cmd) {
    // 处理事件...
    return m, waitForOutput(m.outputCh)  // 重新订阅
}
```

### 4.3 使用 tea.Batch 并行执行多个 Cmd

```go
return m, tea.Batch(
    waitForOutput(m.outputCh),
    spinnerTick(),
    fetchAPI(url),
)
```

## 5. Init() 的正确实现

### 5.1 在 Init() 中启动必要的监听器

**错误示例：**
```go
func (m *Model) Init() tea.Cmd {
    return nil  // 错误！可能错过初始化期间的事件
}
```

**正确示例：**
```go
func (m *Model) Init() tea.Cmd {
    return tea.Batch(
        waitForOutput(m.outputCh),  // 立即开始监听
        spinnerTick(),              // 启动 spinner
    )
}
```

## 6. 消息类型设计

### 6.1 使用专用消息类型

```go
// 定义专用消息类型
type channelClosedMsg struct{}
type inputSentMsg struct{}
type loadingTickMsg struct{}
type errorMsg struct{ err error }

// 在 Update 中处理
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case channelClosedMsg:
        m.statusText = "Agent disconnected"
        return m, nil
    case errorMsg:
        m.statusText = "Error: " + msg.err.Error()
        return m, nil
    }
}
```

### 6.2 避免使用 nil 消息

返回 `nil` 作为 `tea.Msg` 是 no-op，会静默丢弃消息。始终返回有意义的消息类型。

## 7. 组件设计模式

### 7.1 组件应遵循 Bubble Tea 模式

可复用组件应该实现自己的 `Update` 和 `View` 方法：

```go
type TodoList struct {
    items    []TodoItem
    selected int
}

func (t *TodoList) Update(msg tea.Msg) (tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "up", "k":
            if t.selected > 0 {
                t.selected--
            }
        case "down", "j":
            if t.selected < len(t.items)-1 {
                t.selected++
            }
        }
    }
    return nil
}

func (t *TodoList) View() string {
    var b strings.Builder
    for i, item := range t.items {
        prefix := "  "
        if i == t.selected {
            prefix = "▸ "
        }
        b.WriteString(prefix + item.Title + "\n")
    }
    return b.String()
}
```

### 7.2 直接修改方法只能在 Update() 中调用

如果组件暴露了直接修改状态的方法（如 `AddItem()`、`Up()`、`Down()`），这些方法**只能**在父 Model 的 `Update()` 中调用：

```go
// 正确：在 Update 中调用
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "tab" {
            m.todoList.Up()  // 正确
        }
    }
    return m, nil
}

// 错误：在 View 中调用
func (m *Model) View() string {
    m.todoList.Up()  // 错误！View 中修改状态
    return m.todoList.View()
}
```

## 8. 鼠标事件处理

### 8.1 正确处理鼠标事件类型

Bubble Tea v2 中，鼠标事件使用 `tea.MouseMsg` 类型，通过 `msg.Type` 判断具体类型：

```go
func (m *Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
    switch msg.Type {
    case tea.MouseWheelUp:
        m.chatViewport.ScrollUp(3)
        return m, nil
    case tea.MouseWheelDown:
        m.chatViewport.ScrollDown(3)
        return m, nil
    }
    return m, nil
}
```

**错误示例：**
```go
func (m *Model) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.MouseWheelMsg:  // 错误！类型不匹配，永远不会匹配
        // ...
    }
    return m, nil
}
```

## 9. 性能优化

### 9.1 避免重复渲染

**错误示例：**
```go
// 在 Update 中渲染一次
m.chatViewport.SetContent(m.renderMessages())
// 在 View 中又渲染一次
func (m *Model) View() string {
    return m.chatViewport.View()  // viewport.View() 会再次渲染
}
```

**正确示例：**
```go
type Model struct {
    renderedContent string
    // ...
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // 更新时缓存渲染结果
    m.renderedContent = m.renderMessages()
    m.chatViewport.SetContent(m.renderedContent)
    return m, nil
}
```

### 9.2 使用 strings.Builder 拼接字符串

```go
func (m *Model) renderMessages() string {
    var b strings.Builder
    for _, msg := range m.messages {
        b.WriteString(msg.Render(m.chatViewport.Width))
    }
    return b.String()
}
```

## 10. 生命周期管理

### 10.1 正确清理资源

```go
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "ctrl+c" {
            m.cleanup()
            return m, tea.Quit
        }
    }
    return m, nil
}

func (m *Model) cleanup() {
    if m.cmd != nil && m.cmd.Process != nil {
        _ = m.cmd.Process.Kill()
    }
    m.cancel()
}
```

### 10.2 处理子进程退出

当监听的 channel 关闭时，应该更新 UI 状态：

```go
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case channelClosedMsg:
        m.promptRunning = false
        m.loading = false
        m.statusText = "Agent disconnected"
        return m, nil
    }
    return m, nil
}
```

## 11. 常见陷阱

### 11.1 忘记重新订阅 channel

```go
// 错误：处理完事件后没有重新订阅
func (m *Model) handleOutput(ev Event) (tea.Model, tea.Cmd) {
    // 处理事件...
    return m, nil  // 错误！监听链断了
}

// 正确：重新订阅
func (m *Model) handleOutput(ev Event) (tea.Model, tea.Cmd) {
    // 处理事件...
    return m, waitForOutput(m.outputCh)  // 正确
}
```

### 11.2 在 Update 中直接修改嵌套组件

```go
// 错误
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    m.viewport.YOffset = 0  // 直接修改
    return m, nil
}

// 正确
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    m.viewport, _ = m.viewport.Update(msg)
    return m, nil
}
```

### 11.3 忽略 tea.Batch 的返回值

```go
// 错误：Cmd 被丢弃
tea.Batch(cmd1, cmd2)
return m, nil

// 正确：返回 Cmd
return m, tea.Batch(cmd1, cmd2)
```

## 12. 调试技巧

### 12.1 使用 slog 记录消息

```go
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    m.logger.Debug("received message", "type", fmt.Sprintf("%T", msg))
    // ...
}
```

### 12.2 使用 --debug 标志

在运行时添加 `--debug` 标志可以输出详细日志到文件，帮助排查问题。

## 13. 总结

| 规则 | 说明 |
|------|------|
| View 是纯函数 | 只读状态，不修改，不执行副作用 |
| Update 不阻塞 | 所有阻塞操作包装在 Cmd 中 |
| Channel 关闭要处理 | 返回专用消息类型，不返回 nil |
| Init 启动监听器 | 不要遗漏初始化期间的事件 |
| 消息类型专用化 | 使用专用 struct 类型，不用 interface{} |
| 组件遵循模式 | 实现 Update/View，或文档说明调用限制 |
| 鼠标事件用 Type | 不用 type switch，用 msg.Type 判断 |
| 缓存渲染结果 | 避免在 Update 和 View 中重复渲染 |
| 清理资源 | 处理 ctrl+c 和 channel 关闭 |
| 重新订阅 | 处理完 channel 事件后重新订阅 |
