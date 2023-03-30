package main

import (
	"log"
	"time"
)

// Investigation is the big data structure
// that holds every bit of data
type Investigation struct {
	InvolvedAddresses []Address     `json:"addresses"`
	Transactions      []Transaction `json:"transactions"`
	Currency          string        `json:"currency"`
	ApiKey            string        `json:"-"` // The api key used to makes requests to the backend api
	logger            log.Logger    `json:"-"`
}

func (i *Investigation) AddAddress(a Address) {
	i.logger.Printf("Adding address %s to the exploration\n", a.Identifier)
	i.InvolvedAddresses = append(i.InvolvedAddresses, a)
}

func (i *Investigation) AddTransaction(t Transaction) {
	i.logger.Printf("Adding %s-(%+f%s)->%s to the exploration\n", t.From, t.Value, i.Currency, t.To)
	i.Transactions = append(i.Transactions, t)
}

func (i *Investigation) HasAddressWithID(id string) bool {
	for _, a := range i.InvolvedAddresses {
		if id == a.Identifier {
			return true
		}
	}
	return false
}

// Nodes/Address of a crypto currency wallet
type Address struct {
	Identifier string   `json:"identifier"` // Identifier
	Tags       []string `json:"tags"`       // Tags, so that the user can identify the addresses
	Balance    float64  `json:"balance"`
}

func (a1 *Address) equals(a2 Address) bool {
	return a1.Identifier == a2.Identifier
}

func (address Address) isNew(invest Investigation) bool {
	for _, a2 := range invest.InvolvedAddresses {
		if address.equals(a2) {
			return false
		}
	}
	return true
}

// Transaction/Edge
type Transaction struct {
	From  string    `json:"from"` // Identifiers of addresses
	To    string    `json:"to"`
	Value float64   `json:"value"`
	Date  time.Time `json:"date"`
}

func (t1 *Transaction) equals(t2 Transaction) bool {
	return t1.From == t2.From && t1.To == t2.To &&
		t1.Value == t2.Value &&
		t1.Date == t2.Date
}

func (t Transaction) isNew(invest Investigation) bool {
	for _, t2 := range invest.Transactions {
		if t.equals(t2) {
			return false
		}
	}
	return true
}
