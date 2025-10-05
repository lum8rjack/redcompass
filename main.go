package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	_ "github.com/lum8rjack/redcompass/migrations"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/jsvm"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/spf13/cobra"
)

//go:embed frontend/dist/*
var frontend embed.FS

var (
	Version string = "2025.10.1"
	app     *pocketbase.PocketBase
	fsys    fs.FS
)

func RouteRoot(e *core.ServeEvent) error {
	dist := http.FileServer(http.FS(fsys))

	// Serve static assets directly
	e.Router.GET("/assets/", apis.WrapStdHandler(dist))
	e.Router.GET("/favicon.ico", apis.WrapStdHandler(dist))

	// For all other routes, always serve index.html
	e.Router.GET("/", func(c *core.RequestEvent) error {
		indexHTML, err := fs.ReadFile(fsys, "index.html")
		if err != nil {
			return err
		}
		c.Response.Header().Set("Content-Type", "text/html")
		_, err = c.Response.Write(indexHTML)
		return err
	})

	return e.Next()
}

func main() {
	app = pocketbase.NewWithConfig(pocketbase.Config{
		HideStartBanner: false,
	})

	app.RootCmd.Version = Version

	// Change default AppName on first boot
	if app.Settings().Meta.AppName == "Acme" {
		app.Settings().Meta.AppName = "RedCompass"
	}

	// Add a secret command to easily generate a 32 character secret that can be used for the encryption key
	app.RootCmd.AddCommand(&cobra.Command{
		Use:   "secret",
		Short: "Generate a 32 character secret you can use for the encryption key",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("secret: %s\n", security.RandomString(32))
		},
	})

	// Setup JavaScript enginge to use with hooks
	jsvm.MustRegister(app, jsvm.Config{
		HooksWatch: true,
	})

	// Setup Migration command
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{})

	// Setup hooks
	SetupHooks()

	// Frontend website
	var err error
	fsys, err = fs.Sub(frontend, "frontend/dist")
	if err != nil {
		log.Fatalf("Failed to create sub filesystem: %v", err)
		return
	}

	// Verify that index.html exists in the embedded filesystem
	_, err = fs.Stat(fsys, "index.html")
	if err != nil {
		log.Fatalf("Frontend files not found in embedded filesystem. Make sure frontend/dist exists and contains files: %v", err)
		return
	}

	// Add additional routes
	app.OnServe().BindFunc(RouteRoot)

	// When a user is created, set their role to "viewer" if its not already set
	app.OnRecordAfterCreateSuccess("users").BindFunc(func(e *core.RecordEvent) error {
		role := e.Record.GetString("role")
		if role == "" {
			e.Record.Set("role", "viewer")
			e.App.Save(e.Record)
		}
		return e.Next()
	})

	// Start the server
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
