package internal

import "github.com/linn221/myfuzzer/requests"

func NewFuzzer(baseRequest *requests.Request, fuzzs []requests.FuzzFunc, antiFuzzs []requests.FuzzFunc) func() ([]*requests.Request, error) {
	var allRequests = make([]*requests.Request, 0)
	var f func(*requests.Request, int) error
	f = func(r *requests.Request, i int) error {
		if i == -1 {
			allRequests = append(allRequests, r)
			return nil
		}
		r1 := r.Clone()
		r2 := r.Clone()
		r3 := r.Clone()

		// not fuzzing
		err := f(r1, i-1)
		if err != nil {
			return err
		}
		// fuzzing
		fuzzs[i](r2)
		err = f(r2, i-1)
		if err != nil {
			return err
		}
		// anti fuzzing
		antiFuzzs[i](r3)
		err = f(r3, i-1)
		if err != nil {
			return err
		}
		return nil
	}

	return func() ([]*requests.Request, error) {
		err := f(baseRequest, len(fuzzs)-1)
		if err != nil {
			return nil, err
		}

		return allRequests, nil
	}
}
