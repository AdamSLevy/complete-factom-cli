// +build ignore

package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/FactomProject/go-simplejson"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"netki"
	"strings"
	"time"
)

const ecPrivKeyVersion = 1

var (
	oidNamedCurveP224 = asn1.ObjectIdentifier{1, 3, 132, 0, 33}
	oidNamedCurveP256 = asn1.ObjectIdentifier{1, 2, 840, 10045, 3, 1, 7}
	oidNamedCurveP384 = asn1.ObjectIdentifier{1, 3, 132, 0, 34}
	oidNamedCurveP521 = asn1.ObjectIdentifier{1, 3, 132, 0, 35}
)

func namedCurveFromOID(oid asn1.ObjectIdentifier) elliptic.Curve {
	switch {
	case oid.Equal(oidNamedCurveP224):
		return elliptic.P224()
	case oid.Equal(oidNamedCurveP256):
		return elliptic.P256()
	case oid.Equal(oidNamedCurveP384):
		return elliptic.P384()
	case oid.Equal(oidNamedCurveP521):
		return elliptic.P521()
	}
	return nil
}

type ecPrivateKey struct {
	Version       int
	PrivateKey    []byte
	NamedCurveOID asn1.ObjectIdentifier `asn1:"optional,explicit,tag:0"`
	PublicKey     asn1.BitString        `asn1:"optional,explicit,tag:1"`
}

func parseECPrivateKey(namedCurveOID *asn1.ObjectIdentifier, der []byte) (key *ecdsa.PrivateKey, err error) {
	var privKey ecPrivateKey
	if _, err := asn1.Unmarshal(der, &privKey); err != nil {
		return nil, errors.New("x509: failed to parse EC private key: " + err.Error())
	}
	if privKey.Version != ecPrivKeyVersion {
		return nil, fmt.Errorf("x509: unknown EC private key version %d", privKey.Version)
	}

	var curve elliptic.Curve
	if namedCurveOID != nil {
		curve = namedCurveFromOID(*namedCurveOID)
	} else {
		curve = namedCurveFromOID(privKey.NamedCurveOID)
	}
	if curve == nil {
		return nil, errors.New("x509: unknown elliptic curve")
	}

	k := new(big.Int).SetBytes(privKey.PrivateKey)
	if k.Cmp(curve.Params().N) >= 0 {
		return nil, errors.New("x509: invalid elliptic curve private key value")
	}
	priv := new(ecdsa.PrivateKey)
	priv.Curve = curve
	priv.D = k
	priv.X, priv.Y = curve.ScalarBaseMult(privKey.PrivateKey)

	return priv, nil
}

/*
KSK Control Working Example Using Hardcoded EC Keys
*/
func main() {
	// Instantiate Partner Signing Key
	partnerKeyRaw, _ := hex.DecodeString("3077020101042095abcdbf646efc799c9b08866c185dd53cbe6ac4ceeed4ed33e2d27641c23a77a00a06082a8648ce3d030107a14403420004b220745b4195fbe16b55d578d347295c6ca9ab9b42720f794ff70d2dab732bc55049f8de66c37248fed05faaca7b12ac0924d3fb6f8a67e5166a8430f1c860b4")
	partnerKey, _ := parseECPrivateKey(nil, partnerKeyRaw)

	// Instantiate User Key and Sign DER Encoded Public Key with Partner Signing Key
	userKeyRaw, _ := hex.DecodeString("307702010104208977d5312f492e7b04cbd8cf11334ed6ccb3ab1c530252a37dc6a1fad46f8b7ba00a06082a8648ce3d030107a1440342000472866325007154426e80e40039a6414a27db6c09b536f3bb79712020f505e1ff91daff344a73eae5b96535c6864277eeb389f1d1f632d0174b77b67b4627d6b2")
	userKey, _ := parseECPrivateKey(nil, userKeyRaw)

	// Hash (SHA256) User's Public Key (DER-format) for KeySigning + Sign
	userDER, _ := x509.MarshalPKIXPublicKey(&userKey.PublicKey)
	hasher := sha256.New()
	hasher.Write(userDER)
	userDERHash := hasher.Sum(nil)
	r, s, _ := ecdsa.Sign(rand.Reader, partnerKey, userDERHash)
	sig, _ := asn1.Marshal(EcdsaSig{r, s})

	// Initialize a partner for use with key signed control
	partner := netki.NewNetkiRemotePartner("http://localhost:5000", userKey, &partnerKey.PublicKey, sig)

	d := &Domain{DomainName: "mydomain.com"}
	wallet := &Wallet{Currency: "btc", WalletAddress: "1hsdflkjghsfdlkjghfsdlkgjh"}
	submitWallets := make([]Wallet, 0)
	submitWallets = append(wallets, *wallet)

	walletName := partner.CreateNewWalletName(*d, "supertest1", submitWallets, "ext_id")
	walletName.Save(partner)

	wallets, _ := partner.GetWalletNames(*d, "")
}
