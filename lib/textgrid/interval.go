package textgrid
import (
	"fmt"
)

type Interval struct {
	MinTime float64
	MaxTime float64
	Mark string
}
func (I *Interval) Str () string {
	return fmt.Sprintf("Interval (%f, %f, %s)", I.MinTime, I.MaxTime, I.Mark)
}

func (I *Interval) Duration () float64 {
	return I.MaxTime - I.MinTime
}

func (I *Interval) Overlaps (other *Interval) bool {
	return other.MinTime < I.MaxTime && I.MinTime < other.MaxTime
}

func (I *Interval) Lt (other *Interval) bool {
	if I.Overlaps(other) {
		return false
	}
	return I.MinTime < other.MinTime
} 

func (I *Interval) Gt (other *Interval) bool {
	if I.Overlaps(other) {
		return false
	}
	return I.MaxTime > other.MaxTime
} 

func (I *Interval) Eq (other *Interval) bool {
	return I.MinTime == other.MinTime && I.MaxTime == other.MaxTime
} 

func (I *Interval) Lte (other *Interval) bool {
	return I.Lt(other) || I.Eq(other)
} 

func (I *Interval) Gte (other *Interval) bool {
	return I.Gt(other) || I.Eq(other)
}

func (I *Interval) Add (sec float64) {
	I.MaxTime += sec
	I.MinTime += sec
}

func (I *Interval) Sub (sec float64) {
	I.MaxTime -= sec
	I.MinTime -= sec
}

func (I *Interval) Contains (other *Interval) bool {
	return I.MinTime <= other.MinTime && other.MaxTime <= I.MaxTime
}

