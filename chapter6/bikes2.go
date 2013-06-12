// ############## Page 110 ##############
package main

import "fmt"

type style int

const (
  Road style = iota
  Mountain
)

type Bicycle struct {
  Style                 style
  Size, TapeColor       string
  FrontShock, RearShock string
}

/* Checking "style" starts down a slippery slope */
func (bike Bicycle) Spares() map[string]string {
  if bike.Style == Road {
    return map[string]string{
      "chain":      "10-speed",
      "tire_size":  "23", // milimeters
      "tape_color": bike.TapeColor,
    }
  }
  return map[string]string{
    "chain":      "10-speed",
    "tire_size":  "2.1", // inches
    "rear_shock": bike.RearShock,
  }
}

func main() {
  bike := Bicycle{Style: Mountain, Size: "S", FrontShock: "Manitou", RearShock: "Fox"}
  fmt.Println(bike.Spares())
}
