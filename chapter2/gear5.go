// ############## Page 32 ##############
package main

import (
  "fmt"
  "math"
)

/*
  Gear
*/
type Gear struct {
  chainring float64 // number of teeth
  cog       float64
  wheel     *Wheel
}

// note: Go doesn't have default arguments, but we could provide another "constructor"
func NewGear(chainring, cog float64, wheel *Wheel) *Gear {
  return &Gear{chainring, cog, wheel}
}

func (gear Gear) Ratio() float64 {
  return gear.Chainring() / gear.Cog()
}

func (gear Gear) GearInches() float64 {
  return gear.Ratio() * gear.Wheel().Diameter()
}

// getters
func (gear Gear) Chainring() float64 {
  return gear.chainring
}

func (gear Gear) Cog() float64 {
  return gear.cog
}

func (gear Gear) Wheel() Wheel {
  return *gear.wheel
}

/*
  Wheel
*/
type Wheel struct {
  rim, tire float64
}

func NewWheel(rim, tire float64) *Wheel {
  return &Wheel{rim, tire}
}

func (wheel Wheel) Diameter() float64 {
  return wheel.Rim() + (wheel.Tire() * 2)
}

func (wheel Wheel) Circumference() float64 {
  return wheel.Diameter() * math.Pi
}

// getters
func (wheel Wheel) Rim() float64 {
  return wheel.rim
}

func (wheel Wheel) Tire() float64 {
  return wheel.tire
}

/*
  Main
*/
func main() {
  wheel := NewWheel(26, 1.5)
  fmt.Println(wheel.Circumference()) // => 91.106186954104

  fmt.Println(NewGear(52, 11, wheel).GearInches()) // => 137.0909090909091

  fmt.Println(NewGear(52, 11, nil).Ratio()) // => 4.7272727272727275
}
