package main

import (
  "bytes"
  "fmt"
  "github.com/zond/god/murmur"
  "github.com/zond/god/radix"
  "os"
  "runtime/pprof"
)

func main() {
  f, err := os.Create("cpuprofile")
  if err != nil {
    panic(err.Error())
  }
  f2, err := os.Create("memprofile")
  if err != nil {
    panic(err.Error())
  }

  pprof.StartCPUProfile(f)
  defer pprof.StopCPUProfile()
  defer pprof.WriteHeapProfile(f2)

  m := radix.NewTree()
  var k []byte
  for i := 0; i < 100000; i++ {
    k = murmur.HashString(fmt.Sprint(i))
    m.Put(k, k, 1)
    x, _, _ := m.Get(k)
    if bytes.Compare(x, k) != 0 {
      panic("humbug")
    }
  }
}
