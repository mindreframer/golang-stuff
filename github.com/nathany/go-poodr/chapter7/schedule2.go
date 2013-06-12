// ############## Page 150-151 ##############
package main

import (
  "fmt"
  "time"
)

/*
  Schedule
*/
type Schedule struct{}

func (schedule Schedule) IsScheduled(schedulable Schedulable, startDate, endDate time.Time) bool {
  fmt.Printf("This %T is not scheduled\n between %s and %s", schedulable, startDate, endDate)
  return false
}

func NewSchedule() *Schedule {
  return &Schedule{}
}

/*
  Schedulable
*/
type Schedulable struct {
  LeadDays Days
  schedule *Schedule
}

func (schedulable Schedulable) Schedule() *Schedule {
  if schedulable.schedule == nil {
    schedulable.schedule = NewSchedule()
  }
  return schedulable.schedule
}

func (schedulable Schedulable) IsSchedulable(startDate, endDate time.Time) bool {
  leadDuration := schedulable.LeadDays.Duration()
  start := startDate.Add(-leadDuration)
  return !schedulable.IsScheduled(start, endDate)
}

func (schedulable Schedulable) IsScheduled(startDate, endDate time.Time) bool {
  // Right now schedule refers to main.Schedulable rather than Bicycle.
  return schedulable.Schedule().IsScheduled(schedulable, startDate, endDate)
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
  Without inheritance we would need to pass Bicycle to Schedulable to have a
  LeadDays() hook, perhaps via a SchedulableHooks interface. Far simpler to
  just initalize Schedulable with LeadDays here.
*/
type Bicycle struct {
  Schedulable
  Size, Chain, TireSize string
}

func NewBicycle() *Bicycle {
  return &Bicycle{Schedulable: Schedulable{LeadDays: 1}}
}

/*
  Main
*/
func main() {
  const dateLayout = "2006/1/2"
  starting, _ := time.Parse(dateLayout, "2015/09/04")
  ending, _ := time.Parse(dateLayout, "2015/09/10")

  b := NewBicycle()
  b.IsSchedulable(starting, ending)
}
