// ############## Page 46 ##############
package main

import "fmt"

/*
  Options
*/
type Options map[string]interface{}

// merge options, but don't overwrite what's there (unlike Ruby merge)
func (options Options) merge(defaults Options) {
  for k, v := range defaults {
    if _, present := options[k]; !present {
      options[k] = v
    }
  }
}

/*
  Gear
*/
type Gear struct {
  Chainring, Cog float64
  Wheel          Diameter
}

type Diameter interface {
  Diameter() float64
}

/*
  This is how one could do hash-based initialization to remove argument order
  dependence like in the Ruby example.
*/
func NewGear(options Options) *Gear {
  options.merge(defaults())

  chainring, _ := options["chainring"].(float64)
  cog, _ := options["cog"].(float64)
  wheel, _ := options["wheel"].(Diameter) // default nil

  /*
     The key/value syntax of composite literals would be much simpler and more
     idiomatic. It could be used directly by getting rid of NewGear altogether,
     though some setup would be necessary if 0-value defaults are not enough.
  */
  return &Gear{Chainring: chainring, Cog: cog, Wheel: wheel}
}

func defaults() Options {
  return Options{
    "chainring": 40.0, // these must be floats for the type assertion!
    "cog":       18.0,
  }
}

func (gear Gear) GearInches() float64 {
  return gear.Ratio() * gear.diameter()
}

func (gear Gear) diameter() float64 {
  if gear.Wheel != nil {
    return gear.Wheel.Diameter()
  }
  return 0
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
  gear := NewGear(Options{
    "chainring": 52.0, // these must be floats for the type assertion!
    "cog":       11.0,
    "wheel":     NewWheel(26, 1.5),
  })
  fmt.Println(gear.GearInches()) // => 137.0909090909091
  fmt.Println(gear.Ratio())      // => 4.7272727272727275
}
