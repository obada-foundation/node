package obit

import (
	"github.com/obada-foundation/node/business/tests"
	"testing"
)

func TestObit(t *testing.T) {
	log, teardown := tests.NewUnit(t)

	log.Println("dd")

	t.Cleanup(teardown)
}
