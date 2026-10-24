package main

import "testing"

func TestCompleteJobNotifiesOnlyFinishedMatchingAsset(t *testing.T) {
	cases := []struct{ name, status, jobAsset, want string }{{"complete", "complete", "a-7", "a-7"}, {"queued", "queued", "a-7", ""}, {"wrong asset", "complete", "a-8", ""}}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			called := ""
			got, err := completeJob(Asset{ID: "a-7", Phone: "+15550001"}, Job{AssetID: tc.jobAsset, Status: tc.status}, func(to, content, key string) (string, error) { called = key; return "msg-1", nil })
			if tc.want == "" {
				if err == nil && got != "" {
					t.Fatalf("got %q", got)
				}
				if called != "" {
					t.Fatalf("unexpected send %q", called)
				}
			} else if err != nil || got != "msg-1" || called != "media-job-a-7" {
				t.Fatalf("got=%q err=%v key=%q", got, err, called)
			}
		})
	}
}
