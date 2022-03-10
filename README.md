# Go-pl0-compiler
Golang 实现 pl0 编译器

## 目录架构

```
GO-PL0-COMPILER
│   .gitignore
│   go.mod
│   LICENSE
│   README.md
│
├───.idea // Goland 配置文件
│
├───assets
│       a.txt // pl0代码样例
│
├───fp
│       fp.go // 文件读取器
│
├───lexer
│       lexer.go // 词法分析器
│       lexer_test.go // 样例测试 包含exp1和exp2
│       state.go // DFA自动机状态声明
│
├───token
│       token.go // 符号枚举声明
│
└───utils
        charparse.go // 字符处理相关函数
```

## 词法分析器

`Lexer` 结构体封装了文件读指针、当前扫描行号（便于错误处理时的定位）和符号表，其方法 `GetSym` 扫描代码文件，并采用 DFA 方法获取符号信息：

```go
// Lexer 词法分析器
type Lexer struct {
	file    *fp.File
	line    int      // 当前所在行号
	symbols []Symbol // 符号数组
}
```

`Symbol` 结构体维护了一个符号的分类信息以及字面量信息

```go
// Symbol 符号
type Symbol struct {
	Id    token.Token // 符号枚举编号
	Value []rune      // 用户自定义的标识符值(若有)
	Num   int         // 用户自定义的数值(若有)
}
```

### exp1 识别标识符

**输出源程序中所有标识符出现的次数**

键入：

```shell
cd lexer
go test -v -run TestExp1
```

<img src="https://markdown-1303167219.cos.ap-shanghai.myqcloud.com/image-20220310140915531.png" alt="image-20220310140915531" style="zoom:67%;" />

### exp2 词法分析

**输出程序中各个单词符号（关键字、专用符号以及其他标记）**

键入：

```shell
cd lexer
go test -v -run TestExp2
```

<img src="https://markdown-1303167219.cos.ap-shanghai.myqcloud.com/image-20220310141031932.png" alt="image-20220310141031932" style="zoom:67%;" />

<img src="https://markdown-1303167219.cos.ap-shanghai.myqcloud.com/image-20220310141305925.png" alt="image-20220310141305925" style="zoom:67%;" />