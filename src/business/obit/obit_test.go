package obit

import (
	"github.com/obada-foundation/node/business/tests"
	"testing"

	"github.com/obada-foundation/sdkgo"
)

func TestObit(t *testing.T) {
	logger, teardown := tests.NewUnit(t)
	sdk, err := sdkgo.NewSdk(logger, true)

	if err != nil {
		t.Fatal("Cannot construct sdk")
	}

	service := NewObitService(sdk)

	t.Log("Given need to check obit service functionality")
	service.Create()

	t.Cleanup(teardown)
}
