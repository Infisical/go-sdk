package util

import (
	"strings"
	"testing"
)

func TestBuildAzureMetadataServiceURL_DefaultResource(t *testing.T) {
	got := buildAzureMetadataServiceURL("", AzureIMDSIdentitySelector{})
	want := AZURE_METADATA_SERVICE_URL + AZURE_DEFAULT_RESOURCE
	if got != want {
		t.Fatalf("default resource URL: got %q, want %q", got, want)
	}
}

func TestBuildAzureMetadataServiceURL_WithClientID(t *testing.T) {
	got := buildAzureMetadataServiceURL("", AzureIMDSIdentitySelector{ClientID: "11111111-1111-1111-1111-111111111111"})
	if !strings.Contains(got, "&client_id=11111111-1111-1111-1111-111111111111") {
		t.Fatalf("expected client_id query param, got %q", got)
	}
	if strings.Contains(got, "object_id=") {
		t.Fatalf("did not expect object_id when client_id is set, got %q", got)
	}
}

func TestBuildAzureMetadataServiceURL_ClientIDWinsOverObjectID(t *testing.T) {
	got := buildAzureMetadataServiceURL("", AzureIMDSIdentitySelector{
		ClientID: "client-uuid",
		ObjectID: "object-uuid",
	})
	if !strings.Contains(got, "&client_id=client-uuid") {
		t.Fatalf("expected client_id, got %q", got)
	}
	if strings.Contains(got, "object_id=object-uuid") {
		t.Fatalf("client_id should win over object_id, got %q", got)
	}
}

func TestBuildAzureMetadataServiceURL_ObjectIDFallback(t *testing.T) {
	got := buildAzureMetadataServiceURL("", AzureIMDSIdentitySelector{ObjectID: "object-uuid"})
	if !strings.Contains(got, "&object_id=object-uuid") {
		t.Fatalf("expected object_id query param, got %q", got)
	}
}

func TestBuildAzureMetadataServiceURL_QueryEscape(t *testing.T) {
	got := buildAzureMetadataServiceURL("", AzureIMDSIdentitySelector{ClientID: "id with spaces"})
	if !strings.Contains(got, "&client_id=id+with+spaces") &&
		!strings.Contains(got, "&client_id=id%20with%20spaces") {
		t.Fatalf("expected url-escaped client_id, got %q", got)
	}
}

func TestBuildAzureMetadataServiceURL_NoSelectorPreservesLegacyURL(t *testing.T) {
	got := buildAzureMetadataServiceURL("", AzureIMDSIdentitySelector{})
	if strings.Contains(got, "client_id=") || strings.Contains(got, "object_id=") {
		t.Fatalf("expected legacy URL without selectors, got %q", got)
	}
}

func TestNormaliseAzureScope(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{in: "", want: "https://management.azure.com/.default"},
		{in: "https://vault.azure.net/", want: "https://vault.azure.net/.default"},
		{in: "https://vault.azure.net", want: "https://vault.azure.net/.default"},
		{in: "https://vault.azure.net/.default", want: "https://vault.azure.net/.default"},
		{in: "https%3A%2F%2Fmanagement.azure.com/", want: "https://management.azure.com/.default"},
	}
	for _, tc := range cases {
		got := normaliseAzureScope(tc.in)
		if got != tc.want {
			t.Errorf("normaliseAzureScope(%q): got %q, want %q", tc.in, got, tc.want)
		}
	}
}
