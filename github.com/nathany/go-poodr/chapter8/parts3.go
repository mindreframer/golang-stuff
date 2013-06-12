// ############## Page 176-182 ##############
package main

import "fmt"

/*
  Bicycle
  Automatic delegation to Spares() is provided by Parts.
*/
type Bicycle struct {
  Size string
  Parts
}

/*
  Parts
  Parts is a slice vs. a struct containing one, giving it slice like behavior.
  (like on Page 173)
*/
type Parts []Part

func (parts Parts) Spares() (spares Parts) {
  for _, part := range parts {
    if part.NeedsSpare {
      spares = append(spares, part)
    }
  }
  return spares
}

/*
  Part
  I don't see much need for a PartsFactory to convert an array to objects. Go's
  composite literal syntax can use either name:values or just values, with the one
  gotcha that all parameters must be specified with the later syntax.
*/
type Part struct {
  Name        string
  Description string
  NeedsSpare  bool
}

var (
  RoadBikeParts = Parts{
    {"chain", "10-speed", true},
    {"tire_size", "23", true},
    {"tape_color", "red", true},
  }

  MountainBikeParts = Parts{
    {"chain", "10-speed", true},
    {"tire_size", "2.1", true},
    {"front_shock", "Manitou", false},
    {"rear_shock", "Fox", true},
  }

  RecumbentBikeParts = Parts{
    {"chain", "9-speed", true},
    {"tire_size", "28", true},
    {"flag", "tall and orange", true},
  }
)

func main() {
  roadBike := Bicycle{Size: "L", Parts: RoadBikeParts}
  mountainBike := Bicycle{Size: "L", Parts: MountainBikeParts}
  recumbentBike := Bicycle{Size: "L", Parts: RecumbentBikeParts}

  fmt.Println(roadBike.Spares())
  fmt.Println(mountainBike.Spares())
  fmt.Println(recumbentBike.Spares())

  // We can combine Parts and still call Spares, unlike the Ruby example. (page 173)
  comboParts := Parts{}
  comboParts = append(comboParts, mountainBike.Parts...)
  comboParts = append(comboParts, roadBike.Parts...)
  comboParts = append(comboParts, recumbentBike.Parts...)

  fmt.Println(len(comboParts))
  fmt.Println(comboParts.Spares())
}
