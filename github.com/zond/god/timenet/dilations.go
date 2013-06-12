package timenet

import (
	"time"
)

const (
	dilationFactor = 5
)

type dilation struct {
	delta int64
	from  int64
}

func newDilation(delta int64) dilation {
	return dilation{delta, time.Now().UnixNano()}
}
func (self dilation) effect() (effect int64, done bool) {
	absDelta := self.delta
	if absDelta < 0 {
		absDelta *= -1
	}
	passed := float64(time.Now().UnixNano() - self.from)
	duration := float64(dilationFactor * absDelta)
	if passed > duration {
		effect = self.delta
		done = true
	} else {
		effect = int64(float64(self.delta) * (passed / duration))
		done = false
	}
	return
}

type dilations struct {
	content []dilation
}

func (self *dilations) delta() (sum int64) {
	for _, dilation := range self.content {
		sum += dilation.delta
	}
	return
}
func (self *dilations) effect() (temporaryEffect, permanentEffect int64) {
	newContent := make([]dilation, 0, len(self.content))
	for _, dilation := range self.content {
		thisEffect, done := dilation.effect()
		if done {
			permanentEffect += thisEffect
		} else {
			temporaryEffect += thisEffect
			newContent = append(newContent, dilation)
		}
	}
	self.content = newContent
	return
}
func (self *dilations) add(delta int64) {
	self.content = append(self.content, newDilation(delta))
}
