package main

import (
    "crypto/rand"
    "crypto/rc4"
    "encoding/pem"
    "flag"
    "io/ioutil"
    "log"
)

const (
    EncryptedFile = "rc4.enc"
    KeyFile       = "rc4.key"
)

var (
    do      = flag.String("do", "encrypt", "The operation to perform, decrypt or encrypt (default)")
    message = flag.String("message", "Wolverines attack at dawn. Red Dawn.", "The message to encrypt")
    keySize = flag.Int("keysize", 256, "Key size in bytes")
)

func MakeKey() []byte {
    key := make([]byte, *keySize)
    n, err := rand.Read(key)
    if err != nil {
        log.Fatalf("failed to read new random key: %s", err)
    }
    if n < *keySize {
        log.Fatalf("failed to read entire key, only read %d out of %d", n, *keySize)
    }
    return key
}

func SaveKey(filename string, key []byte) {
    block := &pem.Block{
        Type:  "RC4 KEY",
        Bytes: key,
    }
    err := ioutil.WriteFile(filename, pem.EncodeToMemory(block), 0644)
    if err != nil {
        log.Fatalf("failed saving key to %s: %s", filename, err)
    }
}

func ReadKey(filename string) ([]byte, error) {
    key, err := ioutil.ReadFile(filename)
    if err != nil {
        return key, err
    }
    block, _ := pem.Decode(key)
    return block.Bytes, nil
}

func Key() []byte {
    key, err := ReadKey(KeyFile)
    if err != nil {
        log.Println("failed reading key, making a new one...")
        key = MakeKey()
        SaveKey(KeyFile, key)
    }
    return key
}

func Cipher() *rc4.Cipher {
    key := Key()
    cipher, err := rc4.NewCipher(key)
    if err != nil {
        log.Fatalf("failed to make RC4 cipher: %s", err)
    }
    return cipher
}

func Encrypt() {
    cipher := Cipher()
    text := []byte(*message)
    cipher.XORKeyStream(text, text)
    err := ioutil.WriteFile(EncryptedFile, text, 0644)
    if err != nil {
        log.Fatalf("failed to write encrypted file: %s", err)
    }
}

func Decrypt() {
    cipher := Cipher()
    bytes, err := ioutil.ReadFile(EncryptedFile)
    if err != nil {
        log.Fatalf("failed to read encrypted file. Did you encrypt first? %s", err)
    }
    cipher.XORKeyStream(bytes, bytes)
    log.Printf("decrypted message: %s", bytes)
}

func main() {
    flag.Parse()
    switch *do {
    case "encrypt":
        Encrypt()
    case "decrypt":
        Decrypt()
    default:
        log.Fatalf("%s not a valid operation. Must be one of encrypt or decrypt", *do)
    }
}
