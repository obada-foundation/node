package obit

import (
	"github.com/obada-foundation/node/business/tests"
	"testing"

	sdk "github.com/obada-foundation/sdk-go"
)

func TestObit(t *testing.T) {
	log, teardown := tests.NewUnit(t)
	sdk, err := sdk.NewSdk(log, true)

	if err != nil {
		t.Fatal("Cannot construct sdk")
	}

	service := NewObitService(sdk)

	t.Log("Given need to check obit service functionality")
	{
		service.Create()
	}

	t.Cleanup(teardown)
}
