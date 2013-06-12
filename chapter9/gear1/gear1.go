// ############## Page 36 ##############
package gear1

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
