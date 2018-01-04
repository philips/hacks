package main

import (
	"fmt"
	"github.com/likexian/whois-parser-go"

	"github.com/likexian/whois-go"
)

func main() {
	domains := []string{"google.com"}

	for _, domain := range domains {
		result, err := whois.Whois(domain)
		if err != nil {
			fmt.Println(err)
			continue
		}

		out, err := whois_parser.Parser(result)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// Print the domain expiration date
		fmt.Printf("%s: %v\n", domain, out.Registrar.ExpirationDate)
	}
}
