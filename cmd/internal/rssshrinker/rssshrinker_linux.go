// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.

//go:build linux

// Package rssshrinker provides functions to reduce the process’ RSS
package rssshrinker

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

/*
#include <malloc.h>
*/
import "C"

//revive:disable:var-naming

// MADV_PAGEOUT Reclaim these pages.
const MADV_PAGEOUT = 21

//revive:enable:var-naming

// Setup starts timers to regularly release memory to the kernel.
func Setup() {
	// Release the memory garbage collected by the Go runtime to the kernel every 2 minutes.
	go func() {
		c := time.Tick(2 * time.Minute)
		for range c {
			debug.FreeOSMemory()
		}
	}()

	// Release the memory allocated by C malloc to the kernel every 2 minutes.
	// This is for the native libraries pulled by Python.
	go func() {
		c := time.Tick(2 * time.Minute)
		for range c {
			C.malloc_trim(0)
		}
	}()

	// Release file-backed memory to the kernel every 30 minutes.
	// This is for the GO code that isn’t actively used.
	go func() {
		if err := pageOutFileBackedMemory(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to release memory: %s\n", err)
		}

		c := time.Tick(30 * time.Minute)
		for range c {
			if err := pageOutFileBackedMemory(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to release memory: %s\n", err)
			}
		}
	}()
}

// pageOutFileBackedMemory releases file-backed memory by advising the kernel to page out the memory.
func pageOutFileBackedMemory() error {
	selfMap, err := os.Open("/proc/self/maps")
	if err != nil {
		return fmt.Errorf("failed to open /proc/self/maps: %w", err)
	}
	defer selfMap.Close()

	scanner := bufio.NewScanner(selfMap)
	for scanner.Scan() {
		// Each line in /proc/self/maps has the following format:
		// address           perms offset  dev   inode       pathname
		// 00400000-0291f000 r--p 00000000 00:53 5009645     /opt/datadog-agent/bin/agent/agent
		// 0291f000-08038000 r-xp 0251f000 00:53 5009645     /opt/datadog-agent/bin/agent/agent
		// 08038000-0eb30000 r--p 07c38000 00:53 5009645     /opt/datadog-agent/bin/agent/agent
		// 0eb30000-0eb31000 r--p 0e730000 00:53 5009645     /opt/datadog-agent/bin/agent/agent
		// 0eb31000-0efed000 rw-p 0e731000 00:53 5009645     /opt/datadog-agent/bin/agent/agent
		fields := strings.Fields(scanner.Text())

		// If the 6th column is missing, the line is about an anonymous mapping.
		// We ignore it as we want to page out only file-backed memory.
		if len(fields) != 6 {
			continue
		}

		address, perms, _ /* offset */, _ /* device */, _ /* inode */, pathname := fields[0], fields[1], fields[2], fields[3], fields[4], fields[5]

		// Ignore pseudo-paths about stack, heap, vdso, named anonymous mapping, etc.
		if strings.HasPrefix(pathname, "[") {
			continue
		}

		// We only want to page out read-only memory.
		if len(perms) != 4 || perms[0] != 'r' || perms[1] != '-' {
			continue
		}

		beginStr, endStr, found := strings.Cut(address, "-")
		if !found {
			continue
		}

		begin, err := strconv.ParseUint(beginStr, 16, int(8*unsafe.Sizeof(uintptr(0))))
		if err != nil {
			return fmt.Errorf("failed to parse begin address %q: %w", beginStr, err)
		}

		end, err := strconv.ParseUint(endStr, 16, int(8*unsafe.Sizeof(uintptr(0))))
		if err != nil {
			return fmt.Errorf("failed to parse end address %q: %w", endStr, err)
		}

		// nolint:govet
		if err := syscall.Madvise(unsafe.Slice((*byte)(unsafe.Pointer(uintptr(begin))), end-begin), MADV_PAGEOUT); err != nil {
			return fmt.Errorf("failed to madvise: %w", err)
		}
	}

	return nil
}
