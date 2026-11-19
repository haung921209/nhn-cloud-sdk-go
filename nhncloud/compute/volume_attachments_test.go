package compute

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Attach and detach go through nova os-volume_attachments, not the cinder
// os-attach/os-detach bookkeeping actions.
func TestAttachDetachVolumeUseNovaAttachments(t *testing.T) {
	var seen []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		seen = append(seen, r.Method+" "+r.URL.Path+" "+string(b))
		if r.Method == http.MethodPost {
			io.WriteString(w, `{"volumeAttachment":{"id":"v1","device":"/dev/vdb","serverId":"s1","volumeId":"v1"}}`)
		}
	}))
	defer srv.Close()
	c := newTestClient(srv.URL)

	out, err := c.AttachVolume(context.Background(), "s1", "v1", "")
	if err != nil {
		t.Fatal(err)
	}
	if out.VolumeAttachment.Device != "/dev/vdb" || out.VolumeAttachment.ServerID != "s1" {
		t.Errorf("decoded %+v", out.VolumeAttachment)
	}
	if err := c.DetachVolume(context.Background(), "s1", "v1"); err != nil {
		t.Fatal(err)
	}
	want := []string{
		`POST /servers/s1/os-volume_attachments {"volumeAttachment":{"volumeId":"v1"}}`,
		`DELETE /servers/s1/os-volume_attachments/v1 `,
	}
	if len(seen) != 2 || seen[0] != want[0] || seen[1] != want[1] {
		t.Errorf("requests:\n%q\nwant\n%q", seen, want)
	}
}
