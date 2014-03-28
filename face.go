package ndn

import (
	//"fmt"
	"net"
	"net/url"
	"strings"
	"time"
)

/*
   Define ndn face
*/

type Face struct {
	*url.URL
	Id uint64
}

func NewFace(raw string) *Face {
	u, _ := url.Parse(raw)
	// assume tcp
	if len(u.Scheme) == 0 {
		u.Scheme = "tcp"
	}
	// assume port 6363
	if !strings.Contains(u.Host, ":") && (u.Scheme == "tcp" || u.Scheme == "udp") {
		u.Host += ":6363"
	}
	return &Face{
		URL: u,
	}
}

func readChunk(conn net.Conn) (b []byte, err error) {
	fixed := make([]byte, 1024)
	for {
		var n int
		n, err = conn.Read(fixed)
		if err != nil {
			return
		}
		b = append(b, fixed[:n]...)
		if n < len(fixed) {
			break
		}
	}
	return
}

func readData(conn net.Conn) (d *Data, err error) {
	d = &Data{}
	b, err := readChunk(conn)
	if err != nil {
		return
	}
	err = d.Decode(b)
	return
}

func readInterest(conn net.Conn) (i *Interest, err error) {
	i = &Interest{}
	b, err := readChunk(conn)
	if err != nil {
		return
	}
	err = i.Decode(b)
	return
}

func (this *Face) Dial(i *Interest) (d *Data, err error) {
	// interest encode
	ib, err := i.Encode()
	if err != nil {
		return
	}

	// dial
	conn, err := net.Dial(this.Scheme, this.Host)
	if err != nil {
		return
	}
	defer conn.Close()
	// write interest
	conn.Write(ib)
	if i.InterestLifeTime == 0 {
		// default timeout 10s
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	} else {
		// use interestLifeTime
		conn.SetReadDeadline(time.Now().Add(time.Duration(i.InterestLifeTime) * time.Millisecond))
	}

	d, err = readData(conn)
	return
}

func (this *Face) Listen(name string, callback func(*Interest) *Data) error {
	ln, err := net.Listen(this.Scheme, this.Host)
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go func(conn net.Conn) {
			defer conn.Close()
			// wait client for 10s if there is no response
			conn.SetReadDeadline(time.Now().Add(10 * time.Second))
			i, err := readInterest(conn)
			if err != nil {
				return
			}
			d := callback(i)

			if d != nil {
				b, err := d.Encode()
				if err != nil {
					return
				}
				conn.Write(b)
			}
		}(conn)
	}
	return nil
}

func (this *Face) Close() {

}
