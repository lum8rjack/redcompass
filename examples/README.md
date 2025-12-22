# Examples

This directory contains example configurations and hooks for RedCompass.

> Some of the production deployment configurations and recommendations are adapted from the [PocketBase Going to Production](https://pocketbase.io/docs/going-to-production/) documentation.

## block-emails.pb.js

A PocketBase JavaScript hook that blocks all email sending functionality. This includes:
- Password reset emails
- Email verification
- Authentication alerts
- OTP emails
- Custom emails

Useful for development environments or internal deployments where email functionality is not needed or SMTP is not configured. The hook logs warnings when email sending is attempted.

## Caddyfile

An example Caddy server configuration for setting up a reverse proxy to RedCompass. Features:
- Disables Caddy admin interface for security
- Configures TLS with self-signed certificates
- Sets maximum request body size to 10MB
- Implements caching for static assets (1 hour cache duration)
- Configures proper host header forwarding
- Sets extended read timeout (360s) for long-running requests

Use this as a starting point for deploying RedCompass behind a Caddy reverse proxy in production.

## cleanup-domain-ideas.pb.js

A PocketBase Javascript cron that removes a domain idea if the domain has already been purchased. Simple script to clean up the domain ideas automatically.


## redcompass.service

A systemd service unit file for running RedCompass as a system service on Linux. Features:
- Runs as a simple service type
- Automatically restarts on failure (5s delay)
- Sets appropriate file limits (4096)
- Logs both stdout and stderr to `/root/pb/std.log`
- Configures the working directory and domain
- Enables the service to start on system boot

Use this file to set up RedCompass as a systemd service for production deployments on Linux systems.

## webhook.pb.js

A PocketBase JavaScript hook that checks if domains will expire in the next 30 days and if so will send a webhook via Slack or Discord at 8am every Friday.
