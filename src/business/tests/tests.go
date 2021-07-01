package tests

import (
	"context"
	"github.com/google/uuid"
	"github.com/obada-foundation/node/foundation/web"
	"log"
	"os"
	"testing"
	"time"
)

const (
	Success = "\u2713"
	Failed  = "\u2717"
)

var (
	dbImage = "mysql:8"
	dbPort  = "3306/tcp"
	dbArgs  = []string{"-e", "MYSQL_ROOT_PASSWORD=secret"}
)

func NewUnit(t *testing.T) (*log.Logger, func()) {
	c := startContainer(t, dbImage, dbPort, dbArgs...)

	t.Log("waiting for database to be ready ...")

	t.Log(c.Host)

	var pingError error
	maxAttempts := 20

	for attempts := 1; attempts <= maxAttempts; attempts++ {
		time.Sleep(time.Second * 5)
		pingError = nil
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
	}

	if pingError != nil {
		dumpContainerLogs(t, c.ID)
		stopContainer(t, c.ID)
		t.Fatalf("database is never ready: %v", pingError)
	}

	// Schema migrate

	teardown := func() {
		t.Helper()

		t.Logf("Stopping container")
		// database close
		stopContainer(t, c.ID)
	}

	log := log.New(os.Stdout, "TEST :", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	return log, teardown
}

func Context() context.Context {
	values := web.Values{
		TraceID: uuid.New().String(),
		Now:     time.Now(),
	}

	return context.WithValue(context.Background(), web.KeyValues, &values)
}

func StringPointer(s string) *string {
	return &s
}

func IntPointer(i int) *int {
	return &i
}
