package main

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/tls"
    "crypto/x509"
    "crypto/x509/pkix"
    "encoding/pem"
    "flag"
    "io"
    "io/ioutil"
    "log"
    "math/big"
    "net"
    "time"
)

const (
    CertFile = "tls.crt"
    KeyFile  = "tls.key"
)

var (
    do      = flag.String("do", "serve", "The operation to perform, key, cert, or serve (default)")
    keySize = flag.Int("keysize", 2048, "The RSA keysize to use")
)

func MakeKey() *rsa.PrivateKey {
    key, err := rsa.GenerateKey(rand.Reader, *keySize)
    if err != nil {
        log.Fatalf("failed to create RSA key: %s", err)
    }
    return key
}

func PemEncodeKey(key *rsa.PrivateKey) []byte {
    block := &pem.Block{
        Type:  "RSA PRIVATE KEY",
        Bytes: x509.MarshalPKCS1PrivateKey(key),
    }
    return pem.EncodeToMemory(block)
}

func SaveKey(filename string, key *rsa.PrivateKey) {
    err := ioutil.WriteFile(filename, PemEncodeKey(key), 0644)
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

func SaveCert(filename string, cert []byte) []byte {
    block := &pem.Block{
        Type:  "CERTIFICATE",
        Bytes: cert,
    }
    bytes := pem.EncodeToMemory(block)
    err := ioutil.WriteFile(filename, bytes, 0644)
    if err != nil {
        log.Fatalf("failed saving cert to %s: %s", filename, err)
    }
    return bytes
}

func MakeCert() tls.Certificate {
    key := Key()
    now := time.Now()
    template := &x509.Certificate{
        SerialNumber: big.NewInt(1),
        Subject: pkix.Name{
            Country:            []string{"CA"},
            Province:           []string{"Alberta"},
            Locality:           []string{"Edmonton"},
            Organization:       []string{"The Standard Library"},
            OrganizationalUnit: []string{"Go, The Standard Library"},
            CommonName:         "localhost",
        },
        NotBefore: now,
        NotAfter:  now.Add(24 * 365 * time.Hour), // 1 year
        KeyUsage:  0,
    }
    cert, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
    if err != nil {
        log.Fatalf("failed creating certificate: %s", err)
    }
    cert = SaveCert(CertFile, cert)
    c, err := tls.X509KeyPair(cert, PemEncodeKey(key))
    if err != nil {
        log.Fatalf("failed to load certificate: %s", err)
    }
    return c
}

func Cert() tls.Certificate {
    cert, err := tls.LoadX509KeyPair(CertFile, KeyFile)
    if err != nil {
        log.Printf("failed loading certificate, generating a new one: %s", err)
        cert = MakeCert()
    }
    return cert
}

func Config() *tls.Config {
    return &tls.Config{
        Certificates: []tls.Certificate{Cert()},
    }
}

func Serve() {
    addr := "localhost:4443"
    conn, err := net.Listen("tcp", addr)
    if err != nil {
        log.Fatalf("failed to listen on %s: %s", addr, err)
    }

    config := Config()
    listener := tls.NewListener(conn, config)
    log.Printf("listening on %s, connect with 'openssl s_client -tls1 -connect %s'", addr, addr)
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Fatalf("failed to accept: %s", err)
        }
        log.Printf("connection accepted from %s", conn.RemoteAddr())
        go func(c net.Conn) {
            _, err := io.Copy(c, c)
            if err != nil {
                log.Printf("error copying: %s", err)
            }
            log.Println("closing connection")
            c.Close()
        }(conn)
    }
}

func main() {
    flag.Parse()
    switch *do {
    case "serve":
        Serve()
    case "cert":
        Cert()
    case "key":
        Key()
    default:
        log.Fatalf("%s is not a valid operation, must be one of serve, cert, or key", *do)
    }
}
