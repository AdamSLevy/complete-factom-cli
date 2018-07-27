package main

import "github.com/posener/complete"

func main() {
	// addchain [-fq] [-n NAME1 -n NAME2 -h HEXNAME3 ] [-CET] ECADDRESS <STDIN>
	addchain := complete.Command{
		Flags: complete.Flags{
			"-f": complete.PredictNothing,
			"-q": complete.PredictNothing,

			"-n": complete.PredictAnything,
			"-h": complete.PredictAnything,

			"-C": complete.PredictNothing,
			"-E": complete.PredictNothing,
			"-T": complete.PredictNothing,
		},
		Args: complete.PredictAnything,
	}
	// addentry [-fq] [-n NAME1 -h HEXNAME2 ...|-c CHAINID] [-e EXTID1 -e EXTID2 -x HEXEXTID ...] [-CET] ECADDRESS <STDIN>
	addentry := complete.Command{
		Flags: complete.Flags{
			"-f": complete.PredictNothing,
			"-q": complete.PredictNothing,

			"-n": complete.PredictAnything,
			"-h": complete.PredictAnything,
			"-c": complete.PredictAnything,

			"-e": complete.PredictAnything,
			"-x": complete.PredictAnything,

			"-C": complete.PredictNothing,
			"-E": complete.PredictNothing,
			"-T": complete.PredictNothing,
		},
		Args: complete.PredictAnything,
	}
	// addtxecoutput [-rq] TXNAME ADDRESS AMOUNT
	addtxecoutput := complete.Command{
		Flags: complete.Flags{
			"-r": complete.PredictNothing,
			"-q": complete.PredictNothing,
		},
		Args: complete.PredictAnything,
	}
	// addtxfee [-q] TXNAME ADDRESS
	addtxfee := complete.Command{
		Flags: complete.Flags{
			"-q": complete.PredictNothing,
		},
		Args: complete.PredictAnything,
	}
	// addtxinput [-q] TXNAME ADDRESS AMOUNT
	addtxinput := complete.Command{
		Flags: complete.Flags{
			"-q": complete.PredictNothing,
		},
		Args: complete.PredictAnything,
	}
	// addtxoutput [-rq] TXNAME ADDRESS AMOUNT
	addtxoutput := complete.Command{
		Flags: complete.Flags{
			"-r": complete.PredictNothing,
			"-q": complete.PredictNothing,
		},
		Args: complete.PredictAnything,
	}
	// backupwallet
	backupwallet := complete.Command{
		Args: complete.PredictNothing,
	}
	// balance [-r] ADDRESS
	balance := complete.Command{
		Flags: complete.Flags{
			"-r": complete.PredictNothing,
		},
		Args: complete.PredictAnything,
	}
	// buyec [-fqrT] FCTADDRESS ECADDRESS ECAMOUNT
	buyec := complete.Command{
		Flags: complete.Flags{
			"-f": complete.PredictNothing,
			"-r": complete.PredictNothing,
			"-q": complete.PredictNothing,
			"-T": complete.PredictNothing,
		},
		Args: complete.PredictAnything,
	}
	// composechain [-f] [-n NAME1 -n NAME2 -h HEXNAME3 ] ECADDRESS <STDIN>
	composechain := complete.Command{
		Flags: complete.Flags{
			"-f": complete.PredictNothing,

			"-n": complete.PredictAnything,
			"-h": complete.PredictAnything,
		},
		Args: complete.PredictAnything,
	}
	// composeentry [-f] [-n NAME1 -h HEXNAME2 ...|-c CHAINID]  [-e EXTID1 -e EXTID2 -x HEXEXTID ...] ECADDRESS <STDIN>
	composeentry := complete.Command{
		Flags: complete.Flags{
			"-f": complete.PredictNothing,

			"-n": complete.PredictAnything,
			"-h": complete.PredictAnything,
			"-c": complete.PredictAnything,

			"-e": complete.PredictAnything,
			"-x": complete.PredictAnything,
		},
		Args: complete.PredictAnything,
	}
	// composetx TXNAME
	composetx := complete.Command{
		Args: complete.PredictAnything,
	}
	// ecrate
	ecrate := complete.Command{
		Args: complete.PredictNothing,
	}
	// exportaddresses
	exportaddresses := complete.Command{
		Args: complete.PredictNothing,
	}
	// help
	help := complete.Command{
		Sub: complete.Commands{
			"addchain":        complete.Command{},
			"addentry":        complete.Command{},
			"addtxecoutput":   complete.Command{},
			"addtxfee":        complete.Command{},
			"addtxinput":      complete.Command{},
			"addtxoutput":     complete.Command{},
			"backupwallet":    complete.Command{},
			"balance":         complete.Command{},
			"buyec":           complete.Command{},
			"composechain":    complete.Command{},
			"composeentry":    complete.Command{},
			"composetx":       complete.Command{},
			"ecrate":          complete.Command{},
			"exportaddresses": complete.Command{},
			"get":             complete.Command{},
			"importaddress":   complete.Command{},
			"importkoinify":   complete.Command{},
			"listaddresses":   complete.Command{},
			"listtxs":         complete.Command{},
			"newecaddress":    complete.Command{},
			"newfctaddress":   complete.Command{},
			"newtx":           complete.Command{},
			"properties":      complete.Command{},
			"receipt":         complete.Command{},
			"rmaddress":       complete.Command{},
			"rmtx":            complete.Command{},
			"sendfct":         complete.Command{},
			"sendtx":          complete.Command{},
			"signtx":          complete.Command{},
			"status":          complete.Command{},
			"subtxfee":        complete.Command{},
		},
	}
	// get abheight HEIGHT -r (to suppress Raw Data)
	get_abheight := complete.Command{
		Flags: complete.Flags{
			"-r": complete.PredictNothing,
		},
		Args: complete.PredictAnything,
	}
	// get allentries [-n NAME1 -h HEXNAME2 ...|CHAINID] [-E]
	get_allentries := complete.Command{
		Flags: complete.Flags{
			"-n": complete.PredictAnything,
			"-h": complete.PredictAnything,

			"-E": complete.PredictNothing,
		},
		Args: complete.PredictAnything,
	}
	// get chainhead [-n NAME1 -h HEXNAME2 ...|CHAINID] [-K]
	get_chainhead := complete.Command{
		Flags: complete.Flags{
			"-n": complete.PredictAnything,
			"-h": complete.PredictAnything,

			"-K": complete.PredictNothing,
		},
		Args: complete.PredictAnything,
	}
	// get dbheight HEIGHT -r (to suppress Raw Data)
	get_dbheight := complete.Command{
		Flags: complete.Flags{
			"-r": complete.PredictNothing,
		},
		Args: complete.PredictAnything,
	}
	// get dblock KEYMR
	get_dblock := complete.Command{
		Args: complete.PredictAnything,
	}
	// get eblock KEYMR
	get_eblock := complete.Command{
		Args: complete.PredictAnything,
	}
	// get ecbheight HEIGHT -r (to suppress Raw Data)
	get_ecbheight := complete.Command{
		Flags: complete.Flags{
			"-r": complete.PredictNothing,
		},
		Args: complete.PredictAnything,
	}
	// get entry HASH
	get_entry := complete.Command{
		Args: complete.PredictAnything,
	}
	// get fbheight HEIGHT -r (to suppress Raw Data)
	get_fbheight := complete.Command{
		Flags: complete.Flags{
			"-r": complete.PredictNothing,
		},
		Args: complete.PredictAnything,
	}
	// get firstentry [-n NAME1 -h HEXNAME2 ...|CHAINID] [-E]
	get_firstentry := complete.Command{
		Flags: complete.Flags{
			"-n": complete.PredictAnything,
			"-h": complete.PredictAnything,

			"-E": complete.PredictNothing,
		},
		Args: complete.PredictAnything,
	}
	// get head [-K]
	get_head := complete.Command{
		Flags: complete.Flags{
			"-K": complete.PredictNothing,
		},
		Args: complete.PredictAnything,
	}
	// get heights
	get_heights := complete.Command{
		Args: complete.PredictNothing,
	}
	// get pendingentries [-E]
	get_pendingentries := complete.Command{
		Flags: complete.Flags{
			"-E": complete.PredictNothing,
		},
		Args: complete.PredictNothing,
	}
	// get pendingtransactions [-T]
	get_pendingtransactions := complete.Command{
		Flags: complete.Flags{
			"-T": complete.PredictNothing,
		},
		Args: complete.PredictNothing,
	}
	// get raw HASH
	get_raw := complete.Command{
		Args: complete.PredictAnything,
	}
	// get walletheight
	get_walletheight := complete.Command{
		Args: complete.PredictNothing,
	}
	// get allentries|chainhead|dblock|eblock|entry|firstentry|head|heights|walletheight|pendingentries|pendingtransactions|raw|dbheight|abheight|fbheight|ecbheight
	get := complete.Command{
		Sub: complete.Commands{
			"abheight":            get_abheight,
			"allentries":          get_allentries,
			"chainhead":           get_chainhead,
			"dbheight":            get_dbheight,
			"dblock":              get_dblock,
			"eblock":              get_eblock,
			"ecbheight":           get_ecbheight,
			"entry":               get_entry,
			"fbheight":            get_fbheight,
			"firstentry":          get_firstentry,
			"head":                get_head,
			"heights":             get_heights,
			"pendingentries":      get_pendingentries,
			"pendingtransactions": get_pendingtransactions,
			"raw":          get_raw,
			"walletheight": get_walletheight,
		},
	}
	// importaddress ADDRESS [ADDRESS...]
	importaddress := complete.Command{
		Args: complete.PredictAnything,
	}
	// importkoinify '12WORDS'
	importkoinify := complete.Command{
		Args: complete.PredictAnything,
	}
	// listaddresses
	listaddresses := complete.Command{
		Args: complete.PredictNothing,
	}
	// listtxs address [-T] ECADDRESS|FCTADDRESS
	listtxs_address := complete.Command{
		Flags: complete.Flags{
			"-T": complete.PredictNothing,
		},
		Args: complete.PredictAnything,
	}
	// listtxs [all] [-T]
	listtxs_all := complete.Command{
		Flags: complete.Flags{
			"-T": complete.PredictNothing,
		},
		Args: complete.PredictNothing,
	}
	// listtxs id TXID
	listtxs_id := complete.Command{
		Args: complete.PredictAnything,
	}
	// listtxs name TXNAME
	listtxs_name := complete.Command{
		Args: complete.PredictAnything,
	}
	// listtxs range [-T] START END
	listtxs_range := complete.Command{
		Flags: complete.Flags{
			"-T": complete.PredictNothing,
		},
		Args: complete.PredictAnything,
	}
	// listtxs tmp
	listtxs_tmp := complete.Command{
		Flags: complete.Flags{
			"-N": complete.PredictNothing,
		},
		Args: complete.PredictNothing,
	}
	// listtxs [address|all|id|name|tmp|range]
	listtxs := complete.Command{
		Flags: complete.Flags{
			"-T": complete.PredictNothing,
		},
		Sub: complete.Commands{
			"address": listtxs_address,
			"all":     listtxs_all,
			"id":      listtxs_id,
			"name":    listtxs_name,
			"tmp":     listtxs_tmp,
			"range":   listtxs_range,
		},
	}
	// newecaddress
	newecaddress := complete.Command{
		Args: complete.PredictNothing,
	}
	// newfctaddress
	newfctaddress := complete.Command{
		Args: complete.PredictNothing,
	}
	// newtx [-q] TXNAME
	newtx := complete.Command{
		Flags: complete.Flags{
			"-q": complete.PredictNothing,
		},
		Args: complete.PredictAnything,
	}
	// properties
	properties := complete.Command{
		Args: complete.PredictNothing,
	}
	// receipt ENTRYHASH
	receipt := complete.Command{
		Args: complete.PredictAnything,
	}
	// rmaddress ADDRESS
	rmaddress := complete.Command{
		Args: complete.PredictAnything,
	}
	// rmtx TXNAME
	rmtx := complete.Command{
		Args: complete.PredictAnything,
	}
	// sendfct [-fqrT] FROMADDRESS TOADDRESS AMOUNT
	sendfct := complete.Command{
		Flags: complete.Flags{
			"-f": complete.PredictNothing,
			"-q": complete.PredictNothing,
			"-r": complete.PredictNothing,
			"-T": complete.PredictNothing,
		},
		Args: complete.PredictAnything,
	}
	// sendtx [-fqT] TXNAME
	sendtx := complete.Command{
		Flags: complete.Flags{
			"-f": complete.PredictNothing,
			"-q": complete.PredictNothing,
			"-T": complete.PredictNothing,
		},
		Args: complete.PredictAnything,
	}
	// signtx [-fqT] TXNAME
	signtx := complete.Command{
		Flags: complete.Flags{
			"-f": complete.PredictNothing,
			"-q": complete.PredictNothing,
			"-T": complete.PredictNothing,
		},
		Args: complete.PredictAnything,
	}
	// status TxID|FullTx
	status := complete.Command{
		Args: complete.PredictAnything,
	}
	// subtxfee [-q] TXNAME ADDRESS
	subtxfee := complete.Command{
		Flags: complete.Flags{
			"-q": complete.PredictNothing,
		},
		Args: complete.PredictAnything,
	}

	cli := complete.Command{
		Sub: complete.Commands{
			"addchain":        addchain,
			"addentry":        addentry,
			"addtxecoutput":   addtxecoutput,
			"addtxfee":        addtxfee,
			"addtxinput":      addtxinput,
			"addtxoutput":     addtxoutput,
			"backupwallet":    backupwallet,
			"balance":         balance,
			"buyec":           buyec,
			"composechain":    composechain,
			"composeentry":    composeentry,
			"composetx":       composetx,
			"ecrate":          ecrate,
			"exportaddresses": exportaddresses,
			"help":            help,
			"get":             get,
			"importaddress":   importaddress,
			"importkoinify":   importkoinify,
			"listaddresses":   listaddresses,
			"listtxs":         listtxs,
			"newecaddress":    newecaddress,
			"newfctaddress":   newfctaddress,
			"newtx":           newtx,
			"properties":      properties,
			"receipt":         receipt,
			"rmaddress":       rmaddress,
			"rmtx":            rmtx,
			"sendfct":         sendfct,
			"sendtx":          sendtx,
			"signtx":          signtx,
			"status":          status,
			"subtxfee":        subtxfee,
		},
		Flags: complete.Flags{
			"-factomdcert":     complete.PredictFiles("*"),
			"-factomdpassword": complete.PredictAnything,
			"-factomdtls":      complete.PredictNothing,
			"-factomduser":     complete.PredictAnything,
			"-s":               complete.PredictAnything,
			"-w":               complete.PredictAnything,
			"-walletcert":      complete.PredictFiles("*"),
			"-walletpassword":  complete.PredictAnything,
			"-wallettls":       complete.PredictNothing,
			"-walletuser":      complete.PredictAnything,

			"-test.bench":                complete.PredictAnything,
			"-test.benchmem":             complete.PredictNothing,
			"-test.benchtime":            complete.PredictAnything,
			"-test.blockprofile":         complete.PredictFiles("*"),
			"-test.blockprofilerate":     complete.PredictAnything,
			"-test.count":                complete.PredictAnything,
			"-test.coverprofile":         complete.PredictFiles("*"),
			"-test.cpu":                  complete.PredictAnything,
			"-test.cpuprofile":           complete.PredictFiles("*"),
			"-test.failfast":             complete.PredictNothing,
			"-test.list":                 complete.PredictAnything,
			"-test.memprofile":           complete.PredictFiles("*"),
			"-test.memprofilerate":       complete.PredictAnything,
			"-test.mutexprofile":         complete.PredictAnything,
			"-test.mutexprofilefraction": complete.PredictAnything,
			"-test.outputdir":            complete.PredictDirs("*"),
			"-test.parallel":             complete.PredictAnything,
			"-test.run":                  complete.PredictAnything,
			"-test.short":                complete.PredictNothing,
			"-test.testlogfile":          complete.PredictFiles("*"),
			"-test.timeout":              complete.PredictAnything,
			"-test.trace":                complete.PredictFiles("*"),
			"-test.v":                    complete.PredictNothing,
		},
	}
	complete.New("factom-cli", cli).Run()
}

// -factomdcert string
// -factomdpassword string
// -factomdtls
// -factomduser string
// -s string
// -test.bench regexp
// -test.benchmem
// -test.benchtime d
// -test.blockprofile file
// -test.blockprofilerate rate
// -test.count n
// -test.coverprofile file
// -test.cpu list
// -test.cpuprofile file
// -test.failfast
// -test.list regexp
// -test.memprofile file
// -test.memprofilerate rate
// -test.mutexprofile string
// -test.mutexprofilefraction int
// -test.outputdir dir
// -test.parallel n
// -test.run regexp
// -test.short
// -test.testlogfile file
// -test.timeout d
// -test.trace file
// -test.v
// -w string
// -walletcert string
// -walletpassword string
// -wallettls
// -walletuser string

// addchain [-fq] [-n NAME1 -n NAME2 -h HEXNAME3 ] [-CET] ECADDRESS <STDIN>
// addentry [-fq] [-n NAME1 -h HEXNAME2 ...|-c CHAINID] [-e EXTID1 -e EXTID2 -x HEXEXTID ...] [-CET] ECADDRESS <STDIN>
// addtxecoutput [-rq] TXNAME ADDRESS AMOUNT
// addtxfee [-q] TXNAME ADDRESS
// addtxinput [-q] TXNAME ADDRESS AMOUNT
// addtxoutput [-rq] TXNAME ADDRESS AMOUNT
// backupwallet
// balance [-r] ADDRESS
// buyec [-fqrT] FCTADDRESS ECADDRESS ECAMOUNT
// composechain [-f] [-n NAME1 -n NAME2 -h HEXNAME3 ] ECADDRESS <STDIN>
// composeentry [-f] [-n NAME1 -h HEXNAME2 ...|-c CHAINID]  [-e EXTID1 -e EXTID2 -x HEXEXTID ...] ECADDRESS <STDIN>
// composetx TXNAME
// ecrate
// exportaddresses
// get allentries|chainhead|dblock|eblock|entry|firstentry|head|heights|walletheight|pendingentries|pendingtransactions|raw|dbheight|abheight|fbheight|ecbheight
// get abheight HEIGHT -r (to suppress Raw Data)
// get allentries [-n NAME1 -h HEXNAME2 ...|CHAINID] [-E]
// get chainhead [-n NAME1 -h HEXNAME2 ...|CHAINID] [-K]
// get dbheight HEIGHT -r (to suppress Raw Data)
// get dblock KEYMR
// get eblock KEYMR
// get ecbheight HEIGHT -r (to suppress Raw Data)
// get entry HASH
// get fbheight HEIGHT -r (to suppress Raw Data)
// get firstentry [-n NAME1 -h HEXNAME2 ...|CHAINID] [-E]
// get head [-K]
// get heights
// get pendingentries [-E]
// get pendingtransactions [-T]
// get raw HASH
// get walletheight
// importaddress ADDRESS [ADDRESS...]
// importkoinify '12WORDS'
// listaddresses
// listtxs [address|all|id|name|tmp|range]
// listtxs address [-T] ECADDRESS|FCTADDRESS
// listtxs [all] [-T]
// listtxs id TXID
// listtxs name TXNAME
// listtxs range [-T] START END
// listtxs tmp
// newecaddress
// newfctaddress
// newtx [-q] TXNAME
// properties
// receipt ENTRYHASH
// rmaddress ADDRESS
// rmtx TXNAME
// sendfct [-fqrT] FROMADDRESS TOADDRESS AMOUNT
// sendtx [-fqT] TXNAME
// signtx [-fqT] TXNAME
// status TxID|FullTx
// subtxfee [-q] TXNAME ADDRESS
