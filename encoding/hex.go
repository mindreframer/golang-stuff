package main

import (
    "encoding/hex"
    "io/ioutil"
    "log"
    "os"
)

func dumpFile() {
    data, err := ioutil.ReadFile("hex.go")
    if err != nil {
        log.Fatalf("failed reading file: %s", err)
    }
    dumper := hex.Dumper(os.Stdout)
    defer dumper.Close()
    log.Println("dumping hex.go to stdout")
    dumper.Write(data)
}

func main() {
    hero := []byte("Batman and Robin")
    log.Printf("hero: %s", hero)
    encoded := hex.EncodeToString(hero)
    log.Printf("encoded: %s", encoded)
    decoded, _ := hex.DecodeString(encoded)
    log.Printf("decoded: %s", decoded)

    dumpFile()
}
