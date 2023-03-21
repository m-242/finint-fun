// Everything related to a general investigation
package main

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
	}

	// TODO

	return invest
}

// The recursive exploration function
func (invest *Investigation) explore(startingPoint string, depth int) {
	// TODO !!!!!!
}
