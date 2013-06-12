package main

import (
    "crypto/md5"
    "crypto/sha1"
    "crypto/sha256"
    "crypto/sha512"
    "flag"
    "hash"
    "io"
    "log"
)

var (
    algorithm = flag.String("algorithm", "md5", "The algorithm to use. Must be one of {md5,sha1,sha256,sha512}")
    message   = flag.String("message", "Go, The Standard Library", "The message to hash")
)

func GetHash() hash.Hash {
    switch *algorithm {
    case "md5":
        return md5.New()
    case "sha1":
        return sha1.New()
    case "sha256":
        return sha256.New()
    case "sha512":
        return sha512.New()
    default:
        log.Fatalf("No hash algorithm %s found", *algorithm)
    }
    panic("unreachable")
}

func main() {
    flag.Parse()
    hash := GetHash()
    io.WriteString(hash, *message)
    log.Printf("%x", hash.Sum(nil))
}
