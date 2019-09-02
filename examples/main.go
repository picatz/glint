package main

import (
	"crypto/md5"
	"crypto/rsa"
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"

	_ "unsafe"

	"golang.org/x/crypto/md4"
)

func main() {
	var example = Example{
		ReadTimeout: 0,
	}

	// and this
	var example2 = ExampleServer{}

	var httpServer = http.Server{}

	var wiggle = Ignored{
		Wiggle: true,
	}

	var wiggle2 = Ignored{}

	pk, err := rsa.GenerateKey(rand.New(rand.NewSource(0)), 2047)
	if err != nil {
		wrapperErr := fmt.Errorf("wrapped: %v", err.Error())
		wrappedErrUpcased := fmt.Errorf("UPPERCASE WRAP: %v", wrapperErr.Error())
		panic(wrappedErrUpcased.Error())
	}
	fmt.Println(pk)
	h4 := md4.New()
	io.WriteString(h4, fmt.Sprintf("%v", pk))
	fmt.Printf("%x", h4.Sum(nil))

	h5 := md5.New()
	io.WriteString(h5, "wiggle")
	fmt.Printf("%x", h5.Sum(nil))

	http.Handle("/", nil)
	http.HandleFunc("/", nil)

	var ex uint16

	_ = tls.Config{
		CipherSuites: []uint16{tls.TLS_AES_128_GCM_SHA256, ex},
	}

	if h5 == nil {
		// Listen on a random tcp port on all ipv4 interfaces
		listener, err := net.Listen("tcp", "0.0.0.0:0")
		if err != nil {
			panic(err)
		}
		err = listener.Close()
		if err != nil {
			panic(err)
		}
	}
}
