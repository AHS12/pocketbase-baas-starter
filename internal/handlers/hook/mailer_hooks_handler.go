package hook

import (
	"net/mail"
	"strings"

	log "ims-pocketbase-baas-starter/pkg/logger"

	"github.com/pocketbase/pocketbase/core"
)

// Helper function to convert mail.Address slice to string slice
func addressesToStrings(addresses []mail.Address) []string {
	result := make([]string, len(addresses))
	for i, addr := range addresses {
		result[i] = addr.String()
	}
	return result
}

// HandleMailerSend handles email send events
func HandleMailerSend(e *core.MailerEvent) error {

	log.Info("Email being sent",
		"to", strings.Join(addressesToStrings(e.Message.To), ", "),
		"subject", e.Message.Subject,
		"from", e.Message.From.String(),
	)

	// Add your custom logic here

	// Example: Add custom headers
	// if e.Message.Headers == nil {
	// 	e.Message.Headers = make(map[string]string)
	// }
	// e.Message.Headers["X-App-Name"] = "IMS PocketBase BaaS Starter"
	// e.Message.Headers["X-Environment"] = "development"

	// Continue with the execution chain
	return e.Next()
}

// HandleMailerBeforeSend handles pre-send email events
func HandleMailerBeforeSend(e *core.MailerEvent) error {
	// This would be called before the email is actually sent
	log.Debug("Preparing to send email",
		"to", strings.Join(addressesToStrings(e.Message.To), ", "),
		"subject", e.Message.Subject,
	)
	// Add pre-send logic here

	return e.Next()
}

// HandleMailerAfterSend handles post-send email events
func HandleMailerAfterSend(e *core.MailerEvent) error {
	log.Info("Email sent successfully",
		"to", strings.Join(addressesToStrings(e.Message.To), ", "),
		"subject", e.Message.Subject,
	)
	// Add post-send logic here

	return e.Next()
}
