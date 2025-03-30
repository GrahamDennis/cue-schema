/*
Copyright © 2025 Graham Dennis <graham.dennis@gmail.com>
*/
package cmd

import (
	"cmp"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/load"
)

var (
	oldSchemaFilename string
	newSchemaFilename string
	cuePath           string
)

// breakingCmd represents the breaking command
var breakingCmd = &cobra.Command{
	Use:           "breaking",
	Short:         "Validate if a schema change is backwards-compatible",
	Long:          `Validate if a schema change is backwards-compatible.`,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		err := RunBreakingChangeDetection(oldSchemaFilename, newSchemaFilename, cuePath)
		if err != nil {
			printError(err)
			os.Exit(1)
		}
	},
}

func RunBreakingChangeDetection(oldSchemaFilename string, newSchemaFilename string, cuePath string) error {
	ctx := cuecontext.New()

	oldValue := loadSchema(ctx, oldSchemaFilename)
	if err := oldValue.Err(); err != nil {
		return err
	}
	newValue := loadSchema(ctx, newSchemaFilename)
	if err := newValue.Err(); err != nil {
		return err
	}

	return IsBackwardsCompatible(oldValue, newValue)
}

func IsBackwardsCompatible(oldValue cue.Value, newValue cue.Value) error {
	return newValue.Subsume(oldValue)
}

func init() {
	rootCmd.AddCommand(breakingCmd)

	breakingCmd.Flags().StringVar(&oldSchemaFilename, "old", "", "old CUE schema file")
	breakingCmd.MarkFlagRequired("old")
	breakingCmd.MarkFlagFilename("old", "cue")

	breakingCmd.Flags().StringVar(&newSchemaFilename, "new", "", "new CUE schema file")
	breakingCmd.MarkFlagRequired("new")
	breakingCmd.MarkFlagFilename("new", "cue")

	breakingCmd.Flags().StringVar(&cuePath, "path", "", "CUE path that contains the schema to validate in the CUE files")
}

func loadSchema(ctx *cue.Context, filename string) cue.Value {
	insts := load.Instances([]string{filename}, &load.Config{
		Dir: filepath.Join(),
		Env: []string{}, // or nil to use os.Environ
	})

	rootValue := ctx.BuildInstance(insts[0])

	return rootValue.LookupPath(cue.ParsePath(cuePath))
}

func printError(err error) {
	if err == nil {
		return
	}

	// Link x/text as our localizer.
	p := message.NewPrinter(getLang())
	format := func(w io.Writer, format string, args ...any) {
		p.Fprintf(w, format, args...)
	}
	pwd, _ := os.Getwd()
	errors.Print(os.Stderr, err, &errors.Config{
		Format: format,
		Cwd:    pwd,
	})
}

func getLang() language.Tag {
	loc := cmp.Or(os.Getenv("LC_ALL"), os.Getenv("LANG"))
	loc, _, _ = strings.Cut(loc, ".")
	return language.Make(loc)
}
