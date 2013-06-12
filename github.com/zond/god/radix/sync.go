package radix

import (
	"bytes"
	"github.com/zond/god/common"
)

const (
	subTreeError = "Illegal, only one level of sub trees supported"
)

// HashTree is something that can be synchronized using Sync.
type HashTree interface {
	Hash() []byte

	Configuration() (conf map[string]string, timestamp int64)
	Configure(conf map[string]string, timestamp int64)

	Finger(key []Nibble) *Print
	GetTimestamp(key []Nibble) (byteValue []byte, timestamp int64, present bool)
	PutTimestamp(key []Nibble, byteValue []byte, present bool, expected, timestamp int64) bool
	DelTimestamp(key []Nibble, expected int64) bool

	SubConfiguration(key []byte) (conf map[string]string, timestamp int64)
	SubConfigure(key []byte, conf map[string]string, timestamp int64)

	SubFinger(key, subKey []Nibble) (result *Print)
	SubGetTimestamp(key, subKey []Nibble) (byteValue []byte, timestamp int64, present bool)
	SubPutTimestamp(key, subKey []Nibble, byteValue []byte, present bool, subExpected, subTimestamp int64) bool
	SubDelTimestamp(key, subKey []Nibble, subExpected int64) bool
	SubClearTimestamp(key []Nibble, expected, timestamp int64) (deleted int)
	SubKillTimestamp(key []Nibble, expected int64) (deleted int)
}

// Sync synchronizes HashTrees using their fingerprints and mutators.
type Sync struct {
	source      HashTree
	destination HashTree
	from        []Nibble
	to          []Nibble
	destructive bool
	putCount    int
	delCount    int
}

func NewSync(source, destination HashTree) *Sync {
	return &Sync{
		source:      source,
		destination: destination,
	}
}

// From defines from what key, inclusive, that this Sync will synchronize.
func (self *Sync) From(from []byte) *Sync {
	self.from = Rip(from)
	return self
}

// To defines from what key, exclusive, that this Sync will synchronize.
func (self *Sync) To(to []byte) *Sync {
	self.to = Rip(to)
	return self
}

// Destroy defines that this Sync will delete whatever it copies from the source Tree.
func (self *Sync) Destroy() *Sync {
	self.destructive = true
	return self
}

// PutCount returns the number of entries this Sync has inserted into the destination Tree.
func (self *Sync) PutCount() int {
	return self.putCount
}

// DelCount returns the number of entries this Sync has deleted from the source Tree.
func (self *Sync) DelCount() int {
	return self.delCount
}

// Run will start this Sync and return when it is finished.
func (self *Sync) Run() *Sync {
	// If we have from and to, and they are equal, that means this sync is over an empty set... just ignore it
	if self.from != nil && self.to != nil && nComp(self.from, self.to) == 0 {
		return self
	}
	if self.destructive || bytes.Compare(self.source.Hash(), self.destination.Hash()) != 0 {
		sourceConf, sourceTs := self.source.Configuration()
		_, destTs := self.destination.Configuration()
		if sourceTs > destTs {
			self.destination.Configure(sourceConf, sourceTs)
		}
		self.synchronize(self.source.Finger(nil), self.destination.Finger(nil))
	}
	return self
}

// potentiallyWithinLimits will check if the given key can contain children between the from and to Nibbles
// for this Sync when considering the namespace circular.
func (self *Sync) potentiallyWithinLimits(key []Nibble) bool {
	if self.from == nil || self.to == nil {
		return true
	}
	cmpKey := toBytes(key)
	cmpFrom := toBytes(self.from)
	cmpTo := toBytes(self.to)
	m := len(cmpKey)
	if m > len(cmpFrom) {
		m = len(cmpFrom)
	}
	if m > len(cmpTo) {
		m = len(cmpTo)
	}
	return common.BetweenII(cmpKey[:m], cmpFrom[:m], cmpTo[:m])
}

// withinLimits will check i the given key is actually between the from and to Nibbles for this Sync.
func (self *Sync) withinLimits(key []Nibble) bool {
	if self.from == nil || self.to == nil {
		return true
	}
	return common.BetweenIE(toBytes(key), toBytes(self.from), toBytes(self.to))
}

// synchronize will recursively run the actual synchronization.
func (self *Sync) synchronize(sourcePrint, destinationPrint *Print) {
	// If there is a source key
	if sourcePrint.Exists {
		// If it represents a node containing synchronizable data, and it is within our limits
		if !sourcePrint.Empty && self.withinLimits(sourcePrint.Key) {
			// If it contains a sub tree
			if sourcePrint.SubTree {
				// If the sub tree in the destination is not equal to the sub tree in the source
				if bytes.Compare(sourcePrint.TreeHash, destinationPrint.TreeHash) != 0 {
					// If the source is empty, but not the destination, and the source is newer than the destination
					if sourcePrint.TreeSize == 0 && destinationPrint.TreeSize > 0 && sourcePrint.TreeDataTimestamp > destinationPrint.TreeDataTimestamp {
						// Clear the destination and count the number of removed keys.
						self.putCount += self.destination.SubClearTimestamp(sourcePrint.Key, destinationPrint.TreeDataTimestamp, sourcePrint.TreeDataTimestamp)
					}
					// Synchronize the sub trees. If the destination is previously cleared, this will copy the configuration.
					subSync := NewSync(&subTreeWrapper{
						self.source,
						sourcePrint.Key,
					}, &subTreeWrapper{
						self.destination,
						sourcePrint.Key,
					})
					if self.destructive {
						subSync.Destroy()
					}
					subSync.Run()
					self.putCount += subSync.PutCount()
					self.delCount += subSync.DelCount()
				} else if self.destructive {
					// If the trees are equal, but this Sync is destructive, just remove the source sub tree.
					self.delCount += self.source.SubKillTimestamp(sourcePrint.Key, sourcePrint.TreeDataTimestamp)
				}
			}
			// If the source has a byte value.
			if sourcePrint.Timestamp > 0 {
				// If the destination print is not covered by the source print (it is not equal and it is older)
				if !sourcePrint.coveredBy(destinationPrint) {
					// If the source still contains the same timestamp
					if value, timestamp, present := self.source.GetTimestamp(sourcePrint.Key); timestamp == sourcePrint.timestamp() {
						// Put the found data in the destination
						if self.destination.PutTimestamp(sourcePrint.Key, value, present, destinationPrint.timestamp(), sourcePrint.timestamp()) {
							self.putCount++
						}
					}
				}
				// If we are destructive and contain something
				if self.destructive && !sourcePrint.Empty {
					// Remove the byte value in the source.
					if self.source.DelTimestamp(sourcePrint.Key, sourcePrint.timestamp()) {
						self.delCount++
					}
				}
			}
		}
		// For each child of the source print
		for index, subPrint := range sourcePrint.SubPrints {
			// If there is a child node there, and it might contain matching keys
			if subPrint.Exists && self.potentiallyWithinLimits(subPrint.Key) {
				// If we are destructive, or if the prints are dissimilar
				if self.destructive || (!destinationPrint.Exists || !subPrint.equals(destinationPrint.SubPrints[index])) {
					// Synchronize the children
					self.synchronize(
						self.source.Finger(subPrint.Key),
						self.destination.Finger(subPrint.Key),
					)
				}
			}
		}
	}
}
