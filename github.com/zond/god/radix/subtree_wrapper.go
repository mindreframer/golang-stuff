package radix

type subTreeWrapper struct {
	parentTree HashTree
	key        []Nibble
}

func (self *subTreeWrapper) Configuration() (map[string]string, int64) {
	return self.parentTree.SubConfiguration(Stitch(self.key))
}
func (self *subTreeWrapper) Configure(conf map[string]string, ts int64) {
	self.parentTree.SubConfigure(Stitch(self.key), conf, ts)
}
func (self *subTreeWrapper) Hash() (hash []byte) {
	if p := self.parentTree.Finger(self.key); p != nil {
		hash = p.TreeHash
	}
	return
}
func (self *subTreeWrapper) Finger(subKey []Nibble) *Print {
	return self.parentTree.SubFinger(self.key, subKey)
}
func (self *subTreeWrapper) GetTimestamp(subKey []Nibble) (byteValue []byte, version int64, present bool) {
	return self.parentTree.SubGetTimestamp(self.key, subKey)
}
func (self *subTreeWrapper) PutTimestamp(subKey []Nibble, byteValue []byte, present bool, expected, version int64) bool {
	return self.parentTree.SubPutTimestamp(self.key, subKey, byteValue, present, expected, version)
}
func (self *subTreeWrapper) DelTimestamp(subKey []Nibble, expected int64) bool {
	return self.parentTree.SubDelTimestamp(self.key, subKey, expected)
}
func (self *subTreeWrapper) SubConfiguration(key []byte) (map[string]string, int64) {
	panic(subTreeError)
}
func (self *subTreeWrapper) SubConfigure(key []byte, conf map[string]string, ts int64) {
	panic(subTreeError)
}
func (self *subTreeWrapper) SubFinger(key, subKey []Nibble) (result *Print) {
	panic(subTreeError)
}
func (self *subTreeWrapper) SubGetTimestamp(key, subKey []Nibble) (byteValue []byte, version int64, present bool) {
	panic(subTreeError)
}
func (self *subTreeWrapper) SubPutTimestamp(key, subKey []Nibble, byteValue []byte, present bool, subExpected, subTimestamp int64) bool {
	panic(subTreeError)
}
func (self *subTreeWrapper) SubDelTimestamp(key, subKey []Nibble, subExpected int64) bool {
	panic(subTreeError)
}
func (self *subTreeWrapper) SubClearTimestamp(key []Nibble, expected, timestamp int64) int {
	panic(subTreeError)
}
func (self *subTreeWrapper) SubKillTimestamp(key []Nibble, expected int64) int {
	panic(subTreeError)
}
