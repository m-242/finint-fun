// Everything related to a general investigation
package main

import (
	"log"
)

func Investigate(
	startingPoint string,
	startingTags []string,
	depth int,
	currency, apiKey string,
) *Investigation {
	invest := &Investigation{
		InvolvedAddresses: make([]*Address, 0),
		Transactions:      make([]*Transaction, 0),
		Currency:          currency,
		ApiKey:            apiKey,
		logger:            *log.Default(),
	}

	addressLookup, ok := addressFunctions[currency]
	if !ok {
		invest.logger.Fatalf("No Address Lookup function for currency %s\n", currency)
	}

	transactionsLookup, ok := transactionFunctions[currency]
	if !ok {
		invest.logger.Fatalf("No Transactions Lookup function for currency %s\n", currency)
	}

	invest.explore(startingPoint, startingTags, depth, addressLookup, transactionsLookup)

	return invest
}

// The recursive exploration function
func (invest *Investigation) explore(
	startingPoint string,
	tags []string,
	depth int,
	addressLookup addressFunc,
	transactionsLookup transactionFunc,
) error {
	invest.logger.Printf("Investigating from %s\n", startingPoint)

	addr, err := addressLookup(invest, startingPoint, tags)
	if err != nil {
		// TODO traitement d'erreur
		invest.logger.Printf("Could'nt investigate node %s\n", startingPoint)
		return err
	}

	transactions, err := transactionsLookup(invest, &addr)
	if err != nil {
		invest.logger.Printf("Couldn't get node %s's transactions\n", startingPoint)
		return err
	}

	if addr.isNew(*invest) {
		invest.AddAddress(&addr)
	}

	var newAddressID string

	for _, trans := range transactions {
		if trans.isNew(*invest) {
			if trans.From == startingPoint {
				newAddressID = trans.To
			} else if trans.To == startingPoint {
				newAddressID = trans.To
			}

			if invest.HasAddressWithID(newAddressID) {
				invest.AddTransaction(&trans)
			} else if depth > 0 {
				err = invest.explore(newAddressID, []string{}, depth-1, addressLookup, transactionsLookup)
				if err != nil {
					invest.logger.Printf("Couldn't explore from %s\n", newAddressID)
				}
			}
		}
	}

	return nil
}
