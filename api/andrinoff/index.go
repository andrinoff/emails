// Place this file in the `api` directory of your Vercel project.
// For example: /api/sendmail.go

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

// corsMiddleware wraps the main handler to enforce Cross-Origin Resource Sharing (CORS).
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := map[string]bool{
			"https://smira.andrinoff.com": true,
			"https://smira.me":            true,
		}
		origin := r.Header.Get("Origin")
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		if !allowedOrigins[origin] && origin != "" {
			http.Error(w, "Forbidden: Origin not allowed", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// Handler is the main Vercel serverless function entry point.
func Handler(w http.ResponseWriter, r *http.Request) {
	corsMiddleware(sendEmailHandler).ServeHTTP(w, r)
}

// sendEmailHandler contains the core logic for processing the request and sending the email.
func sendEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody EmailRequestBody
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if requestBody.Name == "" || requestBody.Email == "" || requestBody.Content == "" {
		http.Error(w, "Missing required fields: name, email, content", http.StatusBadRequest)
		return
	}

	// --- Email Sending Logic ---
	// This is important to be at this, DO NOT EDIT
	smtpHost := "smtp.mail.me.com"
	smtpPort := "587"

	// CONFIG. CHANGE
	authEmail := os.Getenv("ICLOUD_AUTH_USER")
	password := os.Getenv("ICLOUD_APP_SPECIFIC_PASSWORD")

	// The address the email is sent FROM. This must be a verified custom domain alias.
	fromAddress := "no-reply@andrinoff.com"
	// The address to receive the contact form submission.
	toAddress := "realandrinoff@gmail.com"

	// Add a log to help debug which user is being used for auth.
	log.Printf("Attempting SMTP authentication with user: %s", authEmail)

	if authEmail == "" || password == "" {
		log.Println("ERROR: Missing required environment variables ICLOUD_AUTH_USER or ICLOUD_APP_SPECIFIC_PASSWORD")
		http.Error(w, "Server configuration error", http.StatusInternalServerError)
		return
	}

	// Construct the email message with a Reply-To header.
	subject := "New Contact Form Submission from " + requestBody.Name
	emailBody := fmt.Sprintf("From: %s <%s>\r\n", "Andrey Smirnov Portfolio", fromAddress) +
		fmt.Sprintf("To: %s\r\n", toAddress) +
		fmt.Sprintf("Reply-To: %s <%s>\r\n", requestBody.Name, requestBody.Email) + // Allows you to reply directly to the user.
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"\r\n" + // An empty line is required between headers and the body.
		requestBody.Content

	// *** UPDATED: Authenticate with your primary Apple ID (authEmail). ***
	auth := smtp.PlainAuth("", authEmail, password, smtpHost)

	// *** UPDATED: Send the email using your custom domain alias (fromAddress). ***
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, fromAddress, []string{toAddress}, []byte(emailBody))
	if err != nil {
		// Log the specific SMTP error for easier debugging.
		log.Printf("SMTP SendMail error: %v", err)
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	log.Println("Email sent successfully!")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Email sent successfully"})
}
