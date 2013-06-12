// ############## Page 24 ##############

/*
  "Hide the variables, even from the class that defines them, by wrapping
  them in methods" (not like this)
*/
package road_to_ruin

type Gear struct {
  Chainring float64
  Cog       float64
}

func NewGear(chainring, cog float64) Gear {
  return Gear{chainring, cog}
}

func (gear Gear) Ratio() float64 {
  return gear.Chainring / gear.Cog // <== road to ruin
}
