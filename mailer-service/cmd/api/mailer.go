package main

import (
	"bytes"
	"html/template"
	"log"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

// Definicja struktury obiektu Mail
type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}

// Definicja struktury wiadomości
type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}

// Funkcja wysyłania wiadomości przez serwer SMTP
func (m *Mail) SendSMTPMessage(msg Message) error {

	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	data := map[string]any{
		"message": msg.Data,
	}

	msg.DataMap = data

	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}

	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		log.Println(err)
		return err
	}

	email := mail.NewMSG()
	email.SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject)

	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML, formattedMessage)

	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}

	err = email.Send(smtpClient)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// Tworzenie wiadomości w formacie HTML
func (m *Mail) buildHTMLMessage(msg Message) (string, error) {

	templateToRender := "./templates/mail.html.gohtml"

	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	formattedMessage := tpl.String()
	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}

	return formattedMessage, nil
}

// Tworzenie wiadomości w formacie tekstowym
func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {

	templateToRender := "./templates/mail.plain.gohtml"

	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	plainMessage := tpl.String()

	return plainMessage, nil
}

// Funkcja konwertuje przekazany ciąg znaków zawierający kod HTML i CSS
// na nowy ciąg znaków HTML, w którym style CSS zostały osadzone bezpośrednio w HTML
func (m *Mail) inlineCSS(s string) (string, error) {
	// Tworzymy nowy obiekt "options" z ustawieniami dla biblioteki "premailer".
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	// Tworzymy nowy obiekt "prem" z wykorzystaniem funkcji NewPremailerFromString(),
	// która przyjmuje ciąg znaków "s" i obiekt "options".
	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}

	// Wywołujemy funkcję Transform() na obiekcie "prem",
	// aby przetworzyć ciąg znaków HTML i CSS.
	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	// Zwracamy nowy ciąg znaków HTML, który zawiera osadzone style CSS.
	return html, nil
}

// Funkcja zwraca typ szyfrowania połączenia z serwerem pocztowym
// na podstawie przekazanego ciągu znaków
func (m *Mail) getEncryption(s string) mail.Encryption {
	// Używając instrukcji switch sprawdzamy wartość argumentu s.
	switch s {
	// Jeśli wartość s to "tls",
	// zwracamy wartość mail.EncryptionSTARTTLS.
	case "tls":
		return mail.EncryptionSTARTTLS
	// Jeśli wartość s to "ssl",
	// zwracamy wartość mail.EncryptionSSLTLS.
	case "ssl":
		return mail.EncryptionSSLTLS
	// Jeśli wartość s to "none" lub pusty ciąg znaków,
	// zwracamy wartość mail.EncryptionNone.
	case "none", "":
		return mail.EncryptionNone
	// W przypadku, gdy wartość s jest inna niż powyższe,
	// zwracamy domyślną wartość mail.EncryptionSTARTTLS.
	default:
		return mail.EncryptionSTARTTLS
	}
}
