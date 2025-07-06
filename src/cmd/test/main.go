package main

import (
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/spirozh/timr/internal"
)

func main() {
	a := internal.NewApp()
	a.Run()
	defer func() {
		if err := a.Close(); err != nil {
			log.Fatal("closing app: ", err)
		}
	}()

	if err := testing(); err != nil {
		panic(err)
	}
}

func testing() error {
	r, err := http.Get("http://localhost:8080/selma")
	if err != nil {
		return err
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	expected := "Hello, Selma!\n"
	if string(b) != expected {
		return errors.New("expected '" + expected + "', got: " + string(b))
	}

	return nil
}
