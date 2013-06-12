package main

type mower struct {
	direction   direction
	x, y        int
	actionCount int
}

func NewMower(x, y int, d direction) *mower {
	return &mower{x: x, y: y, direction: d}
}

func (m *mower) mow(f *field) {
	m.move(1, 0)
	f.cut(m.x, m.y)
}

func (m *mower) turn() {
	m.direction = m.direction.turn()
	m.actionCount++
}

func (m *mower) jump(f *field, forward, right int) {
	m.move(forward, right)
	f.cut(m.x, m.y)
}

func (m *mower) move(forward, right int) {
	m.x, m.y = m.direction.move(m.x, m.y, forward, right)
	m.actionCount++
}

func (m *mower) isOutOfFuel() bool {
	return m.actionCount == maxMowerActions
}
