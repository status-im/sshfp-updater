package main

import (
	"bufio"
	"bytes"
	"infra-sshfp-cf/cloudflare"
	"infra-sshfp-cf/config"
	"infra-sshfp-cf/consul"
	"infra-sshfp-cf/sshfp"
	"infra-sshfp-cf/statestore"
	"os"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

func main() {
	//Debug loglevel
	logrus.SetLevel(logrus.InfoLevel)
	//Create configuration components
	// Get config file name from args, if empty - try to configure from ENVs
	var configFilename string = ""
	if len(os.Args) > 1 {
		configFilename = os.Args[1]
	}

	cfgService := config.NewService(config.NewFileRepository())
	config, err := cfgService.LoadConfig(configFilename)
	if err != nil {
		logrus.Fatal(err)
	}

	switch strings.ToLower(config.LogLevel) {
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	//Create cloudflare components
	cloudflare := cloudflare.NewService(cloudflare.NewRepository(config.CloudflareToken, config.DomainName))

	//Create buffer to catch output from consul
	var buf bytes.Buffer

	//Try run consul binary and catch output
	logrus.Debug("Calling binary")
	output, err := exec.Command("consul", "watch", "-type=service", "-service=sshd", "-token", config.ConsulToken, "-http-addr=http://127.0.0.1:8500").Output()

	if err != nil {
		logrus.Fatalf("Cannot run consul, error: %s", err)
	}

	//Store output to the buffer
	buf.Write(output)

	//Use bufio Reader to turn bytes buffer into reader
	bufferedReader := bufio.NewReader(&buf)

	//Create reader listener and start reading on (blocking operation)
	consul := consul.NewService(consul.NewStdinRepository(bufferedReader))
	err = consul.LoadData()

	//Code below is executed upon data receipt

	if err != nil {
		logrus.Fatal(err)
	}

	hosts := consul.GetHostnames()

	//Open statestore
	statestore := statestore.NewService(statestore.NewMapRepository(config.StorageFilePath))

	//Iterate over hosts and check modify indexes
	for _, hostname := range hosts {
		modifiedIndex := consul.GetModifiedIndex(hostname)
		modified, _ := statestore.CheckIfModified(hostname, modifiedIndex)

		if modified {
			//Check if hosts has A record in CF - if not ignore
			exists, err := cloudflare.FindHostByName(hostname)
			if err != nil || !exists {
				continue
			}

			//Generate DNS records based on metadata
			sshfp := sshfp.NewService()
			consulKeys := sshfp.ParseConsulSSHRecords(consul.GetMetaData(hostname))

			//GetCurrent configuration
			cloudflareKeys, err := cloudflare.GetSSHFPRecordsForHost(hostname)
			if err != nil {
				logrus.Error(err)
			}

			configPlan := sshfp.PrepareConfiguration(hostname, cloudflareKeys, consulKeys)

			//ConfigPlan is empty, but host was flagged as modified. It's very likely that is a new host, or db is corrupted
			if len(configPlan) == 0 {
				statestore.SaveState(hostname, modifiedIndex)
			}

			if len(configPlan) > 0 {
				sshfp.PrintConfigPlan(configPlan)
				itemsApplied, err := cloudflare.ApplyConfigPlan(configPlan)
				if err != nil {
					logrus.Error(err)
				}
				if itemsApplied == len(configPlan) {
					statestore.SaveState(hostname, modifiedIndex)
				}
			}
		}
	}

	//Purge old hosts
	hosts, _ = statestore.GetStalledHosts(config.HostTimeout)
	for _, host := range hosts {
		err := cloudflare.DeleteSSHFPRecordsForHost(host)
		if err != nil {
			logrus.Fatalf("Cannot delete records for host: %s", host)
		}

	}
	statestore.PurgeStalledHosts(config.HostTimeout)

}
