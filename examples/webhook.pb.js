// Check weekly if domains will expire in the next 30 days
// If so, send a webhook via Slack or Discord at 8am every Friday

cronAdd("expiring-domains-webhook", "0 8 * * 5", () => {
    const logname = "CRON: expiring-domains-webhook";
    try {
        // Get current date
        const currentDate = new Date();

        // Create a new date 30 days in the future
        const futureDate = new Date();
        futureDate.setDate(currentDate.getDate() + 30);

        // Get Domains that will expire in the next 30 days
        let allrecords = $app.findAllRecords("Domains",
            $dbx.between("Expiration_Date", currentDate, futureDate),
        )

        // If there is at least one domain that will expire in the next 30 days
        if (allrecords.length > 0) {
            let message = "Upcoming Domain Expirations:\n\n";

            // Loop through the result and send a webhook to Slack and/or Discord
            for (let i = 0; i < allrecords.length; i++) {
                let expirationDate = allrecords[i].get("Expiration_Date");
                let autoRenew = allrecords[i].get("Auto_Renew");
                let isLocked = allrecords[i].get("Is_Locked");

                expirationDate = new Date(expirationDate);
                expirationDate = expirationDate.toISOString().split('T')[0];
                message += `Domain: ${allrecords[i].get("Name")} | Expires: ${expirationDate} | Auto-Renew: ${autoRenew} | Is Locked: ${isLocked}\n`;
            }

            // Message for Slack
            const slackWebhookUrl = "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"
            const jsonMessageSlack = JSON.stringify({
                text: message,
            });

            // Message for Discord
            const discordWebhookUrl = "https://discord.com/api/webhooks/1234567890/XXXXXXXXXXXXXXXXXXXXXXXX"
            const jsonMessageDiscord = JSON.stringify({
                content: message,
            });

            // Send a webhook to Slack and/or Discord
            // Replace URL and BODY based on Slack or Discord
            const res = $http.send({
                url: slackWebhookUrl,
                headers: {"content-type": "application/json"},
                method: "POST",
                body: jsonMessageSlack,
            })

            if (res.status === 200) {
                $app.logger().info(logname + " - webhook sent successfully", "expiring_domains", allrecords.length);
            } else {
                $app.logger().error(logname + " - webhook failed to send", "status_code", res.status);
            }
        }

    } catch (error) {
        $app.logger().error(logname + " - error", "error_message", error);
    }
});
