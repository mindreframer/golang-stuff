package main

type bee struct {
	x, y    int
	actions []beeAction
}

func NewBee(x, y int) *bee {
	return &bee{x: x, y: y, actions: make([]beeAction, 0, 10)}
}

func (b *bee) fly(direction direction, distance int) beeAction {
	xStart := b.x
	yStart := b.y
	b.x, b.y = direction.move(b.x, b.y, distance, distance)
	action := beeAction{min(xStart, b.x), min(yStart, b.y), max(xStart, b.x), max(yStart, b.y)}
	b.actions = append(b.actions, action)
	return action
}

func (b *bee) isTired() bool {
	return len(b.actions) == maxBeeActions
}

type beeAction struct {
	xStart, yStart, xEnd, yEnd int
}
