package main

import (
    "crypto/ecdsa"
    "crypto/elliptic"
    "crypto/rand"
    "crypto/sha1"
    "flag"
    "io"
    "log"
)

var message = flag.String("message", "Nuke the site from orbit, it's the only way to be sure.", "The message to sign")

func HashMessage() []byte {
    h := sha1.New()
    _, err := io.WriteString(h, *message)
    if err != nil {
        log.Fatalf("failed to hash message: %s", err)
    }
    return h.Sum(nil)
}

func Key() *ecdsa.PrivateKey {
    key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
    if err != nil {
        log.Fatalf("failed to generate key: %s", err)
    }
    return key
}

func main() {
    flag.Parse()

    key := Key()
    hash := HashMessage()
    r, s, err := ecdsa.Sign(rand.Reader, key, hash)
    if err != nil {
        log.Fatalf("failed to sign message: %s", err)
    }
    log.Printf("r: %s", r)
    log.Printf("s: %s", s)

    if ecdsa.Verify(&key.PublicKey, hash, r, s) {
        log.Println("message is valid!")
    } else {
        log.Println("message invalid :(")
    }
}
