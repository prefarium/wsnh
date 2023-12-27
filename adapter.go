package main

import "wsnh/adapters"

type adapter interface {
	ReadLast() (*adapters.Entry, error)
	ReadAll() ([]*adapters.Entry, error)
	Write(*adapters.Entry) error
}
