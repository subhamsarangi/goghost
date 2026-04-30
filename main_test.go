package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nyaruka/phonenumbers"
)

// --- IP Track ---

func TestIPInfoUnmarshal(t *testing.T) {
	raw := `{
		"ip":"8.8.8.8","type":"IPv4","country":"United States","country_code":"US",
		"city":"Mountain View","continent":"North America","continent_code":"NA",
		"region":"California","region_code":"CA","latitude":37.386,"longitude":-122.083,
		"is_eu":false,"postal":"94035","calling_code":"+1","capital":"Washington D.C.",
		"borders":"CAN,MEX","flag":{"emoji":"🇺🇸"},
		"connection":{"asn":15169,"org":"Google LLC","isp":"Google LLC","domain":"google.com"},
		"timezone":{"id":"America/Los_Angeles","abbr":"PDT","is_dst":true,"offset":-25200,
		"utc":"-07:00","current_time":"2024-01-01T12:00:00"}
	}`
	var info IPInfo
	if err := json.Unmarshal([]byte(raw), &info); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if info.City != "Mountain View" {
		t.Errorf("expected City=Mountain View, got %s", info.City)
	}
	if info.Connection.ISP != "Google LLC" {
		t.Errorf("expected ISP=Google LLC, got %s", info.Connection.ISP)
	}
	if info.Timezone.ID != "America/Los_Angeles" {
		t.Errorf("expected TZ=America/Los_Angeles, got %s", info.Timezone.ID)
	}
	if info.Flag.Emoji != "🇺🇸" {
		t.Errorf("expected flag emoji, got %s", info.Flag.Emoji)
	}
}

func TestIPTrackHTTP(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(IPInfo{
			IP: "1.2.3.4", City: "TestCity", Country: "Testland",
		})
	}))
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/1.2.3.4")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	var info IPInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		t.Fatalf("decode failed: %v", err)
	}
	if info.City != "TestCity" {
		t.Errorf("expected TestCity, got %s", info.City)
	}
}

// --- Show My IP ---

func TestShowMyIPHTTP(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("1.2.3.4"))
	}))
	defer srv.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(srv.URL)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

// --- Phone Track ---

func TestPhoneParse_Valid(t *testing.T) {
	parsed, err := phonenumbers.Parse("+12025550147", "US")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if !phonenumbers.IsValidNumber(parsed) {
		t.Error("expected valid number")
	}
	if phonenumbers.GetRegionCodeForNumber(parsed) != "US" {
		t.Error("expected region US")
	}
	e164 := phonenumbers.Format(parsed, phonenumbers.E164)
	if e164 != "+12025550147" {
		t.Errorf("expected E164 +12025550147, got %s", e164)
	}
}

func TestPhoneParse_Invalid(t *testing.T) {
	_, err := phonenumbers.Parse("notanumber", "US")
	if err == nil {
		t.Error("expected error for invalid number")
	}
}

func TestPhoneParse_WithoutCountryCode(t *testing.T) {
	parsed, err := phonenumbers.Parse("2025550147", "US")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	if !phonenumbers.IsPossibleNumber(parsed) {
		t.Error("expected possible number")
	}
}

// --- Username Track ---

func TestUsernameTrack_Found(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Head(srv.URL + "/testuser")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestUsernameTrack_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Head(srv.URL + "/ghostuser")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		t.Error("expected non-200 for missing user")
	}
}

func TestUsernameTrack_Concurrent(t *testing.T) {
	hits := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	sites := []SocialMediaSite{
		{"Site1", srv.URL + "/user", ""},
		{"Site2", srv.URL + "/user", ""},
		{"Site3", srv.URL + "/user", ""},
	}

	client := &http.Client{Timeout: 5 * time.Second}
	results := make(chan CheckResult, len(sites))

	for _, s := range sites {
		go func(site SocialMediaSite) {
			resp, err := client.Head(site.URL)
			found := err == nil && resp.StatusCode == http.StatusOK
			if resp != nil {
				resp.Body.Close()
			}
			results <- CheckResult{Name: site.Name, URL: site.URL, Found: found}
		}(s)
	}

	for i := 0; i < len(sites); i++ {
		r := <-results
		if r.Found {
			hits++
		}
	}

	if hits != len(sites) {
		t.Errorf("expected %d hits, got %d", len(sites), hits)
	}
}

func TestUsernameTrack_Timeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := &http.Client{Timeout: 300 * time.Millisecond}
	_, err := client.Head(srv.URL + "/user")
	if err == nil {
		t.Error("expected timeout error")
	}
}