package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/miekg/dns"
)

type ClientConfig struct {
	Servers  []string // servers to use
	Search   []string // suffixes to append to local name
	Port     string   // what port to use
	Ndots    int      // number of dots in name to trigger absolute lookup
	Timeout  int      // seconds before giving up on packet
	Attempts int      // lost packets before giving up on server, not used in the package dns
}

func main() {
	test, _ := ClientConfigFromFile("/etc/resolv.conf")

	// fmt.Println(test.Attempts)
	// fmt.Println(test.Ndots)
	// fmt.Println(test.Port)
	// fmt.Println(test.Search)
	// fmt.Println(test.Servers)
	// fmt.Println(test.Timeout)

	m := new(dns.Msg)
	dnstype := "A"
	domain := "google.com"

	if dnstype == "A" {
		m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	}
	if dnstype == "TXT" {
		m.SetQuestion(dns.Fqdn(domain), dns.TypeTXT)
	}
	c := new(dns.Client)
	c.Timeout = time.Duration(5) * time.Second
	c.ReadTimeout = time.Duration(5) * time.Second
	c.DialTimeout = time.Duration(5) * time.Second

	response, _, err := c.Exchange(m, test.Servers[0]+":"+test.Port)
	if err != nil {
		fmt.Printf("Error %s: %s\n", domain, err)
	}

	if len(response.Answer) == 0 {
		fmt.Printf("%s %s %s %v - %d", domain, " Had no results at", test.Servers[0]+":"+test.Port, response.Answer, len(response.Answer))
	}

	fmt.Println(response)

	if err != nil {
		fmt.Println(err.Error())
	}
}

func ClientConfigFromFile(resolvconf string) (*ClientConfig, error) {
	file, err := os.Open(resolvconf)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ClientConfigFromReader(file)
}

// ClientConfigFromReader works like ClientConfigFromFile but takes an io.Reader as argument
func ClientConfigFromReader(resolvconf io.Reader) (*ClientConfig, error) {
	c := new(ClientConfig)
	scanner := bufio.NewScanner(resolvconf)
	c.Servers = make([]string, 0)
	c.Search = make([]string, 0)
	c.Port = "53"
	c.Ndots = 1
	c.Timeout = 5
	c.Attempts = 2

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		line := scanner.Text()
		f := strings.Fields(line)
		if len(f) < 1 {
			continue
		}
		switch f[0] {
		case "nameserver": // add one name server
			if len(f) > 1 {
				// One more check: make sure server name is
				// just an IP address.  Otherwise we need DNS
				// to look it up.
				name := f[1]
				c.Servers = append(c.Servers, name)
			}

		case "domain": // set search path to just this domain
			if len(f) > 1 {
				c.Search = make([]string, 1)
				c.Search[0] = f[1]
			} else {
				c.Search = make([]string, 0)
			}

		case "search": // set search path to given servers
			c.Search = cloneSlice(f[1:])

		case "options": // magic options
			for _, s := range f[1:] {
				switch {
				case len(s) >= 6 && s[:6] == "ndots:":
					n, _ := strconv.Atoi(s[6:])
					if n < 0 {
						n = 0
					} else if n > 15 {
						n = 15
					}
					c.Ndots = n
				case len(s) >= 8 && s[:8] == "timeout:":
					n, _ := strconv.Atoi(s[8:])
					if n < 1 {
						n = 1
					}
					c.Timeout = n
				case len(s) >= 9 && s[:9] == "attempts:":
					n, _ := strconv.Atoi(s[9:])
					if n < 1 {
						n = 1
					}
					c.Attempts = n
				case s == "rotate":
					/* not imp */
				}
			}
		}
	}
	return c, nil
}

func cloneSlice[E any, S ~[]E](s S) S {
	if s == nil {
		return nil
	}
	return append(S(nil), s...)
}
