package main

import (
	"crypto/md5"
	"crypto/rsa"
	"fmt"
	"io"
	"math/rand"

	"golang.org/x/crypto/md4"
)

func main() {
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
	io.WriteString(h5, "The fog is getting thicker!")
	io.WriteString(h5, "And Leon's getting laaarger!")
	fmt.Printf("%x", h5.Sum(nil))
}
