package vena

import (
	"circuit/kit/xor"
	"path"
)

type Config struct {
	Anchor string // Root anchor for the shards
	Shard  []*ShardConfig
}

type ShardConfig struct {
	Key   xor.Key
	Host  string
	Dir   string
	Cache int
}

func (c *Config) ShardAnchor(key xor.Key) string {
	return path.Join(c.Anchor, key.String())
}
