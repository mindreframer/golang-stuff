package main

type direction struct {
	move   func(currentX, currentY, offsetX, offsetY int) (int, int)
	String func() string
}

var north = direction{
	move: func(currentX, currentY, offsetX, offsetY int) (int, int) {
		return currentX, currentY - offsetY
	},
	String: func() string { return "N" },
}

var west = direction{
	move: func(currentX, currentY, offsetX, offsetY int) (int, int) {
		return currentX - offsetX, currentY
	},
	String: func() string { return "W" },
}

var south = direction{
	move: func(currentX, currentY, offsetX, offsetY int) (int, int) {
		return currentX, currentY + offsetY
	},
	String: func() string { return "S" },
}

var east = direction{
	move: func(currentX, currentY, offsetX, offsetY int) (int, int) {
		return currentX + offsetX, currentY
	},
	String: func() string { return "E" },
}
