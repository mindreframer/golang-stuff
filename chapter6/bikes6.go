// ############## Page 136 ##############
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

/* Presumably the hook (subtype) could be set during construction. */
func (bike Bicycle) Spares(hook Hooks) Parts {
  return Parts{
    "chain":     bike.Chain,
    "tire_size": bike.TireSize,
  }.merge(hook.localSpares())
}

/*
  Road Bike
*/
type RoadBike struct {
  Bicycle
  TapeColor string
}

func (bike RoadBike) localSpares() Parts {
  return Parts{
    "tape_color": bike.TapeColor,
  }
}

/*
  Mountain Bike
*/
type MountainBike struct {
  Bicycle
  FrontShock, RearShock string
}

func (bike MountainBike) localSpares() Parts {
  return Parts{
    "rear_shock": bike.RearShock,
  }
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

type Hooks interface {
  localSpares() Parts
}

/*
  Main
*/
func main() {
  roadBike := RoadBike{
    Bicycle:   Bicycle{Size: "M", Chain: "10-speed", TireSize: "23"},
    TapeColor: "red"}
  fmt.Println(roadBike.Spares(roadBike))

  mountainBike := MountainBike{
    Bicycle:    Bicycle{Size: "S", Chain: "10-speed", TireSize: "2.1"},
    FrontShock: "Manitou", RearShock: "Fox"}
  fmt.Println(mountainBike.Spares(mountainBike))
}
