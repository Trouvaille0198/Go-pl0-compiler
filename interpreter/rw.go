package interpreter

import (
	"fmt"
	"log"
)

// read 读取一个整数
func read() (i int) {
	log.Printf("请输入一个无符号整数：")
	_, err := fmt.Scanln(&i)
	for err != nil {
		log.Printf("输入的不是无符号整数，请重新输入：")
		_, err = fmt.Scanln(&i)
	}
	return
}

// write 打印一个整数
func write(i int) {
	log.Printf("Write: %d", i)
}
