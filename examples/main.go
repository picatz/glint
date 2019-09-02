package main

import (
	"archive/zip"
	"crypto/md5"
	"crypto/rsa"
	"crypto/tls"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/md4"
)

func unzip(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0777); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}

type Example struct {
	ReadTimeout time.Duration
}

type Example2 struct {
	Wiggle bool
}

func main() {
	var example = Example{
		ReadTimeout: 0,
	}

	var httpServer = http.Server{
		ReadTimeout: 0,
	}

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
