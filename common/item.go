package common

type Item struct {
	Key       []byte
	SubKey    []byte
	Value     []byte
	Exists    bool
	Timestamp int64
	TTL       int
	Index     int
	Sync      bool
}
