Glamour
=======

A casual introduction. 你好世界!

## Let’s talk about artichokes

The _artichoke_ is mentioned as a garden plant in the 8th century BC by Homer
**and** Hesiod. The naturally occurring variant of the artichoke, the cardoon,
which is native to the Mediterranean area, also has records of use as a food
among the ancient Greeks and Romans. Pliny the Elder mentioned growing of
_carduus_ in Carthage and Cordoba.

> He holds him with a skinny hand,
> ‘There was a ship,’ quoth he.
> ‘Hold off! unhand me, grey-beard loon!’
> An artichoke, dropt he.

--Samuel Taylor Coleridge, [The Rime of the Ancient Mariner][rime]

[rime]: https://poetryfoundation.org/poems/43997/

## Other foods worth mentioning

1. Carrots
1. Celery
1. Tacos
    * Soft
    * Hard
1. Cucumber

## Things to eat today

* [x] Carrots
* [x] Ramen
* [ ] Currywurst

### Power levels of the aforementioned foods

| Name       | Power | Comment          |
| ---        | ---   | ---              |
| Carrots    | 9001  | It’s over 9000?! |
| Ramen      | 9002  | Also over 9000?! |
| Currywurst | 10000 | What?!           |

## Currying Artichokes

Here’s a bit of code in [Haskell](https://haskell.org), because we are fancy.
Remember that to compile Haskell you’ll need `ghc`.

```haskell
module Main where

import Data.List (intercalate)

hello :: String -> String
hello s = "Hello, " <> s <> "."

main :: IO ()
main = putStrLn
     $ intercalate "\n"
     $ hello <$> [ "artichoke", "alcachofa" ]
```

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "runtime"
)

func main() {
    // 获取当前源文件所在目录
    _, filename, _, ok := runtime.Caller(0)
    if !ok {
        fmt.Println("无法获取当前文件路径")
        return
    }
    dir := filepath.Dir(filename)

    // 构造同目录下的目标文件路径
    filePath := filepath.Join(dir, "example.txt")

    // 读取文件内容
    data, err := os.ReadFile(filePath)
    if err != nil {
        fmt.Printf("读取文件失败: %v\n", err)
        return
    }
    fmt.Println(string(data))
}
```

***

_Alcachofa_, if you were wondering, is artichoke in Spanish.