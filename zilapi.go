package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Zilliqa/gozilliqa-sdk/bech32"
	"github.com/Zilliqa/gozilliqa-sdk/keytools"
	"github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/jmoiron/sqlx"
)

// this function takes in a ds block number and store the result
func InsertDSBlockMinerInfo(dsblocknumber int, db *sqlx.DB) (result string) {
	row := db.QueryRow("select count(*) as numrecords from public.\"DSBlockMinerInfo\" where dsblock=$1", dsblocknumber)
	var numrecords uint
	row.Scan(&numrecords)
	if numrecords != 0 {
		println("Already prepared dsblock - doing nothing for ds block " + fmt.Sprintf("%d", dsblocknumber))
		return "Already prepared dsblock - doing nothing for ds block " + fmt.Sprintf("%d", dsblocknumber)
	}

	MinerWallets := []map[string]interface{}{}

	provider := provider.NewProvider("https://api.zilliqa.com/")
	response, _ := provider.GetMinerInfo(fmt.Sprint(dsblocknumber))
	DSCommittee := response.DsCommittee
	Shards := response.Shards
	// Get Nodes
	for i, si := range Shards {
		for i2, v := range si.Nodes {
			walletaddressbech32 := GetBech32AddressFromPublicKey(v)
			walletNonBech32 := strings.ToLower(ConvertToBech32Address(walletaddressbech32))
			mp1 := map[string]interface{}{
				"dsblock":             dsblocknumber,
				"dscommittee":         false,
				"node":                true,
				"walletpublickey":     v,
				"walletaddress":       walletNonBech32,
				"walletaddressbech32": "0x" + walletaddressbech32,
			}
			MinerWallets = append(MinerWallets, mp1)
			_ = i2
		}
		_ = i
	}
	// Get DS Committee
	for i, s := range DSCommittee {

		walletaddressbech32 := GetBech32AddressFromPublicKey(s)
		walletNonBech32 := strings.ToLower(ConvertToBech32Address(walletaddressbech32))
		mp1 := map[string]interface{}{
			"dsblock":             dsblocknumber,
			"dscommittee":         true,
			"node":                false,
			"walletpublickey":     s,
			"walletaddress":       walletNonBech32,
			"walletaddressbech32": "0x" + walletaddressbech32,
		}
		MinerWallets = append(MinerWallets, mp1)
		_ = i
	}

	dbresult, err := db.NamedExec("INSERT INTO public.\"DSBlockMinerInfo\" (dsblock, dscommittee, node, walletpublickey, walletaddress, walletaddressbech32) VALUES (:dsblock, :dscommittee, :node, :walletpublickey, :walletaddress, :walletaddressbech32)", MinerWallets)
	_ = dbresult
	_ = err
	println("Successfully ran InsertDSBlockMinerInfo for dx block " + fmt.Sprintf("%d", dsblocknumber))
	return "Successfully ran InsertDSBlockMinerInfo for dx block " + fmt.Sprintf("%d", dsblocknumber)
}

// this function takes in a tx block number and stores the transactions
func InsertTXBlockTransactions(txblocknumber int, db *sqlx.DB) (result string) {
	row := db.QueryRow("select count(*) as numrecords from public.\"Transactions\" where txblock=$1", txblocknumber)
	var numrecords uint
	row.Scan(&numrecords)
	if numrecords != 0 {
		println("Already prepared txblock - doing nothing for tx block " + fmt.Sprintf("%d", txblocknumber))
		return "Already prepared txblock - doing nothing for tx block " + fmt.Sprintf("%d", txblocknumber)
	}

	Transactions := []map[string]interface{}{}

	provider := provider.NewProvider("https://api.zilliqa.com/")
	response, _ := provider.GetTxnBodiesForTxBlock(fmt.Sprint(txblocknumber))
	// Iterate through the transactions
	for i, si := range response {
		// Skip any smart contract calls as for mining we dont care
		if si.Data == nil {
			senderaddressbech32 := GetBech32AddressFromPublicKey(si.SenderPubKey)
			senderaddressnonnbech32 := strings.ToLower(ConvertToBech32Address(senderaddressbech32))
			toaddressbech32 := strings.ToLower(ConvertToBech32Address(si.ToAddr))
			amount, err := strconv.ParseFloat(si.Amount, 64)
			_ = err
			mp1 := map[string]interface{}{
				"dsblock":             (txblocknumber / 100) - 1,
				"txblock":             txblocknumber,
				"txid":                "0x" + si.ID,
				"amount":              amount / 1000000000000,
				"success":             si.Receipt.Success,
				"senderpublickey":     si.SenderPubKey,
				"senderaddressbech32": "0x" + senderaddressbech32,
				"senderaddress":       senderaddressnonnbech32,
				"toaddress":           toaddressbech32,
				"toaddressbech32":     "0x" + si.ToAddr,
			}
			Transactions = append(Transactions, mp1)
		}

		_ = i
	}

	dbresult, err := db.NamedExec("INSERT INTO public.\"Transactions\" (dsblock, txblock, txid, amount, success, senderpublickey, senderaddressbech32, senderaddress, toaddress, toaddressbech32) VALUES (:dsblock, :txblock, :txid, :amount, :success, :senderpublickey, :senderaddressbech32, :senderaddress, :toaddress, :toaddressbech32)", Transactions)
	_ = dbresult
	_ = err
	println("Successfully ran InsertTXBlockTransactions for txblock " + fmt.Sprintf("%d", txblocknumber))
	return "Successfully ran InsertTXBlockTransactions for txblock " + fmt.Sprintf("%d", txblocknumber)
}

// converts the 0x address to zil1
func ConvertToBech32Address(address string) (Bech32Address string) {
	Bech32Address, _ = bech32.ToBech32Address(address)
	return Bech32Address
}

// converts the zil1 address to 0x
func ConvertFromBech32Address(Bech32Address string) (address string) {
	address, _ = bech32.FromBech32Addr(Bech32Address)
	return address
}

// GetBech32AddressFromPublicKey
func GetBech32AddressFromPublicKey(PubKey string) (address string) {
	b := util.DecodeHex(PubKey)
	address = keytools.GetAddressFromPublic(b)
	return address
}
