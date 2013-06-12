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
  Now a Trip trusts the other objects to do their job, and doesn't know the
  details of how they do it.
*/
type Preparer interface {
  PrepareTrip(*Trip)
}

func (trip Trip) Prepare(preparers []Preparer) {
  for _, preparer := range preparers {
    preparer.PrepareTrip(&trip)
  }
}

/*
  Mechanics prepare the bicycles
*/
type Mechanic string

func (mechanic Mechanic) PrepareTrip(trip *Trip) {
  bicycles := trip.Bicycles

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

func (coordinator TripCoordinator) PrepareTrip(trip *Trip) {
  coordinator.BuyFood(trip.Customers)
}

func (coordinator TripCoordinator) BuyFood(customers []Customer) {
  fmt.Printf("%s is buying food for %s.\n", coordinator, customers)
}

/*
  Drivers gas up the vehicle and fill the water tank
*/
type Driver string

func (driver Driver) PrepareTrip(trip *Trip) {
  vehicle := trip.Vehicle
  driver.GasUp(vehicle)
  driver.FillWaterTank(vehicle)
}

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
  trip.Prepare([]Preparer{
    Mechanic("Joe"), TripCoordinator("Kim"), Driver("Dave"),
  })
}
