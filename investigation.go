// Everything related to a general investigation
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
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
	addr, _ := addressLookup(invest, startingPoint, tags)

	if addr.Identifier == "" {
		return nil
	}

	invest.logger.Println(addr)

	invest.AddAddress(addr)

	if depth == 0 { // Stop condition
		return nil
	}

	transactions, _ := transactionsLookup(invest, &addr)
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
}

func (invest *Investigation) toJSON(path string) error {
	file, err := json.MarshalIndent(invest, "", " ")
	if err != nil {
		invest.logger.Printf("Couldn't marshall invest to json.")
		return err
	}

	return ioutil.WriteFile(path, file, 0644)
}

func fromJSON(path string) (Investigation, error) {
	var iv Investigation

	jsonFile, _ := os.Open(path)
	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &iv)

	return iv, nil
}
