package main

import (
	"strings"
	"net/url"
	"net/http"
	"time"
	"strconv"
	"io/ioutil"
	"github.com/PuerkitoBio/goquery"
	"log"
	"errors"
	"os"
	"path"
)

const searchUrl = "https://www.pexels.com/search/"
const indexUrl = "https://www.pexels.com"
const startFlag = `<article`
const endFlag = `<\/article>`

func download(dir string, count int, keyword string, resolution string) error {
	list, err := getUrl(count, keyword)
	if err != nil {
		return err
	}

	if stat, err := os.Stat(dir); err != nil {
		return err
	} else {
		if !stat.IsDir() {
			os.MkdirAll(dir, os.ModeDir)
		}
	}

	ch := make(chan error, 5)
	for i, link := range list {
		u, err := url.Parse(link)
		if err != nil {
			return err
		}
		fileName := u.Query().Get("dl")
		if fileName == "" {
			return errors.New("can't get fileName(param 'dl')")
		}

		u.RawQuery = "?dl&fit=crop&crop=entropy&w=" + strings.Split(resolution, "x")[0] + "&h=" + strings.Split(resolution, "x")[1]
		log.Println(u)
		go func() { ch <- dl(u, dir, fileName, i) }()
	}
	for i := 0; i < len(list); i++ {
		if err = <-ch; err != nil {
			log.Println(err)
		}
	}
	return nil
}

func dl(u *url.URL, dir string, fileName string, i int) error {
	log.Printf("[%d]downloading: %s\n", i, fileName)
	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	ioutil.WriteFile(path.Clean(dir)+string(os.PathSeparator)+fileName, buf, 0666)
	log.Printf("[%d]download okay\t%s", i, u)
	return nil
}

func getUrl(pageCount int, keyword string) ([]string, error) {
	baseUrl := searchUrl
	if strings.Trim(keyword, " ") == "" {
		baseUrl = indexUrl
	}
	dl := baseUrl + url.QueryEscape(keyword)
	list := make([]string, 0, pageCount*15)
	for page := 1; page <= pageCount; page++ {
		println("get page ", page)
		req, err := http.NewRequest(http.MethodGet, dl, nil)
		if err != nil {
			return nil, err
		}

		q := req.URL.Query()
		q.Add("format", "js")
		q.Add("seed", time.Now().Format("2018-06-23 15:45:58")+"  0000")
		q.Add("page", strconv.Itoa(page))
		req.URL.RawQuery = q.Encode()

		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			return nil, err
		}
		content, err := ioutil.ReadAll(resp.Body)
		contentStr := string(content)
		startIdx := strings.Index(string(contentStr), startFlag)

		endIdx := strings.LastIndex(string(contentStr), endFlag) + len(endFlag)
		htmlContent := contentStr[startIdx:endIdx]

		htmlContent = strings.NewReplacer(`\n`, "\n", `\"`, "\"", `\/`, "/").Replace(htmlContent)

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
		doc.Find("a[download]").Each(func(i int, selection *goquery.Selection) {
			for _, attr := range selection.Nodes[0].Attr {
				if attr.Key == "href" {
					list = append(list, attr.Val)
					log.Println("find link: " + attr.Val)
				}
			}
		})
	}
	return list, nil
}
