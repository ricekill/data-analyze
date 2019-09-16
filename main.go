package main

import (
	"crawler-project/common"
	"crawler-project/module"
	"fmt"
	"time"
)

func main()  {
	//Load Config
	common.CheckErr(common.LoadConfig())
	common.CheckErr(common.OpenDb())

	fmt.Println("开始:", time.Now().Unix())
	//爬虫开始
	module.Start()
	fmt.Println("结束:", time.Now().Unix())
}
