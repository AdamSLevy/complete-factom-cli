package netki

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/hex"
	"fmt"
	"github.com/FactomProject/go-simplejson"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Define our Error Type
type NetkiError struct {
	ErrorString string
	Failures    []string
}

func (e NetkiError) Error() string {
	var buffer bytes.Buffer
	buffer.WriteString(e.ErrorString)
	if len(e.Failures) > 0 {
		buffer.WriteString(": ")
		buffer.WriteString(strings.Join(e.Failures, ", "))
	}
	return buffer.String()
}

// WalletNameLookup resolves an address from a netki address and currency.
func WalletNameLookup(uri, currency string) (string, error) {
	apimethod := "https://pubapi.netki.com/api/wallet_lookup"
	resp, err := http.Get(fmt.Sprintf("%s/%s/%s", apimethod, uri, currency))
	if err != nil {
		return "", err
	}

	j, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		return "", err
	} else if resp.StatusCode != 200 {
		return "", fmt.Errorf("Could not resolve netki address")
	}

	if msg, err := j.Get("message").String(); msg != "" {
		return "", fmt.Errorf(msg)
	} else if err != nil {
		return "", err
	}

	return j.Get("wallet_address").String()
}

type Partner struct {
	id, partnerName string
}

type Domain struct {
	DomainName        string
	Status            string
	DelegationStatus  bool
	DelegationMessage string
	WalletNameCount   int
	Namesevers        []string
	NextRollDate      time.Time
	PublicSigningKey  string
	DsRecords         []string
}

type Wallet struct {
	Currency, WalletAddress string
}

type WalletName struct {
	Id         string
	DomainName string
	Name       string
	Wallets    []Wallet
	ExternalId string
}

type NetkiRequest interface {
	ProcessRequest(partner *NetkiPartner, uri string, method string, bodyData string) (*simplejson.Json, error)
}

type NetkiRequester struct {
	HTTPClient *http.Client
}

type NetkiPartner struct {
	Requester     NetkiRequest
	PartnerId     string
	ApiKey        string
	ApiUrl        string
	UserKey       *ecdsa.PrivateKey
	KeySigningKey *ecdsa.PublicKey
	KeySignature  []byte
}

type EcdsaSig struct {
	R, S *big.Int
}

// Utility Functions
func urlEncode(text string) string {
	result, err := url.Parse(text)
	if err != nil {
		return ""
	}
	return result.String()
}

// Sign Request
func (n NetkiRequester) SignRequest(uri string, bodyData string, key *ecdsa.PrivateKey) (string, error) {
	h := sha256.New()
	r := big.NewInt(0)
	s := big.NewInt(0)

	h.Write([]byte(uri + bodyData))
	signDataHash := h.Sum(nil)

	r, s, err := ecdsa.Sign(rand.Reader, key, signDataHash)
	if err != nil {
		return "", &NetkiError{"Unable to Sign Data", make([]string, 0)}
	}

	sequence := EcdsaSig{r, s}
	encoding, _ := asn1.Marshal(sequence)
	return hex.EncodeToString(encoding), nil
}

// Generic Request Handling
func (n NetkiRequester) ProcessRequest(partner *NetkiPartner, uri string, method string, bodyData string) (*simplejson.Json, error) {
	var supported_methods = [...]string{"GET", "POST", "PUT", "DELETE"}
	var isSupportedMethod = false

	// make sure we have a supported HTTP method
	for _, val := range supported_methods {
		if val == method {
			isSupportedMethod = true
		}
	}

	if !isSupportedMethod {
		return &simplejson.Json{}, &NetkiError{fmt.Sprintf("Unsupported HTTP Method: %s", method), make([]string, 0)}
	}

	buf := new(bytes.Buffer)
	if bodyData != "" {
		_, err := buf.WriteString(bodyData)
		if err != nil {
			return &simplejson.Json{}, &NetkiError{"Unable to Write Request Data to Buffer", make([]string, 0)}
		}
	}

	// Create Our Request
	buffer := new(bytes.Buffer)
	buffer.WriteString(partner.ApiUrl)
	if !strings.HasSuffix(partner.ApiUrl, "/") {
		buffer.UnreadRune()
	}
	buffer.WriteString(uri)

	req, err := http.NewRequest(method, buffer.String(), buf)
	req.Header.Set("Content-Type", "application/json")
	if partner.PartnerId == "" && partner.UserKey != nil {
		sig, err := n.SignRequest(buffer.String(), bodyData, partner.UserKey)
		if err != nil {
			return &simplejson.Json{}, err
		}
		req.Header.Set("X-Identity", partner.GetUserPublicKey())
		req.Header.Set("X-Signature", sig)
		req.Header.Set("X-Partner-Key", partner.GetKeySigningKey())
		req.Header.Set("X-Partner-KeySig", hex.EncodeToString(partner.KeySignature))
	} else {
		req.Header.Set("X-Partner-ID", partner.PartnerId)
		req.Header.Set("Authorization", partner.ApiKey)
	}

	// See if we have an injected HTTPClient
	var client *http.Client
	if n.HTTPClient == nil {
		client = http.DefaultClient
	} else {
		client = n.HTTPClient
	}

	// Send Our Request
	resp, err := client.Do(req)
	if err != nil {
		return &simplejson.Json{}, &NetkiError{fmt.Sprintf("HTTP Request Failed: %s", err), make([]string, 0)}
	}

	// DELETE with 204 Response, Don't Care About Response Data
	if method == "DELETE" && resp.StatusCode == http.StatusNoContent {
		return &simplejson.Json{}, nil
	}

	// Validate Content-Type
	if resp.Header.Get("Content-Type") != "application/json" {
		return &simplejson.Json{}, &NetkiError{fmt.Sprintf("HTTP Response Contains Invalid Content-Type: %s", resp.Header.Get("Content-Type")), make([]string, 0)}
	}

	// Close the body when we're done with the function
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &simplejson.Json{}, &NetkiError{fmt.Sprintf("HTTP Body Read Failed: %s", err), make([]string, 0)}
	}

	// Get Our JSON Data
	js, err := simplejson.NewJson(body)
	if err != nil {
		return &simplejson.Json{}, &NetkiError{fmt.Sprintf("Error Retrieving JSON Data: %s", err), make([]string, 0)}
	}

	// Return message if success is false
	if !js.Get("success").MustBool(false) {
		errMsg := new(bytes.Buffer)
		errMsg.WriteString(js.Get("message").MustString())
		resp, _ := js.Get("failures").Array()
		if resp != nil {
			errMsg.WriteString(" [FAILURES: ")
			failures := make([]string, 0)
			for i := 0; i < len(js.Get("failures").MustArray()); i++ {
				failures = append(failures, js.Get("failures").GetIndex(i).Get("message").MustString())
			}
			errMsg.WriteString(strings.Join(failures, ", "))
			errMsg.WriteString("]")
		}
		return &simplejson.Json{}, &NetkiError{errMsg.String(), make([]string, 0)}
	}

	return js, nil

}

// Defined WalletName Methods
func (w WalletName) GetAddress(currency string) string {
	for _, wallet := range w.Wallets {
		if wallet.Currency == currency {
			return wallet.WalletAddress
		}
	}
	return ""
}

func (w WalletName) UsedCurrencies() []string {
	currencies := make([]string, 0)
	for _, wallet := range w.Wallets {
		currencies = append(currencies, wallet.Currency)
	}
	return currencies
}

func (w *WalletName) SetCurrencyAddress(currency string, address string) {

	for index, wallet := range w.Wallets {
		if wallet.Currency == currency {
			w.Wallets[index].WalletAddress = address
			return
		}
	}

	// Wallet Doesn't Already Exist, Create It & Extend the Array
	wallet := Wallet{Currency: currency, WalletAddress: address}
	w.Wallets = append(w.Wallets, wallet)
}

func (w *WalletName) RemoveCurrency(currency string) {
	for index, wallet := range w.Wallets {
		if wallet.Currency == currency {
			w.Wallets = append(w.Wallets[:index], w.Wallets[index+1:]...)
			return
		}
	}
}

func (w *WalletName) Save(partner *NetkiPartner) error {
	// Set Default HTTP Method
	httpMethod := "POST"

	wallets := make([]simplejson.Json, 0)
	for _, wallet := range w.Wallets {
		_w := simplejson.New()
		_w.Set("currency", wallet.Currency)
		_w.Set("wallet_address", wallet.WalletAddress)
		wallets = append(wallets, *_w)
	}

	d := simplejson.New()
	d.Set("domain_name", w.DomainName)
	d.Set("name", w.Name)
	d.Set("wallets", wallets)
	d.Set("external_id", w.ExternalId)
	if w.Id != "" {
		d.Set("id", w.Id)
		httpMethod = "PUT"
	}

	wnArray := make([]simplejson.Json, 0)
	wnArray = append(wnArray, *d)
	req := simplejson.New()
	req.Set("wallet_names", wnArray)

	jsondata, err := req.MarshalJSON()
	if err != nil {
		return &NetkiError{fmt.Sprintf("Unable to Marshall JSON Data: %s", err), make([]string, 0)}
	}

	resp, err := partner.Requester.ProcessRequest(partner, "/v1/partner/walletname", httpMethod, string(jsondata[:len(jsondata)]))
	if err != nil {
		return err
	}

	w.Id = resp.Get("wallet_names").GetIndex(0).Get("id").MustString()
	return nil
}

func (w WalletName) Delete(partner *NetkiPartner) error {

	if w.Id == "" {
		return &NetkiError{"WalletName has no ID! Cannot Delete!", make([]string, 0)}
	}

	d := simplejson.New()
	d.Set("domain_name", w.DomainName)
	d.Set("id", w.Id)

	wnArray := make([]simplejson.Json, 0)
	wnArray = append(wnArray, *d)
	req := simplejson.New()
	req.Set("wallet_names", wnArray)

	jsondata, err := req.MarshalJSON()
	if err != nil {
		return &NetkiError{fmt.Sprintf("Unable to Marshall JSON Data: %s", err), make([]string, 0)}
	}

	_, err = partner.Requester.ProcessRequest(partner, "/v1/partner/walletname", "DELETE", string(jsondata[:len(jsondata)]))
	if err != nil {
		return err
	}

	return nil
}

// Define NetkiPartner Utility methods
func (n NetkiPartner) GetUserPublicKey() string {
	derkey, err := x509.MarshalPKIXPublicKey(&n.UserKey.PublicKey)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(derkey)
}

func (n NetkiPartner) GetKeySigningKey() string {
	derkey, err := x509.MarshalPKIXPublicKey(n.KeySigningKey)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(derkey)
}

func (n NetkiPartner) SetUserKey(userKey *ecdsa.PrivateKey) {
	n.UserKey = userKey
}

func (n NetkiPartner) SetKeySigningKey(signingKey *ecdsa.PublicKey) {
	n.KeySigningKey = signingKey
}

func (n NetkiPartner) SetKeySignature(sig []byte) {
	n.KeySignature = sig
}

// Define NetkiPartner methods
func (n NetkiPartner) CreateNewPartner(partnerName string) (Partner, error) {
	uri := new(bytes.Buffer)
	uri.WriteString("/v1/admin/partner/")
	uri.WriteString(urlEncode(partnerName))

	resp, err := n.Requester.ProcessRequest(&n, uri.String(), "POST", "")
	if err != nil {
		return Partner{}, err
	}

	return Partner{resp.Get("partner").Get("id").MustString(), resp.Get("partner").Get("name").MustString()}, nil
}

func (n NetkiPartner) GetPartners() ([]Partner, error) {
	resp, err := n.Requester.ProcessRequest(&n, "/v1/admin/partner", "GET", "")
	if err != nil {
		return make([]Partner, 0), err
	}

	ps := resp.Get("partners").MustArray()

	partners := make([]Partner, 0)
	for i := 0; i < len(ps); i++ {
		intPartner := resp.Get("partners").GetIndex(i)
		newObj := Partner{intPartner.Get("id").MustString(), intPartner.Get("name").MustString()}
		partners = append(partners, newObj)
	}
	return partners, nil
}

func (n NetkiPartner) DeletePartner(partner Partner) error {
	uri := new(bytes.Buffer)
	uri.WriteString("/v1/admin/partner/")
	uri.WriteString(urlEncode(partner.partnerName))

	_, err := n.Requester.ProcessRequest(&n, uri.String(), "DELETE", "")
	if err != nil {
		return err
	}

	return nil
}

// Domain Handlers
func (n NetkiPartner) CreateNewDomain(domainName string, partner Partner) (Domain, error) {
	uri := new(bytes.Buffer)
	uri.WriteString("/v1/partner/domain/")
	uri.WriteString(urlEncode(domainName))

	d := simplejson.New()
	if partner.id != "" {
		d.Set("partner_id", partner.id)
	}

	jsondata, err := d.MarshalJSON()
	if err != nil {
		return Domain{}, err
	}

	resp, err := n.Requester.ProcessRequest(&n, uri.String(), "POST", string(jsondata[:len(jsondata)]))
	if err != nil {
		return Domain{}, err
	}

	returnDomain := Domain{}
	returnDomain.DomainName = resp.Get("domain_name").MustString()
	returnDomain.Status = resp.Get("status").MustString()
	returnDomain.Namesevers = resp.Get("nameservers").MustStringArray()

	return returnDomain, nil
}

func (n NetkiPartner) GetDomains() ([]Domain, error) {
	resp, err := n.Requester.ProcessRequest(&n, "/api/domain", "GET", "")
	if err != nil {
		return make([]Domain, 0), err
	}

	returnDomains := make([]Domain, 0)
	for i := 0; i < len(resp.Get("domains").MustArray()); i++ {
		newDomain := Domain{}
		newDomain.DomainName = resp.Get("domains").GetIndex(i).Get("domain_name").MustString()
		returnDomains = append(returnDomains, newDomain)
	}
	return returnDomains, nil
}

func (n NetkiPartner) GetDomainStatus(domain Domain) (returnDomain Domain, err error) {
	resp, err := n.Requester.ProcessRequest(&n, "/v1/partner/domain/"+urlEncode(domain.DomainName), "GET", "")
	if err != nil {
		return Domain{}, err
	}

	returnDomain = Domain{}
	returnDomain.DomainName = domain.DomainName
	returnDomain.Status = resp.Get("status").MustString()
	returnDomain.DelegationStatus = resp.Get("delegation_status").MustBool(false)
	returnDomain.DelegationMessage = resp.Get("delegation_message").MustString()
	returnDomain.WalletNameCount = resp.Get("wallet_name_count").MustInt()
	return returnDomain, nil
}

func (n NetkiPartner) GetDomainDnssec(domain Domain) (returnDomain Domain, err error) {
	resp, err := n.Requester.ProcessRequest(&n, "/v1/partner/domain/dnssec/"+urlEncode(domain.DomainName), "GET", "")
	if err != nil {
		return Domain{}, err
	}

	returnDomain = Domain{}
	returnDomain.DomainName = fmt.Sprint(domain.DomainName)
	returnDomain.NextRollDate, _ = time.Parse("2006-01-02T15:04:05.000Z", resp.Get("nextroll_date").MustString())
	returnDomain.DsRecords = resp.Get("ds_records").MustStringArray()
	returnDomain.PublicSigningKey = resp.Get("public_key_signing_key").MustString()

	return returnDomain, nil
}

func (n NetkiPartner) DeleteDomain(domain Domain) error {
	_, err := n.Requester.ProcessRequest(&n, "/v1/partner/domain/"+urlEncode(domain.DomainName), "DELETE", "")
	if err != nil {
		return err
	}
	return nil
}

func (n NetkiPartner) CreateNewWalletName(domain Domain, name string, wallets []Wallet, externalId string) WalletName {
	wn := WalletName{}
	wn.DomainName = domain.DomainName
	wn.Name = name
	wn.Wallets = wallets
	wn.ExternalId = externalId
	return wn
}

func (n NetkiPartner) GetWalletNames(domain Domain, externalId string) ([]WalletName, error) {
	uri := new(bytes.Buffer)
	uri.WriteString("/v1/partner/walletname")

	argSlice := make([]string, 0)
	if domain.DomainName != "" {
		argSlice = append(argSlice, "domain_name="+domain.DomainName)
	}
	if externalId != "" {
		argSlice = append(argSlice, "external_id="+externalId)
	}

	if len(argSlice) > 0 {
		uri.WriteString("?" + strings.Join(argSlice, "&"))
	}

	resp, err := n.Requester.ProcessRequest(&n, uri.String(), "GET", "")
	if err != nil {
		return make([]WalletName, 0), err
	}

	if resp.Get("wallet_name_count").MustInt(0) == 0 {
		return make([]WalletName, 0), nil
	}

	walletNames := make([]WalletName, 0)
	for i := 0; i < len(resp.Get("wallet_names").MustArray()); i++ {
		wn := resp.Get("wallet_names").GetIndex(i)
		wallets := make([]Wallet, 0)
		for j := 0; j < len(wn.Get("wallets").MustArray()); j++ {
			wallet := wn.Get("wallets").GetIndex(j)
			newWallet := Wallet{wallet.Get("currency").MustString(), wallet.Get("wallet_address").MustString()}
			wallets = append(wallets, newWallet)
		}
		newWalletName := WalletName{}
		newWalletName.Id = wn.Get("id").MustString()
		newWalletName.DomainName = wn.Get("domain_name").MustString()
		newWalletName.Name = wn.Get("name").MustString()
		newWalletName.ExternalId = wn.Get("external_id").MustString()
		newWalletName.Wallets = wallets
		walletNames = append(walletNames, newWalletName)
	}

	return walletNames, nil

}

// Constructor / NetkiPartner Factory
func NewNetkiPartner(partnerId string, apiKey string, apiUrl string) *NetkiPartner {
	return &NetkiPartner{Requester: new(NetkiRequester), PartnerId: partnerId, ApiKey: apiKey, ApiUrl: apiUrl}
}

func NewNetkiRemotePartner(apiUrl string, userKey *ecdsa.PrivateKey, keySigningKey *ecdsa.PublicKey, keySignature []byte) *NetkiPartner {
	return &NetkiPartner{Requester: new(NetkiRequester), ApiUrl: apiUrl, UserKey: userKey, KeySigningKey: keySigningKey, KeySignature: keySignature}
}
