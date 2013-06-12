package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/hmac"
    "crypto/sha256"
    "flag"
    "log"
)

var (
    // 32 byte key for AES256, made from crypto/rand
    key     = []byte{0x98, 0x39, 0xea, 0x42, 0xd0, 0x3e, 0x36, 0x6b, 0xe3, 0x7b, 0x91, 0x6, 0x50, 0x5b, 0x7f, 0xc9, 0x93, 0x56, 0xaa, 0xa8, 0x96, 0x33, 0x7, 0xd7, 0xf7, 0x50, 0xa5, 0x3a, 0xdc, 0x8e, 0xe2, 0x9f}
    iv      = []byte("batman and robin") // 16 bytes
    message = flag.String("message", "Batman and Robin are coming", "The message to use")
)

func main() {
    flag.Parse()
    block, err := aes.NewCipher(key)
    if err != nil {
        log.Fatalf("failed making AES block cipher: %s", err)
    }
    bytes := []byte(*message)
    stream := cipher.NewCTR(block, iv)
    stream.XORKeyStream(bytes, bytes)
    hash := hmac.New(sha256.New, key)
    hash.Write(bytes)
    log.Printf("message: %s", *message)
    log.Printf("encrypted message (raw bytes): %v", bytes)
    log.Printf("HMAC: %x", hash.Sum(nil))
}
