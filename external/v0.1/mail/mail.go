package mail

import (
	"fmt"
	"net/smtp"
	"os"
)

var (
	senderEmail = os.Getenv("AUTHEMAILUSERNAME")
	pwd         = os.Getenv("AUTHEMAILPWD")
)

func SendEmail(mail, token, kind, code, host string) {
	auth := smtp.PlainAuth("", senderEmail, pwd, "smtp.gmail.com")
	to := []string{mail}
	msg := []byte{}
	toBody := fmt.Sprintf("To: %s\r\n", mail)
	subject := ""
	body := ""
	switch kind {
	case "recover":
		subject = "Subject: Recover your account \r\n"
		body = recoverEmailBody(token, mail, code, host)
		fmt.Println(body)
	case "activate":
		subject = "Subject: Activate your account \r\n"
		body = activateEmailBody(token, host)
	}
	msg = formatEmail(toBody, subject, body)
	//msg := []byte("To: maximotejeda.com\r\n" + "Subject: Testing the send\r\n" + "\r\n" + "sended or what??" + message + "\r\n")
	//newMsg := formatEmail("To: maximotejeda.com\r\n", "Subject: Testing the send\r\n")
	err := smtp.SendMail("smtp.gmail.com:587", auth, senderEmail, to, msg)
	if err != nil {
		panic(err)
	}
}

func formatMailHTML() string {
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	return mime
}

func recoverEmailBody(token, email, code, host string) (body string) {
	body = fmt.Sprintf(`
	<html>
		<body>
			<h1>
				Recover Account
			</h1>
			<p>A new recovery process for the account was issued for the email %s</p>
			<p>to continue with the process click on the link <a href="https://%s/recover?token=%s">Recover My account</a> and introduce the next code</p>
			<h3>%s</h3>
			<p>If you did not ask for recovery please dont click any of the links here</p>
		</body>
	</html>
	`, email, host, token, code)
	return body
}

func activateEmailBody(token, host string) (body string) {
	body = fmt.Sprintf(`
	<html>
		<body>
			<h1>
				Activate Account
			</h1>
			<p>Please activate your new account</p>
			<p>Click the link to activate your account <a href="http://%s/activate?token=%s">Activate </a></p>
		</body>
	</html>
	`, host, token)
	return body
}

func formatEmail(to, subject, msg string) []byte {
	mimeType := formatMailHTML()
	body := msg
	msg1 := to + subject + mimeType + body
	return []byte(msg1)
}
