package common

type ConfItem struct {
	TreeKey   []byte
	Key       string
	Value     string
	Timestamp int64
	TTL       int
}

type Conf struct {
	TreeKey   []byte
	Data      map[string]string
	Timestamp int64
}
