package skyd

// A slice of Property objects.
type PropertyList []*Property

// Determines the length of an event slice.
func (s PropertyList) Len() int {
	return len(s)
}

// Compares two properties in a list.
func (s PropertyList) Less(i, j int) bool {
	return s[i].Id < s[j].Id
}

// Swaps two properties in a list
func (s PropertyList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
