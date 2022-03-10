package fp

import (
	"bufio"
	"log"
	"os"
)

type File struct {
	file    *os.File
	scanner *bufio.Scanner
}

// NewFile 创建File对象
func NewFile(path string) *File {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanRunes) // 使用utf-8编码 故支持写入中文
	return &File{file: file, scanner: scanner}
}

// Read 读取一个字符
func (f *File) Read() (r rune, isEnd bool) {
	if f.scanner.Scan() {
		r = []rune(f.scanner.Text())[0] // 转换成rune
		isEnd = false
	} else {
		isEnd = true // 文件结束
	}
	return
}
