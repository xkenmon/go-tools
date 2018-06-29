package main

import (
	"flag"
	"log"
)

func main() {

	path := flag.String("path", ".", "指定存储路径,默认为当前工作路径")
	count := flag.Int("count", 1, "指定下载总页数(每页15)")
	keyword := flag.String("search", "", "指定搜索关键字,为空从首页下载")
	resolution := flag.String("resolution", "1920x1080", "指定分辨率")
	poolSize := flag.Int("n",5,"线程池大小")

	flag.Parse()
	err := download(*path, *count, *keyword, *resolution,*poolSize)
	if err != nil {
		log.Fatal(err)
	}
}
