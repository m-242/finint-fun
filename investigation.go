// Everything related to a general investigation
package main

import (
	"fmt"
	"log"
)

func Investigate(
	startingPoint string,
	startingTags []string,
	depth int,
	currency, apiKey string,
) *Investigation {
	invest := &Investigation{
		InvolvedAddresses: []Address{},
		Transactions:      []Transaction{},
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
		fmt.Printf("Bad aaa")
	}

	if addr.Identifier == "" {
		return nil
	}

	invest.logger.Println(addr)

	invest.AddAddress(addr)

	if depth == 0 { // Stop condition
		return nil
	}

	transactions, err := transactionsLookup(invest, &addr)
	for _, t := range transactions {
		if t.From == "" || t.To == "" {
			continue
		}
		invest.logger.Printf("%s - (%.8f) -> %s\n", t.From, t.Value, t.To)

		if invest.HasAddressWithID(t.From) && !invest.HasAddressWithID(t.To) {
			invest.explore(t.To, []string{}, depth-1, addressLookup, transactionsLookup)
		}

		if invest.HasAddressWithID(t.To) && !invest.HasAddressWithID(t.From) {
			invest.explore(t.From, []string{}, depth-1, addressLookup, transactionsLookup)
		}

		invest.AddTransaction(t)
	}
	return nil

	// TODO
	//	invest.logger.Printf("Investigating from %s\n", startingPoint)
	//
	//	addr, err := addressLookup(invest, startingPoint, tags)
	//	if err != nil {
	//		// TODO traitement d'erreur
	//		invest.logger.Printf("Could'nt investigate node %s\n", startingPoint)
	//		return err
	//	}
	//
	//	transactions, err := transactionsLookup(invest, &addr)
	//	if err != nil {
	//		invest.logger.Printf("Couldn't get node %s's transactions\n", startingPoint)
	//		return err
	//	}
	//
	//	if addr.isNew(*invest) {
	//		invest.AddAddress(&addr)
	//	}
	//
	//	var newAddressID string
	//
	//	for _, trans := range transactions {
	//		if trans.isNew(*invest) {
	//			if trans.From == startingPoint {
	//				newAddressID = trans.To
	//			} else if trans.To == startingPoint {
	//				newAddressID = trans.To
	//			}
	//
	//			if invest.HasAddressWithID(newAddressID) {
	//				invest.AddTransaction(&trans)
	//			} else if depth > 0 {
	//				err = invest.explore(newAddressID, []string{}, depth-1, addressLookup, transactionsLookup)
	//				if err != nil {
	//					invest.logger.Printf("Couldn't explore from %s\n", newAddressID)
	//				}
	//			} else {
	//				addr, err := addressLookup(invest, newAddressID, []string{})
	//				if err != nil {
	//					invest.logger.Printf("Couldn't investigate border node %s\n", newAddressID)
	//				}
	//				invest.AddAddress(&addr)
	//			}
	//		}
	//	}

	return nil
}
