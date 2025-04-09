// Block all emails from being sent

onBootstrap((e) => {
    e.next()
    $app.logger().info("block-emails.pb.js hook loaded",
        "details", "Block all emails from being sent")
})

onMailerSend((e) => {
    $app.logger().warn("Blocking email from being sent",
        "function", "MailerSend",
        "hook", "block-emails.pb.js")
})

onMailerRecordAuthAlertSend((e) => {
    $app.logger().warn("Blocking email from being sent",
        "function", "MailerRecordAuthAlertSend",
        "hook", "block-emails.pb.js")
})

onMailerRecordPasswordResetSend((e) => {
    $app.logger().warn("Blocking email from being sent",
        "function", "MailerRecordPasswordResetSend",
        "hook", "block-emails.pb.js")
})

onMailerRecordVerificationSend((e) => {
    $app.logger().warn("Blocking email from being sent",
        "function", "MailerRecordVerificationSend",
        "hook", "block-emails.pb.js")
})

onMailerRecordEmailChangeSend((e) => {
    $app.logger().warn("Blocking email from being sent",
        "function", "MailerRecordEmailChangeSend",
        "hook", "block-emails.pb.js")
})

onMailerRecordOTPSend((e) => {
    $app.logger().warn("Blocking email from being sent",
        "function", "MailerRecordOTPSend",
        "hook", "block-emails.pb.js")
})
