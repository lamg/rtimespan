package rtimespan

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRSpan(t *testing.T) {
	// ts[0] is the date the RSpan starts
	// ts[1:] are the dates to be tested whether
	// they are inside the RSpan with an hour of
	// active period, a day of duration and repeated 10 times.

	// The boolean values in the following comments indicate
	// whether the associated date belongs the mentioned
	// repeated time span.
	ts := []string{
		// RSpan start date - 0
		"2006-01-02T15:04:05-04:00",
		// One minute later - 1 true
		"2006-01-02T15:05:05-04:00",
		// 59 minutes with 59 seconds later - 2 true
		"2006-01-02T16:04:04-04:00",
		// 1 hour later - 3 false
		"2006-01-02T16:04:05-04:00",
		// One day later - 4 true
		"2006-01-03T15:04:05-04:00",
		// One day later minus one minute - 5 false
		"2006-01-03T15:04:04-04:00",
		// 9 days later plus one minute - 6 true
		"2006-01-11T15:05:05-04:00",
		// 10 days later plus one minute - 7 false
		"2006-01-12T15:05:05-04:00",
		// 11 days later plus one minute - 8 false
		"2006-01-13T15:05:05-04:00",
	}
	tp := make([]time.Time, len(ts))
	for i, j := range ts {
		var e error
		tp[i], e = time.Parse(time.RFC3339, j)
		require.NoError(t, e)
	}
	rs := &RSpan{
		Active:   time.Hour,
		Infinite: false,
		Start:    tp[0],
		Times:    10,
		Total:    24 * time.Hour,
		AllTime:  false,
	}
	tss := []struct {
		rs *RSpan
		t  time.Time
		y  bool
	}{
		{rs, tp[1], true},
		{rs, tp[2], true},
		{rs, tp[3], false},
		{rs, tp[4], true},
		{rs, tp[5], false},
		{rs, tp[6], true},
		{rs, tp[7], false},
		{rs, tp[8], false},
	}
	for i, j := range tss {
		a := j.rs.ContainsTime(j.t)
		require.True(t, a == j.y, "At %d %t != %t", i, a, j.y)
	}
}

func TestJSONMarshal(t *testing.T) {
	d, e := time.Parse(time.RFC3339, "2006-01-02T15:04:05-04:00")
	require.NoError(t, e)
	r, s := &RSpan{
		Active:   time.Hour,
		Infinite: false,
		Start:    d,
		Times:    10,
		Total:    24 * time.Hour,
		AllTime:  false,
	},
		`{
		"active":3600000000000,
		"infinite":false,
		"start":"2006-01-02T15:04:05-04:00",
		"times":10,
		"total":86400000000000,
		"allTime":false
	}`
	x := new(RSpan)
	e = json.Unmarshal([]byte(s), x)
	require.NoError(t, e)
	require.True(t, r.Active == x.Active, "%d != %d",
		r.Active, x.Active)
	require.True(t, r.AllTime == x.AllTime)
	require.True(t, r.Infinite == x.Infinite)
	require.True(t, r.Start.Equal(x.Start))
	require.True(t, r.Times == x.Times)
	require.True(t, r.Total == x.Total)
}

func Example() {
	t0, e := time.Parse(time.RFC3339,
		"2006-01-02T15:04:05-04:00")
	var t1 time.Time
	if e == nil {
		// a day later inside the active time span
		// which is 2006-01-03T15:04:05-04:00 ->
		// 2006-01-03T16:04:05-04:00
		t1, e = time.Parse(time.RFC3339,
			"2006-01-03T15:46:05-04:00")
	}
	if e == nil {
		rs := &RSpan{
			Active:   time.Hour,
			Infinite: false,
			Start:    t0,
			Times:    10,
			Total:    24 * time.Hour,
		}
		// A repeated time span of on hour each day,
		// during 10 days, starting on t0
		y := rs.ContainsTime(t1)
		fmt.Println(y)
		// Output: true
	} else {
		println(e.Error())
	}
}
