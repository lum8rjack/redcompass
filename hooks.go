package main

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

func SetupHooks() {
	bootstrapHook()
	createHook()
	updateHook()
	deleteHook()
}

func bootstrapHook() {
	app.OnBootstrap().BindFunc(func(e *core.BootstrapEvent) error {
		if err := e.Next(); err != nil {
			return err
		}

		checkAllServices()
		return nil
	})
}

// On startup, loop through all services and add the cron job
func checkAllServices() {
	services, err := app.FindAllRecords("Services")
	if err != nil {
		app.Logger().Error("Error finding services", "error", err.Error())
		return
	}

	for _, service := range services {
		serviceName := service.GetString("Provider")
		if serviceName == "" {
			continue
		}

		msg := "CRON:" + serviceName + " startup hook"
		err := AddDomainsCronJob(service)
		if err != nil {
			app.Logger().Error(msg, "function", "AddCronJob", "error", err.Error())
			continue
		}
		app.Logger().Info(msg, "status", "created", "jobID", serviceName)
	}

}

func createHook() {
	app.OnRecordCreate("Services").BindFunc(func(e *core.RecordEvent) error {
		msg := "CRON:" + e.Record.GetString("Provider") + " create hook"
		err := AddDomainsCronJob(e.Record)
		if err != nil {
			app.Logger().Error(msg, "function", "AddCronJob", "error", err.Error())
			return e.Next()
		}
		app.Logger().Info(msg, "status", "created", "jobID", e.Record.GetString("Provider"))
		return e.Next()
	})
}

func updateHook() {
	// When a service is updated, remove the cron job and add a new and updated cron job
	app.OnRecordAfterUpdateSuccess("Services").BindFunc(func(e *core.RecordEvent) error {
		msg := "CRON:" + e.Record.GetString("Provider") + " update hook"

		err := RemoveCronJob(e.Record.GetString("Provider"))
		if err != nil {
			app.Logger().Error(msg, "function", "RemoveCronJob", "error", err.Error())
			return e.Next()
		}

		err = AddDomainsCronJob(e.Record)
		if err != nil {
			app.Logger().Error(msg, "function", "AddCronJob", "error", err.Error())
			return e.Next()
		}

		app.Logger().Info(msg, "status", "updated", "jobID", e.Record.GetString("Provider"))
		return e.Next()
	})

	// When a project is completed, unassign all domains from the project and set the last used project to the project id
	app.OnRecordAfterUpdateSuccess("Projects").BindFunc(func(e *core.RecordEvent) error {
		if e.Record.GetString("Completed") == "true" {
			msg := "PROJECT:" + e.Record.GetString("Name") + " project completed update hook"
			assignedDomains, err := app.FindAllRecords("Domains",
				dbx.NewExp("Assigned_Project = {:project}", dbx.Params{"project": e.Record.GetString("id")}),
			)

			if err != nil {
				app.Logger().Error(msg, "function", "FindAllRecords", "error", err.Error())
				return e.Next()
			}
			for _, domain := range assignedDomains {
				domain.Set("Assigned_Project", "")
				domain.Set("Last_Used", e.Record.GetString("id"))
				err = app.Save(domain)
				if err != nil {
					app.Logger().Error(msg, "projectID", e.Record.GetString("id"), "function", "Save", "error", err.Error())
					return e.Next()
				}
			}
			app.Logger().Info(msg, "projectID", e.Record.GetString("id"), "status", "completed", "note", "unassigned all domains from the project")

		}
		return e.Next()
	})
}

func deleteHook() {
	app.OnRecordDelete("Services").BindFunc(func(e *core.RecordEvent) error {
		msg := "CRON:" + e.Record.GetString("Provider") + " delete hook"
		err := RemoveCronJob(e.Record.GetString("Provider"))
		if err != nil {
			app.Logger().Error(msg, "function", "RemoveCronJob", "error", err.Error())
			return e.Next()
		}
		app.Logger().Info(msg, "status", "deleted", "jobID", e.Record.GetString("Provider"))

		return e.Next()
	})
}
