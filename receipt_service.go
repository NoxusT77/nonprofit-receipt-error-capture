package main

import "fmt"

type Donation struct {
	DonorID     string `json:"donor_id"`
	AmountCents int    `json:"amount_cents"`
}
type ReceiptResult struct {
	DonorID string
	Status  string
}

func SendReceipts(client *ErrorClient, donations []Donation, send func(Donation) error) []ReceiptResult {
	results := make([]ReceiptResult, 0, len(donations))
	for _, donation := range donations {
		if err := send(donation); err != nil {
			if captureErr := client.Capture(map[string]any{"title": "donor receipt failed", "message": err.Error(), "level": "error", "fingerprint": []string{donation.DonorID, "receipt"}, "exception": fmt.Sprintf("%T: %v", err, err), "context": map[string]any{"donor_id": donation.DonorID, "amount_cents": donation.AmountCents}}); captureErr != nil {
				results = append(results, ReceiptResult{DonorID: donation.DonorID, Status: "capture_failed"})
				continue
			}
			results = append(results, ReceiptResult{DonorID: donation.DonorID, Status: "queued_for_review"})
			continue
		}
		results = append(results, ReceiptResult{DonorID: donation.DonorID, Status: "sent"})
	}
	return results
}
