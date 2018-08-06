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

// Parse any previously specified factom-cli options required for connecting to
// factom-walletd
func parseWalletFlags() {
	if flags != nil {
		// We already parsed the flags.
		return
	}
	// Using flag.FlagSet allows us to parse a custom array of flags
	// instead of this programs args.
	flags = flag.NewFlagSet("", flag.ContinueOnError)
	flags.StringVar(&factom.RpcConfig.WalletServer, "w", "localhost:8089", "")
	flags.StringVar(&factom.RpcConfig.WalletTLSCertFile, "walletcert",
		"~/.factom/walletAPIpub.cert", "")
	flags.StringVar(&factom.RpcConfig.WalletRPCUser, "walletuser", "", "")
	flags.StringVar(&factom.RpcConfig.WalletRPCPassword, "walletpassword", "", "")
	flags.BoolVar(&factom.RpcConfig.WalletTLSEnable, "wallettls", false, "")

	// flags.Parse will print warnings if it comes across an unrecognized
	// flag. We don't want this so we temprorarily redirect everything to
	// /dev/null before we call flags.Parse().
	stdout := os.Stdout
	stderr := os.Stderr
	os.Stdout, _ = os.Open(os.DevNull)
	os.Stderr = os.Stdout

	// The current command line being typed is stored in the environment
	// variable COMP_LINE. We split on spaces and discard the first in the
	// list because it is the program name `factom-cli`.
	flags.Parse(strings.Fields(os.Getenv("COMP_LINE"))[1:])

	// Restore stdout and stderr.
	os.Stdout = stdout
	os.Stderr = stderr

	// We want need factom-walletd to timeout or the CLI completion will
	// hang and never return. This is the whole reason we use AdamSLevy's
	// fork of factom.
	factom.SetWalletTimeout(1 * time.Second)
}
