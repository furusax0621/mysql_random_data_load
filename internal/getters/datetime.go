package getters

import (
	"fmt"
	"math/rand"
	"time"
)

type RandomDateTimeInRange struct {
	name      string
	min       string
	max       string
	allowNull bool
}

// Value returns a random time.Time in the range specified by the New method
func (r *RandomDateTimeInRange) Value() interface{} {
	randomSeconds := rand.Int63n(oneYear)
	d := time.Now().Add(-1 * time.Duration(randomSeconds) * time.Second)
	return d
}

func (r *RandomDateTimeInRange) String() string {
	d := r.Value().(time.Time)
	return d.Format("2006-01-02 15:03:04")
}

// Quote returns the value quoted for MySQL
func (r *RandomDateTimeInRange) Quote() string {
	d := r.Value().(time.Time)
	return fmt.Sprintf("'%s'", d.Format("2006-01-02 15:03:04"))
}

// NewRandomDateTimeInRange returns a new random date in the specified range
func NewRandomDateTimeInRange(name string, min, max string, allowNull bool) *RandomDateTimeInRange {
	if min == "" {
		t := time.Now().Add(-1 * time.Duration(oneYear) * time.Second)
		min = t.Format("2006-01-02 15:03:04")
	}
	return &RandomDateTimeInRange{name, min, max, allowNull}
}

// NewRandomDateTime returns a new random datetime between Now() and Now() - 1 year
func NewRandomDateTime(name string, allowNull bool) *RandomDateTimeInRange {
	return &RandomDateTimeInRange{name, "", "", allowNull}
}
