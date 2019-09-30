package task5

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func newJobFunc(dur time.Duration, err bool) func() error {
	return func() error {
		if dur > 0 {
			time.Sleep(dur)
		}
		if err {
			return errors.New("job_error")
		}
		return nil
	}
}

func TestFun1(t *testing.T) {
	cases := []struct {
		jobs       []func() error
		workersLen int
		errLimit   int
		successCnt int           // expected
		errCnt     int           // expected
		duration   time.Duration // expected (in Millisecond)
	}{
		{
			jobs: []func() error{
				newJobFunc(50*time.Millisecond, false),
				newJobFunc(50*time.Millisecond, false),
				newJobFunc(50*time.Millisecond, false),
			},
			workersLen: 1,
			errLimit:   2,
			successCnt: 3,
			errCnt:     0,
			duration:   150 * time.Millisecond,
		},
		{
			jobs: []func() error{
				newJobFunc(0, true),
				newJobFunc(0, true),
				newJobFunc(50*time.Millisecond, false),
			},
			workersLen: 1,
			errLimit:   2,
			successCnt: 0,
			errCnt:     2,
			duration:   0,
		},
		{
			jobs: []func() error{
				newJobFunc(50*time.Millisecond, false),
				newJobFunc(50*time.Millisecond, false),
				newJobFunc(50*time.Millisecond, false),
				newJobFunc(50*time.Millisecond, false),
			},
			workersLen: 2,
			errLimit:   0,
			successCnt: 4,
			errCnt:     0,
			duration:   100 * time.Millisecond,
		},
		{
			jobs: []func() error{
				newJobFunc(0, true),
				newJobFunc(0, true),
				newJobFunc(50*time.Millisecond, false),
				newJobFunc(50*time.Millisecond, false),
			},
			workersLen: 2,
			errLimit:   3,
			successCnt: 2,
			errCnt:     2,
			duration:   50 * time.Millisecond,
		},
		{
			jobs: []func() error{
				newJobFunc(20*time.Millisecond, true),
				newJobFunc(30*time.Millisecond, true),
				newJobFunc(50*time.Millisecond, false),
				newJobFunc(50*time.Millisecond, false),
			},
			workersLen: 2,
			errLimit:   2,
			successCnt: 1,
			errCnt:     2,
			duration:   70 * time.Millisecond,
		},
		{
			jobs: []func() error{
				newJobFunc(20*time.Millisecond, true),
				newJobFunc(50*time.Millisecond, true),
				newJobFunc(0, true),
				newJobFunc(50*time.Millisecond, false),
			},
			workersLen: 2,
			errLimit:   2,
			successCnt: 0,
			errCnt:     3,
			duration:   50 * time.Millisecond,
		},
	}

	for csI, cs := range cases {
		startTime := time.Now()
		successCnt, errCnt := Fun1(cs.jobs, cs.workersLen, cs.errLimit)
		elapsedDuration := time.Since(startTime)
		require.Equal(t, cs.successCnt, successCnt, "successCnt, case %d", csI+1)
		require.Equal(t, cs.errCnt, errCnt, "errCnt, case %d", csI+1)
		durDiff := elapsedDuration - cs.duration
		if durDiff < 0 { // just for abs
			durDiff = -durDiff
		}
		require.Less(t, int64(durDiff), int64(10*time.Millisecond), "durDiff, case %d", csI+1)
	}
}
