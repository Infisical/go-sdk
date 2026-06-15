package util

import "testing"

func TestGetAwsRegion_PrefersAwsRegionEnv(t *testing.T) {
	t.Setenv("AWS_REGION", "us-east-1")
	t.Setenv("AWS_DEFAULT_REGION", "eu-west-1")

	region, err := GetAwsRegion()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if region != "us-east-1" {
		t.Fatalf("expected AWS_REGION to win, got %q", region)
	}
}

func TestGetAwsRegion_FallsBackToDefaultRegionEnv(t *testing.T) {
	t.Setenv("AWS_REGION", "")
	t.Setenv("AWS_DEFAULT_REGION", "ap-south-1")

	region, err := GetAwsRegion()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if region != "ap-south-1" {
		t.Fatalf("expected AWS_DEFAULT_REGION fallback, got %q", region)
	}
}
