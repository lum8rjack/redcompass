package main

import (
	"encoding/json"

	"github.com/lum8rjack/redcompass/services/namecheap"
)

func AddNamecheapCron(jobID string, cron string) {
	app.Cron().MustAdd(jobID, cron, func() { // "*/15 * * * *"
		msg := "Namecheap cron"

		app.Logger().Info(msg, "message", "Namecheap cron started")
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
				app.Logger().Error(msg, "function", "AddDomainRecord", "error", err.Error())
				continue
			}

			// Delete all existing records for this domain
			err = DeleteAllDomainRecords(d.Name)
			if err != nil {
				app.Logger().Error(msg, "function", "DeleteAllDomainRecords", "error", err.Error())
				continue
			}

			// Get records for each domain if:
			// - The domain is using our DNS
			// - The domain is not expired
			// - The domain is not locked
			if d.IsOurDNS && !d.IsExpired && !d.IsLocked {
				ncRecords, err := client.NamecheapGetDomainRecords(d.Name)
				if err != nil {
					app.Logger().Error(msg, "function", "NamecheapGetDomainRecords", "domain", d.Name)
					continue
				}

				// Add new records
				for _, r := range ncRecords {
					err = AddDomainRecord(d.Name, r.Name, r.Type, r.Address)
					if err != nil {
						app.Logger().Error(msg, "function", "AddDomainRecord", "error", err.Error())
						continue
					}
				}
			} else {
				app.Logger().Info(msg, "message", "using a custom DNS", "domain", d.Name)
			}
		}
	})
}

// Remove the Namecheap cron
func RemoveNamecheapCron(jobID string) {
	app.Cron().Remove(jobID)
}
