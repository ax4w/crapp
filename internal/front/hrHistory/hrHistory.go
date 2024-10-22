package hrHitory

import (
	"crapp/internal/middle"
	"time"
)

type History struct {
	Heart_rates []int
	T           time.Time
	Size        int
	Index       int
	End         bool
	R           int
}

func New(e bool, s, r int) *History {
	return &History{
		End:  e,
		Size: s,
		R:    r,
	}
}
func (h *History) IsToday() bool {
	t1 := time.Now()
	return t1.Year() == h.T.Year() && h.T.Month() == h.T.Month() && t1.Day() == h.T.Day()
}

func (h History) AddTimes() {
	var objs []middle.HRObj
	m := time.Date(h.T.Year(), h.T.Month(), h.T.Day(), 0, 0, 0, 0, h.T.Location())
	for _, v := range h.Heart_rates {
		if v != 0 {
			objs = append(objs, middle.HRObj{
				T:  m,
				HR: v,
			})
		}
		m = m.Add(5 * time.Minute)
	}
	middle.HeartRateLogChannel <- objs
}

func (h *History) minSinceMidnight() int {
	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	delta := now.Sub(midnight).Seconds()
	return int(delta/60) + 1
}

func (h *History) Normalize() {
	if len(h.Heart_rates) > 288 {
		h.Heart_rates = h.Heart_rates[0:288]
	}
	if h.IsToday() {
		h.Heart_rates = h.Heart_rates[:h.minSinceMidnight()]
	}

}
