// ############## Page 87 ##############
package main

import "fmt"

type Bicycle string

type Trip struct {
  bicycles []Bicycle
}

/*
  For the first example, we use a BicyclePreparer. This isn't as restrictive
  as requiring a Mechanic, but PrepareBicycles is still too specific for what
  Trip wants to achieve.
*/
type BicyclePreparer interface {
  PrepareBicycles([]Bicycle)
}

func (trip Trip) Prepare(mechanic BicyclePreparer) {
  mechanic.PrepareBicycles(trip.bicycles)
}

type Mechanic string

func (mechanic Mechanic) PrepareBicycles(bicycles []Bicycle) {
  fmt.Printf("Preparing %d bicycles...\n", len(bicycles))
  for _, bicycle := range bicycles {
    mechanic.PrepareBicycle(bicycle)
  }
}

func (mechanic Mechanic) PrepareBicycle(bicycle Bicycle) {
  fmt.Println("Preparing bicycle...", bicycle)
}

func main() {
  mechanic := new(Mechanic)
  trip := Trip{[]Bicycle{"my bike", "your bike"}}
  trip.Prepare(mechanic)
}
