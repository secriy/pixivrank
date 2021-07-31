package pixivrank

import (
	"testing"
)

func TestPixivTask(t *testing.T) {
	pt := &PixivTask{
		Client:     NewClientWithPorxy("http://127.0.0.1:1080"),
		R18Enabled: true,
	}
	pt.Run()
}
