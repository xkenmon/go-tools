package download

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	pic_info_url     = "https://www.mystart.com/api/get_background_stats_extension/"
	ext              = "569e9c14a6fde15e017c4eb5"
	lang             = "en"
	pic_download_url = "https://gallery.mystartcdn.com/mycats/images/"
)

type PicInfo struct {
	Success            bool   `json:"success"`
	StartDate          string `json:"startDate"`
	EndDate            string `json:"endDate"`
	TotalNbLike        int    `json:"total_nb_like"`
	TotalNbUnlike      int    `json:"total_nb_unlike"`
	TotalLikeNocumul   int    `json:"total_like_nocumul"`
	TotalUnlikeNocumul int    `json:"total_unlike_nocumul"`
	List               []struct {
		Ext         string `json:"ext"`
		ImageID     string `json:"image_id"`
		CumulLike   int    `json:"cumul_like"`
		CumulUnlike int    `json:"cumul_unlike"`
		NoOfLike    int    `json:"no_of_like"`
		NoOfUnlike  int    `json:"no_of_unlike"`
	} `json:"list"`
}

func GetPicInfo(startdate string) *PicInfo {
	if startdate == "" {
		startdate = "1970-01-01"
	}
	req, _ := http.NewRequest(http.MethodGet, pic_info_url, nil)
	q := req.URL.Query()
	q.Add("ext", ext)
	q.Add("lang", lang)
	q.Add("startdate", startdate)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	json_content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	info := &PicInfo{}
	err = json.Unmarshal(json_content, info)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return info
}

func DownloadPic(id, path string) (*string, error) {
	if id == "" {
		return nil, errors.New("id can't empty")
	}
	resp, err := http.Get(pic_download_url + id + ".jpeg")
	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Status code is %d, not OK(200)!", resp.StatusCode))
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	filepath := path + "/" + id + ".jpeg"
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(filepath, buf, 0666)
	return &filepath, err
}
