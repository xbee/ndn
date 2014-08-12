package ndn

import (
	"bufio"
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"math/big"
	"time"
)

var (
	SignKey   Key
	VerifyKey Key
)

type Key struct {
	Name       Name
	privateKey crypto.PrivateKey
}

func (this *Key) LocatorName() (name Name) {
	name.Components = append(this.Name.Components, []byte("KEY"), []byte("ID-CERT"))
	return
}

func (this *Key) Decode(pemData []byte) (err error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		err = errors.New("not pem data")
		return
	}
	this.Name.Set(block.Headers["NAME"])
	switch block.Type {
	case "RSA PRIVATE KEY":
		this.privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	case "ECDSA PRIVATE KEY":
		this.privateKey, err = x509.ParseECPrivateKey(block.Bytes)
	default:
		err = errors.New("unsupported key type")
	}
	return
}

func (this *Key) Encode() (pemData []byte, err error) {
	var b []byte
	var keyType string
	switch this.privateKey.(type) {
	case *rsa.PrivateKey:
		b = x509.MarshalPKCS1PrivateKey(this.privateKey.(*rsa.PrivateKey))
		keyType = "RSA PRIVATE KEY"
	case *ecdsa.PrivateKey:
		b, err = x509.MarshalECPrivateKey(this.privateKey.(*ecdsa.PrivateKey))
		if err != nil {
			return
		}
		keyType = "ECDSA PRIVATE KEY"
	default:
		err = errors.New("unsupported key type")
		return
	}
	pemData = pem.EncodeToMemory(&pem.Block{
		Type: keyType,
		Headers: map[string]string{
			"NAME": this.Name.String(),
		},
		Bytes: b,
	})
	return
}

var (
	oidRsa   = asn1.ObjectIdentifier{1, 2, 840, 113549, 1, 1, 1}
	oidEcdsa = asn1.ObjectIdentifier{1, 2, 840, 10045, 2, 1}
)

func (this *Key) EncodeCertificate() (raw []byte, err error) {
	var sigType uint64
	var publicKeyBytes []byte
	var oidSig asn1.ObjectIdentifier
	switch this.privateKey.(type) {
	case *rsa.PrivateKey:
		publicKeyBytes, err = asn1.Marshal(this.privateKey.(*rsa.PrivateKey).PublicKey)
		if err != nil {
			return
		}
		oidSig = oidRsa
		sigType = SignatureTypeSha256WithRsa
	case *ecdsa.PrivateKey:
		publicKeyBytes, err = asn1.Marshal(this.privateKey.(*ecdsa.PrivateKey).PublicKey)
		if err != nil {
			return
		}
		oidSig = oidEcdsa
		sigType = SignatureTypeSha256WithEcdsa
	default:
		err = errors.New("unsupported key type")
		return
	}

	d := Data{
		Name: this.LocatorName(),
		MetaInfo: MetaInfo{
			ContentType: 2, //key
		},
		SignatureInfo: SignatureInfo{
			SignatureType: sigType,
			KeyLocator: KeyLocator{
				Name: this.LocatorName(),
			},
		},
	}
	d.Content, err = asn1.Marshal(certificate{
		Validity: validity{
			NotBefore: time.Now(),
			NotAfter:  time.Date(2049, 12, 31, 23, 59, 59, 0, time.UTC), // end of asn.1
		},
		Subject: []pkix.AttributeTypeAndValue{{
			Type:  asn1.ObjectIdentifier{2, 5, 4, 41},
			Value: this.Name.String(),
		}},
		SubjectPubKeyInfo: subjectPubKeyInfo{
			AlgorithmIdentifier: pkix.AlgorithmIdentifier{
				Algorithm: oidSig,
				// This is a NULL parameters value which is technically
				// superfluous, but most other code includes it and, by
				// doing this, we match their public key hashes.
				Parameters: asn1.RawValue{
					Tag: 5,
				},
			},
			Bytes: asn1.BitString{
				Bytes:     publicKeyBytes,
				BitLength: 8 * len(publicKeyBytes),
			},
		},
	})
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)
	enc := base64.NewEncoder(base64.StdEncoding, buf)
	err = d.WriteTo(enc)
	if err != nil {
		return
	}
	enc.Close()
	for buf.Len() != 0 {
		raw = append(raw, buf.Next(64)...)
		raw = append(raw, 0xA)
	}
	return
}

func NewKey(name string, privateKey crypto.PrivateKey) (key Key, err error) {
	key.Name.Set(name)
	switch privateKey.(type) {
	case *rsa.PrivateKey:
	case *ecdsa.PrivateKey:
	default:
		err = errors.New("unsupported key type")
		return
	}
	key.privateKey = privateKey
	return
}

type certificate struct {
	Validity          validity
	Subject           []pkix.AttributeTypeAndValue
	SubjectPubKeyInfo subjectPubKeyInfo
}

type validity struct {
	NotBefore time.Time
	NotAfter  time.Time
}

type subjectPubKeyInfo struct {
	AlgorithmIdentifier pkix.AlgorithmIdentifier
	Bytes               asn1.BitString
}

func PrintCertificate(raw []byte) (err error) {
	// newline does not matter
	dec := base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer(raw))
	var d Data
	err = d.ReadFrom(bufio.NewReader(dec))
	if err != nil {
		return
	}
	cert := &certificate{}
	_, err = asn1.Unmarshal(d.Content, cert)
	if err != nil {
		return
	}
	Print(d, cert)
	return
}

type ecdsaSignature struct {
	R, S *big.Int
}

func (this *Key) Sign(digest []byte) (signature []byte, err error) {
	switch this.privateKey.(type) {
	case *rsa.PrivateKey:
		signature, err = rsa.SignPKCS1v15(rand.Reader, this.privateKey.(*rsa.PrivateKey), crypto.SHA256, digest)
	case *ecdsa.PrivateKey:
		var sig ecdsaSignature
		sig.R, sig.S, err = ecdsa.Sign(rand.Reader, this.privateKey.(*ecdsa.PrivateKey), digest)
		if err != nil {
			return
		}
		signature, err = asn1.Marshal(sig)
	default:
		err = errors.New("unsupported key type")
	}
	return
}

func (this *Key) Verify(digest, signature []byte) error {
	switch this.privateKey.(type) {
	case *rsa.PrivateKey:
		return rsa.VerifyPKCS1v15(&this.privateKey.(*rsa.PrivateKey).PublicKey, crypto.SHA256, digest, signature)
	case *ecdsa.PrivateKey:
		var sig ecdsaSignature
		_, err := asn1.Unmarshal(signature, &sig)
		if err != nil {
			return err
		}
		if ecdsa.Verify(&this.privateKey.(*ecdsa.PrivateKey).PublicKey, digest, sig.R, sig.S) {
			return nil
		} else {
			return errors.New("crypto/ecdsa: verification error")
		}
	default:
		return errors.New("unsupported key type")
	}
}
