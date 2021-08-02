package main

import (
	"flag"

	pr "github.com/secriy/pixivrank"
)

var (
	proxy      = flag.String("p", "", "Proxy address, such as 'http://127.0.0.1:1080'")
	r18Enabled = flag.Bool("r", false, "Enable for R18 mode")
	numbers    = flag.Int("n", 0, "The number of images to download.")
)

func main() {
	flag.Parse()
	pt := new(pr.PixivTask)
	pt.RunTask(*proxy, *r18Enabled, *numbers)
}
