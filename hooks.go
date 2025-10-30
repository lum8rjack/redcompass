package main

import (
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
