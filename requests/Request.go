package requests

type Request struct {
	Base string
}

type FuzzFunc func(Request) Request
