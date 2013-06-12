package main

type direction struct {
	move func(currentX, currentY, forward, right int) (int, int)
	turn func() direction
}

var north = direction{
	move: func(currentX, currentY, forward, right int) (int, int) {
		return wrap(currentX+right, fieldWidth), wrap(currentY-forward%fieldHeight, fieldHeight)
	},
}

var west = direction{
	move: func(currentX, currentY, forward, right int) (int, int) {
		return wrap(currentX-forward%fieldWidth, fieldWidth), wrap(currentY-right%fieldHeight, fieldHeight)
	},
}

var south = direction{
	move: func(currentX, currentY, forward, right int) (int, int) {
		return wrap(currentX-right%fieldWidth, fieldWidth), wrap(currentY+forward, fieldHeight)
	},
}

var east = direction{
	move: func(currentX, currentY, forward, right int) (int, int) {
		return wrap(currentX+forward, fieldWidth), wrap(currentY+right, fieldHeight)
	},
}

func init() {
	north.turn = func() direction { return west }
	west.turn = func() direction { return south }
	south.turn = func() direction { return east }
	east.turn = func() direction { return north }
}

func wrap(value, max int) int {
	if value < 0 {
		value += max
	}
	return value % max
}
