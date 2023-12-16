package main

type adapter interface {
	readLast() (*entry, error)
	write(*entry) error
}
