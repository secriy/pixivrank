# Pixiv-Rank

用于下载 Pixiv 排行榜图片的小工具，包含 Python 版本和 Go 版本。

## Usage

1. 在浏览器上登录 Pixiv 账号，复制 cookie
2. 在*pixiv-rank/*目录下创建文件*.cookies*，并用文本编辑器打开，把 cookie 粘贴进去

### Python

```shell
python3 pixiv.py
```

### Go

#### Build

- make

  ```shell
  make build
  ```

- go build

  ```go
  go build -o pixivrank.exe ./cmd/
  ```

#### Execute

```shell
Usage of pixivrank.exe:
  -n int
        The number of images to download.
  -p string
        Proxy address, such as 'http://127.0.0.1:1080'
  -r    Enable for R18 mode
```

- 使用`-r`参数指定是否爬取 R18 分类图片
- 使用`-n [number]`参数指定每个排行榜需要爬取的图片数量，默认为全部爬取，最大为 50
- 使用`-p`参数指定使用的代理，如`http://127.0.0.1:1080`，默认不使用代理（不使用代理在中国大陆地区通常无法正常爬取）

#### Example

```shell
make build
pixivrank.exe -p http://127.0.0.1:1080 -r -n 10	// 爬取非R18图片10张，R18图片10张，共20张
pixivrank.exe -p http://127.0.0.1:1080 // 爬取非R18图片50张
```

## Document

```
package pixivrank // import "github.com/secriy/pixivrank"


FUNCTIONS

func NewClient() *http.Client
    NewClient return the http client instance.
    返回携带cookie存储的 http client

func NewClientWithPorxy(proxyUrl string) *http.Client
    NewClientWithProxy return the http client instance with a proxy server.
    返回包含代理的 http client


TYPES

type PixivTask struct {
        Client     *http.Client
        Cookie     string
        Proxy      string
        R18Enabled bool
        Numbers    int // numbers of images to download.
}
    PixivTask is use to store the task information.
    用于记录任务的部分参数

func (pt *PixivTask) DownloadImage(url string, id int) error
    DownloadImage create the image files and download from pixiv by URL.
    下载指定URL的单张图片到img文件夹

func (pt *PixivTask) ImageURL(id int) string
    ImageURL return the original URL of the image which find by illust_id.
    通过插画ID获取插画图片的原图下载地址

func (pt *PixivTask) RankIDList() ([]int, error)
    RankIDList return the list of all illust_id which in daily rank.
    获取当日的插画排行榜列表

func (pt *PixivTask) RunTask(proxy string, r18Enabled bool, numbers int)
    RunTask is to initialize the task instance.
    用于创建Task实例，如果开启R18，则分别爬取R18和非R18的图片

func (pt *PixivTask) Task()
    Task is used to excute one task.
    用于执行一次单一模式（R18/非R18）的爬取任务
```
