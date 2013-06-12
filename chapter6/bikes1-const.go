// ############## Page 107 ##############
package main

import "fmt"

type Bicycle struct {
  Size      string
  TapeColor string
}

type Part int

/*
  It probably makes sense to use constants where the Ruby code was symbols.
*/
const (
  Chain Part = iota
  TireSize
  TapeColor
)

/* What the mechanic needs to bring along */
func (bike Bicycle) Spares() map[Part]string {
  return map[Part]string{
    Chain:     "10-speed",
    TireSize:  "23",
    TapeColor: bike.TapeColor,
  }
}

func main() {
  bike := Bicycle{Size: "M", TapeColor: "Red"}
  fmt.Println(bike.Size)
  fmt.Println(bike.Spares())
}
