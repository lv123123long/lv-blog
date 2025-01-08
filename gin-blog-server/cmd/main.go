package main

import (
	"flag"
	"lv-blog/internal/global"
)

func main() {
	configPath := flag.String("c", "../config.yml", "配置文件路径")
	flag.Parse()

	// 根据命令行读取配置文件，其他变量的初始化依赖于配置文件对象
	conf := global.ReadConfig(*configPath)

}
