package main

import (
	"crypto/sha256"
	"log"
)

func encryptSessionID(sessionID string) []byte {
	//get the master key for encrypting
	//this will be a combination of the mac address
	//if we fail to get the hostname or the mac address
	machineID, err := getMachineID()
	if err != nil {
		log.Println("Failed to get machine id")
		return nil
	}
	macAddress, err := getMacAddr()
	hash := sha256.New()
	combined := macAddress + machineID
	hash.Write([]byte(combined))
	return hash.Sum(nil)
}
