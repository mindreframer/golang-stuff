package main

import (
    "crypto/cipher"
    "crypto/des"
    "crypto/rand"
    "encoding/pem"
    "flag"
    "io/ioutil"
    "log"
)

const (
    KeyFile       = "des.key"
    EncryptedFile = "des.enc"
)

var (
    IV      = []byte("superman") // 8 bytes
    triple  = flag.Bool("3", false, "Use 3DES")
    message = flag.String("message", "Batman is Bruce Wayne", "The message to encrypt")
    do      = flag.String("do", "encrypt", "The operation to perform: decrypt or encrypt (default) ")
)

func MakeKey() []byte {
    size := 8
    if *triple {
        size *= 3
    }
    key := make([]byte, size)
    n, err := rand.Read(key)
    if err != nil {
        log.Fatalf("failed to read new random key: %s", err)
    }
    if n < size {
        log.Fatalf("failed to read entire key, only read %d out of %d", n, size)
    }
    return key
}

func SaveKey(filename string, key []byte) {
    block := &pem.Block{
        Type:  "DES KEY",
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
        log.Println("failed reading keyfile, making a new one...")
        key = MakeKey()
        SaveKey(KeyFile, key)
    }
    return key
}

func MakeCipher() cipher.Block {
    var c cipher.Block
    var err error
    if *triple {
        c, err = des.NewTripleDESCipher(Key())
    } else {
        c, err = des.NewCipher(Key())
    }
    if err != nil {
        log.Fatalf("failed making the DES cipher: %s", err)
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
    err := ioutil.WriteFile(EncryptedFile, encrypted, 0644)
    if err != nil {
        log.Fatalf("failed writing encrypted file: %s", err)
    }
}

func Decrypt() {
    bytes, err := ioutil.ReadFile(EncryptedFile)
    if err != nil {
        log.Fatalf("failed reading encrypted file: %s", err)
    }
    plaintext := Crypt(bytes)
    log.Printf("decrypted message: %s", plaintext)
}

func main() {
    flag.Parse()
    switch *do {
    case "encrypt":
        Encrypt()
    case "decrypt":
        Decrypt()
    default:
        log.Fatalf("%s is not a valid operation. Must be one of encrypt or decrypt", *do)
    }
}
