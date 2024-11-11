// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build linux_bpf

// Package eventparser is used for parsing raw bytes from bpf code into events
package eventparser

import (
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
	"unsafe"

	"golang.org/x/sys/unix"

	"github.com/DataDog/datadog-agent/pkg/dynamicinstrumentation/ditypes"
	"github.com/DataDog/datadog-agent/pkg/dynamicinstrumentation/ratelimiter"
	"github.com/DataDog/datadog-agent/pkg/util/log"
	"github.com/kr/pretty"
)

// MaxBufferSize is the maximum size of the output buffer from bpf which is read by this package
const MaxBufferSize = 10000

var (
	byteOrder = binary.LittleEndian
)

// ParseEvent takes the raw buffer from bpf and parses it into an event. It also potentially
// applies a rate limit
func ParseEvent(record []byte, ratelimiters *ratelimiter.MultiProbeRateLimiter) *ditypes.DIEvent {
	event := ditypes.DIEvent{}

	if len(record) < ditypes.SizeofBaseEvent {
		log.Tracef("malformed event record (length %d)", len(record))
		return nil
	}
	baseEvent := *(*ditypes.BaseEvent)(unsafe.Pointer(&record[0]))
	event.ProbeID = unix.ByteSliceToString(baseEvent.Probe_id[:])

	allowed, droppedEvents, successfulEvents := ratelimiters.AllowOneEvent(event.ProbeID)
	if !allowed {
		log.Tracef("event dropped by rate limit. Probe %s\t(%d dropped events out of %d)\n",
			event.ProbeID, droppedEvents, droppedEvents+successfulEvents)
		return nil
	}

	event.PID = baseEvent.Pid
	event.UID = baseEvent.Uid
	event.StackPCs = baseEvent.Program_counters[:]
	event.Argdata = readParams(record[ditypes.SizeofBaseEvent:])
	pretty.Log("Event:", event.Argdata)

	return &event
}

// ParseParams extracts just the parsed parameters from the full event record
func ParseParams(record []byte) ([]*ditypes.Param, error) {
	if len(record) < 392 {
		return nil, fmt.Errorf("malformed event record (length %d)", len(record))
	}
	return readParams(record[392:]), nil
}

func readParams(values []byte) []*ditypes.Param {
	fmt.Println(values[0:100])
	outputParams := []*ditypes.Param{}
	for i := 0; i < MaxBufferSize; {
		if i+3 >= len(values) {
			break
		}
		paramTypeDefinition := parseTypeDefinition(values[i:])
		pretty.Log(paramTypeDefinition)
		if paramTypeDefinition == nil {
			break
		}
		sizeOfTypeDefinition := countBufferUsedByTypeDefinition(paramTypeDefinition)
		i += sizeOfTypeDefinition
		val, numBytesRead := parseParamValue(paramTypeDefinition, values[i:])
		if val == nil {
			return outputParams
		}

		i += numBytesRead
		outputParams = append(outputParams, val)
	}
	return outputParams
}

// parseParamValue takes the representation of the param type's definition and the
// actual values in the buffer and populates the definition with the value parsed
// from the byte buffer. It returns the resulting parameter and an indication of
// how many bytes were read from the buffer
func parseParamValue(definition *ditypes.Param, buffer []byte) (*ditypes.Param, int) {
	var bufferIndex int = 0

	// Start by creating a stack with each layer of the definition
	// which will correspond with the layers of the values read from buffer.
	// This is done using a temporary stack to reverse the order.
	tempStack := newParamStack()
	definitionStack := newParamStack()
	tempStack.push(definition)
	for !tempStack.isEmpty() {
		current := tempStack.pop()

		if reflect.Kind(current.Kind) == reflect.Slice {
			len, err := readRuntimeSizedLength(buffer[bufferIndex : bufferIndex+2])
			if err != nil {
				log.Error(err)
			}
			bufferIndex += 2
			//TODO: Limit `len` to max slice elements
			current.Size = len
			if len == 0 {
				current.Fields = []*ditypes.Param{}
				_ = definitionStack.pop()
			} else if len > 1 {
				for i := 0; i < int(len)-1; i++ {
					copiedSliceElementDefinition := &ditypes.Param{}
					deepCopyParam(copiedSliceElementDefinition, current.Fields[0])
					current.Fields = append(current.Fields, copiedSliceElementDefinition)
				}
			}
		}
		copiedParam := copyParam(current)
		definitionStack.push(copiedParam)
		for n := 0; n < len(current.Fields); n++ {
			tempStack.push(current.Fields[n])
		}
	}

	valueStack := newParamStack()
	for bufferIndex+8 < len(buffer) {
		paramDefinition := definitionStack.pop()
		if paramDefinition == nil {
			break
		}

		if reflect.Kind(paramDefinition.Kind) == reflect.String {
			// We read the length first
			size, err := readRuntimeSizedLength(buffer[bufferIndex : bufferIndex+2])
			if err != nil {
				log.Error(err)
				break
			}
			bufferIndex += 2
			paramDefinition.Size = size
			paramDefinition.ValueStr = string(buffer[bufferIndex : bufferIndex+int(size)])
			bufferIndex += int(ditypes.StringMaxSize)
			valueStack.push(paramDefinition)
		} else if !isTypeWithHeader(paramDefinition.Kind) {
			if bufferIndex+int(paramDefinition.Size) >= len(buffer) {
				break
			}
			// This is a regular value (no sub-fields).
			// We parse the value of it from the buffer and push it to the value stack
			paramDefinition.ValueStr = parseIndividualValue(paramDefinition.Kind, buffer[bufferIndex:bufferIndex+int(paramDefinition.Size)])
			bufferIndex += int(paramDefinition.Size)
			valueStack.push(paramDefinition)
		} else if reflect.Kind(paramDefinition.Kind) == reflect.Pointer {
			if bufferIndex+int(paramDefinition.Size) >= len(buffer) {
				break
			}
			// Pointers are unique in that they have their own value, and sub-fields.
			// We parse the value of it from the buffer, place it in the value for
			// the pointer itself, then pop the next value and place it as a sub-field.
			paramDefinition.ValueStr = parseIndividualValue(paramDefinition.Kind, buffer[bufferIndex:bufferIndex+int(paramDefinition.Size)])
			bufferIndex += int(paramDefinition.Size)
			paramDefinition.Fields = append(paramDefinition.Fields, valueStack.pop())
			valueStack.push(paramDefinition)
		} else {
			// This is a type with sub-fields which have already been parsed and push
			// onto the value stack. We pop those and set them as fields in this type.
			// We then push this type onto the value stack as it may also be a sub-field.
			// In header types like this, paramDefinition.Size corresponds with the number of
			// fields under it.
			for n := 0; n < int(paramDefinition.Size); n++ {
				paramDefinition.Fields = append([]*ditypes.Param{valueStack.pop()}, paramDefinition.Fields...)
			}
			valueStack.push(paramDefinition)
		}
	}
	return valueStack.pop(), bufferIndex
}

func deepCopyParam(dst, src *ditypes.Param) {
	dst.Type = src.Type
	dst.Kind = src.Kind
	dst.Size = src.Size
	dst.Fields = make([]*ditypes.Param, len(src.Fields))
	for i, field := range src.Fields {
		dst.Fields[i] = &ditypes.Param{}
		deepCopyParam(dst.Fields[i], field)
	}
}

func copyParam(p *ditypes.Param) *ditypes.Param {
	return &ditypes.Param{
		Type: p.Type,
		Kind: p.Kind,
		Size: p.Size,
	}
}

func parseKindToString(kind byte) string {
	if kind == 255 {
		return "Unsupported"
	} else if kind == 254 {
		return "reached field limit"
	}

	return reflect.Kind(kind).String()
}

// parseTypeDefinition is given a buffer which contains the header type definition
// for basic/complex types, and the actual content of those types.
// It returns a fully populated tree of `ditypes.Param` which will be used for parsing
// the actual values
func parseTypeDefinition(b []byte) *ditypes.Param {
	stack := newParamStack()
	i := 0
	for {
		if len(b) < 3 {
			return nil
		}

		kind := b[i]
		newParam := &ditypes.Param{
			Kind: kind,
			Type: parseKindToString(kind),
		}
		if reflect.Kind(kind) == reflect.Slice {
			stack.push(newParam)
			i += 1
			continue
		}
		if reflect.Kind(kind) == reflect.String {
			i += 1
			goto stackCheck
		}

		newParam.Size = binary.LittleEndian.Uint16(b[i+1 : i+3])
		if newParam.Kind == 0 && newParam.Size == 0 {
			break
		}
		i += 3
		if isTypeWithHeader(newParam.Kind) {
			stack.push(newParam)
			continue
		}

	stackCheck:
		if stack.isEmpty() {
			return newParam
		}
		top := stack.peek()
		top.Fields = append(top.Fields, newParam)
		if len(top.Fields) == int(top.Size) ||
			(reflect.Kind(top.Kind) == reflect.Slice) ||
			(reflect.Kind(top.Kind) == reflect.Pointer && len(top.Fields) == 1) {
			newParam = stack.pop()
			goto stackCheck
		}

	}
	return nil
}

// countBufferUsedByTypeDefinition is used to determine that amount of bytes
// that were used to read the type definition. Each individual element of the
// definition uses 3 bytes (1 for kind, 2 for size). This is a needed calculation
// so we know where we should read the actual values in the buffer.
func countBufferUsedByTypeDefinition(root *ditypes.Param) int {
	queue := []*ditypes.Param{root}
	counter := 0
	for len(queue) != 0 {
		front := queue[0]
		queue = queue[1:]

		if reflect.Kind(front.Kind) == reflect.String ||
			reflect.Kind(front.Kind) == reflect.Slice {
			counter += 1
		} else {
			counter += 3
		}
		queue = append(queue, front.Fields...)
	}
	return counter
}

func isTypeWithHeader(pieceType byte) bool {
	return reflect.Kind(pieceType) == reflect.Struct ||
		reflect.Kind(pieceType) == reflect.Array ||
		reflect.Kind(pieceType) == reflect.Slice ||
		reflect.Kind(pieceType) == reflect.Pointer
}

func isRuntimeSizedType(pieceType byte) bool {
	return reflect.Kind(pieceType) == reflect.Slice ||
		reflect.Kind(pieceType) == reflect.String
}

func readRuntimeSizedLength(lengthBytes []byte) (uint16, error) {
	if len(lengthBytes) != 2 {
		return 0, errors.New("malformed bytes for runtime sized length")
	}
	return binary.NativeEndian.Uint16(lengthBytes), nil
}

func parseIndividualValue(paramType byte, paramValueBytes []byte) string {
	switch reflect.Kind(paramType) {
	case reflect.Uint8:
		return fmt.Sprintf("%d", uint8(paramValueBytes[0]))
	case reflect.Int8:
		return fmt.Sprintf("%d", int8(paramValueBytes[0]))
	case reflect.Uint16:
		return fmt.Sprintf("%d", byteOrder.Uint16(paramValueBytes))
	case reflect.Int16:
		return fmt.Sprintf("%d", int16(byteOrder.Uint16(paramValueBytes)))
	case reflect.Uint32:
		return fmt.Sprintf("%d", byteOrder.Uint32(paramValueBytes))
	case reflect.Int32:
		return fmt.Sprintf("%d", int32(byteOrder.Uint32(paramValueBytes)))
	case reflect.Uint64:
		return fmt.Sprintf("%d", byteOrder.Uint64(paramValueBytes))
	case reflect.Int64:
		return fmt.Sprintf("%d", int64(byteOrder.Uint64(paramValueBytes)))
	case reflect.Uint:
		return fmt.Sprintf("%d", byteOrder.Uint64(paramValueBytes))
	case reflect.Int:
		return fmt.Sprintf("%d", int(byteOrder.Uint64(paramValueBytes)))
	case reflect.Pointer:
		return fmt.Sprintf("0x%X", byteOrder.Uint64(paramValueBytes))
	case reflect.String:
		return string(paramValueBytes)
	case reflect.Bool:
		if paramValueBytes[0] == 1 {
			return "true"
		} else {
			return "false"
		}
	case ditypes.KindUnsupported:
		return "UNSUPPORTED"
	default:
		return ""
	}
}
