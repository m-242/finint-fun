// Contains everything related to API queries
// Uses tonapi.io
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
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
		"BTC": lookupBTCAddress,
	}
	transactionFunctions = map[string]transactionFunc{
		"TON": lookupTONtransactions,
		"BS":  lookupBStransactions,
		"BTC": lookupBTCtransaction,
	}
)

// BTC Investigation
type btcAddrAndTxs struct {
	Hash160       string `json:"hash160"`
	Address       string `json:"address"`
	NTx           int    `json:"n_tx"`
	NUnredeemed   int    `json:"n_unredeemed"`
	TotalReceived int    `json:"total_received"`
	TotalSent     int    `json:"total_sent"`
	FinalBalance  int    `json:"final_balance"`
	Txs           []struct {
		Hash        string `json:"hash"`
		Ver         int    `json:"ver"`
		VinSz       int    `json:"vin_sz"`
		VoutSz      int    `json:"vout_sz"`
		LockTime    int    `json:"lock_time"`
		Size        int    `json:"size"`
		RelayedBy   string `json:"relayed_by"`
		BlockHeight int    `json:"block_height"`
		TxIndex     int    `json:"tx_index"`
		Inputs      []struct {
			PrevOut struct {
				Hash  string `json:"hash"`
				Value int    `json:"value"`
			} `json:"prev_out"`
			Script string `json:"script"`
		} `json:"inputs"`
		Out []struct {
			Value  int    `json:"value"`
			Hash   string `json:"hash"`
			Script string `json:"script"`
		} `json:"out"`
	} `json:"txs"`
}

func btcHelper(invest *Investigation, identifier string) (btcAddrAndTxs, error) {
	rc := 429
	for rc != 200 {
		url := fmt.Sprintf("https://blockchain.info/rawaddr/%s", identifier)
		res, err := http.Get(url)
		if err != nil {
			invest.logger.Fatalln(err)
		}

		rc = res.StatusCode

		if rc == 429 {
			time.Sleep(time.Second * 5)
			invest.logger.Println("NOPE")
		} else {
			body, _ := ioutil.ReadAll(res.Body)
			var tmpAddress btcAddrAndTxs
			if err := json.Unmarshal(body, &tmpAddress); err != nil { // Parse []byte to go struct pointer
				invest.logger.Println("Can not unmarshal JSON")
				invest.logger.Println(string(body))
			}

			return tmpAddress, nil
		}

	}
	return btcAddrAndTxs{}, fmt.Errorf("Baad")

}

func lookupBTCAddress(invest *Investigation, identifier string, tags []string) (Address, error) {
	// https://blockchain.info/rawaddr/$bitcoin_address
	tmpAddress, _ := btcHelper(invest, identifier)

	return Address{
		Identifier: tmpAddress.Hash160,
		Tags:       tags,
		Balance:    float64(tmpAddress.FinalBalance),
	}, nil
}

func lookupBTCtransaction(
	invest *Investigation,
	identifier *Address,
) ([]Transaction, error) {
	tmpAddress, _ := btcHelper(invest, identifier.Identifier)

	txs := []Transaction{}
	for _, t := range tmpAddress.Txs {
		// Extract biggest outpout
		to := ""
		bb := 0.0
		for _, r := range t.Out {
			bb2 := float64(r.Value)
			if bb2 > bb {
				bb = bb2
				to = r.Hash
			}
		}

		v := float64(t.Inputs[0].PrevOut.Value)
		txs = append(txs, Transaction{
			From:  t.Inputs[0].PrevOut.Hash,
			Value: v,
			To:    to,
		})
	}

	return []Transaction{}, nil
}

// TON investigation
type tonAddress struct {
	Address struct {
		Bounceable    string `json:"bounceable"`
		NonBounceable string `json:"non_bounceable"`
		Raw           string `json:"raw"`
	} `json:"address"`
	Balance      int64    `json:"balance"`
	Interfaces   []string `json:"interfaces"`
	IsScam       bool     `json:"is_scam"`
	LastUpdate   int64    `json:"last_update"`
	MemoRequired bool     `json:"memo_required"`
	Status       string   `json:"status"`
}

func lookupTONAddress(invest *Investigation, identifier string, tags []string) (Address, error) {
	rc := 429
	for rc == 429 {
		url := fmt.Sprintf("https://tonapi.io/v1/account/getInfo?account=%s", identifier)
		res, err := http.Get(url)
		if err != nil {
			invest.logger.Fatalln(err)
		}

		body, err := ioutil.ReadAll(res.Body)
		var tmpAddress tonAddress
		if err := json.Unmarshal(body, &tmpAddress); err != nil { // Parse []byte to go struct pointer
			invest.logger.Println("Can not unmarshal JSON")
			invest.logger.Println(string(body))
			return Address{}, err
		}

		rc = res.StatusCode
		if rc == 429 {
			invest.logger.Println("Being rate limited, waiting")
			time.Sleep(time.Second * 5)
		} else {
			addr := Address{
				Identifier: tmpAddress.Address.Raw,
				Tags:       tags,
				Balance:    float64(tmpAddress.Balance),
			}

			return addr, nil
		}
	}
	panic("WTF")
}

type tonTransactions struct {
	Transactions []struct {
		InMsg struct {
			CreatedLt   int64 `json:"created_lt"`
			Destination struct {
				Address string `json:"address"`
				Icon    string `json:"icon"`
				IsScam  bool   `json:"is_scam"`
				Name    string `json:"name"`
			} `json:"destination"`
			FwdFee int `json:"fwd_fee"`
			IhrFee int `json:"ihr_fee"`
			Source struct {
				Address string `json:"address"`
				Icon    string `json:"icon"`
				IsScam  bool   `json:"is_scam"`
				Name    string `json:"name"`
			} `json:"source"`
			Value int `json:"value"`
		} `json:"in_msg"`
	}
}

func lookupTONtransactions(invest *Investigation, address *Address) ([]Transaction, error) {
	rc := 429
	for rc == 429 {
		transactions := []Transaction{}
		if address.Identifier == "" {
			return transactions, nil
		}

		url := fmt.Sprintf(
			"https://tonapi.io/v1/blockchain/getTransactions?account=%s&limit=100",
			address.Identifier,
		)
		// fmt.Println(url)
		res, err := http.Get(url)
		if err != nil {
			invest.logger.Print(err)
		}

		rc = res.StatusCode
		if rc == 429 {
			invest.logger.Println("Being rate limited, waiting")
			time.Sleep(time.Second * 5)
		} else {

			body, _ := ioutil.ReadAll(res.Body)
			var tmp tonTransactions
			if err := json.Unmarshal(body, &tmp); err != nil { // Parse []byte to go struct pointer
				invest.logger.Println("Can not unmarshal JSON")
				invest.logger.Println(err)
				// return transactions, err
			}

			for _, tt := range tmp.Transactions {
				trans := Transaction{
					From:  tt.InMsg.Source.Address,
					To:    tt.InMsg.Destination.Address,
					Value: float64(tt.InMsg.Value),
					// Date:  time.Unix(tt.InMsg.CreatedLt, 0),
				}
				transactions = append(transactions, trans)
			}

			return transactions, nil
		}
	}
	panic("WTF2")
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
		b := rand.Int() % 2

		t = Transaction{
			To:    identifiers[rand.Intn(len(identifiers))],
			Value: rand.Float64() * 150.0,
			Date:  randate(),
		}

		if b == 1 {
			t.From = identifiers[rand.Intn(len(identifiers))]
			t.To = address.Identifier
		} else {
			t.From = address.Identifier
			t.To = identifiers[rand.Intn(len(identifiers))]
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
