package radix

// each will iterate over the tree in order
func (self *node) each(prefix []Nibble, use int, f nodeIterator) (cont bool) {
	cont = true
	if self != nil {
		prefix = append(prefix, self.segment...)
		if !self.empty && (use == 0 || self.use&use != 0) {
			cont = f(Stitch(prefix), self.byteValue, self.treeValue, self.use, self.timestamp)
		}
		if cont {
			for _, child := range self.children {
				cont = child.each(prefix, use, f)
				if !cont {
					break
				}
			}
		}
	}
	return
}

// reverseEach will iterate over the tree in reverse order
func (self *node) reverseEach(prefix []Nibble, use int, f nodeIterator) (cont bool) {
	cont = true
	if self != nil {
		prefix = append(prefix, self.segment...)
		for i := len(self.children) - 1; i >= 0; i-- {
			cont = self.children[i].reverseEach(prefix, use, f)
			if !cont {
				break
			}
		}
		if cont {
			if !self.empty && (use == 0 || self.use&use != 0) {
				cont = f(Stitch(prefix), self.byteValue, self.treeValue, self.use, self.timestamp)
			}
		}
	}
	return
}

// eachBetween will iterate between min and max, including each depending on mincmp and maxcmp, in order
func (self *node) eachBetween(prefix, min, max []Nibble, mincmp, maxcmp, use int, f nodeIterator) (cont bool) {
	cont = true
	prefix = append(prefix, self.segment...)
	if !self.empty && (use == 0 || self.use&use != 0) && (min == nil || nComp(prefix, min) > mincmp) && (max == nil || nComp(prefix, max) < maxcmp) {
		cont = f(Stitch(prefix), self.byteValue, self.treeValue, self.use, self.timestamp)
	}
	if cont {
		for _, child := range self.children {
			if child != nil {
				childKey := make([]Nibble, len(prefix)+len(child.segment))
				copy(childKey, prefix)
				copy(childKey[len(prefix):], child.segment)
				mmi := len(childKey)
				if mmi > len(min) {
					mmi = len(min)
				}
				mma := len(childKey)
				if mma > len(max) {
					mma = len(max)
				}
				if (min == nil || nComp(childKey[:mmi], min[:mmi]) > -1) && (max == nil || nComp(childKey[:mma], max[:mma]) < 1) {
					cont = child.eachBetween(prefix, min, max, mincmp, maxcmp, use, f)
				}
				if !cont {
					break
				}
			}
		}
	}
	return
}

// eachBetween will iterate between min and max, including each depending on mincmp and maxcmp, in reverse order
func (self *node) reverseEachBetween(prefix, min, max []Nibble, mincmp, maxcmp, use int, f nodeIterator) (cont bool) {
	cont = true
	prefix = append(prefix, self.segment...)
	var child *node
	for i := len(self.children) - 1; i >= 0; i-- {
		child = self.children[i]
		if child != nil {
			childKey := make([]Nibble, len(prefix)+len(child.segment))
			copy(childKey, prefix)
			copy(childKey[len(prefix):], child.segment)
			mmi := len(childKey)
			if mmi > len(min) {
				mmi = len(min)
			}
			mma := len(childKey)
			if mma > len(max) {
				mma = len(max)
			}
			if (min == nil || nComp(childKey[:mmi], min[:mmi]) > -1) && (max == nil || nComp(childKey[:mma], max[:mma]) < 1) {
				cont = child.reverseEachBetween(prefix, min, max, mincmp, maxcmp, use, f)
			}
			if !cont {
				break
			}
		}
	}
	if cont {
		if !self.empty && (use == 0 || self.use&use != 0) && (min == nil || nComp(prefix, min) > mincmp) && (max == nil || nComp(prefix, max) < maxcmp) {
			cont = f(Stitch(prefix), self.byteValue, self.treeValue, self.use, self.timestamp)
		}
	}
	return
}

// sizeBetween will count values between min and max, including each depending on mincmp and maxcmp, counting values of types included in use (byteValue and/or treeValue)
func (self *node) sizeBetween(prefix, min, max []Nibble, mincmp, maxcmp, use int) (result int) {
	prefix = append(prefix, self.segment...)
	if !self.empty && (use == 0 || self.use&use != 0) && (min == nil || nComp(prefix, min) > mincmp) && (max == nil || nComp(prefix, max) < maxcmp) {
		if use == 0 || self.use&use&byteValue != 0 {
			result++
		}
		if use == 0 || self.use&use&treeValue != 0 {
			result += self.treeValue.Size()
		}
	}
	for _, child := range self.children {
		if child != nil {
			childKey := make([]Nibble, len(prefix)+len(child.segment))
			copy(childKey, prefix)
			copy(childKey[len(prefix):], child.segment)
			mmi := len(childKey)
			if mmi > len(min) {
				mmi = len(min)
			}
			mma := len(childKey)
			if mma > len(max) {
				mma = len(max)
			}
			mires := nComp(childKey[:mmi], min[:mmi])
			mares := nComp(childKey[:mma], max[:mma])
			if (min == nil || mires > -1) && (max == nil || mares < 1) {
				if (min == nil || mires > 0) && (max == nil || mares < 0) {
					if use == 0 {
						result += child.realSize
					} else {
						if use&byteValue != 0 {
							result += child.byteSize
						}
						if use&treeValue != 0 {
							result += child.treeSize
						}
					}
				} else {
					result += child.sizeBetween(prefix, min, max, mincmp, maxcmp, use)
				}
			}
		}
	}
	return
}

// eachBetweenIndex will iterate over the tree between index min and max, inclusive.
// Missing min or max will mean 'from the start' or 'to the end' respectively.
func (self *node) eachBetweenIndex(prefix []Nibble, count int, min, max *int, use int, f nodeIndexIterator) (cont bool) {
	cont = true
	prefix = append(prefix, self.segment...)
	if !self.empty && (use == 0 || self.use&use != 0) && (min == nil || count >= *min) && (max == nil || count <= *max) {
		cont = f(Stitch(prefix), self.byteValue, self.treeValue, self.use, self.timestamp, count)
		if use == 0 || self.use&use&byteValue != 0 {
			count++
		}
		if use == 0 || self.use&use&treeValue != 0 {
			count += self.treeValue.Size()
		}
	}
	if cont {
		relevantChildSize := 0
		for _, child := range self.children {
			if child != nil {
				relevantChildSize = 0
				if use == 0 {
					relevantChildSize = child.realSize
				} else {
					if use&byteValue != 0 {
						relevantChildSize += child.byteSize
					}
					if use&treeValue != 0 {
						relevantChildSize += child.treeSize
					}
				}
				if (min == nil || relevantChildSize+count > *min) && (max == nil || count <= *max) {
					cont = child.eachBetweenIndex(prefix, count, min, max, use, f)
				}
				count += relevantChildSize
				if !cont {
					break
				}
			}
		}
	}
	return
}

// reverseEachBetweenIndex is like eachBetweenIndex, but iterates in reverse.
func (self *node) reverseEachBetweenIndex(prefix []Nibble, count int, min, max *int, use int, f nodeIndexIterator) (cont bool) {
	cont = true
	prefix = append(prefix, self.segment...)
	var child *node
	relevantChildSize := 0
	for i := len(self.children) - 1; i >= 0; i-- {
		child = self.children[i]
		if child != nil {
			relevantChildSize = 0
			if use == 0 {
				relevantChildSize = child.realSize
			} else {
				if use&byteValue != 0 {
					relevantChildSize += child.byteSize
				}
				if use&treeValue != 0 {
					relevantChildSize += child.treeSize
				}
			}
			if (min == nil || relevantChildSize+count > *min) && (max == nil || count <= *max) {
				cont = child.reverseEachBetweenIndex(prefix, count, min, max, use, f)
			}
			count += relevantChildSize
			if !cont {
				break
			}
		}
	}
	if cont {
		if !self.empty && (use == 0 || self.use&use != 0) && (min == nil || count >= *min) && (max == nil || count <= *max) {
			cont = f(Stitch(prefix), self.byteValue, self.treeValue, self.use, self.timestamp, count)
			if use == 0 || self.use&use&byteValue != 0 {
				count++
			}
			if use == 0 || self.use&use&treeValue != 0 {
				count += self.treeValue.Size()
			}
		}
	}
	return
}
