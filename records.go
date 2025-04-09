package main

import (
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

// Add domain record to the database
func AddDomain(domainProvider string, domainName string, purchasedDate time.Time, expirationDate time.Time, isExpired bool, autoRenew bool, isLocked bool, customDNS bool) error {
	domainsCollection, err := app.FindCollectionByNameOrId("Domains")
	if err != nil {
		return err
	}

	record, err := app.FindFirstRecordByData("Domains", "Name", domainName)
	if err != nil {
		record = core.NewRecord(domainsCollection)
		record.Set("Name", domainName)
		record.Set("Healthy", true)
	}

	// Set record details
	record.Set("Purchased_Date", purchasedDate)
	record.Set("Expiration_Date", expirationDate)
	record.Set("Is_Expired", isExpired)
	record.Set("Auto_Renew", autoRenew)
	record.Set("Is_Locked", isLocked)
	record.Set("Domain_Provider", domainProvider)
	record.Set("Custom_DNS", customDNS)
	err = app.Save(record)
	if err != nil {
		return err
	}

	return nil
}

// Add a new record to a domain
func AddDomainRecord(domainName string, recordName string, recordType string, address string) error {
	// Get domain record id
	domainRecord, err := app.FindFirstRecordByData("Domains", "Name", domainName)
	if err != nil {
		return err
	}

	recordsCollection, err := app.FindCollectionByNameOrId("Domain_Records")
	if err != nil {
		return err
	}

	record := core.NewRecord(recordsCollection)
	record.Set("Domain", domainRecord.Get("id").(string))
	record.Set("Record_Name", recordName)
	record.Set("Record_Type", recordType)
	record.Set("Address", address)
	err = app.Save(record)
	if err != nil {
		return err
	}

	return nil
}

// Delete all records for a domain
func DeleteAllDomainRecords(domainName string) error {
	// Get domain record id
	domainRecord, err := app.FindFirstRecordByData("Domains", "Name", domainName)
	if err != nil {
		return err
	}

	records, _ := app.FindAllRecords("Domain_Records",
		dbx.NewExp("Domain = {:domain}", dbx.Params{"domain": domainRecord.Get("id").(string)}),
	)

	if len(records) > 0 {
		for _, record := range records {
			err = app.Delete(record)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
