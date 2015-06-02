package ndn

import (
	"bytes"
	"os"
	"reflect"
	"testing"
)

var (
	rsaKey   = readKey("key/default.pri")
	ecdsaKey = readKey("key/ecdsa.pri")
)

func readKey(file string) Key {
	pem, err := os.Open(file)
	if err != nil {
		return nil
	}
	defer pem.Close()
	key, err := DecodePrivateKey(pem)
	if err != nil {
		return nil
	}
	return key
}

func TestPrivateKey(t *testing.T) {
	for _, key1 := range []Key{rsaKey, ecdsaKey} {
		buf := new(bytes.Buffer)
		err := EncodePrivateKey(key1, buf)
		if err != nil {
			t.Fatal(err)
		}

		key2, err := DecodePrivateKey(buf)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(key1, key2) {
			t.Fatal("not equal", key1, key2)
		}
	}
}

func TestCertificate(t *testing.T) {
	for _, key := range []Key{rsaKey, ecdsaKey} {
		buf := new(bytes.Buffer)
		err := EncodeCertificate(key, buf)
		if err != nil {
			t.Fatal(err)
		}

		_, err = DecodeCertificate(buf)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestSignVerify(t *testing.T) {
	d := new(Data)
	for _, key := range []Key{rsaKey, ecdsaKey} {
		err := SignData(key, d)
		if err != nil {
			t.Fatal(err)
		}
		err = key.Verify(d, d.SignatureValue)
		if err != nil {
			t.Fatal(err)
		}
	}
}
