package main

import (
	"encoding/json"

	"github.com/lum8rjack/redcompass/services/namecheap"
)

func AddNamecheapCron(jobID string, cron string) {
	msg := "CRON:Namecheap"

	if cron == "" {
		app.Logger().Error(msg+" error cron is empty", "method", "AddNamecheapCron")
		return
	}

	if jobID == "" {
		app.Logger().Error(msg+" error jobID is empty", "method", "AddNamecheapCron")
		return
	}

	app.Cron().MustAdd(jobID, cron, func() {
		app.Logger().Info(msg+" started", "status", "started")
		record, err := app.FindFirstRecordByData("Services", "Provider", "Namecheap")
		if err != nil {
			app.Logger().Error(msg, "function", "app.FindFirstRecordByData", "error", err.Error())
			return
		}

		recordSettings := record.GetString("Settings")
		if recordSettings == "" {
			app.Logger().Error(msg, "function", "record.GetString", "error", "no settings")
			return
		}

		var settings namecheap.Settings
		err = json.Unmarshal([]byte(recordSettings), &settings)
		if err != nil {
			app.Logger().Error(msg, "function", "json.Unmarshal", "error", err.Error())
			return
		}

		if settings.ApiKey == "" || settings.IP == "" || settings.Username == "" {
			app.Logger().Error(msg, "error", "invalid settings")
			return
		}

		client := namecheap.NewClient(settings.Username, settings.ApiKey, settings.IP)

		// Get domains
		response, err := client.NamecheapGetDomains()
		if err != nil {
			app.Logger().Error(msg, "function", "namecheap.NamecheapGetDomains", "error", err.Error())
			return
		}

		for _, d := range response {
			err = AddDomain("Namecheap", d.Name, d.Created, d.Expires, d.IsExpired, d.AutoRenew, d.IsLocked, !d.IsOurDNS)
			if err != nil {
				app.Logger().Error(msg, "function", "AddDomainRecord", "domain", d.Name, "error", err.Error())
				continue
			}

			// Delete all existing records for this domain
			err = DeleteAllDomainRecords(d.Name)
			if err != nil {
				app.Logger().Error(msg, "function", "DeleteAllDomainRecords", "domain", d.Name, "error", err.Error())
				continue
			}

			// Get records for each domain if:
			// - The domain is using our DNS
			// - The domain is not expired
			// - The domain is not locked
			if d.IsOurDNS && !d.IsExpired && !d.IsLocked {
				ncRecords, err := client.NamecheapGetDomainRecords(d.Name)
				if err != nil {
					app.Logger().Error(msg, "function", "NamecheapGetDomainRecords", "domain", d.Name, "error", err.Error())
					continue
				}

				// Add new records
				for _, r := range ncRecords {
					err = AddDomainRecord(d.Name, r.Name, r.Type, r.Address)
					if err != nil {
						app.Logger().Error(msg, "function", "AddDomainRecord", "domain", d.Name, "error", err.Error())
						continue
					}
				}
			} else {
				app.Logger().Debug(msg+" skipped record retrieval", "status", "skipped", "domain", d.Name, "isOurDNS", d.IsOurDNS, "isExpired", d.IsExpired, "isLocked", d.IsLocked)
			}
		}

		app.Logger().Info(msg+" completed", "status", "completed")
	})
}

// Remove the Namecheap cron
func RemoveNamecheapCron(jobID string) {
	if jobID == "" {
		app.Logger().Error("CRON:Namecheap error jobID is empty", "method", "RemoveNamecheapCron")
		return
	}

	app.Cron().Remove(jobID)
}
