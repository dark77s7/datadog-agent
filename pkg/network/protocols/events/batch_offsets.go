// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build linux_bpf

// Package events contains implementation to unify perf-map communication between kernel and user space.
package events

import (
	"sync"

	"github.com/DataDog/datadog-agent/pkg/util/log"
)

// offsetManager is responsible for keeping track of which chunks of data we
// have consumed from each batch object
type offsetManager struct {
	mux        sync.Mutex
	stateByCPU []*cpuReadState
}

type cpuReadState struct {
	// this is the nextBatchID we're expecting for a particular CPU core. we use
	// this when we attempt to retrieve data that hasn't been sent from kernel space
	// yet because it belongs to an incomplete batch.
	nextBatchID int

	// information associated to partial batch reads
	partialBatchID int
	partialOffset  int
}

func newOffsetManager(numCPUS int) *offsetManager {
	stateByCPU := make([]*cpuReadState, numCPUS)
	for i := range stateByCPU {
		stateByCPU[i] = new(cpuReadState)
	}

	return &offsetManager{stateByCPU: stateByCPU}
}

// Get returns the data offset that hasn't been consumed yet for a given batch
func (o *offsetManager) Get(cpu int, batch *batch, syncing bool, id string) (begin, end int) {
	o.mux.Lock()
	defer o.mux.Unlock()
	state := o.stateByCPU[cpu]
	batchID := int(batch.Idx)

	// Log state and batch info
	log.Info("[USM] Get called: cpu=%d, batchID=%d, batchLen=%d, syncing=%v, id=%s", cpu, batchID, batch.Len, syncing, id)
	log.Info("[USM] State for cpu %d: nextBatchID=%d, partialBatchID=%d, partialOffset=%d, id=%s", cpu, state.nextBatchID, state.partialBatchID, state.partialOffset, id)

	if batchID < state.nextBatchID {
		// metric
		// we have already consumed this data
		log.Info("[USM] Skipping batch: batchID %d is less than nextBatchID %d, id=%s", batchID, state.nextBatchID, id)
		return 0, 0
	}

	if batchComplete(batch) {
		state.nextBatchID = batchID + 1
		log.Info("[USM] Batch complete: updating nextBatchID to %d, id %s", state.nextBatchID, id)
	}

	// determining the begin offset
	// usually this is 0, but if we've done a partial read of this batch
	// we need to take that into account
	if int(batch.Idx) == state.partialBatchID {
		begin = state.partialOffset
		log.Info("[USM] Partial batch detected: setting begin offset to %d, for id=%s", begin, id)
	}

	// determining the end offset
	// usually this is the full batch size but it can be less
	// in the context of a forced (partial) read
	end = int(batch.Len)
	log.Info("[USM] End offset set to: %d for id=%s", end, id)

	// if this is part of a forced read (that is, we're reading a batch before
	// it's complete) we need to keep track of which entries we're reading
	// so we avoid reading the same entries again
	if syncing {
		state.partialBatchID = int(batch.Idx)
		state.partialOffset = end
		log.Info("[USM] Syncing: updating partialBatchID to %d and partialOffset to %d, for id=%s", state.partialBatchID, state.partialOffset, id)
	}

	return
}

func (o *offsetManager) NextBatchID(cpu int) int {
	o.mux.Lock()
	defer o.mux.Unlock()

	return o.stateByCPU[cpu].nextBatchID
}

func batchComplete(b *batch) bool {
	return b.Cap > 0 && b.Len == b.Cap
}
