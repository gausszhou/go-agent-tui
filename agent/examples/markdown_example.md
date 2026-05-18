# Markdown 示例文档

## 1. 文本格式

这是**粗体文本**，这是*斜体文本*，这是~~删除线文本~~。

也可以使用 `_斜体_` 和 `__粗体__` 语法。

## 2. 代码

### 行内代码

在 Go 中使用 `fmt.Println("Hello, World!")` 输出内容。

### 代码块

```go
package main

import "fmt"

func main() {
    // 这是一个注释
    message := "Hello, World!"
    fmt.Println(message)

    // 循环示例
    for i := 0; i < 10; i++ {
        fmt.Printf("Count: %d\n", i)
    }
}
```

### 其他语言

```javascript
const greeting = "Hello, World!";
console.log(greeting);

function add(a, b) {
    return a + b;
}
```

```python
def greet(name):
    return f"Hello, {name}!"

result = greet("World")
print(result)
```

```json
{
    "name": "bubblecode",
    "version": "1.0.0",
    "description": "A TUI application for AI agents",
    "dependencies": {
        "bubbletea": "^4.0.0",
        "lipgloss": "^2.0.0"
    }
}
```

## 3. 列表

### 无序列表

- 第一项
- 第二项
  - 子项 A
  - 子项 B
- 第三项

### 有序列表

1. 步骤一
2. 步骤二
   2.1 子步骤
   2.2 子步骤
3. 步骤三

## 4. 表格

| 左对齐 | 居中对齐 | 右对齐 |
|:-------|:--------:|-------:|
| 内容1  |  内容2   |   内容3|
| 内容4  |  内容5   |   内容6|

## 5. 引用

> 这是一段引用文本。
> 可以跨越多行。

> ## 引用中也可以使用标题
> - 列表
> - 也可以使用

## 6. 链接和图片

[访问 GitHub](https://github.com)

![替代文本](https://example.com/image.png)

## 7. 分割线

---

上面是分割线

***

下面也是分割线

## 8. 标题层级

### 三级标题
#### 四级标题
##### 五级标题
###### 六级标题

## 9. 任务列表

- [x] 已完成的任务
- [ ] 未完成的任务
- [x] 另一个已完成

## 10. 脚注

这里有一个脚注[^1]。

[^1]: 这是脚注的内容。

## 11. emoji 表情

😀 😃 😄 😁 😆 😅 😂 🙂 🙃 😉 😊 😏 😒 😓 😔 😖 😘 😗 😙 😚 😜 😝 😛 😐 😑 😶 😏 😒 😓 😔 😖 😘 😗 😙 😚 😜 😝 😛 😗

## 12. 数学公式

行内公式：$E = mc^2$

块级公式：

$$
\frac{d}{dx} \left( \int_{0}^{x} f(t) dt \right) = f(x)
$$

## 13. 长文本测试

这是一段非常长的文本，用于测试界面在处理大量内容时的渲染效果。Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.

Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.

继续添加更多内容以确保文本足够长，能够充分测试渲染器的性能。ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()_+-=[]{}|;':",./<>?`~