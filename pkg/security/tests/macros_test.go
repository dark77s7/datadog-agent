// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build linux && functionaltests

// Package tests holds tests related files
package tests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DataDog/datadog-agent/pkg/security/secl/model"
	"github.com/DataDog/datadog-agent/pkg/security/secl/rules"
)

func TestMacros(t *testing.T) {
	SkipIfNotAvailable(t)

	macros := []*rules.MacroDefinition{
		{
			ID:         "testmacro",
			Expression: `"{{.Root}}/test-macro"`,
		},
		{
			ID:         "testmacro2",
			Expression: `["{{.Root}}/test-macro"]`,
		},
	}

	ruleDefs := []*rules.RuleDefinition{
		{
			ID:         "test_rule",
			Expression: `testmacro in testmacro2 && mkdir.file.path in testmacro2`,
		},
	}

	test, err := newTestModule(t, macros, ruleDefs)
	if err != nil {
		t.Fatal(err)
	}
	defer test.Close()

	testFile, _, err := test.Path("test-macro")
	if err != nil {
		t.Fatal(err)
	}

	test.WaitSignal(t, func() error {
		if err = os.Mkdir(testFile, 0777); err != nil {
			return err
		}
		return os.Remove(testFile)
	}, func(event *model.Event, _ *rules.Rule) {
		assert.Equal(t, "mkdir", event.GetType(), "wrong event type")
	})
}
