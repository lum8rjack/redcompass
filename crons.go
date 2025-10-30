package main

import (
	"errors"

	"github.com/lum8rjack/redcompass/services"
	"github.com/pocketbase/pocketbase/core"
)

// Add a cron job
func AddDomainsCronJob(record *core.Record) error {
	// Check if the record has a provider, settings, and cron
	if record.GetString("Provider") == "" {
		return errors.New("provider is empty")
	}

	if record.GetString("Settings") == "" {
		return errors.New("settings is empty")
	}

	if record.GetString("Cron") == "" {
		return errors.New("cron is empty")
	}

	jobID := record.GetString("Provider")
	cron := record.GetString("Cron")

	app.Cron().MustAdd(jobID, cron, func() {
		msg := "CRON:" + jobID + " cron job"
		app.Logger().Info(msg, "status", "started")

		// Get the service client
		service, err := services.NewService(jobID, record.GetString("Settings"))
		if err != nil {
			app.Logger().Error(msg, "function", "services.NewService", "error", err.Error())
			return
		}

		// Get the domains
		domains, err := service.GetDomains()
		if err != nil {
			app.Logger().Error(msg, "function", "service.GetDomains", "error", err.Error())
			return
		}

		// Loop through the domains
		for _, d := range domains {
			err = AddDomain(jobID, d.Name, d.Created, d.Expires, d.IsExpired, d.AutoRenew, d.IsLocked, !d.IsOurDNS)
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
				domainRecords, err := service.GetDomainRecords(d.Name)
				if err != nil {
					app.Logger().Error(msg, "function", "GetDomainRecords", "domain", d.Name, "error", err.Error())
					continue
				}

				// Add new records
				for _, r := range domainRecords {
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
		app.Logger().Info(msg, "status", "completed")
	})

	return nil
}

// Remove a cron job
func RemoveCronJob(jobID string) error {
	if jobID == "" {
		return errors.New("jobID is empty")
	}

	app.Cron().Remove(jobID)
	return nil
}
