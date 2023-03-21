// Contains everything related to API queries
// Uses tonapi.io
package main

import (
	"math/rand"
	"time"
)

type (
	addressFunc     func(*Investigation, string, []string) (Address, error)
	transactionFunc func(*Investigation, *Address) ([]Transaction, error)
)

// Lookup table for function pointers
var (
	addressFunctions = map[string]addressFunc{
		"TON": lookupTONAddress,
		"BS":  lookupBSAddress,
	}
	transactionFunctions = map[string]transactionFunc{
		"TON": lookupTONtransactions,
		"BS":  lookupBStransactions,
	}
)

// TON investigation
// Basic Brick of the recursion
func lookupTONAddress(invest *Investigation, identifier string, tags []string) (Address, error) {
	// TODO
	return Address{}, nil
}

func lookupTONtransactions(invest *Investigation, address *Address) ([]Transaction, error) {
	// TODO
	return []Transaction{}, nil
}

// BS crypto, used for tests
// For the test to make sense, we use addresses that use the alphabet
// as identifier (meaning max 24 addresses, which is fine for testing purposes)
func lookupBSAddress(invest *Investigation, identifier string, tags []string) (Address, error) {
	addr := &Address{
		Identifier: identifier,
		Tags:       tags,
		Balance:    0.0 + rand.Float64()*200.0, // Balance is limited to 200 BS
	}

	return *addr, nil
}

func lookupBStransactions(invest *Investigation, address *Address) ([]Transaction, error) {
	identifiers := []string{
		"A",
		"B",
		"C",
		"D",
		"E",
		"F",
		"G",
		"H",
		"I",
		"J",
		"K",
		"L",
		"M",
		"N",
		"O",
		"P",
		"Q",
		"R",
		"S",
		"T",
		"U",
		"V",
		"W",
		"X",
		"Y",
		"Z",
	}

	var t Transaction

	ret := []Transaction{}

	for i := 0; i < 16; i++ {
		t = Transaction{
			From:  identifiers[rand.Intn(len(identifiers))],
			To:    identifiers[rand.Intn(len(identifiers))],
			Value: rand.Float64() * 150.0,
			Date:  randate(),
		}

		ret = append(ret, t)
	}

	return ret, nil
}

// Helper for BS crypto
func randate() time.Time {
	min := time.Date(1970, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(2070, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min

	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0)
}
