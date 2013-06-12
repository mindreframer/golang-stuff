// ############## Page 110 ##############
package main

import "fmt"

type Bicycle struct {
  Size      bikeSize
  TapeColor string
}

/*
  Though I didn't implement an accessor for Size and TapeColor,
  the fact that they are strings preserves our ability to change.
  (in contrast with gear3.go)
*/
type bikeSize string

func (size bikeSize) String() string {
  return fmt.Sprintf("(%s)", string(size)) // return a different string
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
