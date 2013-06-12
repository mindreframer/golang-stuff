package main

import (
    "crypto/dsa"
    "crypto/rand"
    "crypto/sha1"
    "encoding/asn1"
    "encoding/pem"
    "flag"
    "io"
    "io/ioutil"
    "log"
    "math/big"
)

const (
    KeyFile = "dsa.key"
)

var (
    message = flag.String("message", "Nuke the site from orbit, it's the only way to be sure.", "The message to sign")
    do      = flag.String("do", "sign", "The operation to do, verify or sign (default)")
    rc      = flag.String("r", "", "The r to use when verifying")
    sc      = flag.String("s", "", "The s to use when verifying")
)

func HashMessage() []byte {
    h := sha1.New()
    _, err := io.WriteString(h, *message)
    if err != nil {
        log.Fatalf("failed to hash message: %s", err)
    }
    return h.Sum(nil)
}

type DsaKeyFormat struct {
    Version       int
    P, Q, G, Y, X *big.Int
}

func SaveKey(key *dsa.PrivateKey) {
    val := DsaKeyFormat{
        P:  key.P, Q: key.Q, G: key.G,
        Y:  key.Y, X: key.X,
    }
    bytes, err := asn1.Marshal(val)
    if err != nil {
        log.Fatalf("failed marshalling key to asn1: %s", err)
    }
    block := &pem.Block{
        Type:  "DSA PRIVATE KEY",
        Bytes: bytes,
    }
    err = ioutil.WriteFile(KeyFile, pem.EncodeToMemory(block), 0644)
    if err != nil {
        log.Fatalf("failed saving key to file %s: %s", KeyFile, err)
    }
}

func ReadKey() (*dsa.PrivateKey, error) {
    bytes, err := ioutil.ReadFile(KeyFile)
    if err != nil {
        return nil, err
    }
    block, _ := pem.Decode(bytes)
    val := new(DsaKeyFormat)
    _, err = asn1.Unmarshal(block.Bytes, val)
    if err != nil {
        return nil, err
    }
    key := &dsa.PrivateKey{
        PublicKey: dsa.PublicKey{
            Parameters: dsa.Parameters{
                P:  val.P,
                Q:  val.Q,
                G:  val.G,
            },
            Y:  val.Y,
        },
        X:  val.X,
    }
    return key, nil
}

func MakeKey() *dsa.PrivateKey {
    key := new(dsa.PrivateKey)
    err := dsa.GenerateParameters(&key.Parameters, rand.Reader, dsa.L2048N256)
    if err != nil {
        log.Fatalf("failed to parameters: %s", err)
    }
    err = dsa.GenerateKey(key, rand.Reader)
    if err != nil {
        log.Fatalf("failed to generate key: %s", err)
    }
    return key
}

func Key() *dsa.PrivateKey {
    key, err := ReadKey()
    if err != nil {
        log.Printf("failed reading keyfile, making a new one: %s", err)
        key = MakeKey()
        SaveKey(key)
    }
    return key
}

func Sign() {
    key := Key()
    hash := HashMessage()
    r, s, err := dsa.Sign(rand.Reader, key, hash)
    if err != nil {
        log.Fatalf("failed to sign message: %s", err)
    }
    log.Printf("r: %v", r)
    log.Printf("s: %v", s)
}

func Verify() {
    r := big.NewInt(0)
    r.SetString(*rc, 10)

    s := big.NewInt(0)
    s.SetString(*sc, 10)

    hash := HashMessage()
    key := Key()
    if dsa.Verify(&key.PublicKey, hash, r, s) {
        log.Println("message is valid!")
    } else {
        log.Println("message is invalid :(")
        log.Println("did you use the -r and -s flags to pass the r and s values?")
    }
}

func main() {
    flag.Parse()
    switch *do {
    case "sign":
        Sign()
    case "verify":
        Verify()
    default:
        log.Fatalf("%s is not a valid operation, must be one of sign or verify", *do)
    }
}
