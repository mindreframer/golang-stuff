// ############## Page 32 ##############
package main

import "fmt"

type Gear struct {
  chainring float64 // number of teeth
  cog       float64
  wheel     *wheel
}

func NewGear(chainring, cog, rim, tire float64) Gear {
  return Gear{chainring, cog, &wheel{rim, tire}}
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

func (gear Gear) Wheel() wheel {
  return *gear.wheel
}

// Wheel (internal)
type wheel struct {
  rim, tire float64
}

func (wheel wheel) Diameter() float64 {
  return wheel.Rim() + (wheel.Tire() * 2)
}

// getters
func (wheel wheel) Rim() float64 {
  return wheel.rim
}

func (wheel wheel) Tire() float64 {
  return wheel.tire
}

// Main
func main() {
  fmt.Println(NewGear(52, 11, 26, 1.5).GearInches())  // => 137.0909090909091
  fmt.Println(NewGear(52, 11, 24, 1.25).GearInches()) // => 125.27272727272728
}
