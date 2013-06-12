// ############## Page 52 ##############
package main

import "fmt"

/*
  Gear
*/
type Gear struct {
  Chainring, Cog float64
}

func NewGear(chainring, cog float64) *Gear {
  return &Gear{chainring, cog}
}

func (gear Gear) GearInches(diameter float64) float64 {
  return gear.Ratio() * diameter
}

func (gear Gear) Ratio() float64 {
  return gear.Chainring / gear.Cog
}

/*
  Wheel
*/
type Wheel struct {
  Rim, Tire float64
  Gear      *Gear
}

func NewWheel(rim, tire, chainring, cog float64) *Wheel {
  return &Wheel{rim, tire, NewGear(chainring, cog)}
}

func (wheel Wheel) Diameter() float64 {
  return wheel.Rim + (wheel.Tire * 2)
}

func (wheel Wheel) GearInches() float64 {
  return wheel.Gear.GearInches(wheel.Diameter())
}

/*
  Main
*/
func main() {
  wheel := NewWheel(26, 1.5, 52, 11)
  fmt.Println(wheel.GearInches()) // => 137.0909090909091
}
