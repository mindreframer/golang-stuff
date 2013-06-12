// ############## Page 107 ##############
package main

import "fmt"

type Bicycle struct {
  Size      string
  TapeColor string
}

/* What the mechanic needs to bring along */
func (bike Bicycle) Spares() map[string]string {
  return map[string]string{
    "chain":      "10-speed",
    "tire_size":  "23",
    "tape_color": bike.TapeColor,
  }
}

func main() {
  bike := Bicycle{Size: "M", TapeColor: "Red"}
  fmt.Println(bike.Size)
  fmt.Println(bike.Spares())
}
