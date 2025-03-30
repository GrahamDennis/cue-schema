package cmd_test

import (
	"strconv"
	"testing"

	"cuelang.org/go/cue/cuecontext"
	"github.com/GrahamDennis/cue-schema/cmd"
)

func TestBreakingChange(t *testing.T) {
	type breakingChangeTest struct {
		compatible bool
		old        string
		new        string
	}

	testCases := []breakingChangeTest{
		// Add new message is compatible
		0: {compatible: true, old: `#schema: messages: {}`, new: `#schema: messages: foo?: string`},
		// Remove message is not compatible
		1: {compatible: false, old: `#schema: messages: foo?: string`, new: `#schema: messages: {}`},
		// Add new enum option is compatible
		2: {compatible: true, old: `#schema: enums: enum1: { value1?: 1}`, new: `#schema: enums: enum1: {value1?: 1, value2?: 2}`},
		// Removing an enum option is not compatible
		3: {compatible: false, old: `#schema: enums: enum1: { value1?: 1, value2?: 2}`, new: `#schema: enums: enum1: {value1?: 1}`},
		// Adding an optional field to a message is compatible
		4: {compatible: true, old: `#schema: messages: message1?: { field1: int }`, new: `#schema: messages: message1?: { field1: int, field2?: int}`},
		// Adding a required field to a message is not compatible
		5: {compatible: false, old: `#schema: messages: message1?: { field1: int }`, new: `#schema: messages: message1?: { field1: int, field2: int}`},
		// Removing an optional field from a message is not compatible
		6: {compatible: false, old: `#schema: messages: message1?: { field1: int, field2?: int}`, new: `#schema: messages: message1?: { field1: int }`},
		// Removing a required field from a message is not compatible
		7: {compatible: false, old: `#schame: messages: message1?: { field1: int, field2: int}`, new: `#schema: messages: message1?: { field1: int }`},
		// Defining enums when they weren't defined before is compatible
		8: {compatible: true, old: `#schema: messages: {}`, new: `#schema: { messages: {}, enums?: enum1: {value1?: 1 } }`},
	}

	for i, tc := range testCases {
		if tc.old == "" {
			continue
		}

		key := tc.old + " ⊑ " + tc.new
		t.Run(strconv.Itoa(i)+"/"+key, func(t *testing.T) {
			ctx := cuecontext.New()
			oldValue := ctx.CompileString(tc.old)
			newValue := ctx.CompileString(tc.new)

			err := cmd.IsBackwardsCompatible(oldValue, newValue)
			got := err == nil

			if got != tc.compatible {
				t.Errorf(`IsBackwardsCompatible(%q, %q) = %v, want %v; (err = %v)`, tc.old, tc.new, got, tc.compatible, err)
			}
		})
	}
}
