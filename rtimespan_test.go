package rtimespan

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRSpan(t *testing.T) {
	// t0 is the date the RSpan starts
	t0, e := time.Parse(time.RFC3339,
		"2006-01-02T15:04:05-04:00")
	require.NoError(t, e)
	rs := &RSpan{
		Active:   time.Hour,
		Infinite: false,
		Start:    t0,
		Times:    10,
		Total:    24 * time.Hour,
		AllTime:  false,
	}
	// ts are the dates to be tested whether
	// they are inside the RSpan with an hour of
	// active period, a day of duration and repeated 10 times.

	// The boolean values in the following comments indicate
	// whether the associated date belongs the mentioned
	// repeated time span.
	ts := [][3]string{
		// One minute later - true
		{"2006-01-02T15:05:05-04:00",
			"2006-01-02T15:04:05-04:00",
			"2006-01-02T16:04:05-04:00"},
		// 59 minutes with 59 seconds later - true
		{"2006-01-02T16:04:04-04:00",
			"2006-01-02T15:04:05-04:00",
			"2006-01-02T16:04:05-04:00"},
		// One day later - true
		{"2006-01-03T15:04:05-04:00",
			"2006-01-03T15:04:05-04:00",
			"2006-01-03T16:04:05-04:00"},
		// 9 days later plus one minute - true
		{"2006-01-11T15:05:05-04:00",
			"2006-01-11T15:04:05-04:00",
			"2006-01-11T16:04:05-04:00"},
		// One day later minus one minute - false
		{"2006-01-03T15:04:04-04:00",
			"2006-01-02T15:04:05-04:00",
			"2006-01-02T16:04:05-04:00"},
		// 1 hour later - false
		{"2006-01-02T16:04:05-04:00",
			"2006-01-02T15:04:05-04:00",
			"2006-01-02T16:04:05-04:00"},
		// 10 days later plus one minute - false
		{"2006-01-12T15:05:05-04:00",
			"2006-01-11T15:04:05-04:00",
			"2006-01-11T16:04:05-04:00"},
		// 11 days later plus one minute - false
		{"2006-01-13T15:05:05-04:00",
			"2006-01-11T15:04:05-04:00",
			"2006-01-11T16:04:05-04:00"},
	}
	tp := make([][3]time.Time, len(ts))
	for i, j := range ts {
		for k := 0; k != 3; k++ {
			var e error
			tp[i][k], e = time.Parse(time.RFC3339, j[k])
			require.NoError(t, e, "At %d,%d", i, k)
		}
	}
	tss := make([]*rTst, len(tp))
	for i := 0; i != len(tss); i++ {
		tss[i] = &rTst{
			rs: rs,
			t:  tp[i][0],
			a:  tp[i][1],
			b:  tp[i][2],
			y:  i < 4,
		}
	}
	for i, j := range tss {
		y := j.rs.ContainsTime(j.t)
		require.True(t, y == j.y, "At %d %t != %t", i, y, j.y)
		a, b := j.rs.CurrActIntv(j.t)
		require.Equal(t, j.a, a, "At %d", i)
		require.Equal(t, j.b, b, "At %d", i)
	}
}

type rTst struct {
	rs   *RSpan
	t    time.Time
	a, b time.Time
	y    bool
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

func TestRSpan0(t *testing.T) {
	ts := [2]string{"2018-01-22T00:01:00Z",
		"2018-01-21T01:50:00Z"}
	tp := make([]time.Time, len(ts))
	for i, j := range ts {
		var e error
		tp[i], e = time.Parse(time.RFC3339, j)
		require.NoError(t, e)
	}
	s := &RSpan{
		Start:    tp[0],
		Active:   time.Hour,
		Total:    48 * time.Hour,
		Times:    1,
		AllTime:  false,
		Infinite: false,
	}
	require.False(t, s.ContainsTime(tp[1]))
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
