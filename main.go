package main

import (
	"fmt"
	"log"
)

func main() {
	client := NewErrorClient()
	donations := []Donation{{DonorID: "donor-42", AmountCents: 2500}, {DonorID: "donor-73", AmountCents: 5000}}
	results := SendReceipts(client, donations, func(d Donation) error {
		if d.DonorID == "donor-73" {
			return fmt.Errorf("mailbox rejected receipt")
		}
		return nil
	})
	for _, result := range results {
		log.Printf("%s: %s", result.DonorID, result.Status)
	}
}
