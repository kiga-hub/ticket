package textgrid
import (
	"errors"
	"fmt"
)

type IntervalTier struct {
	Name 		string
	MaxTime		float64
	MinTime 	float64
	Intervals 	[]*Interval
}

func (T *IntervalTier) BisectLeft (interval *Interval) int {
	if T.Intervals == nil {
		return 0
	}
	for index, I := range T.Intervals {
		if I.Gte(interval) {
			return index
		}
	}
	return len(T.Intervals)
}

func (T *IntervalTier) Str () string {
	return fmt.Sprintf("<IntervalTier %s, %d intervals>", T.Name, len(T.Intervals))
}

func (T *IntervalTier) AddInterval (interval *Interval) error {
	if interval.MinTime < T.MinTime {
		return errors.New(fmt.Sprintf("%f", T.MinTime))
	}
	if T.MaxTime > 0 && interval.MaxTime > T.MaxTime {
		return errors.New(fmt.Sprintf("%f", T.MaxTime))
	}
	index := T.BisectLeft(interval)
	if index < len(T.Intervals) && T.Intervals[index].Eq(interval) {
		return errors.New(T.Intervals[index].Str())
	}
	rear := append([]*Interval{},T.Intervals[index:]...)
	T.Intervals = append(T.Intervals[0:index], interval)
	T.Intervals = append(T.Intervals, rear...)
	return nil
}

func (T *IntervalTier) Add (MinTime, MaxTime float64, Mark string) error {
	return T.AddInterval(&Interval{MinTime: MinTime, MaxTime: MaxTime, Mark: Mark})
}

func (T *IntervalTier) RemoveInterval (interval *Interval) {
	if T.Intervals == nil {
		return
	}
	for index, I := range T.Intervals {
		if I.Eq(interval) {
			T.Intervals = append(T.Intervals[:index], T.Intervals[index+1:]...)
		}
	}
}

func (T *IntervalTier) Remove (MinTime, MaxTime float64, Mark string) {
	T.RemoveInterval(&Interval{MinTime: MinTime, MaxTime: MaxTime, Mark: Mark})
}

