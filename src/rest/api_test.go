package rest

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

func TestRest_Shutdown(t *testing.T) {
	logger := log.New(os.Stdout, "TEST-REST ", 0)
	srv := Rest{Logger: logger}
	done := make(chan bool)

	// without waiting for channel close at the end goroutine will stay alive after test finish
	// which would create data race with next test
	go func() {
		time.Sleep(200 * time.Millisecond)
		srv.Shutdown()
		close(done)
	}()

	st := time.Now()
	srv.Run("127.0.0.1", 0)
	assert.True(t, time.Since(st).Seconds() < 1, "should take about 100ms")
	<-done
}
