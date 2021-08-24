// Copyright (c) 2021 deadc0de6

package transport

// Transport the interface to transports
type Transport interface {
	Execute(string) (string, string, error)
	Copy(string, string, string) error
	Mkdir(string) error
	Close()
}
