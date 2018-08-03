package retry

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRetryTimes(t *testing.T) {
	action := New(func() error {
		return fmt.Errorf("test")
	}, WithTimes(3), WithBackoff(time.Millisecond*100, 4))

	defer action.Close()

	times := 0

	for {
		select {
		case status := <-action.Do():
			require.False(t, status)
			goto FINISH
		case <-action.Error():
			times++
		}
	}

FINISH:

	require.Equal(t, times, 3)
}
