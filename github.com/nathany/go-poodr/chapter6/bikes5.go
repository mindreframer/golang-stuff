// ############## Page 131 ##############
package main

import "fmt"

/*
  Bicycle
*/
type Bicycle struct {
  Size     string
  Chain    string
  TireSize string
}

func (bike Bicycle) Spares() Parts {
  return Parts{
    "chain":     bike.Chain,
    "tire_size": bike.TireSize,
  }
}

/*
  Road Bike
*/
type RoadBike struct {
  Bicycle
  TapeColor string
}

func (bike RoadBike) Spares() Parts {
  return Parts{
    "tape_color": bike.TapeColor,
  }.merge(bike.Bicycle.Spares())
}

/*
  Mountain Bike
*/
type MountainBike struct {
  Bicycle
  FrontShock, RearShock string
}

func (bike MountainBike) Spares() Parts {
  return Parts{
    "rear_shock": bike.RearShock,
  }.merge(bike.Bicycle.Spares())
}

/*
  Parts
*/
type Parts map[string]string

// merge parts but don't overwrite what's there
func (parts Parts) merge(defaults Parts) Parts {
  for k, v := range defaults {
    if _, present := parts[k]; !present {
      parts[k] = v
    }
  }
  return parts
}

/*
  Main
*/
func main() {
  roadBike := RoadBike{
    Bicycle:   Bicycle{Size: "M", Chain: "10-speed", TireSize: "23"},
    TapeColor: "red"}
  fmt.Println(roadBike.Spares())

  mountainBike := MountainBike{
    Bicycle:    Bicycle{Size: "S", Chain: "10-speed", TireSize: "2.1"},
    FrontShock: "Manitou", RearShock: "Fox"}
  fmt.Println(mountainBike.Spares())
}
