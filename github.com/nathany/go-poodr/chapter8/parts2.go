// ############## Page 169-172 ##############
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
  NeedsSpare defaults to false here, the opposite of the Ruby example. We could
  write a NewSparePart initializer that defaults it to true.
*/
type Part struct {
  Name        string
  Description string
  NeedsSpare  bool
}

func main() {
  chain := Part{Name: "chain", Description: "10-speed", NeedsSpare: true}
  roadTire := Part{Name: "tire_size", Description: "23", NeedsSpare: true}
  tape := Part{Name: "tape_color", Description: "red", NeedsSpare: true}
  mountainTire := Part{Name: "tire_size", Description: "2.1", NeedsSpare: true}
  rearShock := Part{Name: "rear_shock", Description: "Fox", NeedsSpare: true}
  frontShock := Part{Name: "front_shock", Description: "Manitou"}

  roadBikeParts := Parts{chain, roadTire, tape}
  roadBike := Bicycle{Size: "L", Parts: roadBikeParts}

  mountainBikeParts := Parts{chain, mountainTire, frontShock, rearShock}
  mountainBike := Bicycle{Size: "L", Parts: mountainBikeParts}

  fmt.Println(roadBike.Size)
  fmt.Println(roadBike.Spares())

  fmt.Println(mountainBike.Size)
  fmt.Println(mountainBike.Spares())

  // We can combine Parts and still call Spares, unlike the Ruby example. (page 173)
  comboParts := Parts{}
  comboParts = append(comboParts, mountainBike.Parts...)
  comboParts = append(comboParts, roadBike.Parts...)
  fmt.Println(len(comboParts))
  fmt.Println(comboParts.Spares())
}
