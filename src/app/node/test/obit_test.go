package handlers

import (
	"encoding/json"
	"github.com/obada-foundation/node/app/node/handlers"
	obitService "github.com/obada-foundation/node/business/obit"
	"github.com/obada-foundation/node/business/sys/validate"
	"github.com/obada-foundation/node/business/tests"
	"github.com/obada-foundation/sdkgo"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// ObitTests holds methods for each obit subtest. This type allows
// passing dependencies for tests while still providing a convenient syntax
// when subtests are registered.
type ObitTests struct {
	app http.Handler
}

func TestAPI(t *testing.T) {
	test := tests.NewIntegration(t)
	t.Cleanup(test.Teardown)

	shutdown := make(chan os.Signal, 1)

	sdk, err := sdkgo.NewSdk(test.Logger, false)

	if err != nil {
		t.Fatalf("cannot initialize SDK: %s", err)
	}

	os := obitService.NewObitService(sdk, test.Logger, test.DB, nil, nil)

	tests := ObitTests{
		app: handlers.API(handlers.APIConfig{
			shutdown,
			test.Logger,
			os,
			nil,
		}),
	}

	t.Run("generateID200", tests.generateID200)
	t.Run("generateID200", tests.generateID422)
}

func (os ObitTests) generateID200(t *testing.T) {
	body := `{"serial_number": "SN123456X", "manufacturer": "SONY", "part_number": "PN123456S"}`

	r := httptest.NewRequest(http.MethodPost, "/obit/id", strings.NewReader(body))
	w := httptest.NewRecorder()

	os.app.ServeHTTP(w, r)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Handler() status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	contentType := resp.Header.Get("Content-Type")
	wantContentType := "application/json"

	if contentType != wantContentType {
		t.Errorf("Handler() Content-Type = %q; want %q", contentType, wantContentType)
	}

	var ID obitService.ID
	if err := json.NewDecoder(resp.Body).Decode(&ID); err != nil {
		t.Fatalf("json.NewDecoder(%v+).Decode(%v+) err = %s", resp.Body, &ID, err)
	}

	wantID := "d7cf869423d12f623f5611e48d6f6665bbc4a270b6e09da2f4c32bcb1b949ecd"
	if ID.ID != wantID {
		t.Errorf("Handler() Body.ID = %q; want %q", ID.ID, wantID)
	}

	wantDID := "did:obada:d7cf869423d12f623f5611e48d6f6665bbc4a270b6e09da2f4c32bcb1b949ecd"
	if ID.DID != wantDID {
		t.Errorf("Handler() Body.DID = %q; want %q", ID.DID, wantDID)
	}

}

func (os ObitTests) generateID422(t *testing.T) {
	testCases := []struct {
		arg string
		want string
	}{
		{
			arg: `{}`,
			want: `[{"field":"serial_number","error":"serial_number is a required field"},{"field":"manufacturer","error":"manufacturer is a required field"},{"field":"part_number","error":"part_number is a required field"}]`,
		},
		{
			arg: `{"serial_number": "", "manufacturer": "", "part_number": ""}`,
			want: `[{"field":"serial_number","error":"serial_number is a required field"},{"field":"manufacturer","error":"manufacturer is a required field"},{"field":"part_number","error":"part_number is a required field"}]`,
		},
		{
			arg: `{"serial_number": "SN123456X", "manufacturer": "", "part_number": ""}`,
			want: `[{"field":"manufacturer","error":"manufacturer is a required field"},{"field":"part_number","error":"part_number is a required field"}]`,
		},
		{
			arg: `{"serial_number": "", "manufacturer": "SONY", "part_number": ""}`,
			want: `[{"field":"serial_number","error":"serial_number is a required field"},{"field":"part_number","error":"part_number is a required field"}]`,
		},
		{
			arg: `{"serial_number": "", "manufacturer": "", "part_number": "PN123456S"}`,
			want: `[{"field":"serial_number","error":"serial_number is a required field"},{"field":"manufacturer","error":"manufacturer is a required field"}]`,
		},
	}

	for _, tc := range testCases {
		r := httptest.NewRequest(http.MethodPost, "/obit/id", strings.NewReader(tc.arg))
		w := httptest.NewRecorder()

		os.app.ServeHTTP(w, r)

		resp := w.Result()
		if resp.StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("Handler() status = %d, want %d", resp.StatusCode, http.StatusUnprocessableEntity)
		}

		contentType := resp.Header.Get("Content-Type")
		wantContentType := "application/json"
		if contentType != wantContentType {
			t.Errorf("Handler() Content-Type = %q; want %q", contentType, wantContentType)
		}

		var er validate.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&er); err != nil {
			t.Fatalf("json.NewDecoder(%v+).Decode(%v+) err = %s", resp.Body, &er, err)
		}

		wantErrMsg := "data validation error"
		if er.Error != wantErrMsg {
			t.Errorf("Handler() Body.Error = %q; want %q", er.Error, wantErrMsg)
		}

		if er.Fields != tc.want {
			t.Errorf("Handler() Body.Fields = %q; want %q", er.Fields, tc.want)
		}
	}
}
