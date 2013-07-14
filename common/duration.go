package common

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	ds := formatDuration(time.Duration(d))
	return []byte(fmt.Sprintf(`"%s"`, ds)), nil
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), "\"")
	u := strings.Split(str, ":")

	ds := u[0]
	di, err := strconv.ParseInt(ds[:len(ds)-1], 10, 64)
	if err != nil {
		return err
	}

	hs := u[1]
	hi, err := strconv.ParseInt(hs[:len(hs)-1], 10, 64)
	if err != nil {
		return err
	}

	hi += di * 24

	u[1] = fmt.Sprintf("%dh", hi)

	dur, err := time.ParseDuration(strings.Join(u[1:], ""))
	if err != nil {
		return err
	}

	*d = Duration(dur)
	return nil
}

func formatDuration(d time.Duration) string {
	t := int64(d.Seconds())
	day := t / (60 * 60 * 24)
	t = t % (60 * 60 * 24)
	hour := t / (60 * 60)
	t = t % (60 * 60)
	min := t / 60
	sec := t % 60

	ds := fmt.Sprintf("%dd:%dh:%dm:%ds", day, hour, min, sec)
	return ds
}

type Time time.Time

func (t Time) MarshalJSON() ([]byte, error) {
	f := "2006-01-02 15:04:05 -0700"
	s := time.Time(t).Format(f)
	return []byte(fmt.Sprintf(`"%s"`, s)), nil
}

func (t *Time) UnmarshalJSON(b []byte) error {
	s := string(b)
	f := `"2006-01-02 15:04:05 -0700"`

	tt, err := time.Parse(f, s)
	if err != nil {
		return err
	}

	*t = Time(tt)
	return nil
}

func (t Time) Elapsed() Duration {
	d := time.Since(time.Time(t))
	return Duration(d)
}
