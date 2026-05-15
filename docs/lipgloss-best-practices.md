# Lip Gloss v2 最佳实践指南

本文档总结了使用 Lip Gloss v2 (`charm.land/lipgloss/v2`) 构建终端 UI 时的最佳实践。

## 1. 样式方法链顺序

### 1.1 推荐顺序

Lip Gloss v2 的样式方法链**顺序不影响最终结果**，因为每个方法都返回一个新的 `Style` 副本，所有属性在 `Render()` 时统一计算。但为了代码可读性和团队一致性，推荐遵循以下顺序：

```go
var style = lipgloss.NewStyle().
    // 1. 文本样式
    Bold(true).
    Italic(true).
    Foreground(color).
    Background(color).
    
    // 2. 尺寸
    Width(24).
    Height(10).
    MaxWidth(40).
    
    // 3. 内边距
    PaddingTop(2).
    PaddingLeft(4).
    
    // 4. 边框
    Border(lipgloss.NormalBorder()).
    BorderForeground(color).
    
    // 5. 对齐
    Align(lipgloss.Center).
    
    // 6. 其他
    Underline(true)
```

### 1.2 常见误区

**误区：必须先设置 Width 再设置 Background/Foreground**

这个说法**不正确**。Lip Gloss v2 中，方法调用顺序不影响渲染结果。以下两种写法完全等价：

```go
// 写法 A
lipgloss.NewStyle().Width(24).Foreground(color).Background(color)

// 写法 B
lipgloss.NewStyle().Foreground(color).Background(color).Width(24)
```

**但有一个例外**：当使用 `SetString()` 时，建议在设置尺寸之后再调用，因为 `SetString()` 只是设置底层字符串，不影响样式计算：

```go
lipgloss.NewStyle().
    Width(24).
    Height(10).
    Foreground(color).
    SetString("Hello")  // 最后设置字符串
```

## 2. 颜色使用

### 2.1 十六进制颜色

```go
lipgloss.Color("#201d1d")  // 推荐：使用十六进制
```

### 2.2 ANSI 颜色

```go
lipgloss.Color("63")       // ANSI 256 色
lipgloss.Magenta           // 标准 16 色
```

### 2.3 自适应颜色（亮/暗背景）

Lip Gloss v2 不再自动检测终端背景色，需要显式查询：

```go
hasDark, _ := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
lightDark := lipgloss.LightDark(hasDark)

fg := lightDark(lipgloss.Color("#333333"), lipgloss.Color("#f1f1f1"))
s := lipgloss.NewStyle().Foreground(fg)
```

### 2.4 颜色降级

Lip Gloss 会自动将颜色降级为终端支持的最佳配置。当输出不是 TTY 时，会自动移除 ANSI 颜色码。

## 3. 尺寸与布局

### 3.1 Width 和 Height

`Width` 和 `Height` 设置的是**内容块的最小尺寸**。如果内容超出，会自动换行或扩展：

```go
lipgloss.NewStyle().
    Width(40).    // 最小宽度 40
    Height(10).   // 最小高度 10
    Render(text)
```

### 3.2 MaxWidth 和 MaxHeight

`MaxWidth` 和 `MaxHeight` 用于**强制限制最大尺寸**，常用于渲染时：

```go
// 渲染时限制最大宽度
style.MaxWidth(16).Render(userInput)
```

### 3.3 测量尺寸

```go
block := style.Render(text)
width := lipgloss.Width(block)
height := lipgloss.Height(block)

// 或简写
w, h := lipgloss.Size(block)
```

## 4. 内边距与边框

### 4.1 Padding 简写

```go
Padding(2)           // 所有方向
Padding(2, 4)        // 垂直 2，水平 4
Padding(1, 2, 3)     // 上 1，水平 2，下 3
Padding(1, 2, 3, 4)  // 上、右、下、左
```

### 4.2 边框

```go
lipgloss.NewStyle().
    Border(lipgloss.NormalBorder()).
    BorderForeground(lipgloss.Color("#007aff")).
    BorderTop(true).      // 只显示上边框
    BorderBottom(true)    // 只显示下边框
```

### 4.3 边框尺寸计算

Lip Gloss v2 中，`Width` 和 `Height` 包含边框、内边距和内容。如果设置 `Width(40)` 并有左右各 1 的边框和 2 的内边距，实际内容宽度为 `40 - 2 - 4 = 34`。

## 5. 文本对齐

### 5.1 水平对齐

```go
Align(lipgloss.Left)
Align(lipgloss.Center)
Align(lipgloss.Right)
```

### 5.2 垂直对齐

```go
AlignVertical(lipgloss.Top)
AlignVertical(lipgloss.Center)
AlignVertical(lipgloss.Bottom)
```

### 5.3 对齐需要设置 Width/Height

对齐只在设置了 `Width` 或 `Height` 时生效，因为需要参考尺寸：

```go
// 正确：设置了 Width，对齐生效
lipgloss.NewStyle().Width(40).Align(lipgloss.Center).Render("Hi")

// 错误：未设置 Width，对齐不生效
lipgloss.NewStyle().Align(lipgloss.Center).Render("Hi")
```

## 6. 样式复用与继承

### 6.1 基础样式函数

推荐在 `theme/` 包中定义基础样式函数：

```go
// theme/theme.go
func BaseStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Background(ThemeBg).
        Foreground(ThemeText)
}

func StatusBar() lipgloss.Style {
    return lipgloss.NewStyle().
        Background(ThemeBgDark).
        Foreground(ThemeMuted)
}
```

### 6.2 Inherit 继承

使用 `Inherit` 可以合并两个样式：

```go
base := lipgloss.NewStyle().Background(bg).Foreground(text)
focused := base.Inherit(lipgloss.NewStyle().Border(lipgloss.NormalBorder()))
```

### 6.3 预定义样式变量

对于不改变的样式，使用包级变量：

```go
var (
    StyleUser    = lipgloss.NewStyle().Foreground(ThemeUser).Bold(true)
    StyleAgent   = lipgloss.NewStyle().Foreground(ThemeAgent).Bold(true)
    StyleThought = lipgloss.NewStyle().Foreground(ThemeThought).Italic(true)
)
```

## 7. 渲染性能

### 7.1 避免在 View() 中创建样式

样式创建有一定开销，建议缓存：

```go
// 错误：每次 View() 都创建新样式
func (m *Model) View() string {
    style := lipgloss.NewStyle().Foreground(color).Bold(true)
    return style.Render(text)
}

// 正确：使用预定义样式
var headerStyle = lipgloss.NewStyle().Foreground(color).Bold(true)

func (m *Model) View() string {
    return headerStyle.Render(text)
}
```

### 7.2 动态样式使用函数

当样式需要动态参数时，使用函数：

```go
func ChatBg(w int) lipgloss.Style {
    return lipgloss.NewStyle().Background(ThemeChatBg).Width(w)
}
```

## 8. 空白区域样式

### 8.1 WhitespaceStyle

Lip Gloss v2 使用 `WithWhitespaceStyle` 替代 v1 的 `WithWhitespaceForeground/Background`：

```go
lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, text,
    lipgloss.WithWhitespaceStyle(
        lipgloss.NewStyle().
            Foreground(lipgloss.Color("#333")).
            Background(lipgloss.Color("#000")),
    ),
)
```

## 9. 下划线样式

### 9.1 基本下划线

```go
lipgloss.NewStyle().Underline(true)
```

### 9.2 下划线样式

```go
lipgloss.NewStyle().
    Underline(true).
    UnderlineStyle(lipgloss.UnderlineCurly).
    UnderlineColor(lipgloss.Color("#007aff"))
```

### 9.3 下划线空格

默认下划线不应用于空白字符，使用 `UnderlineSpaces(true)` 可以改变此行为：

```go
lipgloss.NewStyle().Underline(true).UnderlineSpaces(true)
```

## 10. 复合布局

### 10.1 水平排列

```go
lipgloss.JoinHorizontal(lipgloss.Top, left, right)
```

### 10.2 垂直排列

```go
lipgloss.JoinVertical(lipgloss.Left, top, bottom)
```

### 10.3 列表布局

```go
lipgloss.JoinVertical(lipgloss.Left,
    header,
    content,
    footer,
)
```

## 11. 图层与合成（v2 新特性）

### 11.1 创建图层

```go
box := lipgloss.NewStyle().Width(10).Height(5).Border(lipgloss.NormalBorder())

a := lipgloss.NewLayer(box.Render("Layer A"))
b := lipgloss.NewLayer(box.Render("Layer B"))

canvas := lipgloss.NewCanvas(
    a.X(5).Y(10).Z(1),
    b.X(3).Y(7)
)

lipgloss.Println(canvas.Render())
```

## 12. 常见陷阱

### 12.1 忘记设置 Width 导致对齐不生效

```go
// 错误：对齐不生效
lipgloss.NewStyle().Align(lipgloss.Center).Render("Hi")

// 正确：设置 Width
lipgloss.NewStyle().Width(40).Align(lipgloss.Center).Render("Hi")
```

### 12.2 使用 fmt.Println 而非 lipgloss.Println

```go
// 错误：可能颜色显示不正确
fmt.Println(style.Render(text))

// 正确：使用 lipgloss.Println 进行颜色降级
lipgloss.Println(style.Render(text))
```

### 12.3 在 Bubble Tea View() 中修改样式

```go
// 错误：View() 中修改样式
func (m *Model) View() string {
    m.style = m.style.Width(m.width)  // 不应该在 View 中修改
    return m.style.Render(text)
}

// 正确：在 Update() 中更新样式
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    case tea.WindowSizeMsg:
        m.style = m.style.Width(msg.Width)
}
```

### 12.4 边框与尺寸计算错误

```go
// 设置 Width(40) 包含边框和内边距
style := lipgloss.NewStyle().
    Width(40).
    Border(lipgloss.NormalBorder()).  // 左右各占 1
    Padding(0, 2)                      // 左右各占 2

// 实际内容宽度 = 40 - 2(边框) - 4(内边距) = 34
```

## 13. 总结

| 规则 | 说明 |
|------|------|
| 方法链顺序 | 不影响结果，但推荐：文本→尺寸→内边距→边框→对齐 |
| Width/Background 顺序 | 无先后要求，可任意顺序调用 |
| 对齐需要 Width | 设置 Align 前必须先设置 Width 或 Height |
| 缓存样式 | 避免在 View() 中重复创建样式 |
| 使用 lipgloss.Println | 而非 fmt.Println，确保颜色降级正确 |
| 边框包含在尺寸内 | Width/Height 包含边框和内边距 |
| v2 自适应颜色 | 需要显式调用 HasDarkBackground |
| 空白区域样式 | 使用 WithWhitespaceStyle 而非独立的前景色/背景色 |
