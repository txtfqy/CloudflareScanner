package main

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"sort"
	"sync"
	"time"
	"flag"
)

func handleUserInput() {
	var downloadSecond int64
	flag.IntVar(&pingRoutine, "thread", 400, "请输入扫描协程数(数字越大越快,不能超过1000)")
	flag.IntVar(&pingTime, "tcppings", 10, "tcping次数")
	flag.IntVar(&downloadTestCount, "nodes", 10, "要测试的下载节点个数")
	flag.Int64Var(&downloadSecond, "seconds", 10, "下载测试时间(单位为秒)")
	//解析获取参数
	flag.Parse()
	if pingRoutine <= 0 {
		pingRoutine = 400
	}
	if pingRoutine > 1000 {
		pingRoutine =1000
	}
	if pingTime <= 0 {
		pingTime = 10
	}
	if downloadTestCount <= 0 {
		downloadTestCount = 10
	}
	if downloadSecond <= 0 {
		downloadSecond = 10
	}
	downloadTestTime = time.Duration(downloadSecond) * time.Second
}

func main() {
	initipEndWith()
	handleUserInput()
	ips := loadFirstIPOfRangeFromFile()
	pingCount := len(ips) * pingTime
	bar := pb.StartNew(pingCount)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var data = make([]CloudflareIPData, 0)

	fmt.Println("开始tcping")

	control := make(chan bool, pingRoutine)
	for _, ip := range ips {
		wg.Add(1)
		control <- false
		handleProgress := handleProgressGenerator(bar)
		go tcpingGoroutine(&wg, &mu, ip, pingTime, &data, control, handleProgress)
	}
	wg.Wait()
	bar.Finish()
	bar = pb.StartNew(downloadTestCount)
	sort.Sort(CloudflareIPDataSet(data))
	fmt.Println("开始下载测速")
	for i := 0; i < downloadTestCount; i++ {
		_, speed := DownloadSpeedHandler(data[i].ip)
		data[i].downloadSpeed = speed
		bar.Add(1)
	}
	bar.Finish()
	ExportCsv("./result.csv", data)
}
