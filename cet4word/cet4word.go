package main

import (
	"fmt"
	"github.com/xkenmon/go-tools/cet4word/excel"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().Unix())
	item, err := excel.Reader("/home/bigmeng/Desktop/cet4word.xlsx", rand.Int()%4449)
	if err != nil {
		println(err)
		return
	}
	fmtStr :=
		`[35m==========================================================[0m
[31mEnglish:[0m  [33m%s[0m
[31mphonetic:[0m  [34m[%s][0m
[31mChinese:[0m  [36m%s[0m
[35m==========================================================[0m
`

	fmt.Printf(fmtStr, item[0], item[1], item[2])
}
