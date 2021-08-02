package pixivrank

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"sync"
)

// PixivTask is use to store the task information.
// 用于记录任务的部分参数
type PixivTask struct {
	Client     *http.Client
	Cookie     string
	Proxy      string
	R18Enabled bool
	Numbers    int // numbers of images to download.
}

type illust struct {
	Contents []struct {
		IllustID int `json:"illust_id"`
	} `json:"contents"`
}

// RunTask is to initialize the task instance.
// 用于创建Task实例，如果开启R18，则分别爬取R18和非R18的图片
func (pt *PixivTask) RunTask(proxy string, r18Enabled bool, numbers int) {
	if proxy == "" {
		pt.Client = NewClient()
	} else {
		pt.Client = NewClientWithPorxy(proxy)
	}
	pt.R18Enabled = r18Enabled
	pt.Numbers = numbers
	if pt.R18Enabled {
		pt.Task()
		pt.R18Enabled = false
	}
	pt.Task()
}

// Task is used to excute one task.
// 用于执行一次单一模式（R18/非R18）的爬取任务
func (pt *PixivTask) Task() {
	// read the cookies file
	if err := pt.readCookies(); err != nil {
		fmt.Println("Failed to read cookies: " + err.Error())
	}
	// get the list of illust-id
	ridList, err := pt.RankIDList()
	if err != nil {
		fmt.Println("Failed to get ranklist: " + err.Error())
	}
	// get url of images
	wg := &sync.WaitGroup{}
	// limit the numbers of images to download
	if pt.Numbers > 0 && pt.Numbers <= len(ridList) {
		ridList = ridList[:pt.Numbers]
	}
	// download the images
	for k, v := range ridList {
		wg.Add(1)
		go func(k, v int) {
			defer wg.Done()
			url := pt.ImageURL(v)
			if url != "" {
				err := pt.DownloadImage(url, v)
				if err != nil {
					fmt.Println(err)
				}
			}
		}(k, v)
	}
	wg.Wait()
}

// RankIDList return the list of all illust_id which in daily rank.
// 获取当日的插画排行榜列表
func (pt *PixivTask) RankIDList() ([]int, error) {
	url := fmt.Sprintf("https://www.pixiv.net/ranking.php?mode=%s&content=illust&p=1&format=json", pt.mode())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.96 Safari/537.36")
	req.Header.Add("cookie", pt.Cookie)
	resp, err := pt.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	illust := &illust{}
	if err := json.Unmarshal(body, illust); err != nil {
		return nil, err
	}
	ids := make([]int, len(illust.Contents))
	for k, v := range illust.Contents {
		ids[k] = v.IllustID
	}
	return ids, err
}

// ImageURL return the original URL of the image which find by illust_id.
// 通过插画ID获取插画图片的原图下载地址
func (pt *PixivTask) ImageURL(id int) string {
	url := fmt.Sprintf("https://www.pixiv.net/artworks/%d", id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}
	resp, err := pt.Client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	reg := regexp.MustCompile(`(?:"original":")(.*?)(?:"\})`)
	str := reg.FindStringSubmatch(string(body))
	if len(str) < 2 {
		return ""
	}
	return str[1]
}

// DownloadImage create the image files and download from pixiv by URL.
// 下载指定URL的单张图片到img文件夹
func (pt *PixivTask) DownloadImage(url string, id int) error {
	pathStr := "img/" + pt.mode()
	if err := os.MkdirAll(pathStr, os.ModeDir); err != nil {
		return err
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("referer", fmt.Sprintf("https://www.pixiv.net/artworks/%d", id))
	resp, err := pt.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	file, err := os.Create(pathStr + "/" + path.Base(req.URL.Path))
	if err != nil {
		return err
	}
	wt := bufio.NewWriter(file)
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	wt.Flush()
	return nil
}

func (pt *PixivTask) mode() (mode string) {
	if pt.R18Enabled {
		mode = "daily_r18"
	} else {
		mode = "daily"
	}
	return
}

func (pt *PixivTask) readCookies() error {
	file, err := os.Open("./.cookies")
	if err != nil {
		return err
	}
	defer file.Close()
	text, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	pt.Cookie = string(text)
	return nil
}
