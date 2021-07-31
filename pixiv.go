package pixivrank

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"
)

type PixivTask struct {
	Client     *http.Client
	Cookie     string
	Proxy      string
	R18Enabled bool
	Quantity   int // numbers of images to download.
}

type Illust struct {
	Contents []struct {
		IllustID int `json:"illust_id"`
	} `json:"contents"`
}

func (pt *PixivTask) Run() {
	if err := pt.ReadCookies(); err != nil {
		log.Println("Failed to read cookies: " + err.Error())
	}
	// get the list of illust-id
	ridList, err := pt.GetRankIDList()
	if err != nil {
		log.Println("Failed to get ranklist: " + err.Error())
	}
	// get url of images
	urlList := make([]string, len(ridList))
	wg := &sync.WaitGroup{}
	for k, v := range ridList {
		wg.Add(1)
		go func(k, v int) {
			defer wg.Done()
			urlList[k], _ = pt.GetImageURL(v)
		}(k, v)
	}
	wg.Wait()
	log.Println(urlList)
}

func (pt *PixivTask) GetRankIDList() ([]int, error) {
	url := fmt.Sprintf("https://www.pixiv.net/ranking.php?mode=%s&content=illust&p=1&format=json", pt.GetMode())
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
	illust := &Illust{}
	if err := json.Unmarshal(body, illust); err != nil {
		return nil, err
	}
	ids := make([]int, len(illust.Contents))
	for k, v := range illust.Contents {
		ids[k] = v.IllustID
	}
	return ids, err
}

func (pt *PixivTask) GetImageURL(id int) (string, error) {
	url := fmt.Sprintf("https://www.pixiv.net/artworks/%d", id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	resp, err := pt.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	reg := regexp.MustCompile(`(?:"original":")(.*?)(?:"\})`)
	str := reg.FindStringSubmatch(string(body))
	if len(str) < 2 {
		return "", errors.New("cant find")
	}
	return str[1], nil
}

// func (pt *PixivTask) DownloadImage(url string) error {

// }

// GetMode return the state of R18.
func (pt *PixivTask) GetMode() (mode string) {
	if pt.R18Enabled {
		mode = "daily_r18"
	} else {
		mode = "daily"
	}
	return
}

func (pt *PixivTask) ReadCookies() error {
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
