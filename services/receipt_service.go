package services

import (
	"html/template"
	"log"
	"maxpark_opp_registration/utils"
	"strings"
)

type ReceiptService struct {
}

func NewReceiptService() (service *ReceiptService, err error) {
	return &ReceiptService{}, err
}

func (s *ReceiptService) SendReceipt(from, to, subject, recipientName string) {
	go func() {
		tmpl, err := template.ParseFiles("templates/registration_received.html")
		if err != nil {
			log.Printf("receipt template parse error: %v", err)
			return
		}
		var sb strings.Builder
		if err := tmpl.Execute(&sb, map[string]interface{}{
			"Name": recipientName,
		}); err != nil {
			log.Printf("receipt template execute error: %v", err)
			return
		}
		if err := utils.SendEmail(from, to, subject, sb.String()); err != nil {
			log.Printf("receipt send error to %s: %v", to, err)
		}
	}()
}
