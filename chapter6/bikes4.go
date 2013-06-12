// ############## Page 119-122 ##############
package main

import "fmt"

/*
  Promote abstract behavior up to Bicycle rather than extracting it down.
*/
type Bicycle struct {
  Size string
}

type RoadBike struct {
  Bicycle
  TapeColor string
}

func (bike RoadBike) Spares() map[string]string {
  return map[string]string{
    "chain":      "10-speed",
    "tire_size":  "23", // milimeters
    "tape_color": bike.TapeColor,
  }
}

type MountainBike struct {
  Bicycle
  FrontShock, RearShock string
}

func (bike MountainBike) Spares() map[string]string {
  return map[string]string{
    "chain":      "10-speed",
    "tire_size":  "2.1", // inches
    "rear_shock": bike.RearShock,
  }
}

func main() {
  roadBike := RoadBike{Bicycle: Bicycle{Size: "M"}, TapeColor: "red"}
  fmt.Println(roadBike.Size)
  fmt.Println(roadBike.Spares())

  mountainBike := MountainBike{Bicycle: Bicycle{Size: "S"}, FrontShock: "Manitou", RearShock: "Fox"}
  fmt.Println(mountainBike.Size)
  fmt.Println(mountainBike.Spares())
}
