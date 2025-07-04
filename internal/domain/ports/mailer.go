package ports

type MailerManager interface {
	LoadEmailTemplate(templateName string, data any) (string, error)
	SendEmail(to string, subject string, body string) error
}
