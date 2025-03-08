package main

import (
	"bytes"
	"net"
	"os"
	"strings"
)

const (
	dbusPath         = "/var/lib/dbus/machine-id"
	dbusPathFallback = "/etc/machine-id"
)

func getMachineID() (string, error) {
	id, err := os.ReadFile(dbusPath)
	if err != nil {
		id, err = os.ReadFile(dbusPathFallback)
	}
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(strings.Trim(string(id), "\n")), nil
}

func getMacAddr() (string, error) {
	macAddress := ""
	interfaces, err := net.Interfaces()
	if err != nil {
		return macAddress, err
	}
	for _, i := range interfaces {
		if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
			macAddress = i.HardwareAddr.String()
			break
		}
	}
	return macAddress, nil
}
