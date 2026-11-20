package util

import "testing"

func TestFormatValidationIssues(t *testing.T) {
	cases := []struct {
		name string
		body string
		want string
		ok   bool
	}{
		{
			name: "issues with paths",
			body: `{"reqId":"req-1","statusCode":422,"message":[{"code":"too_small","message":"Number must be greater than or equal to 1","path":["limit"]},{"code":"invalid_union_discriminator","message":"Invalid discriminator value","path":["services",0,"type"]}],"error":"ValidationFailure"}`,
			want: "limit: Number must be greater than or equal to 1; services.0.type: Invalid discriminator value",
			ok:   true,
		},
		{
			name: "issue without a path",
			body: `{"message":[{"code":"custom","message":"Provide at least one setting to update","path":[]}]}`,
			want: "Provide at least one setting to update",
			ok:   true,
		},
		{name: "plain message", body: `{"message":"Invalid request"}`, want: "Invalid request", ok: true},
		{name: "not JSON", body: `<html>Bad Gateway</html>`},
		{name: "no message", body: `{"error":"ValidationFailure"}`},
		{name: "issue without a message", body: `{"message":[{"path":["limit"]}]}`},
		{name: "empty issue list", body: `{"message":[]}`},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, ok := formatValidationIssues(c.body)
			if got != c.want || ok != c.ok {
				t.Fatalf("got (%q, %v), want (%q, %v)", got, ok, c.want, c.ok)
			}
		})
	}
}
