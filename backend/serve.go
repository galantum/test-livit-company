package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/chromedp/chromedp"
)

type RequestBody struct {
	Search string `json:"search"`
	Page   int    `json:"page"`
}

func setupCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func handleOptions(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		setupCORS(w)
		w.WriteHeader(http.StatusNoContent)
	}
}

func scrapeTikTok(w http.ResponseWriter, reqBody RequestBody, searchURL string) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var videoLinks []string
	err := chromedp.Run(ctx,
		chromedp.Navigate(searchURL),
		chromedp.WaitReady("#app", chromedp.ByID),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for i := 0; i < reqBody.Page; i++ {
				if err := chromedp.Evaluate(`window.scrollTo(0, document.body.scrollHeight)`, nil).Do(ctx); err != nil {
					return err
				}
				time.Sleep(3 * time.Second)
			}
			return nil
		}),
		chromedp.Evaluate(`Array.from(document.querySelectorAll('a.css-1g95xhm-AVideoContainer.e19c29qe13')).map(a => a.href)`, &videoLinks),
	)
	if err != nil {
		http.Error(w, "Error running chromedp tasks", http.StatusInternalServerError)
		return
	}

	data := []map[string]string{}
	for _, link := range videoLinks {
		if link != "" {
			data = append(data, map[string]string{"link": link})
		}
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Data found",
		"data":    data,
	}

	setupCORS(w)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}

func explore(w http.ResponseWriter, r *http.Request) {
	handleOptions(w, r)
	scrapeTikTok(w, RequestBody{}, "https://www.tiktok.com/explore?lang=id-ID")
}

func search(w http.ResponseWriter, r *http.Request) {
	setupCORS(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	//log.Println("Request Body (raw):", string(body))

	r.Body = io.NopCloser(bytes.NewBuffer(body))

	var reqBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	//log.Printf("Decoded Request Body: %+v\n", reqBody)

	searchURL := "https://www.tiktok.com/search?lang=id-ID&q=" + reqBody.Search
	scrapeTikTok(w, reqBody, searchURL)
}

func main() {
	http.HandleFunc("/", explore)
	http.HandleFunc("/search", search)
	fmt.Println("Server running on http://127.0.0.1:8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
