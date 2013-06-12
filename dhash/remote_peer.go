package dhash

import (
	"github.com/zond/god/common"
	"time"
)

type remotePeer common.Remote

func (self remotePeer) ActualTime() (result time.Time) {
	if err := (common.Remote)(self).Call("Timenet.ActualTime", 0, &result); err != nil {
		result = time.Now()
	}
	return
}
