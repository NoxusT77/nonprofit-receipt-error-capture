package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendReceiptsQueuesFailedDonor(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":true,"data":{}}`))
	}))
	defer server.Close()
	client := &ErrorClient{BaseURL: server.URL, Key: "", HTTP: server.Client()}
	donations := []Donation{{DonorID: "ok", AmountCents: 100}, {DonorID: "bad", AmountCents: 200}}
	results := SendReceipts(client, donations, func(d Donation) error {
		if d.DonorID == "bad" {
			return errReceipt{}
		}
		return nil
	})
	if results[0].Status != "sent" || results[1].Status != "queued_for_review" {
		t.Fatalf("unexpected results: %#v", results)
	}
}

func TestSendReceiptsReportsCaptureFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok":false,"error":{"code":"unavailable"}}`))
	}))
	defer server.Close()
	client := &ErrorClient{BaseURL: server.URL, Key: "", HTTP: server.Client()}

	results := SendReceipts(client, []Donation{{DonorID: "bad", AmountCents: 200}}, func(Donation) error {
		return errReceipt{}
	})
	if results[0].Status != "capture_failed" {
		t.Fatalf("unexpected results: %#v", results)
	}
}

type errReceipt struct{}

func (errReceipt) Error() string { return "mail service rejected" }
