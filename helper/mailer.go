package helper

import (
	"greenenvironment/configs"
	"log"
	"net/smtp"
)

type MailerInterface interface {
	Send(to string, code string, subject string) error
}

type Mailer struct {
	config configs.SMTPConfig
}

func NewMailer(config configs.SMTPConfig) MailerInterface {
	return &Mailer{
		config: config,
	}
}

func (m *Mailer) Send(to string, code string, subject string) error {
	from := m.config.Username
	pass := m.config.Password

	htmlBody := `<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <title>EcoMate ` + subject + `</title>
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin />
        <link href="https://fonts.googleapis.com/css2?family=Nunito&display=swap" rel="stylesheet" />
        <style>
            p,
            h1 h2 h3 h4 h5 h6,
            div,
            td {
                font-family: Nunito, sans-serif;
            }
            a {
                color: #365cce;
                text-decoration: none;
            }
            .border {
                border-style: solid;
                border-width: 1px;
                border-color: #365cce;
                border-radius: 0.25rem;
            }

            .footertext {
                font-size: 12px;
            }

            @media (min-width: 640px) {
                .footertext {
                    font-size: 16px;
                }
            }
        </style>
    </head>
    <body style="margin: 0; padding: 0">
        <table cellpadding="0" cellspacing="0" border="0" width="100%" style="max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 8px; overflow: hidden">
            <tr>
                <td style="padding: 20px; background-color: #ffffff; text-align: center">
                    <img src="https://res.cloudinary.com/djsrf0egn/image/upload/v1733985438/ecomate/qzdb9lri1tm9dulwm3pt.png" alt="logo" />
                </td>
            </tr>
            <tr>
                <td style="background-color: #2e7d32; padding: 30px 20px; text-align: center">
                    <table cellpadding="0" cellspacing="0" border="0" width="100%">
                        <tr>
                            <td style="text-align: center">
                                <img src="https://res.cloudinary.com/djsrf0egn/image/upload/v1733985367/ecomate/x9jfkqysf1rvnweq1msz.png" alt="arrow" />
                            </td>
                        </tr>
                        <tr>
                            <td style="text-align: center; font-size: 26px; font-weight: normal; color: #fafafa; padding-top: 16px; padding-bottom: 16px">` + subject + `</td>
                        </tr>
                        <tr>
                            <td style="font-size: 24px; font-weight: bold; text-transform: capitalize; text-align: center; color: #fafafa">Verify Your OTP Code</td>
                        </tr>
                    </table>
                </td>
            </tr>
            <tr>
                <td style="padding: 30px 20px">
                    <p style="margin: 0 0 20px 0">Hello ` + to + `,</p>
                    <p style="margin: 0 0 20px 0">Please use the following One Time Password (OTP)</p>
                    <table cellpadding="0" cellspacing="0" border="0" width="100%" style="margin-bottom: 20px">
                        <tr>
                            <td style="width: 16.66%; text-align: center"><div style="border: 1px solid #2e7d32; padding: 10px; font-size: 24px; font-weight: bold">` + string(code[0]) + `</div></td>
                            <td style="width: 16.66%; text-align: center"><div style="border: 1px solid #2e7d32; padding: 10px; font-size: 24px; font-weight: bold">` + string(code[0]) + `</div></td>
                            <td style="width: 16.66%; text-align: center"><div style="border: 1px solid #2e7d32; padding: 10px; font-size: 24px; font-weight: bold">` + string(code[0]) + `</div></td>
                            <td style="width: 16.66%; text-align: center"><div style="border: 1px solid #2e7d32; padding: 10px; font-size: 24px; font-weight: bold">` + string(code[0]) + `</div></td>
                            <td style="width: 16.66%; text-align: center"><div style="border: 1px solid #2e7d32; padding: 10px; font-size: 24px; font-weight: bold">` + string(code[0]) + `</div></td>
                            <td style="width: 16.66%; text-align: center"><div style="border: 1px solid #2e7d32; padding: 10px; font-size: 24px; font-weight: bold">` + string(code[0]) + `</div></td>
                        </tr>
                    </table>
                    <p style="margin: 0 0 20px 0">This passcode will only be valid for the next 2 minutes. If the passcode does not work, you can use this login verification link:</p>
                    <br />
                    <br />
                    <p style="margin: 20px 0 0 0">Thank you,</p>
                    <br />
                    <p style="margin: 0">Ecomate Team.</p>
                </td>
            </tr>
            <tr>
                <td style="background-color: #2e7d32; padding: 20px; text-align: center; color: #ffffff">
                    <p style="color: #fafafa; font-weight: bold; font-size: 20px; letter-spacing: 2px">Get in touch</p>
                    <a href="mailto:greenenvironmentcaps@gmail.com" style="color: #fafafa" alt="greenenvironmentcaps@gmail.com">greenenvironmentcaps@gmail.com</a>
                    <p
                        class="footertext"
                        style="font-size: 12px; color: #fafafa; width: 90%; text-align: center; border-top: 0.5px solid rgba(229, 231, 235, 0.7); padding-top: 13px; margin: 15px auto"
                    >
                        Copyright © 2024 Ecomate. All rights reserved.
                    </p>
                </td>
            </tr>
        </table>
    </body>
</html>
`

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = "Ecomate Verification Code"
	headers["MIME-version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	header := ""
	for k, v := range headers {
		header += k + ": " + v + "\r\n"
	}

	msg := header + "\r\n" + htmlBody

	err := smtp.SendMail(m.config.Host+":"+m.config.Port,
		smtp.PlainAuth("", from, pass, m.config.Host),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}
	return nil
}
