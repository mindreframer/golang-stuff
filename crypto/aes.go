package main

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/pem"
    "flag"
    "fmt"
    "io/ioutil"
    "log"
)

const (
    KeyFile       = "aes.%d.key"
    EncryptedFile = "aes.%d.enc"
)

var (
    IV      = []byte("batman and robin") // 16 bytes
    message = flag.String("message", "Batman is Bruce Wayne", "The message to encrypt")
    keySize = flag.Int("keysize", 32, "The keysize in bytes to use: 16, 24, or 32 (default)")
    do      = flag.String("do", "encrypt", "The operation to perform: decrypt or encrypt (default) ")
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
        Type:  "AES KEY",
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
    file := fmt.Sprintf(KeyFile, *keySize)
    key, err := ReadKey(file)
    if err != nil {
        log.Println("failed reading keyfile, making a new one...")
        key = MakeKey()
        SaveKey(file, key)
    }
    return key
}

func MakeCipher() cipher.Block {
    c, err := aes.NewCipher(Key())
    if err != nil {
        log.Fatalf("failed making the AES cipher: %s", err)
    }
    return c
}

func Crypt(bytes []byte) []byte {
    blockCipher := MakeCipher()
    stream := cipher.NewCTR(blockCipher, IV)
    stream.XORKeyStream(bytes, bytes)
    return bytes
}

func Encrypt() {
    encrypted := Crypt([]byte(*message))
    err := ioutil.WriteFile(fmt.Sprintf(EncryptedFile, *keySize), encrypted, 0644)
    if err != nil {
        log.Fatalf("failed writing encrypted file: %s", err)
    }
}

func Decrypt() {
    bytes, err := ioutil.ReadFile(fmt.Sprintf(EncryptedFile, *keySize))
    if err != nil {
        log.Fatalf("failed reading encrypted file: %s", err)
    }
    plaintext := Crypt(bytes)
    log.Printf("decrypted message: %s", plaintext)
}

func main() {
    flag.Parse()

    switch *keySize {
    case 16, 24, 32:
        // Keep calm and carry on...
    default:
        log.Fatalf("%d is not a valid keysize. Must be one of 16, 24, 32", *keySize)
    }

    switch *do {
    case "encrypt":
        Encrypt()
    case "decrypt":
        Decrypt()
    default:
        log.Fatalf("%s is not a valid operation. Must be one of encrypt or decrypt", *do)
    }
}
