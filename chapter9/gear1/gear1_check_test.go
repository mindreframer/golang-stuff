package gear1_test

import (
  "go-poodr/chapter9/gear1"
  . "launchpad.net/gocheck"
  "math"
  "testing"
)

// Hook up gocheck into the gotest runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) TestCalculatesDiameter(c *C) {
  wheel := gear1.NewWheel(26, 1.5)
  c.Assert(wheel.Diameter(), Within, 0.01, 29.0)
}

func (s *MySuite) TestCalculatesGearInches(c *C) {
  gear := gear1.NewGear(52, 11, 26, 1.5)
  c.Assert(gear.GearInches(), Within, 0.01, 137.1)
}

/*
  Within Delta Custom Checker
  Would be nice if gocheck included something like this.
*/
type withinChecker struct {
  *CheckerInfo
}

var Within Checker = &withinChecker{
  &CheckerInfo{Name: "Within", Params: []string{"obtained", "delta", "expected"}},
}

func (c *withinChecker) Check(params []interface{}, names []string) (result bool, error string) {
  obtained, ok := params[0].(float64)
  if !ok {
    return false, "obtained must be a float64"
  }
  delta, ok := params[1].(float64)
  if !ok {
    return false, "delta must be a float64"
  }
  expected, ok := params[2].(float64)
  if !ok {
    return false, "expected must be a float64"
  }
  return math.Abs(obtained-expected) <= delta, ""
}
