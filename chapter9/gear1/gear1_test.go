package gear1_test

import (
  "go-poodr/chapter9/gear1"
  "testing"
)

func TestCalculatesDiameter(t *testing.T) {
  wheel := gear1.NewWheel(26, 1.5)
  if wheel.Diameter() < 28.99 || wheel.Diameter() > 29.01 {
    t.Errorf("Expected diameter %f to be 29", wheel.Diameter())
  }
}

func TestCalculatesGearInches(t *testing.T) {
  gear := gear1.NewGear(52, 11, 26, 1.5)
  if gear.GearInches() < (137.1-0.1) || gear.GearInches() > (137.1+0.1) {
    t.Errorf("Expected gear inches %f to be 137.1", gear.GearInches())
  }
}
