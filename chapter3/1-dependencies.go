// ############## Page 36 ##############
package main

import "fmt"

/*
  Gear
*/
type Gear struct {
  Chainring, Cog float64
  Rim, Tire      float64
}

func NewGear(chainring, cog, rim, tire float64) *Gear {
  return &Gear{chainring, cog, rim, tire}
}

func (gear Gear) GearInches() float64 {
  /*
     anti-pattern
     Gear knows too much about Wheel
     * only willing to collaborate with Wheels
     * the message sent (Diameter)
     * the arguments
     * the order of the arguments
  */
  return gear.Ratio() * NewWheel(gear.Rim, gear.Tire).Diameter()
}

func (gear Gear) Ratio() float64 {
  return gear.Chainring / gear.Cog
}

/*
  Wheel
*/
type Wheel struct {
  Rim, Tire float64
}

func NewWheel(rim, tire float64) *Wheel {
  return &Wheel{rim, tire}
}

func (wheel Wheel) Diameter() float64 {
  return wheel.Rim + (wheel.Tire * 2)
}

/*
  Main
*/
func main() {
  gear := NewGear(52, 11, 26, 1.5)
  fmt.Println(gear.GearInches()) // => 137.0909090909091
  fmt.Println(gear.Ratio())      // => 4.7272727272727275
}
