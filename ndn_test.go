package ndn

import (
	"bytes"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
)

var (
	key = &rsa.PrivateKey{
		PublicKey: rsa.PublicKey{
			N: fromBase10("14314132931241006650998084889274020608918049032671858325988396851334124245188214251956198731333464217832226406088020736932173064754214329009979944037640912127943488972644697423190955557435910767690712778463524983667852819010259499695177313115447116110358524558307947613422897787329221478860907963827160223559690523660574329011927531289655711860504630573766609239332569210831325633840174683944553667352219670930408593321661375473885147973879086994006440025257225431977751512374815915392249179976902953721486040787792801849818254465486633791826766873076617116727073077821584676715609985777563958286637185868165868520557"),
			E: 3,
		},
		D: fromBase10("9542755287494004433998723259516013739278699355114572217325597900889416163458809501304132487555642811888150937392013824621448709836142886006653296025093941418628992648429798282127303704957273845127141852309016655778568546006839666463451542076964744073572349705538631742281931858219480985907271975884773482372966847639853897890615456605598071088189838676728836833012254065983259638538107719766738032720239892094196108713378822882383694456030043492571063441943847195939549773271694647657549658603365629458610273821292232646334717612674519997533901052790334279661754176490593041941863932308687197618671528035670452762731"),
		Primes: []*big.Int{
			fromBase10("130903255182996722426771613606077755295583329135067340152947172868415809027537376306193179624298874215608270802054347609836776473930072411958753044562214537013874103802006369634761074377213995983876788718033850153719421695468704276694983032644416930879093914927146648402139231293035971427838068945045019075433"),
			fromBase10("109348945610485453577574767652527472924289229538286649661240938988020367005475727988253438647560958573506159449538793540472829815903949343191091817779240101054552748665267574271163617694640513549693841337820602726596756351006149518830932261246698766355347898158548465400674856021497190430791824869615170301029"),
		},
	}
	byteRSA      = []byte{0x6, 0xfd, 0x1, 0x35, 0x7, 0x10, 0x8, 0x6, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x8, 0x6, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x14, 0xf, 0x18, 0x1, 0x2, 0x19, 0x1, 0x3, 0x1a, 0x7, 0x8, 0x5, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x15, 0x3, 0x1, 0x2, 0x3, 0x16, 0x7, 0x1b, 0x1, 0x1, 0x1c, 0x2, 0x7, 0x0, 0x17, 0xfd, 0x1, 0x0, 0x44, 0x0, 0xe1, 0x2a, 0x16, 0x8, 0x45, 0xf, 0xb7, 0x38, 0xfe, 0x8f, 0x7a, 0x2a, 0xc7, 0x79, 0x72, 0xbc, 0x20, 0xfb, 0x70, 0x5, 0xdb, 0xe, 0xf1, 0x35, 0xd1, 0x93, 0xc6, 0xf, 0x9, 0x89, 0xd6, 0xa6, 0x97, 0x27, 0x5c, 0x6b, 0x7c, 0x11, 0x54, 0x5b, 0x48, 0x96, 0x56, 0x4c, 0x5d, 0x42, 0x1e, 0xe5, 0x3f, 0xc4, 0xea, 0xf4, 0xff, 0x7c, 0xe, 0x97, 0x9c, 0xc, 0x8b, 0xa2, 0xa, 0x33, 0x4f, 0x5, 0x48, 0x69, 0x3, 0x99, 0x63, 0xa0, 0xaa, 0xba, 0x8b, 0xf8, 0x15, 0xd2, 0x82, 0xc9, 0x89, 0x2, 0xe1, 0x3c, 0xe, 0x5b, 0x80, 0x98, 0x17, 0x8b, 0x45, 0x4b, 0x57, 0xa1, 0xcf, 0xc6, 0x26, 0x90, 0xbd, 0xae, 0x56, 0xc0, 0xc4, 0x59, 0x28, 0x7f, 0xa4, 0x99, 0xad, 0xed, 0x55, 0xa4, 0xbc, 0x33, 0x2b, 0x82, 0x17, 0xfe, 0xd5, 0x0, 0xe1, 0xa4, 0xa, 0x89, 0x8, 0xed, 0x8b, 0x24, 0xbb, 0xc0, 0x9b, 0x43, 0x97, 0x62, 0xd3, 0x5, 0x38, 0x27, 0xfb, 0x4e, 0x39, 0x99, 0xaf, 0x15, 0xe, 0x12, 0xf, 0x5a, 0x3b, 0x26, 0x87, 0xfe, 0x83, 0x2e, 0x89, 0x6e, 0xa3, 0x70, 0x2f, 0x2e, 0x6e, 0x9a, 0x45, 0x60, 0x46, 0x4e, 0x2a, 0x54, 0x53, 0xf0, 0xa7, 0x56, 0xb6, 0x22, 0xac, 0x6b, 0x8a, 0x7a, 0x6, 0xeb, 0x8d, 0xdd, 0xde, 0xbe, 0xd7, 0x84, 0x68, 0x69, 0x1a, 0x2e, 0x18, 0xe2, 0x47, 0xcd, 0xbe, 0xc9, 0x60, 0x33, 0xaa, 0xd5, 0x10, 0x44, 0x54, 0x57, 0xae, 0x1b, 0x85, 0xe9, 0xc0, 0xff, 0xdf, 0x86, 0x33, 0xea, 0x41, 0x1d, 0xf0, 0xe7, 0x96, 0x11, 0xe9, 0x8b, 0x70, 0x5c, 0xee, 0x53, 0x7e, 0xd7, 0xb6, 0x68, 0x2b, 0x47, 0xe0, 0xf0, 0xea, 0x20, 0xa3, 0x81, 0xb8, 0x3d, 0xe8, 0xbb, 0x43, 0xc1, 0x2b, 0x8, 0x27, 0x6b, 0xfe, 0x85, 0xb2, 0xae, 0x63, 0x87, 0x8d, 0x9a}
	byteSHA256   = []byte{0x6, 0x4f, 0x7, 0x10, 0x8, 0x6, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x8, 0x6, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x14, 0xf, 0x18, 0x1, 0x2, 0x19, 0x1, 0x3, 0x1a, 0x7, 0x8, 0x5, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x15, 0x3, 0x1, 0x2, 0x3, 0x16, 0x3, 0x1b, 0x1, 0x0, 0x17, 0x20, 0x29, 0x5, 0xa7, 0x90, 0x15, 0xfb, 0xd7, 0xe5, 0x66, 0xa1, 0x52, 0x11, 0xf0, 0x2c, 0xbb, 0x4d, 0xb8, 0xc0, 0x8a, 0x9e, 0x5, 0xca, 0x47, 0x82, 0xee, 0x3b, 0x2a, 0xbf, 0x20, 0xc0, 0x73, 0xdc}
	byteInterest = []byte{0x5, 0x2b, 0x7, 0x11, 0x8, 0x8, 0x66, 0x61, 0x63, 0x65, 0x62, 0x6f, 0x6f, 0x6b, 0x8, 0x5, 0x75, 0x73, 0x65, 0x72, 0x73, 0x9, 0xb, 0xd, 0x1, 0x3, 0xe, 0x1, 0x5, 0x11, 0x1, 0x4, 0x12, 0x0, 0xa, 0x3, 0x1, 0x2, 0x3, 0xb, 0x1, 0x8, 0xc, 0x1, 0x9}
)

func fromBase10(base10 string) *big.Int {
	i, ok := new(big.Int).SetString(base10, 10)
	if !ok {
		panic("bad number: " + base10)
	}
	return i
}

func BenchmarkDataSHA256Encode(b *testing.B) {
	data := NewData("/google/search")
	for n := 0; n < b.N; n++ {
		_, err := data.Encode()
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkDataSHA256Decode(b *testing.B) {
	data := NewData("")
	for n := 0; n < b.N; n++ {
		err := data.Decode(byteSHA256)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkDataRSAEncode(b *testing.B) {
	RSAPrivateKey = key
	data := NewData("/google/search")
	data.Signature.Type = 1
	for n := 0; n < b.N; n++ {
		_, err := data.Encode()
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkDataRSADecode(b *testing.B) {
	RSAPrivateKey = key
	data := NewData("")
	for n := 0; n < b.N; n++ {
		err := data.Decode(byteRSA)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkInterestEncode(b *testing.B) {
	interest := NewInterest("/google/search")
	for n := 0; n < b.N; n++ {
		_, err := interest.Encode()
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func BenchmarkInterestDecode(b *testing.B) {
	interest := NewInterest("")
	for n := 0; n < b.N; n++ {
		err := interest.Decode(byteInterest)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

func TestData(t *testing.T) {
	RSAPrivateKey = key
	data := NewData("/google/search")
	data.MetaInfo.ContentType = 2
	data.MetaInfo.FreshnessPeriod = 3
	data.MetaInfo.FinalBlockId = []byte("hello")
	data.Content = []byte{0x1, 0x2, 0x3}

	data.Signature.Type = 0

	b, err := data.Encode()
	if err != nil {
		t.Error(err)
	}

	data_decode := Data{}
	err = data_decode.Decode(b)
	if err != nil {
		t.Error(err)
	}
	// name order changes
	data.Name = nil
	data_decode.Name = nil
	data.Signature.Value = nil
	data_decode.Signature.Value = nil
	data.Signature.Info = nil
	data_decode.Signature.Info = nil
	b1, err := json.Marshal(data)
	if err != nil {
		t.Error(err)
	}
	b2, err := json.Marshal(data_decode)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(b1, b2) {
		t.Errorf("expected %v, got %v", b1, b2)
	}
}

func TestInterest(t *testing.T) {
	interest := NewInterest("/facebook/users")

	interest.Selectors.MinSuffixComponents = 3
	interest.Selectors.MaxSuffixComponents = 5
	interest.Selectors.ChildSelector = 4
	interest.Selectors.MustBeFresh = true
	interest.Scope = 8
	interest.InterestLifeTime = 9
	interest.Nonce = []byte{0x1, 0x2, 0x3}
	b, err := interest.Encode()
	if err != nil {
		t.Error(err)
	}

	interest_decode := Interest{}
	err = interest_decode.Decode(b)
	if err != nil {
		t.Error(err)
	}
	// name order changes
	interest.Name = nil
	interest_decode.Name = nil
	b1, err := json.Marshal(interest)
	if err != nil {
		t.Error(err)
	}
	b2, err := json.Marshal(interest_decode)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(b1, b2) {
		t.Errorf("expected %v, got %v", b1, b2)
	}
}
