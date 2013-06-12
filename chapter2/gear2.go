// ############## Page 20 ##############
package main

import "fmt"

type Gear struct {
  chainring float64 // number of teeth
  cog       float64
  rim       float64 // diameters
  tire      float64
}

func NewGear(chainring, cog, rim, tire float64) Gear {
  return Gear{chainring, cog, rim, tire}
}

func (gear Gear) Chainring() float64 {
  return gear.chainring
}

func (gear Gear) Cog() float64 {
  return gear.cog
}

func (gear Gear) Rim() float64 {
  return gear.rim
}

func (gear Gear) Tire() float64 {
  return gear.tire
}

func (gear Gear) Ratio() float64 {
  return gear.Chainring() / gear.Cog()
}

func (gear Gear) GearInches() float64 {
  return gear.Ratio() * (gear.Rim() + (gear.Tire() * 2))
}

func main() {
  fmt.Println(NewGear(52, 11, 26, 1.5).GearInches())  // => 137.0909090909091
  fmt.Println(NewGear(52, 11, 24, 1.25).GearInches()) // => 125.27272727272728
}
