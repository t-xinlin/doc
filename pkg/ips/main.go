package main

import (
	"log"
	"net"
)

func main() {
	ips, _ := Ips()
	log.Println("ips: ", ips)
}

func Ips() (map[string]string, error) {
	ips := make(map[string]string)

	interfaces, err := net.Interfaces()
	if err != nil {
		log.Printf("Interfaces Error: %s", err.Error())
		return nil, err
	}

	for _, i := range interfaces {
		byName, err := net.InterfaceByName(i.Name)
		if err != nil {
			log.Printf("InterfaceByName Error: %s", err.Error())
			return nil, err
		}
		address, err := byName.Addrs()
		if err != nil {
			log.Printf("byName.Addrs Error: %s", err.Error())
		}

		for _, v := range address {
			ips[byName.Name] = v.String()
		}
	}

	return ips, nil
}
