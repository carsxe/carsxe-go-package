package carsxe

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func testClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return New("test-key", WithBaseURL(srv.URL), WithHTTPClient(srv.Client()))
}

func TestRecallsYmm(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/recalls-ymm" {
			t.Errorf("path = %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("key") != "test-key" || q.Get("year") != "2026" || q.Get("make") != "toyota" {
			t.Errorf("query = %s", r.URL.RawQuery)
		}
		json.NewEncoder(w).Encode(map[string]any{"success": true})
	})
	out := c.RecallsYmm(map[string]string{"year": "2026", "make": "toyota", "model": "corolla"})
	if out["success"] != true {
		t.Fatalf("got %#v", out)
	}
}

func TestSubmitRecallsBatch(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/recalls-batch/submit" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("content-type = %s", r.Header.Get("Content-Type"))
		}
		body, _ := io.ReadAll(r.Body)
		var payload map[string]any
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("body: %v", err)
		}
		vins, _ := payload["vins"].([]any)
		if len(vins) != 2 {
			t.Errorf("vins = %#v", payload["vins"])
		}
		json.NewEncoder(w).Encode(map[string]any{"success": true, "data": map[string]any{"batchId": "brb_test"}})
	})
	out := c.SubmitRecallsBatch(map[string]any{"vins": []string{"1HGBH41JXMN109186", "5YJSA1E26HF000001"}})
	data := out["data"].(map[string]any)
	if data["batchId"] != "brb_test" {
		t.Fatalf("got %#v", out)
	}
}

func TestRecallsBatchStatusResults(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s", r.Method)
		}
		if r.URL.Query().Get("batchId") != "brb_1" {
			t.Errorf("batchId = %s", r.URL.Query().Get("batchId"))
		}
		json.NewEncoder(w).Encode(map[string]any{"path": r.URL.Path})
	})
	if c.RecallsBatchStatus(map[string]string{"batchId": "brb_1"})["path"] != "/v1/recalls-batch/status" {
		t.Fatal("status path")
	}
	if c.RecallsBatchResults(map[string]string{"batchId": "brb_1"})["path"] != "/v1/recalls-batch/results" {
		t.Fatal("results path")
	}
}

func TestRecallsBatchDownloadCSV(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/recalls-batch/download" {
			t.Errorf("path = %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "text/csv")
		io.WriteString(w, "vin,hasRecalls\n1HGBH41JXMN109186,true\n")
	})
	out := c.RecallsBatchDownload(map[string]string{"batchId": "brb_1"})
	csv, _ := out["csv"].(string)
	if csv == "" || csv[:3] != "vin" {
		t.Fatalf("got %#v", out)
	}
}

func TestYmmOptionsAndOwnership(t *testing.T) {
	c := testClient(t, func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"path": r.URL.Path, "q": r.URL.RawQuery})
	})
	cases := []struct {
		got  map[string]any
		want string
	}{
		{c.YmmOptions(map[string]string{"year": "2026"}), "/v1/ymm-options"},
		{c.OwnershipVin(map[string]string{"vin": "1FT8X3BT0BEA61538"}), "/v1/ownership/vin"},
		{c.OwnershipPerson(map[string]string{"first_name": "John", "last_name": "Sample", "address": "123 Example St", "zip": "90210"}), "/v1/ownership/person"},
		{c.OwnershipAddress(map[string]string{"address": "123 Example St", "zip": "90210"}), "/v1/ownership/address"},
		{c.OwnershipZip(map[string]string{"zip": "90210"}), "/v1/ownership/zip"},
	}
	for _, tc := range cases {
		if tc.got["path"] != tc.want {
			t.Errorf("got path %v want %s", tc.got["path"], tc.want)
		}
	}
}
