package authotp

import "net/smtp"

func SendEmail(email string, otp string) error {

	from := "mubashirullatuthodi@gmail.com"
	password := "jhzv jkkb ewgr sbrh"
	to := []string{email}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	msg := []byte("Subject: Your OTP for Sign Up\n\n OTP is: " + otp)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	if err != nil {
		return err
	}
	return nil
}
