package main

import (
	"fmt"

	"github.com/blachniet/go-surfbeam"
)

func main() {
	surfbeamClient := surfbeam.New(surfbeam.DefaultModemURI)
	if s, err := surfbeamClient.ModemStatus(); err != nil {
		fmt.Printf("Err: %v\n", err.Error())
	} else {
		fmt.Printf(`Status: %v
IP Address: %v
Software Version: %v
Hardware Version: %v
Total Traffic (GB): %.2f
`,
			s.Status,
			s.IPAddress,
			s.SoftwareVersion,
			s.HardwareVersion,
			(float64(s.TransmittedBytes)+float64(s.ReceivedBytes))/1024./1024./1024.)
	}
}
