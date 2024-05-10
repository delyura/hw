package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type Client struct {
	address    string
	connection net.Conn
	in         io.ReadCloser
	out        io.Writer
	timeout    time.Duration
}

func (cl *Client) Connect() (err error) {
	cl.connection, err = net.Dial("tcp", cl.address)
	if err != nil {
		return err
	}
	return nil
}

func (cl *Client) Close() (err error) {
	return cl.connection.Close()
}

func (cl *Client) Send() error {
	_, err := io.Copy(cl.connection, cl.in)
	return err
}

func (cl *Client) Receive() error {
	_, err := io.Copy(cl.out, cl.connection)
	return err
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}
