package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/nyaruka/phonenumbers"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[38;5;197m" // hot pink
	Green  = "\033[38;5;51m"  // electric cyan
	Yellow = "\033[38;5;214m" // amber orange
	White  = "\033[38;5;183m" // soft lavender
	Cyan   = "\033[38;5;129m" // deep purple
)

var reader = bufio.NewReader(os.Stdin)

func input(prompt string) string {
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		fmt.Print("\033[H\033[2J")
	}
}

func printBanner() {
	clearScreen()
	fmt.Fprintf(os.Stderr, White+"\n"+
		"     .--.                                      \n"+
		"   .'    `.                                    \n"+
		"  /  .- -.  \\                                  \n"+
		" |  (o) (o)  |        "+White+"--------------------------------\n"+
		" |     __    |        "+White+"| "+Green+"GOGHOST - TRACKER TOOL       "+White+"|\n"+
		" |    (__}   |        "+White+"|    "+White+                                   "|\n"+
		"  \\          /        "+White+"--------------------------------\n"+
		"   `-.____.-'                                  \n"+
		"   :        `.                                 \n"+
		"  :    :.     `.                               \n"+
		" :    :  :      `.                             \n"+
		":    :    :    .  `.                           \n"+
		":   :      :  . `.  `.                         \n"+
		" `.:        `..   `...`                        \n"+
		"  ~  ~   ~    ~  ~   ~                         \n"+
		Reset)
	time.Sleep(500 * time.Millisecond)
}

func withBanner(fn func()) func() {
	return func() {
		printBanner()
		fn()
	}
}

type IPInfo struct {
	IP            string  `json:"ip"`
	Type          string  `json:"type"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"country_code"`
	City          string  `json:"city"`
	Continent     string  `json:"continent"`
	ContinentCode string  `json:"continent_code"`
	Region        string  `json:"region"`
	RegionCode    string  `json:"region_code"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	IsEU          bool    `json:"is_eu"`
	Postal        string  `json:"postal"`
	CallingCode   string  `json:"calling_code"`
	Capital       string  `json:"capital"`
	Borders       string  `json:"borders"`
	Flag          struct {
		Emoji string `json:"emoji"`
	} `json:"flag"`
	Connection struct {
		ASN    int    `json:"asn"`
		Org    string `json:"org"`
		ISP    string `json:"isp"`
		Domain string `json:"domain"`
	} `json:"connection"`
	Timezone struct {
		ID          string `json:"id"`
		Abbr        string `json:"abbr"`
		IsDST       bool   `json:"is_dst"`
		Offset      int    `json:"offset"`
		UTC         string `json:"utc"`
		CurrentTime string `json:"current_time"`
	} `json:"timezone"`
}

func ipTrack() {
	ip := input(White + "\n Enter IP target : " + Green)
	fmt.Println()

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("http://ipwho.is/" + ip)
	if err != nil {
		fmt.Println(Red + " Request failed: " + err.Error() + Reset)
		return
	}
	defer resp.Body.Close()

	var info IPInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		fmt.Println(Red + " Failed to parse response: " + err.Error() + Reset)
		return
	}

	fmt.Printf(" %s============= %sSHOW INFORMATION IP ADDRESS %s=============\n", White, Green, White)
	fmt.Printf("\n %sIP target       :%s %s\n", White, Green, ip)
	fmt.Printf(" %sType IP         :%s %s\n", White, Green, info.Type)
	fmt.Printf(" %sCountry         :%s %s\n", White, Green, info.Country)
	fmt.Printf(" %sCountry Code    :%s %s\n", White, Green, info.CountryCode)
	fmt.Printf(" %sCity            :%s %s\n", White, Green, info.City)
	fmt.Printf(" %sContinent       :%s %s\n", White, Green, info.Continent)
	fmt.Printf(" %sContinent Code  :%s %s\n", White, Green, info.ContinentCode)
	fmt.Printf(" %sRegion          :%s %s\n", White, Green, info.Region)
	fmt.Printf(" %sRegion Code     :%s %s\n", White, Green, info.RegionCode)
	fmt.Printf(" %sLatitude        :%s %f\n", White, Green, info.Latitude)
	fmt.Printf(" %sLongitude       :%s %f\n", White, Green, info.Longitude)
	fmt.Printf(" %sMaps            :%s https://www.google.com/maps/@%.4f,%.4f,8z\n", White, Green, info.Latitude, info.Longitude)
	fmt.Printf(" %sEU              :%s %v\n", White, Green, info.IsEU)
	fmt.Printf(" %sPostal          :%s %s\n", White, Green, info.Postal)
	fmt.Printf(" %sCalling Code    :%s %s\n", White, Green, info.CallingCode)
	fmt.Printf(" %sCapital         :%s %s\n", White, Green, info.Capital)
	fmt.Printf(" %sBorders         :%s %s\n", White, Green, info.Borders)
	fmt.Printf(" %sCountry Flag    :%s %s\n", White, Green, info.Flag.Emoji)
	fmt.Printf(" %sASN             :%s %d\n", White, Green, info.Connection.ASN)
	fmt.Printf(" %sORG             :%s %s\n", White, Green, info.Connection.Org)
	fmt.Printf(" %sISP             :%s %s\n", White, Green, info.Connection.ISP)
	fmt.Printf(" %sDomain          :%s %s\n", White, Green, info.Connection.Domain)
	fmt.Printf(" %sID              :%s %s\n", White, Green, info.Timezone.ID)
	fmt.Printf(" %sABBR            :%s %s\n", White, Green, info.Timezone.Abbr)
	fmt.Printf(" %sDST             :%s %v\n", White, Green, info.Timezone.IsDST)
	fmt.Printf(" %sOffset          :%s %d\n", White, Green, info.Timezone.Offset)
	fmt.Printf(" %sUTC             :%s %s\n", White, Green, info.Timezone.UTC)
	fmt.Printf(" %sCurrent Time    :%s %s\n", White, Green, info.Timezone.CurrentTime)
}

func phoneTrack() {
	num := input(White + "\n Enter phone number target " + Green + "Ex [+12025550123] " + White + ": " + Green)
	defaultRegion := "US"

	parsed, err := phonenumbers.Parse(num, defaultRegion)
	if err != nil {
		fmt.Println(Red + " Failed to parse number: " + err.Error() + Reset)
		return
	}

	regionCode := phonenumbers.GetRegionCodeForNumber(parsed)
	carrier, _ := phonenumbers.GetCarrierForNumber(parsed, "en")
	location, _ := phonenumbers.GetGeocodingForNumber(parsed, "en")
	isValid := phonenumbers.IsValidNumber(parsed)
	isPossible := phonenumbers.IsPossibleNumber(parsed)
	intlFormat := phonenumbers.Format(parsed, phonenumbers.INTERNATIONAL)
	e164Format := phonenumbers.Format(parsed, phonenumbers.E164)
	nationalFormat := phonenumbers.Format(parsed, phonenumbers.NATIONAL)
	numType := phonenumbers.GetNumberType(parsed)

	typeStr := "Other"
	switch numType {
	case phonenumbers.MOBILE:
		typeStr = "Mobile"
	case phonenumbers.FIXED_LINE:
		typeStr = "Fixed Line"
	case phonenumbers.FIXED_LINE_OR_MOBILE:
		typeStr = "Fixed Line or Mobile"
	}

	fmt.Printf("\n %s========== %sSHOW INFORMATION PHONE NUMBERS %s==========\n", White, Green, White)
	fmt.Printf("\n %sLocation             :%s %s\n", White, Green, location)
	fmt.Printf(" %sRegion Code          :%s %s\n", White, Green, regionCode)
	fmt.Printf(" %sOperator             :%s %s\n", White, Green, carrier)
	fmt.Printf(" %sValid number         :%s %v\n", White, Green, isValid)
	fmt.Printf(" %sPossible number      :%s %v\n", White, Green, isPossible)
	fmt.Printf(" %sInternational format :%s %s\n", White, Green, intlFormat)
	fmt.Printf(" %sNational format      :%s %s\n", White, Green, nationalFormat)
	fmt.Printf(" %sE.164 format         :%s %s\n", White, Green, e164Format)
	fmt.Printf(" %sCountry code         :%s %d\n", White, Green, parsed.GetCountryCode())
	fmt.Printf(" %sLocal number         :%s %d\n", White, Green, parsed.GetNationalNumber())
	fmt.Printf(" %sType                 :%s %s\n", White, Green, typeStr)
}

type SocialMediaSite struct {
	Name   string
	URL    string
	Signup string
}

type CheckResult struct {
	Name          string
	URL           string
	Signup        string
	Found         bool
	Unverifiable  bool
}

func usernameTrack() {
	username := input(White + "\n Enter Username : " + Green)

	sites := []SocialMediaSite{
		{"Facebook", "https://www.facebook.com/" + username, "https://www.facebook.com/reg"},
		{"X (Twitter)", "https://x.com/" + username, "https://x.com/i/flow/signup"},
		{"Instagram", "https://www.instagram.com/" + username, "https://www.instagram.com/accounts/emailsignup"},
		{"LinkedIn", "https://www.linkedin.com/in/" + username, "https://www.linkedin.com/signup"},
		{"GitHub", "https://www.github.com/" + username, "https://github.com/signup"},
		{"Pinterest", "https://www.pinterest.com/" + username, "https://www.pinterest.com/join"},
		{"Tumblr", "https://www.tumblr.com/" + username, "https://www.tumblr.com/register"},
		{"YouTube", "https://www.youtube.com/" + username, "https://accounts.google.com/signup"},
		{"SoundCloud", "https://soundcloud.com/" + username, "https://soundcloud.com/signup"},
		{"Snapchat", "https://www.snapchat.com/add/" + username, "https://accounts.snapchat.com/accounts/signup"},
		{"TikTok", "https://www.tiktok.com/@" + username, "https://www.tiktok.com/signup"},
		{"Behance", "https://www.behance.net/" + username, "https://www.behance.net/signup"},
		{"Medium", "https://www.medium.com/@" + username, "https://medium.com/m/signin"},
		{"Quora", "https://www.quora.com/" + username, "https://www.quora.com/"},
		{"Flickr", "https://www.flickr.com/people/" + username, "https://identity.flickr.com/sign-up"},
		{"Twitch", "https://www.twitch.tv/" + username, "https://www.twitch.tv/signup"},
		{"Dribbble", "https://www.dribbble.com/" + username, "https://dribbble.com/signup"},
		{"Product Hunt", "https://www.producthunt.com/@" + username, "https://www.producthunt.com/join"},
		{"Telegram", "https://www.telegram.me/" + username, "https://telegram.org/"},
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// bodyCheck: sites that return 200 for missing profiles (JS-rendered or soft 404).
	// mustContain = at least one phrase must be present for profile to be "found".
	// jsOnly      = fully JS-rendered, cannot verify via raw HTTP.
	type bodyRule struct {
		mustContain []string
		mustAbsent  []string
		jsOnly      bool
	}
	bodyChecks := map[string]bodyRule{
		// X and Facebook are fully React-rendered — raw HTML has zero profile data
		"X (Twitter)": {jsOnly: true},
		"Facebook":    {jsOnly: true},
		// Quora: profile page has og:url; 404 page does not
		"Quora": {mustContain: []string{
			`og:url" content="https://www.quora.com/profile/` + username,
			`og:url" content="https://www.quora.com/` + username,
		}},
		// Dribbble: public profile has og:url; private/missing does not
		"Dribbble": {mustContain: []string{
			`og:url" content="https://dribbble.com/` + username,
		}},
	}

	results := make(chan CheckResult, len(sites))
	var wg sync.WaitGroup

	for _, site := range sites {
		wg.Add(1)
		go func(s SocialMediaSite) {
			defer wg.Done()

			rule, needsBodyCheck := bodyChecks[s.Name]

			// JS-only sites cannot be verified via raw HTTP
			if needsBodyCheck && rule.jsOnly {
				results <- CheckResult{Name: s.Name, URL: s.URL, Signup: s.Signup, Unverifiable: true}
				return
			}

			method := http.MethodHead
			if needsBodyCheck {
				method = http.MethodGet
			}

			req, err := http.NewRequest(method, s.URL, nil)
			if err != nil {
				results <- CheckResult{Name: s.Name, URL: s.URL, Signup: s.Signup, Found: false}
				return
			}
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
			req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
			req.Header.Set("Accept-Language", "en-US,en;q=0.5")

			resp, err := client.Do(req)
			if err != nil {
				results <- CheckResult{Name: s.Name, URL: s.URL, Signup: s.Signup, Found: false}
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusGone {
				results <- CheckResult{Name: s.Name, URL: s.URL, Signup: s.Signup, Found: false}
				return
			}

			if needsBodyCheck {
				// Read up to 64KB of raw HTML — og: meta tags appear in <head>, well within this range
				buf := make([]byte, 65536)
				n, _ := io.ReadFull(resp.Body, buf)
				body := strings.ToLower(string(buf[:n]))

				if len(rule.mustContain) > 0 {
					for _, phrase := range rule.mustContain {
						if strings.Contains(body, strings.ToLower(phrase)) {
							results <- CheckResult{Name: s.Name, URL: s.URL, Signup: s.Signup, Found: true}
							return
						}
					}
					results <- CheckResult{Name: s.Name, URL: s.URL, Signup: s.Signup, Found: false}
					return
				}

				for _, phrase := range rule.mustAbsent {
					if strings.Contains(body, strings.ToLower(phrase)) {
						results <- CheckResult{Name: s.Name, URL: s.URL, Signup: s.Signup, Found: false}
						return
					}
				}
				results <- CheckResult{Name: s.Name, URL: s.URL, Signup: s.Signup, Found: true}
				return
			}

			// For LinkedIn: returns 999 for valid profiles when blocking scrapers, 404 for missing
			io.Copy(io.Discard, resp.Body)
			found := resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusGone
			results <- CheckResult{Name: s.Name, URL: s.URL, Signup: s.Signup, Found: found}
		}(site)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	collected := make([]CheckResult, 0, len(sites))
	for r := range results {
		collected = append(collected, r)
	}

	// Sort: verified results first, unverifiable (Facebook, X) at the end
	sort.Slice(collected, func(i, j int) bool {
		if collected[i].Unverifiable != collected[j].Unverifiable {
			return !collected[i].Unverifiable
		}
		return collected[i].Name < collected[j].Name
	})

	fmt.Printf("\n %s========== %sSHOW INFORMATION USERNAME %s==========\n\n", White, Green, White)
	for _, r := range collected {
		if r.Unverifiable {
			fmt.Printf(" %s[ %s? %s] %s : %sCheck manually → %s%s  %s| Sign up: %s%s\n", White, Cyan, White, r.Name, Cyan, Cyan, r.URL, White, Cyan, r.Signup)
		} else if r.Found {
			fmt.Printf(" %s[ %s+ %s] %s : %s%s\n", White, Green, White, r.Name, Green, r.URL)
		} else {
			fmt.Printf(" %s[ %s- %s] %s : %sNot found  %s→ Sign up: %s%s\n", White, Yellow, White, r.Name, Yellow, White, Cyan, r.Signup)
		}
	}
}

func showMyIP() {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://api.ipify.org/")
	if err != nil {
		fmt.Println(Red + " Request failed: " + err.Error() + Reset)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(Red + " Failed to read response" + Reset)
		return
	}

	fmt.Printf("\n %s========== %sSHOW INFORMATION YOUR IP %s==========\n", White, Green, White)
	fmt.Printf("\n %s[ %s+ %s] Your IP Address : %s%s\n", White, Green, White, Green, string(body))
	fmt.Printf("\n %s==========================================\n", White)
}

func printMenu() {
	clearScreen()
	fmt.Fprintf(os.Stderr, White+`
   ________      ________               __ 
  / ____/ /___  / ____/ /_  ____  _____/ /_
 / / __/ / __ \/ / __/ __ \/ __ \/ ___/ __/
/ /_/ / / /_/ / /_/ / / / / /_/ (__  ) /_  
\____/_/\____/\____/_/ /_/\____/____/\__/  

`+Reset)

	fmt.Fprintf(os.Stderr, `
%s[ 1 ] %sTrace an IP Address
%s[ 2 ] %sReveal My IP
%s[ 3 ] %sInvestigate Phone Number
%s[ 4 ] %sHunt a Username
%s[ 0 ] %sDrop Connection

`, White, Green, White, Green, White, Green, White, Green, White, Green)
}

func main() {
	for {
		printMenu()
		choice := input(White + "\n [ + ] " + Green + "Select Option : " + White)

		switch choice {
		case "1":
			withBanner(ipTrack)()
		case "2":
			withBanner(showMyIP)()
		case "3":
			withBanner(phoneTrack)()
		case "4":
			withBanner(usernameTrack)()
		case "0":
			fmt.Println(Red + "\n [ ! ] Exit" + Reset)
			time.Sleep(500 * time.Millisecond)
			os.Exit(0)
		default:
			fmt.Println(Red + "\n [ ! ] Please input a valid number" + Reset)
			time.Sleep(1500 * time.Millisecond)
			continue
		}

		input(White + "\n[ + ] " + Green + "Press enter to continue")
	}
}