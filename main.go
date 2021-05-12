package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	cloudflare "github.com/cloudflare/cloudflare-go"
)

// API credentials and domain to be updated
type CloudflareOptions struct {
	user   string
	domain string
	record string
	token  string
	newIp  string
}

func main() {
	fmt.Println("Starting …")
	http.HandleFunc("/", handler)

	// -------------------------
	// Start the listening loop
	// -------------------------
	http.ListenAndServe(":4242", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received call from ")
	fmt.Println(r.RemoteAddr)
	fmt.Println(r.RequestURI)

	// -------------------------
	// Parse the received URL and check for the API token
	// -------------------------
	r.ParseForm()
	fmt.Println(r.Form)
	if len(r.Form["token"]) == 0 {
		fmt.Println("token empty")
		return
	}

	co := CloudflareOptions{
		user:   r.FormValue("user"),
		domain: r.FormValue("domain"),
		record: r.FormValue("record"),
		token:  r.FormValue("token"),
		newIp:  r.FormValue("newip"),
	}

	fmt.Println("Supplied new IP is: " + string(co.newIp))
	fmt.Println(co)

	// -------------------------
	// Update the record using the Cloudflare API
	// -------------------------

	UpdateRecord(co)
}

func UpdateRecord(options CloudflareOptions) {

	// -------------------------
	// Get API object
	// -------------------------
	api, err := cloudflare.NewWithAPIToken(options.token)
	if err != nil {
		log.Fatal(err)
	}

	// -------------------------
	// Get Zone ID in order to get the zone's DNS record
	// -------------------------
	zoneId, err := api.ZoneIDByName(options.domain)
	if err != nil {
		log.Fatal(err)
	}

	// -------------------------
	// Get the DNS record to update it later
	// -------------------------
	dnsRecord, err := api.DNSRecords(context.Background(), zoneId, cloudflare.DNSRecord{Name: options.record})
	if err != nil {
		log.Fatal(err)
	}
	if dnsRecord == nil {
		fmt.Println("Record not found")
		os.Exit(1)
	}

	fmt.Println("Trying to change record …")

	// -------------------------
	// Update the record with the new A record
	// -------------------------
	err = api.UpdateDNSRecord(context.Background(), zoneId, dnsRecord[0].ID, cloudflare.DNSRecord{Content: options.newIp})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Record changed.")
}
