// Package rtimespan determines whether a time is
// contained in a time span that repeates several
// times (possibly infinite).
//
// As an example you may consider a time span that
// starts at determined time for an hour, and you
// also want to specify an homologal time spans that
// start a day later, for 10 days (includes current day).
// For representing that situtation you would create an
// RSpan with start time determined by you, active duration
// of one hour, total duration of 24 hours and Times
// field equal 10. Notice that the 10th day counting
// from the start date is excluded.
package rtimespan

import (
	"time"
)

// RSpan is a repeated time span
type RSpan struct {
	// Time the repeated span starts
	Start time.Time `json:"start"`
	// Duration of the active phase
	Active time.Duration `json:"active"`
	// Duration of the cycle
	Total time.Duration `json:"total"`
	// Number of times the span repeats
	Times int `json:"times"`
	// Repeats forever
	Infinite bool `json:"infinite"`
	// AllTime indicates whether all times are
	// contained in this span, i.e. for all times
	// RSpan.ContainsTime returns true
	AllTime bool `json:"allTime"`
}

// ContainsTime returns whether r contains x
func (r *RSpan) ContainsTime(x time.Time) (y bool) {
	p := x.Sub(r.Start)
	d, m := p/r.Total, p%r.Total
	y = r.AllTime ||
		(((r.Times >= 0 && d < time.Duration(r.Times)) ||
			r.Infinite) && m < r.Active)
	return
}
