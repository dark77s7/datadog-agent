// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-present Datadog, Inc.

package common

import (
	"fmt"
	"strings"

	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/components"
)

// GetRegistryValue returns a registry value from a remote host
func GetRegistryValue(host *components.RemoteHost, path string, value string) (string, error) {
	cmd := fmt.Sprintf("Get-ItemPropertyValue -Path '%s' -Name '%s'", path, value)
	out, err := host.Execute(cmd)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// RegistryKeyExists returns true if the registry key exists on the remote host
func RegistryKeyExists(host *components.RemoteHost, path string) (bool, error) {
	cmd := fmt.Sprintf("Test-Path -Path '%s'", path)
	out, err := host.Execute(cmd)
	if err != nil {
		return false, err
	}
	return strings.EqualFold(strings.TrimSpace(out), "True"), nil
}

// DeleteRegistryKey deletes a registry key on the remote host
func DeleteRegistryKey(host *components.RemoteHost, path string) error {
	cmd := fmt.Sprintf("Remove-Item -Path '%s' -Recurse -Force", path)
	_, err := host.Execute(cmd)
	return err
}

// SetRegistryDWORDValue sets, creating if necessary, a DWORD value at the specified path
func SetRegistryDWORDValue(host *components.RemoteHost, path string, name string, value int) error {
	return SetTypedRegistryValue(host, path, name, fmt.Sprintf("%d", value), "DWORD")
}

// SetTypedRegistryValue sets, creating if necessary, the value at the specified path with the specified type
//
// https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.management/set-itemproperty?view=powershell-7.4#-type
func SetTypedRegistryValue(host *components.RemoteHost, path string, name string, value string, typeName string) error {
	cmd := fmt.Sprintf("New-Item -Path '%s' -Force; Set-ItemProperty -Path '%s' -Name '%s' -Value '%s' -Type '%s'", path, path, name, value, typeName)
	_, err := host.Execute(cmd)
	if err != nil {
		return err
	}
	return nil
}
