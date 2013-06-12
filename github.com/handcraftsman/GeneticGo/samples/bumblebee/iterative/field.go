package main

type field struct {
	xFlowerLookup        map[int]*flower
	yFlowerLookup        map[int]*flower
	width, height        int
	numberOfFlowers      int
	numberOfFlowersFound int
}

func NewField(fieldWidth, fieldHeight int, flowerPoints *[]point) *field {
	numberOfFlowers := len(*flowerPoints)
	field := field{
		xFlowerLookup:   make(map[int]*flower, numberOfFlowers),
		yFlowerLookup:   make(map[int]*flower, numberOfFlowers),
		width:           fieldWidth,
		height:          fieldHeight,
		numberOfFlowers: numberOfFlowers,
	}

	for _, point := range *flowerPoints {
		flower := NewFlower(point.x, point.y)
		field.xFlowerLookup[point.x] = flower
		field.yFlowerLookup[point.y] = flower
	}
	return &field
}

func (f *field) markFlowersInXPath(action beeAction) {
	xStart, xEnd := sort(action.xStart, action.xEnd)

	yStart := max(action.yStart-2, 0)
	yStop := min(action.yStart+2, f.height)
	for y := yStart; y < yStop; y++ {
		flower := f.yFlowerLookup[y]
		if flower == nil {
			continue
		}

		if !(*flower).hasPollen {
			continue
		}

		if (*flower).x >= xStart && (*flower).x <= xEnd {
			(*flower).hasPollen = false
			f.numberOfFlowersFound++
		}
	}
}

func (f *field) markFlowersInYPath(action beeAction) {
	yStart, yEnd := sort(action.yStart, action.yEnd)

	xStart := max(action.xStart-2, 0)
	xStop := min(action.xStart+2, f.width)
	for x := xStart; x < xStop; x++ {
		flower := f.xFlowerLookup[x]
		if flower == nil {
			continue
		}

		if !(*flower).hasPollen {
			continue
		}

		if (*flower).y >= yStart && (*flower).y <= yEnd {
			(*flower).hasPollen = false
			f.numberOfFlowersFound++
		}
	}
}

func (f *field) allFlowersFound() bool {
	return f.numberOfFlowersFound == f.numberOfFlowers
}

func (f *field) String() string {
	s := ""
	for _, flower := range f.xFlowerLookup {
		fl := flower.String()
		if fl != "" {
			s += fl + "\n"
		}
	}
	return s
}

type point struct {
	x, y int
}
