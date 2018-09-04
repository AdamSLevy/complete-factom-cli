// +build ignore

package main

import (
	"fmt"
	"netki"
	"strconv"
)

func handleErr(err error) {
	fmt.Println(fmt.Sprintf("Received Error: %s", err))
}

func main() {

	partner := netki.NewNetkiPartner("PARTNER_ID", "APU_KEY", "http://localhost:5000")

	// Test Partners
	newPartner, err := partner.CreateNewPartner("Golang Test Partner")
	if err != nil {
		handleErr(err)
	}

	partners, err := partner.GetPartners()
	fmt.Println(partners)
	if err != nil {
		handleErr(err)
	}

	err = partner.DeletePartner(newPartner)
	if err != nil {
		handleErr(err)
	}

	domain, err := partner.CreateNewDomain("golang-testdomain.com", partners[0])
	if err != nil {
		handleErr(err)
	}

	domains, err := partner.GetDomains()
	if err != nil {
		handleErr(err)
	}
	fmt.Println(domains)

	for _, d := range domains {
		result, err := partner.GetDomainStatus(d)
		if err != nil {
			handleErr(err)
		} else {
			fmt.Println("Domain Status: ")
			fmt.Println(result)
		}

		result, err = partner.GetDomainDnssec(d)
		if err != nil {
			handleErr(err)
		} else {
			fmt.Print("DNSSEC Status: ")
			fmt.Println(result)
		}

	}

	partner.DeleteDomain(domain)

	// Test Walletname GET and Clear
	names, err := partner.GetWalletNames(netki.Domain{}, "")
	if err != nil {
		handleErr(nil)
	}

	// Delete golangtest wallet if it exists
	for _, wn := range names {
		fmt.Println("Found WN: " + wn.Name)
		if wn.Name == "golangtest" {
			fmt.Println("Deleting WN: " + wn.Name)
			wn.Delete(partner)
		}
	}

	// Test WalletName CUD Cycle
	wn := partner.CreateNewWalletName(domains[1], "golangtest", make([]netki.Wallet, 0), "goLangExternalId")
	wn.SetCurrencyAddress("btc", "1dgkjsdfhlkjfsdhkjlsdfhsgdf")
	wn.GetAddress("btc")

	usedCurrencies := wn.UsedCurrencies()
	for index, currency := range usedCurrencies {
		fmt.Println("Currency [" + strconv.Itoa(index) + "]: " + currency)
	}

	err = wn.Save(partner)
	if err != nil {
		handleErr(err)
	}

	wn.SetCurrencyAddress("dgc", "D548376529834756928376523")
	err = wn.Save(partner)
	if err != nil {
		handleErr(err)
	}

	err = wn.Delete(partner)
	if err != nil {
		handleErr(err)
	}

}
