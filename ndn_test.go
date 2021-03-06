package ndn

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/go-ndn/tlv"
)

var (
	interest = &Interest{Name: NewName("/hello")}
	data     = &Data{Name: NewName("/hello")}

	discard = tlv.NewWriter(ioutil.Discard)
)

func BenchmarkDataEncodeRSA(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := SignData(rsaKey, data)
		if err != nil {
			b.Fatal(err)
		}
		err = data.WriteTo(discard)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDataEncodeECDSA(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := SignData(ecdsaKey, data)
		if err != nil {
			b.Fatal(err)
		}
		err = data.WriteTo(discard)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDataEncodeHMAC(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := SignData(hmacKey, data)
		if err != nil {
			b.Fatal(err)
		}
		err = data.WriteTo(discard)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDataEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := data.WriteTo(discard)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDataDecode(b *testing.B) {
	buf := new(bytes.Buffer)
	data.WriteTo(tlv.NewWriter(buf))
	seeker := bytes.NewReader(buf.Bytes())
	r := tlv.NewReader(seeker)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		seeker.Seek(0, 0)
		err := new(Data).ReadFrom(r)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkInterestEncode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := interest.WriteTo(discard)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkInterestDecode(b *testing.B) {
	buf := new(bytes.Buffer)
	interest.WriteTo(tlv.NewWriter(buf))
	seeker := bytes.NewReader(buf.Bytes())
	r := tlv.NewReader(seeker)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		seeker.Seek(0, 0)
		err := new(Interest).ReadFrom(r)
		if err != nil {
			b.Fatal(err)
		}
	}
}
