package main

type adapter interface {
	readLast() (*entry, error)
	readAll() ([]*entry, error)
	write(*entry) error
}
