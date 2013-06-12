// ############## Page 26 ##############
package main

import "fmt"

type ObscuringReferences struct {
  data [][2]int
}

func NewObscuringReferences(data [][2]int) ObscuringReferences {
  return ObscuringReferences{data}
}

func (ref ObscuringReferences) Data() [][2]int {
  return ref.data
}

func (ref ObscuringReferences) Diameters() (diameters []int) {
  for _, cell := range ref.Data() {
    diameters = append(diameters, cell[0]+(cell[1]*2))
  }
  return diameters
}

func main() {
  // rim and tire sizes (now in milimeters!) in a 2d array
  diameters := [][2]int{{622, 20}, {622, 23}, {559, 30}, {559, 40}}
  obscure := NewObscuringReferences(diameters)
  fmt.Println(obscure.Data())
  fmt.Println(obscure.Diameters())
}
