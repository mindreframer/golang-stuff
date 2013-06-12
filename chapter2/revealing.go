// ############## Page 28 ##############
package main

import "fmt"

/*
  "If you can control the input, pass in a useful object, but if you are
  compelled to take a messy structure, hide the mess even from yourself."
*/

type RevealingReferences struct {
  wheels []wheel
}

func NewRevealingReferences(data [][2]int) RevealingReferences {
  return RevealingReferences{wheelify(data)}
}

func (ref RevealingReferences) Wheels() []wheel {
  return ref.wheels
}

func (ref RevealingReferences) Diameters() (diameters []int) {
  for _, wheel := range ref.Wheels() {
    diameters = append(diameters, wheel.rim+(wheel.tire*2))
  }
  return diameters
}

// ... now everyone can send rim/tire to wheel
type wheel struct {
  rim, tire int
}

func wheelify(data [][2]int) (wheels []wheel) {
  for _, cell := range data {
    wheels = append(wheels, wheel{cell[0], cell[1]})
  }
  return wheels
}

func main() {
  // rim and tire sizes (now in milimeters!) in a 2d array
  diameters := [][2]int{{622, 20}, {622, 23}, {559, 30}, {559, 40}}
  revealing := NewRevealingReferences(diameters)
  fmt.Println(revealing.Wheels())
  fmt.Println(revealing.Diameters())
}
