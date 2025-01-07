package main

import (
	"flag"
	"lv-blog/internal/global"
)

func main() {
	configPath := flag.String("c", "../config.yml", "配置文件路径")
	flag.Parse()

	conf := global.ReadConfig(*configPath)
}
