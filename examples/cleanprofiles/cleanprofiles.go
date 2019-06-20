package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"vizzlo.com/mixpanel"
)

func usage() {
	fmt.Printf("usage: %s <API-SECRET> <TOKEN>\n", os.Args[0])
}

const Day = time.Hour * 24
const Month = Day * 30
const Year = Day * 365

func main() {
	if len(os.Args) < 3 {
		usage()
		return
	}

	secret := os.Args[1]
	token := os.Args[2]

	exp := mixpanel.NewExportClient(secret)
	profiles, err := exp.ListProfiles(&mixpanel.ProfileQuery{
		LastSeenBefore:   time.Now().Add(-3 * Month),
		OutputProperties: []string{"$name"},
	})
	if err != nil {
		fmt.Println("Error occurred:", err)
	}

	for _, p := range profiles {
		fmt.Printf("* %s:\n", p.ID)
		for k, v := range p.Properties {
			fmt.Printf("  - %s = %v\n", k, v)
		}
		fmt.Println()
	}

	fmt.Printf("%d profiles.\n", len(profiles))

	mp := mixpanel.New(token)
	mp.Client.Timeout = time.Minute

	for _, p := range profiles {
		if err := mp.DeleteProfile(p.ID); err != nil {
			log.Fatalf("unable to delete %s: %s", p.ID, err)
		}
	}
}
