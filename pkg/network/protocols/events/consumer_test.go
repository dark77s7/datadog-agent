// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build linux_bpf

package events

import (
	"math"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
	"unsafe"

	ebpftelemetry "github.com/DataDog/datadog-agent/pkg/ebpf/telemetry"
	manager "github.com/DataDog/ebpf-manager"
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/features"
	"github.com/cilium/ebpf/perf"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"

	ddebpf "github.com/DataDog/datadog-agent/pkg/ebpf"
	"github.com/DataDog/datadog-agent/pkg/ebpf/bytecode"
	"github.com/DataDog/datadog-agent/pkg/network/config"
	"github.com/DataDog/datadog-agent/pkg/util/kernel"
)

func TestConsumer(t *testing.T) {
	kversion, err := kernel.HostVersion()
	require.NoError(t, err)
	if minVersion := kernel.VersionCode(4, 14, 0); kversion < minVersion {
		t.Skipf("package not supported by kernels < %s", minVersion)
	}

	const numEvents = 100
	c := config.New()
	program, err := newEBPFProgram(c)
	require.NoError(t, err)

	var mux sync.Mutex
	result := make(map[uint64]int)
	callback := func(events []uint64) {
		mux.Lock()
		defer mux.Unlock()
		for _, n := range events {
			result[n] = +1
		}
	}

	consumer, err := NewConsumer("test", program.Manager, callback)
	require.NoError(t, err)
	consumer.Start()

	err = program.Start()
	require.NoError(t, err)

	// generate test events
	generator := newEventGenerator(program.Manager, t)
	for i := 0; i < numEvents; i++ {
		generator.Generate(uint64(i))
	}
	generator.Stop()
	time.Sleep(100 * time.Millisecond)

	// this ensures that any incomplete batch left in eBPF is fully processed
	consumer.Sync()
	program.Stop(manager.CleanAll)
	consumer.Stop()

	// ensure that we have received each event exactly once
	for i := 0; i < numEvents; i++ {
		actual := result[uint64(i)]
		assert.Equalf(t, 1, actual, "eventID=%d should have 1 occurrence. got %d", i, actual)
	}
}

func TestInvalidBatchCountMetric(t *testing.T) {
	kversion, err := kernel.HostVersion()
	require.NoError(t, err)
	if minVersion := kernel.VersionCode(4, 14, 0); kversion < minVersion {
		t.Skipf("package not supported by kernels < %s", minVersion)
	}

	c := config.New()
	program, err := newEBPFProgram(c)
	require.NoError(t, err)
	t.Cleanup(func() { program.Stop(manager.CleanAll) })

	consumer, err := NewConsumer("test", program.Manager, func([]uint64) {})
	require.NoError(t, err)

	// We are creating a raw sample with a data length of 4, which is smaller than sizeOfBatch
	// and would be considered an invalid batch.
	recordSample(c, consumer, []byte("test"))

	consumer.Start()
	t.Cleanup(func() { consumer.Stop() })
	require.Eventually(t, func() bool {
		// Wait for the consumer to process the invalid batch.
		return consumer.invalidBatchCount.Get() == 1
	}, 5*time.Second, 100*time.Millisecond)
}

type eventGenerator struct {
	// map used for coordinating test with eBPF program space
	testMap *ebpf.Map

	// file used for triggering write(2) syscalls
	testFile *os.File
}

// recordSample records a sample using the consumer handler.
func recordSample(c *config.Config, consumer *Consumer[uint64], sampleData []byte) {
	// Ring buffers require kernel version 5.8.0 or higher, therefore, the handler is chosen based on the kernel version.
	if c.EnableUSMRingBuffers && features.HaveMapType(ebpf.RingBuf) == nil {
		handler := consumer.handler.(*ddebpf.RingBufferHandler)
		handler.RecordHandler(&ringbuf.Record{
			RawSample: sampleData,
		}, nil, nil)
	} else {
		handler := consumer.handler.(*ddebpf.PerfHandler)
		handler.RecordHandler(&perf.Record{
			RawSample: sampleData,
		}, nil, nil)
	}
}

func newEventGenerator(program *manager.Manager, t *testing.T) *eventGenerator {
	m, _, _ := program.GetMap("test")
	require.NotNilf(t, m, "couldn't find test map")

	f, err := os.Create(filepath.Join(t.TempDir(), "foobar"))
	require.NoError(t, err)

	return &eventGenerator{
		testMap:  m,
		testFile: f,
	}
}

func (e *eventGenerator) Generate(eventID uint64) error {
	type testCtx struct {
		pid     uint32
		fd      uint64
		eventID uint64
	}

	var (
		fd  = uint64(e.testFile.Fd())
		key = uint32(0)
		val = testCtx{pid: uint32(os.Getpid()), fd: fd, eventID: eventID}
	)

	// this is done so the eBPF program will echo back an event that is equal to
	// eventID once the write syscall is triggered below
	err := e.testMap.Put(unsafe.Pointer(&key), unsafe.Pointer(&val))
	if err != nil {
		return err
	}

	e.testFile.Write([]byte("whatever"))
	return nil
}

func (e *eventGenerator) Stop() {
	e.testFile.Close()
}

func newEBPFProgram(c *config.Config) (*ddebpf.Manager, error) {
	bc, err := bytecode.GetReader(c.BPFDir, "usm_events_test-debug.o")
	if err != nil {
		return nil, err
	}
	defer bc.Close()

	m := &manager.Manager{
		Probes: []*manager.Probe{
			{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "tracepoint__syscalls__sys_enter_write",
				},
			},
		},
	}
	options := manager.Options{
		RLimit: &unix.Rlimit{
			Cur: math.MaxUint64,
			Max: math.MaxUint64,
		},
		ActivatedProbes: []manager.ProbesSelector{
			&manager.ProbeSelector{
				ProbeIdentificationPair: manager.ProbeIdentificationPair{
					EBPFFuncName: "tracepoint__syscalls__sys_enter_write",
				},
			},
		},
		ConstantEditors: []manager.ConstantEditor{
			{
				Name:  "test_monitoring_enabled",
				Value: uint64(1),
			},
		},
	}

	ddEbpfManager := ddebpf.NewManager(m, "usm", &ebpftelemetry.ErrorsTelemetryModifier{})

	Configure(config.New(), "test", ddEbpfManager.Manager, &options)
	err = ddEbpfManager.InitWithOptions(bc, &options)
	if err != nil {
		return nil, err
	}

	return ddEbpfManager, nil
}
