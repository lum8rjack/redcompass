// Cleanup domain ideas every day at 2am
// and delete ideas that are already owned

cronAdd("cleanup-domain-ideas", "0 2 * * *", () => {
    const logname = "CRON: cleanup-domain-ideas";
    try {
        // Get all domain ideas
        let allideas = $app.findAllRecords("Domain_Ideas");

        // If there is at least one domain idea
        if (allideas.length > 0) {
            // Loop through the result and delete the ideas that are already owned
            for (let i = 0; i < allideas.length; i++) {
                let idea = allideas[i];
                let domainName = idea.get("Domain");
                let owned = $app.findAllRecords("Domains", 
                    $dbx.hashExp({"Name": domainName}),
                );
                if (owned.length > 0) {
                    $app.delete(idea);
                    $app.logger().info(logname + " - domain idea is already owned", "domain", domainName, "action", "deleted");
                }
            }
        }
    } catch (error) {
        $app.logger().error(logname + " - error", "error_message", error.Error());
    }
});
