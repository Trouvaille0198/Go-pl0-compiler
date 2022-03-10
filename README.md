# Go-pl0-compiler
Golang 实现 pl0 编译器

## 目录架构

```json
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