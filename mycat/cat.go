package main

import (
	"flag"
	"fmt"
	"github.com/xkenmon/tools/mycat/download"
)

var (
	key = "AIzaSyBLv-O5At7qzBBkpA5PanCoFjf5XrvOvPA"
)

func main() {

	path := flag.String("path", ".", "图片存储路径")

	flag.Parse()

	info := download.GetPicInfo("")
	count := 0

	for idx, pic := range info.List {
		fmt.Printf("[%d] start download:%s\n", idx, pic.ImageID)
		filename, err := download.DownloadPic(pic.ImageID, *path)
		if err != nil {
			fmt.Println(err)
			continue
		}
		count++
		fmt.Println("download complete! " + *filename)
	}
	fmt.Printf("Download complete!!\ttotal download %d\n", count)
	fmt.Println("Enjoy your cat.")
}
