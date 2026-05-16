package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	// Обробник головної сторінки оплати (тепер чітко слухає /payment)
	http.HandleFunc("/payment", func(w http.ResponseWriter, r *http.Request) {
		bookingID := r.URL.Query().Get("booking_id")
		if bookingID == "" {
			http.Error(w, "Missing booking_id parameter", http.StatusBadRequest)
			return
		}

		tmpl, err := template.ParseFiles("templates/payment.html")
		if err != nil {
			http.Error(w, fmt.Sprintf("Template error: %v", err), http.StatusInternalServerError)
			return
		}

		_ = tmpl.Execute(w, map[string]string{"BookingID": bookingID})
	})

	// Обробник успішної оплати (слухає /payment/success)
	http.HandleFunc("/payment/success", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		bookingID := r.FormValue("booking_id")
		if bookingID == "" {
			http.Error(w, "Missing booking_id inside form data", http.StatusBadRequest)
			return
		}

		payload, _ := json.Marshal(map[string]string{
			"order_id": bookingID,
			"status":   "success",
		})

		targetURL := "http://app:8080/api/v1/payments/webhook"
		resp, err := http.Post(targetURL, "application/json", bytes.NewBuffer(payload))
		if err != nil {
			log.Printf("Failed to send webhook: %v", err)
			http.Error(w, "Internal app communication failure", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Main app rejected webhook with status: %d", resp.StatusCode)
			http.Error(w, "Main application rejected the payment", http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "http://localhost/index.html?payment=success", http.StatusSeeOther)
	})

	fmt.Println("Payment Mock service started on :8081...")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
