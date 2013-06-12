package skyd

import (
	"testing"
	"time"
)

func TestShiftEpoch(t *testing.T) {
	timestamp, _ := time.Parse(time.RFC3339, "1970-01-01T00:00:00Z")
	value := ShiftTime(timestamp)
	if value != 0 {
		t.Fatalf("Invalid time shift: %v", value)
	}
}

func TestShiftOneSecondAfterEpoch(t *testing.T) {
	timestamp, _ := time.Parse(time.RFC3339, "1970-01-01T00:00:01Z")
	value := ShiftTime(timestamp)
	if value != 0x100000 {
		t.Fatalf("Invalid time shift: %v", value)
	}
}

func TestShiftOneSecondBeforeEpoch(t *testing.T) {
	timestamp, _ := time.Parse(time.RFC3339, "1969-12-31T23:59:59Z")
	value := ShiftTime(timestamp)
	if value != -0x100000 {
		t.Fatalf("Invalid time shift: %v", value)
	}
}

func TestShiftOneAndAHalfSecondsAfterEpoch(t *testing.T) {
	timestamp, _ := time.Parse(time.RFC3339, "1970-01-01T00:00:01.5Z")
	value := ShiftTime(timestamp)
	if value != 0x17a120 {
		t.Fatalf("Invalid time shift: %v", value)
	}
}

func TestUnshiftEpoch(t *testing.T) {
	value := UnshiftTime(0)
	if value.UTC().Format(time.RFC3339) != "1970-01-01T00:00:00Z" {
		t.Fatalf("Invalid time unshift: %v", value)
	}
}

func TestUnshiftOneSecondAfterEpoch(t *testing.T) {
	value := UnshiftTime(0x100000)
	if value.UTC().Format(time.RFC3339) != "1970-01-01T00:00:01Z" {
		t.Fatalf("Invalid time unshift: %v", value)
	}
}

func TestUnshiftOneSecondBeforeEpoch(t *testing.T) {
	value := UnshiftTime(-0x100000)
	if value.UTC().Format(time.RFC3339) != "1969-12-31T23:59:59Z" {
		t.Fatalf("Invalid time unshift: %v", value)
	}
}

func TestUnshiftOneAndAHalfSecondsAfterEpoch(t *testing.T) {
	value := UnshiftTime(0x17a120)
	if value.UTC().Format(time.RFC3339Nano) != "1970-01-01T00:00:01.5Z" {
		t.Fatalf("Invalid time unshift: %v", value)
	}
}
