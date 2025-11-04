package main

import (
	"errors"
	"fmt"
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

type client struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (c *client) Connect() error {
	if c.in == nil || c.out == nil {
		return errors.New("need check in and out")
	}

	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	c.conn = conn

	return err
}

func (c *client) Close() error {
	if c.conn == nil {
		return errors.New("connection is nil")
	}
	return c.conn.Close()
}

func (c *client) Send() error {
	if _, err := io.Copy(c.conn, c.in); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (c *client) Receive() error {
	if _, err := io.Copy(c.out, c.conn); err != nil {
		return fmt.Errorf("failed to receive message: %w", err)
	}

	return nil
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
