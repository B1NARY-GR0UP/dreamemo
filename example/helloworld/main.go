package main

import (
	"github.com/B1NARY-GR0UP/dreamemo/dream"
)

func main() {
	// 1. 使用 Default 或其他函数配置好所有的选项
	// 2. 使用 go run . 加命令行选项的模式配置节点地址
	// 3. 通过 HTTP 获取值
	dream.Default()
}
