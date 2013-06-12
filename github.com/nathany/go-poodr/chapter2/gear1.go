// ############## Page 19 ##############
package main

import "fmt"

type Gear struct {
  chainring float64
  cog       float64
}

func NewGear(chainring, cog float64) Gear {
  return Gear{chainring, cog}
}

func (gear Gear) Chainring() float64 {
  return gear.chainring
}

func (gear Gear) Cog() float64 {
  return gear.cog
}

func (gear Gear) Ratio() float64 {
  return gear.Chainring() / gear.Cog()
}

func main() {
  fmt.Println(NewGear(52, 11).Ratio()) // => 4.7272727272727275
  fmt.Println(NewGear(30, 27).Ratio()) // => 1.1111111111111112
}
