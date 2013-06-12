// ############## Page 148-149 ##############
package main

import (
  "fmt"
  "time"
)

/*
  Schedule
*/
type Schedule struct{}
type Schedulable interface{}

func (schedule Schedule) IsScheduled(schedulable Schedulable, startDate, endDate time.Time) bool {
  fmt.Printf("This %T is not scheduled\n between %s and %s", schedulable, startDate, endDate)
  return false
}

/*
  Days
*/
type Days int

func (days Days) Duration() time.Duration {
  return time.Duration(days) * time.Hour * 24
}

/*
  Bicycle
*/
type Bicycle struct {
  Schedule              *Schedule
  Size, Chain, TireSize string
}

func (bike Bicycle) IsSchedulable(startDate, endDate time.Time) bool {
  leadDuration := bike.leadDays().Duration()
  start := startDate.Add(-leadDuration)
  return !bike.IsScheduled(start, endDate)
}

func (bike Bicycle) IsScheduled(startDate, endDate time.Time) bool {
  return bike.Schedule.IsScheduled(bike, startDate, endDate)
}

func (bike Bicycle) leadDays() Days {
  return 1
}

/*
  Main
*/
func main() {
  const dateLayout = "2006/1/2"
  starting, _ := time.Parse(dateLayout, "2015/09/04")
  ending, _ := time.Parse(dateLayout, "2015/09/10")

  b := Bicycle{Schedule: &Schedule{}}
  b.IsSchedulable(starting, ending)
}
