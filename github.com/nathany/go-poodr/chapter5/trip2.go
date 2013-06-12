// ############## Page 88 ##############
package main

import "fmt"

type Bicycle string
type Customer string
type Vehicle string

type Trip struct {
  Bicycles  []Bicycle
  Customers []Customer
  Vehicle   Vehicle
}

/*
  Trip preparation becomes more complicated
*/
func (trip Trip) Prepare(preparers []interface{}) {
  for _, preparer := range preparers {
    /*
       Not all type switches are bad, but this is an anti-pattern. Trip knows
       too many of the implementation details needed to prepare for a trip.
    */
    switch p := preparer.(type) {
    case Mechanic:
      p.PrepareBicycles(trip.Bicycles)
    case TripCoordinator:
      p.BuyFood(trip.Customers)
    case Driver:
      p.GasUp(trip.Vehicle)
      p.FillWaterTank(trip.Vehicle)
    }
  }
}

/*
  Mechanics prepare the bicycles
*/
type Mechanic string

func (mechanic Mechanic) PrepareBicycles(bicycles []Bicycle) {
  fmt.Printf("%s is preparing %d bicycles...\n", mechanic, len(bicycles))
  for _, bicycle := range bicycles {
    mechanic.PrepareBicycle(bicycle)
  }
}

func (mechanic Mechanic) PrepareBicycle(bicycle Bicycle) {
  fmt.Printf("Preparing bicycle... %s.\n", bicycle)
}

/*
  Trip Coordinators buy food for customers
*/
type TripCoordinator string

func (coordinator TripCoordinator) BuyFood(customers []Customer) {
  fmt.Printf("%s is buying food for %s.\n", coordinator, customers)
}

/*
  Drivers gas up the vehicle and fill the water tank
*/
type Driver string

func (driver Driver) GasUp(vehicle Vehicle) {
  fmt.Printf("%s is gassing up %s.\n", driver, vehicle)
}

func (driver Driver) FillWaterTank(vehicle Vehicle) {
  fmt.Printf("%s is filling the water tank for %s.\n", driver, vehicle)
}

func main() {
  trip := Trip{
    Bicycles:  []Bicycle{"my bike", "your bike"},
    Customers: []Customer{"me", "you"},
    Vehicle:   "the truck",
  }
  trip.Prepare([]interface{}{
    Mechanic("Joe"), TripCoordinator("Kim"), Driver("Dave"),
  })
}
