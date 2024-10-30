package main

import "log"

type PublicIPFetcher interface {
	FetchPublicIPv4() (string, error)
}

type Updater struct {
	DigitalOceanDNS  DigitalOceanDNS
	Record           string
	Domain           string
	PublicIPFetchers []PublicIPFetcher
}

func (u *Updater) UpdatePublicIP() error {
	var err error
	var publicIPv4 string

	for _, fetcher := range u.PublicIPFetchers {
		publicIPv4, err = fetcher.FetchPublicIPv4()
		if err == nil {
			break
		}
	}

	if err != nil {
		return err
	}

	record, err := u.DigitalOceanDNS.FindDomainRecord(u.Domain, u.Record, "A")
	if err != nil {
		return err
	}

	if record.Data == publicIPv4 {
        log.Println("public ip "+publicIPv4+" not changed")

		return nil
	}

    respRecord, err := u.DigitalOceanDNS.UpdateDomainRecord(u.Domain, record, publicIPv4)
    if err != nil {
        return err
    }

    log.Println("record updated to new public ip: " + respRecord.Data)

	return err
}
