package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
)

// EmailRequestBody defines the structure of the incoming JSON request.
type EmailRequestBody struct {
	Name    string `json:"name"`
	Email   string `json:"email"` // This is the "from" email address for the reply-to header.
	Content string `json:"content"`
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// These are the only domains allowed to make requests.
		allowedOrigins := map[string]bool{
			"https://smira.andrinoff.com": true,
			"https://smira.me":            true,
		}

		origin := r.Header.Get("Origin")

		// Check if the request's origin is in our list of allowed origins.
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}

		// Handle preflight requests (OPTIONS method) which browsers send first
		// to check if the actual request is safe to send.
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// If the origin is not allowed, send a forbidden error.
		if !allowedOrigins[origin] && origin != "" {
			http.Error(w, "Forbidden: Origin not allowed", http.StatusForbidden)
			return
		}

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	corsMiddleware(sendEmailHandler).ServeHTTP(w, r)
}

func sendEmailHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody EmailRequestBody
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if requestBody.Name == "" || requestBody.Email == "" || requestBody.Content == "" {
		http.Error(w, "Missing required fields: name, email, content", http.StatusBadRequest)
		return
	}
	smtpHost := "smtp.mail.me.com"
	smtpPort := "587"
	fromEmail := os.Getenv("ICLOUD_EMAIL")                // Your no-reply@andrinoff.com email
	password := os.Getenv("ICLOUD_APP_SPECIFIC_PASSWORD") // The password you generated in Step 1
	toEmail := "smirnov.andrey@gmail.com"                 // The email address where you want to receive the contact form submissions.

	if fromEmail == "" || password == "" {
		log.Println("Error: Environment variables ICLOUD_EMAIL or ICLOUD_APP_SPECIFIC_PASSWORD are not set.")
		http.Error(w, "Server configuration error", http.StatusInternalServerError)
		return
	}

	// The message body.
	// We format it to be clear and include a "Reply-To" header so you can
	// reply directly to the person who filled out the form.
	subject := "New Contact Form Submission from " + requestBody.Name
	body := fmt.Sprintf("From: %s <%s>\r\n", requestBody.Name, requestBody.Email) +
		fmt.Sprintf("To: %s\r\n", toEmail) +
		fmt.Sprintf("Reply-To: %s\r\n", requestBody.Email) + // This is important!
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"\r\n" + // Empty line separates headers from the body
		requestBody.Content

	// Create the SMTP authentication object.
	auth := smtp.PlainAuth("", fromEmail, password, smtpHost)

	// Send the email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, fromEmail, []string{toEmail}, []byte(body))
	if err != nil {
		log.Printf("Error sending email: %v", err)
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	// Send a success response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Email sent successfully"})
}
