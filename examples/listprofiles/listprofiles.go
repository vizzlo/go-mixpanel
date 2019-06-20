package main

import (
	"fmt"
	"os"
	"time"

	"vizzlo.com/mixpanel"
)

func usage() {
	fmt.Printf("usage: %s <API-SECRET>\n", os.Args[0])
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	secret := os.Args[1]

	mp := mixpanel.NewExportClient(secret)
	profiles, err := mp.ListProfiles(&mixpanel.ProfileQuery{LastSeenAfter: time.Now().Add(-time.Hour)})
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
}
