package main

import "fmt"

type Asset struct{ ID, Creator, Phone string }
type Job struct{ AssetID, Status string }

func completeJob(asset Asset, job Job, send func(string, string, string) (string, error)) (string, error) {
	if job.AssetID != asset.ID {
		return "", fmt.Errorf("job asset mismatch")
	}
	if job.Status != "complete" {
		return "", nil
	}
	message := fmt.Sprintf("Asset %s is ready for review.", asset.ID)
	return send(asset.Phone, message, "media-job-"+asset.ID)
}
