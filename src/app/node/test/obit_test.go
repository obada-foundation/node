package handlers

import (
	"encoding/json"
	"github.com/obada-foundation/node/app/node/handlers"
	obitService "github.com/obada-foundation/node/business/obit"
	searchService "github.com/obada-foundation/node/business/search"
	"github.com/obada-foundation/node/business/search"
	"github.com/obada-foundation/node/business/sys/validate"
	"github.com/obada-foundation/node/business/tests"
	"github.com/obada-foundation/node/business/types"
	"github.com/obada-foundation/sdkgo"
	"net/http"
	"net/http/httptest"
	"net/url"
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

	tests.CreateObit(t, test)
	tests.CreateOwnerObits(t, test)

	shutdown := make(chan os.Signal, 1)

	sdk, err := sdkgo.NewSdk(test.Logger, false)

	if err != nil {
		t.Fatalf("cannot initialize SDK: %s", err)
	}

	os := obitService.NewObitService(sdk, test.Logger, test.DB, nil, nil)

	ss := searchService.NewService(test.Logger, test.DB)

	tests := ObitTests{
		app: handlers.API(handlers.APIConfig{
			Shutdown:    shutdown,
			Logger:      test.Logger,
			ObitService: os,
			SearchService: ss,
		}),
	}

	t.Run("generateID200", tests.generateID200)
	t.Run("generateID200", tests.generateID422)
	t.Run("checksum200", tests.checksum200)
	t.Run("checksum422", tests.checksum422)
	t.Run("checksum500", tests.checksum500)
	t.Run("get404", tests.get404)
	t.Run("get200", tests.get200)
	t.Run("search200", tests.search200)
}

func (obs ObitTests) generateID200(t *testing.T) {
	body := `{"serial_number": "SN123456X", "manufacturer": "SONY", "part_number": "PN123456S"}`

	r := httptest.NewRequest(http.MethodPost, "/obit/id", strings.NewReader(body))
	w := httptest.NewRecorder()

	obs.app.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

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

func (obs ObitTests) generateID422(t *testing.T) {
	testCases := []struct {
		arg  string
		want string
	}{
		{
			arg:  `{}`,
			want: `[{"field":"serial_number","error":"serial_number is a required field"},{"field":"manufacturer","error":"manufacturer is a required field"},{"field":"part_number","error":"part_number is a required field"}]`,
		},
		{
			arg:  `{"serial_number": "", "manufacturer": "", "part_number": ""}`,
			want: `[{"field":"serial_number","error":"serial_number is a required field"},{"field":"manufacturer","error":"manufacturer is a required field"},{"field":"part_number","error":"part_number is a required field"}]`,
		},
		{
			arg:  `{"serial_number": "SN123456X", "manufacturer": "", "part_number": ""}`,
			want: `[{"field":"manufacturer","error":"manufacturer is a required field"},{"field":"part_number","error":"part_number is a required field"}]`,
		},
		{
			arg:  `{"serial_number": "", "manufacturer": "SONY", "part_number": ""}`,
			want: `[{"field":"serial_number","error":"serial_number is a required field"},{"field":"part_number","error":"part_number is a required field"}]`,
		},
		{
			arg:  `{"serial_number": "", "manufacturer": "", "part_number": "PN123456S"}`,
			want: `[{"field":"serial_number","error":"serial_number is a required field"},{"field":"manufacturer","error":"manufacturer is a required field"}]`,
		},
	}

	for _, tc := range testCases {
		r := httptest.NewRequest(http.MethodPost, "/obit/id", strings.NewReader(tc.arg))
		w := httptest.NewRecorder()

		obs.app.ServeHTTP(w, r)

		resp := w.Result()
		defer resp.Body.Close()

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

func (obs ObitTests) checksum200(t *testing.T) {
	body := `
		{
			"serial_number_hash": "SN123456X", 
			"manufacturer": "SONY", 
			"part_number": "PN123456S",
			"owner_did": "did:obada:owner:123456",
			"modified_on": 1624387537
		}
	`

	r := httptest.NewRequest(http.MethodPost, "/obit/checksum", strings.NewReader(body))
	w := httptest.NewRecorder()

	obs.app.ServeHTTP(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Handler() status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func (obs ObitTests) checksum500(t *testing.T) {
	testCases := []struct {
		arg  string
		want string
	}{
		{
			arg:  "",
			want: "Internal Server Error",
		},
	}

	for _, tc := range testCases {
		r := httptest.NewRequest(http.MethodPost, "/obit/checksum", strings.NewReader(tc.arg))
		w := httptest.NewRecorder()

		obs.app.ServeHTTP(w, r)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusInternalServerError {
			t.Fatalf("Handler() status = %d, want %d", resp.StatusCode, http.StatusInternalServerError)
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

		if er.Error != tc.want {
			t.Errorf("Handler() Body.Error = %q; want %q", er.Error, tc.want)
		}
	}
}

func (obs ObitTests) checksum422(t *testing.T) {
	testCases := []struct {
		arg  string
		want string
	}{
		{
			arg:  "{}",
			want: `[{"field":"serial_number_hash","error":"serial_number_hash is a required field"},{"field":"manufacturer","error":"manufacturer is a required field"},{"field":"part_number","error":"part_number is a required field"},{"field":"owner_did","error":"owner_did is a required field"}]`,
		},
		{
			arg: `{
				"serial_number_hash": "cae6b797ae2627d96689fed03adc28311d5f2175253c3a0e375301e225ddf44d"
			}`,
			want: `[{"field":"manufacturer","error":"manufacturer is a required field"},{"field":"part_number","error":"part_number is a required field"},{"field":"owner_did","error":"owner_did is a required field"}]`,
		},
	}

	for _, tc := range testCases {
		r := httptest.NewRequest(http.MethodPost, "/obit/checksum", strings.NewReader(tc.arg))
		w := httptest.NewRecorder()

		obs.app.ServeHTTP(w, r)

		resp := w.Result()
		defer resp.Body.Close()

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

func (obs ObitTests) get404(t *testing.T) {
	testCases := []struct {
		arg  string
		want string
	}{
		{
			arg:  "d7cf869423d12f623f5611e48d6f6665bbc4a270b6e09da2f4c32bcb1b949ecx",
			want: "",
		},
	}

	for _, tc := range testCases {
		r := httptest.NewRequest(http.MethodGet, "/obits/"+tc.arg, nil)
		w := httptest.NewRecorder()

		obs.app.ServeHTTP(w, r)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("Handler() status = %d, want %d", resp.StatusCode, http.StatusNotFound)
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

		wantErrMsg := "not found"
		if er.Error != wantErrMsg {
			t.Errorf("Handler() Body.Error = %q; want %q", er.Error, wantErrMsg)
		}

		if er.Fields != tc.want {
			t.Errorf("Handler() Body.Fields = %q; want %q", er.Fields, tc.want)
		}
	}
}

func (obs ObitTests) get200(t *testing.T) {
	testCases := []struct {
		arg  string
		want string
	}{
		{
			arg: "d7cf869423d12f623f5611e48d6f6665bbc4a270b6e09da2f4c32bcb1b949ecd",
		},
	}

	for _, tc := range testCases {
		r := httptest.NewRequest(http.MethodGet, "/obits/"+tc.arg, nil)
		w := httptest.NewRecorder()

		obs.app.ServeHTTP(w, r)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Handler() status = %d, want %d", resp.StatusCode, http.StatusOK)
		}

		contentType := resp.Header.Get("Content-Type")
		wantContentType := "application/json"
		if contentType != wantContentType {
			t.Errorf("Handler() Content-Type = %q; want %q", contentType, wantContentType)
		}

		var o types.QLDBObit
		if err := json.NewDecoder(resp.Body).Decode(&o); err != nil {
			t.Fatalf("json.NewDecoder(%v+).Decode(%v+) err = %s", resp.Body, &o, err)
		}

		wantDID := tc.arg
		if o.ObitDID != wantDID {
			t.Errorf("Handler() Body.ObitID = %q; want %q", o.ObitDID, wantDID)
		}
	}
}

func (obs ObitTests) search200(t *testing.T) {
	testCases := []struct {
		arg  string
		want int
	}{
		{
			arg: "usn1",
			want: 1,
		},
		{
			arg: "did:obada:owner:678910",
			want: 50,
		},
	}

	for _, tc := range testCases {
		r := httptest.NewRequest(http.MethodGet, "/obits?q=" + url.QueryEscape(tc.arg), nil)
		w := httptest.NewRecorder()

		obs.app.ServeHTTP(w, r)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Handler() status = %d, want %d", resp.StatusCode, http.StatusOK)
		}

		contentType := resp.Header.Get("Content-Type")
		wantContentType := "application/json"
		if contentType != wantContentType {
			t.Errorf("Handler() Content-Type = %q; want %q", contentType, wantContentType)
		}

		var obits search.Obits
		if err := json.NewDecoder(resp.Body).Decode(&obits); err != nil {
			t.Fatalf("json.NewDecoder(%v+).Decode(%v+) err = %s", resp.Body, &obits, err)
		}

		t.Logf("%v", obits)

		wantObitsCount := len(obits.Obits)
		if wantObitsCount != tc.want {
			t.Errorf("Handler(%q, %d) count(Body.Obits) = %d; want %d", tc.arg, 0, wantObitsCount, tc.want)
		}
	}
}
