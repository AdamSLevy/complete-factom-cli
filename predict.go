package main

import (
	"flag"
	"os"
	"strings"
	"time"

	"github.com/AdamSLevy/factom"
	"github.com/posener/complete"
)

var predictSingleTxName = complete.PredictFunc(func(a complete.Args) []string {
	argc := 0
	for _, arg := range a.Completed[1:] {
		argc++
		if string(arg[0]) == "-" {
			argc--
		}
	}
	if argc == 0 {
		return listTxNames()
	}
	return nil
})

var predictTxNameFCTAddress = complete.PredictFunc(func(a complete.Args) []string {
	argc := 0
	for _, arg := range a.Completed[1:] {
		argc++
		if string(arg[0]) == "-" {
			argc--
		}
	}
	switch argc {
	case 0:
		return listTxNames()
	case 1:
		return listFCTAddresses()
	}
	return nil
})

var predictTxNameECAddress = complete.PredictFunc(func(a complete.Args) []string {
	argc := 0
	for _, arg := range a.Completed[1:] {
		argc++
		if string(arg[0]) == "-" {
			argc--
		}
	}
	switch argc {
	case 0:
		return listTxNames()
	case 1:
		return listECAddresses()
	}
	return nil
})

var predictSingleAddress = complete.PredictFunc(func(a complete.Args) []string {
	argc := 0
	for _, arg := range a.Completed[1:] {
		argc++
		if string(arg[0]) == "-" {
			argc--
		}
	}
	if argc == 0 {
		return listAddresses()
	}
	return nil
})

var predictFCTAddressECAddress = complete.PredictFunc(func(a complete.Args) []string {
	argc := 0
	for _, arg := range a.Completed[1:] {
		argc++
		if string(arg[0]) == "-" {
			argc--
		}
	}
	switch argc {
	case 0:
		return listFCTAddresses()
	case 1:
		return listECAddresses()
	}
	return nil
})

var predictFCTAddressFCTAddress = complete.PredictFunc(func(a complete.Args) []string {
	argc := 0
	for _, arg := range a.Completed[1:] {
		argc++
		if string(arg[0]) == "-" {
			argc--
		}
	}
	switch argc {
	case 0:
		fallthrough
	case 1:
		return listFCTAddresses()
	}
	return nil
})

func predictSingleECAddress(optArgs []string) complete.PredictFunc {
	return complete.PredictFunc(func(a complete.Args) []string {
		argc := 0
		for _, arg := range a.Completed[1:] {
			argc++
			if string(arg[0]) == "-" {
				argc--
			}
			for _, optArg := range optArgs {
				if arg == optArg {
					argc--
					break
				}
			}
		}
		if argc == 0 {
			return listECAddresses()
		}
		return nil
	})
}

func listTxNames() []string {
	parseWalletFlags()
	txs, err := factom.ListTransactionsTmp()
	if err != nil {
		complete.Log("error: %v", err)
		return nil
	}
	txNames := make([]string, len(txs))
	for i, tx := range txs {
		txNames[i] = tx.Name
	}
	return txNames
}

func listECAddresses() []string {
	_, ecs := addressPubStrings()
	return ecs
}

func listFCTAddresses() []string {
	fcts, _ := addressPubStrings()
	return fcts
}

func listAddresses() []string {
	fcts, ecs := addressPubStrings()
	return append(fcts, ecs...)
}

func addressPubStrings() ([]string, []string) {
	parseWalletFlags()
	// Fetch all addresses.
	fcts, ecs, err := factom.FetchAddresses()
	if err != nil {
		complete.Log("error: %v", err)
		return nil, nil
	}

	// Create slices of the public address strings.
	fctAddresses := make([]string, len(fcts))
	for i, fct := range fcts {
		fctAddresses[i] = fct.String()
	}
	ecAddresses := make([]string, len(ecs))
	for i, ec := range ecs {
		ecAddresses[i] = ec.PubString()
	}
	return fctAddresses, ecAddresses
}

var flags *flag.FlagSet

func parseWalletFlags() {
	if flags != nil {
		return
	}
	// Parse any previously specified factom-cli options required for
	// connecting to factom-walletd
	flags = flag.NewFlagSet("", flag.ContinueOnError)
	flags.StringVar(&factom.RpcConfig.WalletServer, "w", "localhost:8089", "")
	flags.StringVar(&factom.RpcConfig.WalletTLSCertFile, "walletcert", "~/.factom/walletAPIpub.cert", "")
	flags.StringVar(&factom.RpcConfig.WalletRPCUser, "walletuser", "", "")
	flags.StringVar(&factom.RpcConfig.WalletRPCPassword, "walletpassword", "", "")
	flags.BoolVar(&factom.RpcConfig.WalletTLSEnable, "wallettls", false, "")
	flags.Parse(strings.Fields(os.Getenv("COMP_LINE"))[1:])
	factom.SetWalletTimeout(1 * time.Second)
}
