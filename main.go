package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	client, err := newSMSClient()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/ingest", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST required", http.StatusMethodNotAllowed)
			return
		}
		var asset Asset
		if err := json.NewDecoder(r.Body).Decode(&asset); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		var job Job
		if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
			job = Job{AssetID: asset.ID, Status: "complete"}
		}
		id, err := completeJob(asset, job, client.send)
		if err != nil {
			http.Error(w, err.Error(), 502)
			return
		}
		fmt.Fprintf(w, "{\"message_id\":%q}\n", id)
	})
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
