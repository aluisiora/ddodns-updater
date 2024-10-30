package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func guessDomain(record string) string {
	firstDot := strings.Index(record, ".")
	if firstDot == -1 {
		return record
	}

	return record[firstDot+1:]
}

func main() {
	token := os.Getenv("DO_TOKEN")

	if token == "" {
		log.Fatal("no token set")
	}

	record := os.Getenv("DNS_RECORD")

	if record == "" {
		log.Fatal("no dns record set")
	}

	updater := Updater{
		DigitalOceanDNS: NewDigitalOceanDNS(token),
		Record:          record,
		Domain:          guessDomain(record),
		PublicIPFetchers: []PublicIPFetcher{
			NewIpify(),
			NewIpapi(),
		},
	}

	for {
		go func() {
            log.Print("running updater")
			err := updater.UpdatePublicIP()
            if err != nil {
                log.Print(fmt.Errorf("error updating public ip [%w]", err))
            }
		}()
		time.Sleep(time.Minute * 3)
	}
}
