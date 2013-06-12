package main

import (
    "crypto"
    "crypto/md5"
    "crypto/rand"
    "crypto/rsa"
    "crypto/sha1"
    "crypto/sha256"
    "crypto/sha512"
    "crypto/x509"
    "encoding/pem"
    "flag"
    "hash"
    "io/ioutil"
    "log"
)

const (
    KeyFile       = "rsa.key"
    SignatureFile = "rsa.sig"
    EncryptedFile = "rsa.enc"
)

var (
    keySize       = flag.Int("keysize", 2048, "The size of the key in bits")
    do            = flag.String("do", "encrypt", "The operation to perform, decrypt or encrypt (default)")
    message       = flag.String("message", "The revolution has begun!", "The message to encrypt")
    hashAlgorithm = flag.String("algorithm", "sha256", "The hash algorithm to use. Must be one of md5, sha1, sha256 (default), sha512")
)

func MakeKey() *rsa.PrivateKey {
    key, err := rsa.GenerateKey(rand.Reader, *keySize)
    if err != nil {
        log.Fatalf("failed to create RSA key: %s", err)
    }
    return key
}

func SaveKey(filename string, key *rsa.PrivateKey) {
    block := &pem.Block{
        Type:  "RSA PRIVATE KEY",
        Bytes: x509.MarshalPKCS1PrivateKey(key),
    }
    err := ioutil.WriteFile(filename, pem.EncodeToMemory(block), 0644)
    if err != nil {
        log.Fatalf("failed saving key to %s: %s", filename, err)
    }
}

func ReadKey(filename string) (*rsa.PrivateKey, error) {
    bytes, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    block, _ := pem.Decode(bytes)
    key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        return nil, err
    }
    return key, nil
}

func Key() *rsa.PrivateKey {
    key, err := ReadKey(KeyFile)
    if err != nil {
        log.Printf("failed to read key, creating a new one: %s", err)
        key = MakeKey()
        SaveKey(KeyFile, key)
    }
    return key
}

func HashAlgorithm() (hash.Hash, crypto.Hash) {
    switch *hashAlgorithm {
    case "md5":
        return md5.New(), crypto.MD5
    case "sha1":
        return sha1.New(), crypto.SHA1
    case "sha256":
        return sha256.New(), crypto.SHA256
    case "sha512":
        return sha512.New(), crypto.SHA512
    default:
        log.Fatalf("%s is not a valid hash algorithm. Must be one of md5, sha1, sha256, sha512")
    }
    panic("not reachable")
}

func HashMessage(data []byte) []byte {
    h, _ := HashAlgorithm()
    h.Write(data)
    return h.Sum(nil)
}

func Encrypt() {
    h, ha := HashAlgorithm()
    key := Key()
    encrypted, err := rsa.EncryptOAEP(h, rand.Reader, &key.PublicKey, []byte(*message), nil)
    if err != nil {
        log.Fatalf("encryption failed: %s", err)
    }
    signature, err := rsa.SignPKCS1v15(rand.Reader, key, ha, HashMessage(encrypted))
    if err != nil {
        log.Fatalf("signing failed; %s", err)
    }
    err = ioutil.WriteFile(EncryptedFile, encrypted, 0644)
    if err != nil {
        log.Fatalf("failed saving encrypted data: %s", err)
    }
    err = ioutil.WriteFile(SignatureFile, signature, 0644)
    if err != nil {
        log.Fatalf("failed saving signature data: %s", err)
    }
}

func Decrypt() {
    key := Key()
    h, ha := HashAlgorithm()
    encrypted, err := ioutil.ReadFile(EncryptedFile)
    if err != nil {
        log.Fatalf("failed reading encrypted data: %s", err)
    }

    signature, err := ioutil.ReadFile(SignatureFile)
    if err != nil {
        log.Fatalf("failed saving signature data: %s", err)
    }

    if err = rsa.VerifyPKCS1v15(&key.PublicKey, ha, HashMessage(encrypted), signature); err != nil {
        log.Fatalf("message not valid: %s", err)
    } else {
        log.Printf("message is valid!")
    }

    plaintext, err := rsa.DecryptOAEP(h, rand.Reader, key, encrypted, nil)
    if err != nil {
        log.Fatalf("failed decrypting: %s", err)
    }
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
        log.Fatalf("%s is not a valid operation. Must be one of encrypt or decrypt")
    }
}
