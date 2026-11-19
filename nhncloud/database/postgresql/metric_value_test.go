package postgresql

import (
	"encoding/json"
	"testing"
)

// The live API sends numbers; the guide documents strings. Both decode.
func TestMetricStatisticValueNumberOrString(t *testing.T) {
	var vs []MetricStatisticValue
	if err := json.Unmarshal([]byte(`[{"timestamp":1,"value":2.5},{"timestamp":2,"value":"3.25"}]`), &vs); err != nil {
		t.Fatal(err)
	}
	if vs[0].Value != 2.5 || vs[1].Value != 3.25 || vs[1].Timestamp != 2 {
		t.Fatalf("got %+v", vs)
	}
}
