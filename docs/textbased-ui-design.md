# 🎨 终端 TUI 设计规范 v1.0

## 一、核心数学基础

### 1.1 基本单位

```go
// 物理像素比
const (
    CharWidth  = 1  // 英文字符宽度（单位）
    CharHeight = 2  // 字符高度（单位）= 2 × 宽度
    CJKWidth   = 2  // 中文字符宽度
    CJKHeight  = 2  // 中文字符高度（实际等于行高）
)

// 视觉比例公式
视觉宽高比 = 列数 / (行数 × 2)
```

### 1.2 黄金尺寸表

| 尺寸 | 视觉比例 | 适用场景 | 容量 |
|------|---------|---------|------|
| 40×12 | 1.667 | 弹窗/提示 | 480 |
| 60×18 | 1.667 | 侧边栏 | 1,080 |
| **80×24** | 1.667 | **标准终端** | 1,920 |
| 100×30 | 1.667 | 宽屏终端 | 3,000 |
| **120×36** | 1.667 | **现代宽屏** | 4,320 |
| 140×42 | 1.667 | 超宽屏 | 5,880 |
| 160×48 | 1.667 | 4K 显示 | 7,680 |

### 1.3 最小可行尺寸

```go
const (
    MinCols = 40  // 最小列数（低于此值无法正常显示）
    MinRows = 12  // 最小行数（低于此值内容拥挤）
    MinWidth  = 40   // 40列
    MinHeight = 12   // 12行
)
```

---

## 二、布局规范

### 2.1 黄金分割布局（推荐）

基于 1.618 的视觉黄金比例：

```go
// 布局分配（以 120×36 为例）
const (
    // 水平方向（总列数 = 120）
    MainAreaWidth   = 74  // 黄金份额（61.7%）
    SidebarWidth    = 46  // 白银份额（38.3%）
    
    // 垂直方向（总行数 = 36）
    HeaderHeight    = 4   // 顶部栏（11.1%）
    ContentHeight   = 28  // 主内容区（77.8%）
    FooterHeight    = 4   // 底部栏（11.1%）
)
```

### 2.2 标准布局模板

```go
type Layout struct {
    // 预留边距
    MarginTop    int = 1
    MarginBottom int = 1
    MarginLeft   int = 1
    MarginRight  int = 1
    
    // 组件区域
    Header  Region  // 顶部：标题、标签页、状态
    Sidebar Region  // 侧边栏：导航、列表
    Content Region  // 主内容：核心显示区
    Footer  Region  // 底部：帮助、提示、进度
}

type Region struct {
    X, Y          int  // 起始坐标
    Width, Height int  // 尺寸
}
```

### 2.3 布局分配公式

```go
// 根据总尺寸自动计算黄金布局
func GoldenLayout(totalWidth, totalHeight int) (main, sidebar, header, content, footer int) {
    // 水平方向：黄金分割
    main = int(float64(totalWidth) / 1.618)      // 较大区域
    sidebar = totalWidth - main - 2              // 扣除边距
    
    // 垂直方向：三等分变体
    header = int(float64(totalHeight) * 0.15)    // 约 15%
    footer = int(float64(totalHeight) * 0.10)    // 约 10%
    content = totalHeight - header - footer - 2  // 剩余 75%
    
    return
}
```

---

## 三、间距规范

### 3.1 标准间距单位

```go
const (
    SpacingXS = 0  // 无间距（边框紧贴）
    SpacingS  = 1  // 小间距（1 列/行）
    SpacingM  = 2  // 中间距（推荐）
    SpacingL  = 3  // 大间距
    SpacingXL = 4  // 超大间距
)

// 推荐：使用 SpacingM = 2 作为标准单元
```

### 3.2 边距规范

```go
// 窗口级边距（最外层）
const WindowPadding = 1

// 组件间间距
const ComponentGap = 1

// 内容区内边距
const ContentPadding = 1
```

### 3.3 边框使用规范

```go
// Lip Gloss 边框风格推荐
边框类型         = lipgloss.NormalBorder()  // 日常使用
圆角边框         = lipgloss.RoundedBorder() // 弹窗/对话框
加粗边框         = lipgloss.ThickBorder()   // 强调区域
隐藏边框         = lipgloss.HiddenBorder()  // 沉浸模式
无边框 + 背景色   = lipgloss.Border{}        // 卡片效果
```

---

## 四、色彩规范

### 4.1 色彩系统

```go
// 基础色板（256 色终端安全）
const (
    // 前景色（文字）
    TextPrimary    = lipgloss.Color("15")  // 亮白
    TextSecondary  = lipgloss.Color("7")   // 浅灰
    TextDisabled   = lipgloss.Color("8")   // 暗灰
    TextAccent     = lipgloss.Color("39")  // 亮蓝
    
    // 背景色
    BgPrimary      = lipgloss.Color("0")   // 黑色
    BgSecondary    = lipgloss.Color("236") // 深灰
    BgAccent       = lipgloss.Color("24")  // 深蓝
    
    // 状态色
    ColorSuccess   = lipgloss.Color("10")  // 绿色
    ColorWarning   = lipgloss.Color("11")  // 黄色
    ColorError     = lipgloss.Color("9")   // 红色
    ColorInfo      = lipgloss.Color("39")  // 蓝色
)
```

### 4.2 对比度要求

```go
// 最小对比度：4.5:1（WCAG AA 标准）
// 推荐对比度：7:1（WCAG AAA 标准）

// 安全组合
推荐组合 = []Pair{
    {TextPrimary, BgPrimary},     // 白字黑底
    {TextAccent, BgPrimary},      // 蓝字黑底
}
```

### 4.3 透明度使用

```go
// 仅在支持真彩色的终端使用透明度
if lipgloss.HasDarkBackground() {
    // 使用半透明效果
    style = style.Background(lipgloss.Color("#000000aa"))
}
```

---

## 五、字体与文本规范

### 5.1 字符宽度处理

```go
// 关键函数
func DisplayWidth(s string) int {
    return lipgloss.Width(s)  // 自动处理中英文
}

func Truncate(s string, maxWidth int) string {
    if lipgloss.Width(s) <= maxWidth {
        return s
    }
    return lipgloss.NewStyle().Width(maxWidth - 1).Render(s) + "…"
}

func PadRight(s string, width int) string {
    current := lipgloss.Width(s)
    if current >= width {
        return s
    }
    return s + strings.Repeat(" ", width-current)
}
```

### 5.2 对齐规范

```go
对齐方式：
- 标题文本：居中对齐
- 正文内容：左对齐
- 数字/状态：右对齐
- 混合场景：使用表格布局

// 左对齐（默认）
style.Align(lipgloss.Left)

// 右对齐（用于数字、时间戳）
style.Align(lipgloss.Right)

// 居中对齐（用于标题、弹窗）
style.Align(lipgloss.Center)
```

---

## 六、组件设计规范

### 6.1 Header（顶部栏）

```go
Header 规范：
- 高度：3-5 行
- 内容：标题、当前时间、通知标识
- 样式：加粗 + 背景色

示例：
┌──────────────────────────────────────┐
│  MyApp v1.0                    14:30 │
│  Dashboard > Settings                │
└──────────────────────────────────────┘
```

### 6.2 Sidebar（侧边栏）

```go
Sidebar 规范：
- 宽度：占总宽度的 25-38%
- 内容：导航菜单、列表、状态
- 样式：使用边框或背景色区分

宽度建议（基于总宽）：
总宽 80  → 侧边栏 20-30 列
总宽 100 → 侧边栏 25-38 列
总宽 120 → 侧边栏 30-46 列
```

### 6.3 Content（主内容区）

```go
Content 规范：
- 宽度：占总宽度的 62-75%
- 高度：占总高度的 70-80%
- 滚动：使用 viewport 组件
- 边距：内部保留 1 列内边距

// 内容区模板
contentStyle := lipgloss.NewStyle().
    Width(contentWidth).
    Height(contentHeight).
    Padding(0, 1).  // 左右内边距 1
    Border(lipgloss.NormalBorder())
```

### 6.4 Footer（底部栏）

```go
Footer 规范：
- 高度：2-4 行
- 内容：快捷键提示、进度状态、错误信息
- 样式：使用较暗的颜色或背景

示例：
┌──────────────────────────────────────┐
│  ↑↓ 移动  Enter 选择  q 退出  ? 帮助  │
└──────────────────────────────────────┘
```

### 6.5 Dialog/Modal（弹窗）

```go
弹窗规范：
- 宽度：40-60 列
- 高度：10-20 行
- 位置：屏幕中央
- 样式：圆角边框 + 背景色

// 弹窗居中计算
dialogX := (totalWidth - dialogWidth) / 2
dialogY := (totalHeight - dialogHeight) / 2
```

---

## 七、交互规范

### 7.1 键盘映射

```go
// 通用快捷键
通用键位：
- ? / F1 : 帮助
- q / Ctrl+C : 退出
- ↑/↓/←/→ : 导航
- Tab : 切换焦点
- Enter : 确认/选择
- Esc : 取消/返回

// 编辑键位
- Ctrl+S : 保存
- / : 搜索
- n : 新建
- d : 删除
- e : 编辑
```

### 7.2 焦点管理

```go
// 焦点指示器
焦点样式 = lipgloss.NewStyle().Background(ColorAccent).Foreground(ColorWhite)
非焦点样式 = lipgloss.NewStyle().Background(BgSecondary).Foreground(TextSecondary)

// 焦点环
func FocusIndicator(focused bool, text string) string {
    if focused {
        return "▶ " + text
    }
    return "  " + text
}
```

### 7.3 反馈规范

```go
// 操作反馈
反馈类型：
- 成功：绿色提示，1.5 秒后消失
- 错误：红色提示，3 秒后消失
- 警告：黄色提示，需要用户确认
- 信息：蓝色提示，2 秒后消失

// 进度指示器
func Spinner(frame int) string {
    frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
    return frames[frame%len(frames)]
}
```

---

## 八、性能规范

### 8.1 渲染优化

```go
// 使用 Viewport 减少重绘
viewport := viewport.New(width, height)
viewport.SetContent(content)

// 批量更新
tea.Batch(
    updateViewportCmd,
    updateStatusCmd,
    updateHelpCmd,
)
```

### 8.2 内容分页

```go
// 大文件/长列表处理
const PageSize = 100

func Paginate(items []string, page int) []string {
    start := page * PageSize
    end := start + PageSize
    if start >= len(items) {
        return []string{}
    }
    if end > len(items) {
        end = len(items)
    }
    return items[start:end]
}
```

---

## 九、响应式规范

### 9.1 断点定义

```go
const (
    BreakpointSmall  = 60   // < 60列：紧凑模式
    BreakpointMedium = 100  // 60-100列：标准模式
    BreakpointLarge  = 140  // > 100列：宽屏模式
)

func ResponsiveLayout(width int) LayoutMode {
    switch {
    case width < BreakpointSmall:
        return CompactMode    // 隐藏侧边栏，堆叠布局
    case width < BreakpointMedium:
        return StandardMode   // 显示侧边栏，标准布局
    default:
        return WideMode       // 显示额外信息，宽松布局
    }
}
```

### 9.2 窗口变化处理

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        // 拒绝过小的窗口
        if msg.Width < MinCols || msg.Height < MinRows {
            m.ready = false
            return m, nil
        }
        
        // 更新所有组件尺寸
        m.width = msg.Width
        m.height = msg.Height
        m.updateLayout()
        
        return m, viewport.Sync(m.viewport)
    }
    return m, nil
}
```

---

## 十、检查清单

### 创建新 TUI 前的检查项

- [ ] 设置最小尺寸限制（40×12）
- [ ] 选择黄金比例尺寸（80×24 或 120×36）
- [ ] 实现响应式布局
- [ ] 正确处理 WindowSizeMsg
- [ ] 先换行后加样式的顺序
- [ ] 支持 CJK 字符显示
- [ ] 提供颜色 Fallback 方案
- [ ] 实现基本的键盘导航
- [ ] 添加帮助信息
- [ ] 处理退出快捷键

### 发布前检查

- [ ] 在不同终端测试（iTerm2, Terminal, Alacritty）
- [ ] 测试不同窗口大小
- [ ] 验证 256 色和真彩色支持
- [ ] 测试无颜色模式
- [ ] 验证屏幕阅读器兼容性
- [ ] 测试极长内容和边界情况

---

## 总结

这套规范的核心思想：

1. **数学基础**：2:1 像素比 + 黄金比例 1.618
2. **经典尺寸**：80×24 或 120×36
3. **布局原则**：黄金分割（61.8% / 38.2%）
4. **间距单位**：2 作为标准间距单元
5. **文本处理**：先换行，后样式
6. **响应式**：适配不同窗口尺寸
7. **可访问性**：足够对比度，键盘导航

遵循这套规范，你可以创建出**视觉舒适、功能完整、跨平台一致**的终端 TUI 应用。
