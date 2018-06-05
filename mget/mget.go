package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type dataRange struct {
	start int64
	end   int64
}

type block struct {
	data   []byte
	dRange dataRange
	id     int
	ok     bool
}

type blocksInfo struct {
	downloaded  []bool
	downloading []bool
	unDownload  []bool
	mutx        sync.Mutex
}

type downloadInfo struct {
	all        int64
	downloaded int64
	bytePerS   int64
	mutx       sync.Mutex
}

type threadPool struct {
	count int
	mutx  sync.Mutex
}

var (
	u          string
	thread     int
	blocksize  int64
	flush      int
	cookiesStr string
	userAgent  string
	headerStr  string
)

func main() {
	threadPtr := flag.Int("n", 1, "下载线程数")
	fileNamePtr := flag.String("o", "unamed", "输出文件名")
	pathPtr := flag.String("path", "./", "存储文件夹路径")
	uPtr := flag.String("u", "", "下载地址")
	blockSizePtr := flag.Int("b", 8, "缓冲区大小(MB)")
	fl := flag.Int("s", 1, "速度刷新率(s)")
	cookiesPtr := flag.String("c", "", "设置头部")
	userAgentPtr := flag.String("ua", "", "userAgent")
	headerPtr := flag.String("h", "", "Header")

	flag.Parse()

	headerStr = *headerPtr
	u = *uPtr
	userAgent = *userAgentPtr
	blocksize = int64(*blockSizePtr) * 1024 * 1024
	path := *pathPtr
	fileName := *fileNamePtr
	cookiesStr = *cookiesPtr
	flush = *fl
	if fileName == "unamed" {
		fileName = getFileName(u)
	}
	thread = *threadPtr

	if u == "" {
		fmt.Println("使用-u参数指定下载地址\n--help获取完整帮助")
		return
	}
	if !strings.HasPrefix(u, "http") {
		u = "http://" + u
	}
	// fmt.Println("thread", *thread)
	// fmt.Println("fileName", *fileName)
	// fmt.Println("path", *path)
	// fmt.Println("u", *u)
	// fmt.Println("buf", *buf)

	_, rangeable := getLength(u)

	if !rangeable || thread == 1 {
		fmt.Println("Single thread Download...")
		DownloadSingle(path, fileName)
	} else {
		fmt.Println("Mutil thread Download...")
		MutilDownload(path, fileName)
	}
}

func MutilDownload(path, fName string) {
	length, _ := getLength(u)
	fmt.Println("Content-Length:", length)
	info := &downloadInfo{
		all: length,
	}

	count := length / blocksize //10MB

	blocks := make([]block, count)

	for i := range blocks {
		blocks[i].dRange = dataRange{int64(i) * blocksize, (int64(i) + 1) * blocksize}
		blocks[i].id = i
	}

	fmt.Println("block count: ", count)

	if length%blocksize != 0 {
		blocks = append(blocks, block{
			dRange: dataRange{count * blocksize, length},
			id:     int(count),
		})
		fmt.Println("block count++ ", "{", count*blocksize, "-", length, "}")
		count++
	}

	ch := make(chan *block, 8)
	defer close(ch)

	blockinfos := &blocksInfo{
		downloaded:  make([]bool, count),
		downloading: make([]bool, count),
		unDownload:  make([]bool, count),
	}

	pool := &threadPool{}

	for i := range blockinfos.unDownload {
		blockinfos.unDownload[i] = true
	}

	for i := 0; i < thread && i < len(blocks); i++ {
		pool.mutx.Lock()
		go goDownload(&blocks[i], ch)
		pool.count++
		blockinfos.Downloading(i)
		pool.mutx.Unlock()
	}
	go printSpeed(info, flush)

	f, _ := os.OpenFile(path+fName, os.O_CREATE|os.O_RDWR, 0666)
	defer f.Close()
	for b := range ch {
		if !(b.ok) {
			fmt.Println(b, "Download failed, retry...")
			go goDownload(b, ch)
			continue
		}
		//fmt.Println("Download succeed", b)
		//fmt.Printf("ratio : %s\\%s\n", sizeFormat(info.downloaded), sizeFormat(info.all))
		info.AddDownloaded(b.dRange.end - b.dRange.start)
		blockinfos.Downloaded(b.id)
		//fmt.Println("Write to ", b.dRange.start, ",len:", len(b.data))
		f.WriteAt(b.data, b.dRange.start)
		//释放内存
		b.data = nil
		undown := blockinfos.GetUndownload()
		if undown == -1 {
			if blockinfos.GetDownloading() == -1 {
				if info.DownloadSucceed() {
					fmt.Println("Download Ok!")
					return
				}
				fmt.Println("something Wrong...")
				return
			}
		} else {
			go goDownload(&blocks[undown], ch)
			blockinfos.Downloading(undown)
		}
	}
}

func (info *downloadInfo) AddDownloaded(l int64) {
	info.mutx.Lock()
	info.downloaded += l
	info.bytePerS += l
	defer info.mutx.Unlock()
}

func (info *downloadInfo) DownloadSucceed() bool {
	info.mutx.Lock()
	defer info.mutx.Unlock()
	return info.all == info.downloaded
}

func (info *blocksInfo) GetUndownload() int {
	info.mutx.Lock()
	defer info.mutx.Unlock()
	for i, v := range info.unDownload {
		if v {
			return i
		}
	}
	return -1
}

func (info *blocksInfo) GetDownloading() int {
	info.mutx.Lock()
	defer info.mutx.Unlock()
	for i, v := range info.downloading {
		if v {
			return i
		}
	}
	return -1
}

func (info *blocksInfo) Downloading(i int) {
	info.mutx.Lock()
	info.downloading[i] = true
	info.unDownload[i] = false
	defer info.mutx.Unlock()
}

func (info *blocksInfo) Downloaded(i int) {
	info.mutx.Lock()
	info.downloaded[i] = true
	info.downloading[i] = false
	defer info.mutx.Unlock()
}

func (blk *block) String() string {
	// return fmt.Sprintf("{id:%d\trange{%d - %d}\tisOK? %v\ndata{\n%v\n}}", blk.id, blk.dRange.start, blk.dRange.end, blk.ok, blk.data)
	return fmt.Sprintf("{id:%d\trange{%s - %s}\tisOK? %v}", blk.id, sizeFormat(blk.dRange.start), sizeFormat(blk.dRange.end), blk.ok)
}

func goDownload(blk *block, ch chan *block) {
	//fmt.Println("start Download block : ", blk)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		fmt.Println("解析", u, "时出现错误: ", err)
		blk.ok = false
		ch <- blk
		return
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", blk.dRange.start, blk.dRange.end))
	req.Header.Set("Connection", "close")

	req.ParseForm()

	for k, v := range req.Form {
		fmt.Printf("Param: %s=%s\n", k, v)
	}

	for _, c := range parseCookies(cookiesStr) {
		req.AddCookie(c)
	}

	if userAgent != "" {
		req.Header.Add("User-Agent", userAgent)
	}

	for k, v := range parseHeader(headerStr) {
		req.Header.Add(k, v)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求时出现错误", err)
		blk.ok = false
		ch <- blk
		return
	}
	// l, err := io.Copy(buf, resp.Body)
	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("copy err", err)
		blk.ok = false
		ch <- blk
		return
	}
	blk.data = data
	blk.ok = true
	ch <- blk
}

func DownloadSingle(path, fName string) (length int64) {

	req, err := http.NewRequest(http.MethodGet, path, nil)
	for _, cookie := range parseCookies(cookiesStr) {
		req.AddCookie(cookie)
	}

	if userAgent != "" {
		req.Header.Add("User-Agent", userAgent)
	}

	for k, v := range parseHeader(headerStr) {
		req.Header.Add(k, v)
	}
	req.ParseForm()
	for k, v := range req.Form {
		fmt.Printf("Param: %s=%s\n", k, v)
	}

	resp, err := http.Get(u)
	if "unamed" == fName {
		fName = getFileName(u)
		fmt.Println("Get file name :", fName)
	}
	length = resp.ContentLength
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	f, err := os.Create(path + fName)

	fmt.Println("Start download file :", f.Name())

	io.Copy(f, resp.Body)

	return length
}

func getLength(downloadURL string) (len int64, rangeable bool) {
	resp, err := http.Get(downloadURL)

	if err != nil {
		fmt.Printf("Can't open url : %s\n%v\n", downloadURL, err)
		os.Exit(-1)
	}

	acceptRange := resp.Header.Get("Accept-Ranges")
	fmt.Println("Accept-Ranges", acceptRange)
	if acceptRange == "" || acceptRange == "none" {
		rangeable = false
	} else {
		rangeable = true
	}

	defer resp.Body.Close()
	return resp.ContentLength, rangeable
}

func getFileName(u string) string {
	return u[strings.LastIndex(u, "/")+1:]
}

func printSpeed(info *downloadInfo, ratio int) {
	t := time.Tick(time.Duration(ratio) * time.Second)
	for _ = range t {
		fmt.Printf("speed %s/%d s\t", sizeFormat(info.bytePerS), ratio)
		fmt.Printf("ratio %s/%s\n", sizeFormat(info.downloaded), sizeFormat(info.all))
		info.bytePerS = 0
	}
}

func sizeFormat(s int64) string {
	sizes := []string{"B", "kB", "MB", "GB", "TB", "PB", "EB"}
	return humanateBytes(s, 1000, sizes)
}

func humanateBytes(s int64, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%d B", s)
	}
	e := math.Floor(math.Log(float64(s)) / math.Log(base))
	suffix := sizes[int(e)]
	val := math.Floor(float64(s)/math.Pow(base, e)*10+0.5) / 10
	f := "%.0f %s"
	if val < 10 {
		f = "%.1f %s"
	}

	return fmt.Sprintf(f, val, suffix)
}

func parseCookies(c string) []*http.Cookie {
	cookies := make([]*http.Cookie, 0)
	arr := strings.Split(c, ";")
	for _, s := range arr {
		s = strings.Trim(s, " \r\n\t")
		cookieMap := strings.Split(s, "=")
		if len(cookieMap) == 2 {
			cookie := &http.Cookie{Name: cookieMap[0], Value: cookieMap[1]}
			cookies = append(cookies, cookie)
			fmt.Printf("set cookie: %s=%s\n", cookieMap[0], cookieMap[1])
		}
	}
	return cookies
}

func parseHeader(h string) map[string]string {
	headers := make(map[string]string)
	arr := strings.Split(h, ";")
	for _, s := range arr {
		s = strings.Trim(s, " ")
		m := strings.Split(s, ": ")
		fmt.Println(m)
		if len(m) == 2 {
			headers[m[0]] = m[1]
			fmt.Printf("set Header: %s = %s\n", m[0], m[1])
		}
	}
	return headers
}
