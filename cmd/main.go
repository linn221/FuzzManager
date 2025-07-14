package main

import (
	"fmt"
	"log"

	"github.com/linn221/myfuzzer/internal"
	"github.com/linn221/myfuzzer/requests"
)

func main() {
	base := requests.Request{Base: "example.com?fuzztest"}
	fuzzs := []requests.FuzzFunc{
		func(r requests.Request) requests.Request {
			r.Base = r.Base + "&name=xxx"
			return r
		},
		func(r requests.Request) requests.Request {
			r.Base = r.Base + "&age=xxx"
			return r
		},
		func(r requests.Request) requests.Request {
			r.Base = r.Base + "&price=xxx"
			return r
		},
		func(r requests.Request) requests.Request {
			r.Base = r.Base + "&order=xxx"
			return r
		},
	}
	antiFuzzs := []requests.FuzzFunc{
		func(r requests.Request) requests.Request {
			r.Base = r.Base + "&name=product"
			return r
		},
		func(r requests.Request) requests.Request {
			r.Base = r.Base + "&age=20"
			return r
		},
		func(r requests.Request) requests.Request {
			r.Base = r.Base + "&price=2000"
			return r
		},
		func(r requests.Request) requests.Request {
			r.Base = r.Base + "&order=rate"
			return r
		},
	}
	fuzz := internal.NewFuzzer(base, fuzzs, antiFuzzs)
	requests, err := fuzz()
	if err != nil {
		log.Fatal(err.Error())
	}
	for _, r := range requests {
		fmt.Println(r.Base)
	}
	fmt.Println(len(requests))
}
