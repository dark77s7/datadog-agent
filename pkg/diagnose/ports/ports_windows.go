// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package ports

import (
	"fmt"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	SystemProcessIDInformationClass = 88 // SystemProcessIDInformationClass gives access to process names without elevated privileges on Windows
)

// SystemProcessIDInformation is a struct for Windows API.
type SystemProcessIDInformation struct {
	ProcessID uintptr
	ImageName windows.NTUnicodeString
}

// RetrieveProcessName fetches the process name on Windows using NtQuerySystemInformation
// with SystemProcessIDInformation, which does not require elevated privileges.
func RetrieveProcessName(pid int, _ string) (string, error) {
	var processInfo SystemProcessIDInformation
	processInfo.ProcessID = uintptr(pid)
	ret := windows.NtQuerySystemInformation(SystemProcessIDInformationClass, unsafe.Pointer(&processInfo), uint32(unsafe.Sizeof(processInfo)), nil)

	if ret != nil {
		return "", ret
	}

	// Convert UTF-16 string and handle trailing nulls
	processName := windows.UTF16PtrToString((*uint16)(unsafe.Pointer(processInfo.ImageName.Buffer)))
	processName = strings.TrimRight(processName, "\x00")
	processName = strings.ToLower(strings.TrimSuffix(processName, ".exe"))
	fmt.Printf("PRE-RETURN VALUE: %s", processName)
	return processName, nil
}
