package mariadb

import (
	"encoding/json"
	"testing"
)

// The rule-create API answers with a jobId (guide v4.0), not a ruleId.
func TestCreateSecurityRuleResponseReadsJobID(t *testing.T) {
	var r CreateSecurityRuleResponse
	if err := json.Unmarshal([]byte(`{"header":{"isSuccessful":true,"resultCode":0},"jobId":"job-1"}`), &r); err != nil {
		t.Fatal(err)
	}
	if r.JobID != "job-1" {
		t.Fatalf("jobId lost: %+v", r)
	}
}
