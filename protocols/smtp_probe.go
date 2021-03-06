// SMTP Tester
//
// The SMTP tester connects to a remote host and ensures that a response
// is received that looks like an SMTP-server banner.
//
// This test is invoked via input like so:
//
//    host.example.com must run smtp [with port 25]
//

package protocols

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/skx/overseer/test"
)

// SMTPTest is our object
type SMTPTest struct {
}

// Arguments returns the names of arguments which this protocol-test
// understands, along with corresponding regular-expressions to validate
// their values.
func (s *SMTPTest) Arguments() map[string]string {
	known := map[string]string{
		"port": "^[0-9]+$",
	}
	return known
}

// Example returns sample usage-instructions for self-documentation purposes.
func (s *SMTPTest) Example() string {
	str := `
SMTP Tester
-----------
 The SMTP tester connects to a remote host and ensures that a response
 is received that looks like an SMTP-server banner.

 This test is invoked via input like so:

    host.example.com must run smtp
`
	return str
}

// RunTest is the part of our API which is invoked to actually execute a
// test against the given target.
//
// In this case we make a TCP connection, defaulting to port 25, and
// look for a response which appears to be an SMTP-server.
func (s *SMTPTest) RunTest(tst test.Test, target string, opts test.TestOptions) error {
	var err error

	//
	// The default port to connect to.
	//
	port := 25

	//
	// If the user specified a different port update to use it.
	//
	if tst.Arguments["port"] != "" {
		port, err = strconv.Atoi(tst.Arguments["port"])
		if err != nil {
			return err
		}
	}

	//
	// Set an explicit timeout
	//
	d := net.Dialer{Timeout: opts.Timeout}

	//
	// Default to connecting to an IPv4-address
	//
	address := fmt.Sprintf("%s:%d", target, port)

	//
	// If we find a ":" we know it is an IPv6 address though
	//
	if strings.Contains(target, ":") {
		address = fmt.Sprintf("[%s]:%d", target, port)
	}

	//
	// Make the TCP connection.
	//
	conn, err := d.Dial("tcp", address)
	if err != nil {
		return err
	}

	//
	// Read the banner.
	//
	banner, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return err
	}
	conn.Close()

	if !strings.Contains(banner, "SMTP") {
		return errors.New("Banner doesn't look like an SMTP server")
	}

	return nil
}

//
// Register our protocol-tester.
//
func init() {
	Register("smtp", func() ProtocolTest {
		return &SMTPTest{}
	})
}
