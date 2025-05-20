/*
Copyright © 2025 Graham Dennis <graham.dennis@gmail.com>
*/
package cmd

import (
	"cmp"
	"fmt"
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
	oldSchemaFilenames       []string
	newSchemaFilenames       []string
	cuePath                  string
	maxErrors                int = 0
	suppressErrorsLongerThan int = 0
)

// breakingCmd represents the breaking command
var breakingCmd = &cobra.Command{
	Use:           "breaking",
	Short:         "Validate if a schema change is backwards-compatible",
	Long:          `Validate if a schema change is backwards-compatible.`,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		err := RunBreakingChangeDetection(oldSchemaFilenames, newSchemaFilenames, cuePath)
		if err != nil {
			printError(err)
			os.Exit(1)
		}
	},
}

func RunBreakingChangeDetection(oldSchemaFilenames []string, newSchemaFilenames []string, cuePath string) error {
	ctx := cuecontext.New()

	oldValue := loadSchemas(ctx, oldSchemaFilenames, cuePath)
	if err := oldValue.Err(); err != nil {
		return err
	}

	newValue := loadSchemas(ctx, newSchemaFilenames, cuePath)
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

	breakingCmd.Flags().StringArrayVar(&oldSchemaFilenames, "old", []string{}, "old CUE schema files")
	breakingCmd.MarkFlagRequired("old")
	breakingCmd.MarkFlagFilename("old", "cue")

	breakingCmd.Flags().StringArrayVar(&newSchemaFilenames, "new", []string{}, "new CUE schema files")
	breakingCmd.MarkFlagRequired("new")
	breakingCmd.MarkFlagFilename("new", "cue")

	breakingCmd.Flags().StringVar(&cuePath, "path", "", "CUE path that contains the schema to validate in the CUE files")

	breakingCmd.Flags().IntVar(&maxErrors, "max-errors", 0, "maximum number of errors to report")
	breakingCmd.Flags().IntVar(&suppressErrorsLongerThan, "suppress-errors-longer-than", 0, "suppress errors longer than this length")
}

func loadSchema(ctx *cue.Context, filename string, cuePath string) cue.Value {
	insts := load.Instances([]string{filename}, &load.Config{
		Dir: filepath.Join(),
		Env: []string{}, // or nil to use os.Environ
	})

	rootValue := ctx.BuildInstance(insts[0])

	return rootValue.LookupPath(cue.ParsePath(cuePath))
}

func loadSchemas(ctx *cue.Context, filenames []string, cuePath string) cue.Value {
	value := ctx.CompileString("_")
	for _, filename := range filenames {
		value = value.Unify(loadSchema(ctx, filename, cuePath))
		if err := value.Err(); err != nil {
			return value
		}
	}
	return value
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
	printErrors(os.Stderr, err, &errors.Config{
		Format: format,
		Cwd:    pwd,
	})
}

func getLang() language.Tag {
	loc := cmp.Or(os.Getenv("LC_ALL"), os.Getenv("LANG"))
	loc, _, _ = strings.Cut(loc, ".")
	return language.Make(loc)
}

func printErrors(w io.Writer, err error, cfg *errors.Config) {
	errs := errors.Errors(err)
	maxErrorsToDisplay := len(errs)
	if maxErrors > 0 {
		maxErrorsToDisplay = min(maxErrorsToDisplay, maxErrors)
	}
	for _, e := range errs[0:maxErrorsToDisplay] {
		if suppressErrorsLongerThan > 0 && len(e.Error()) > suppressErrorsLongerThan {
			fmt.Fprintf(w, "Suppressed error longer than %d characters\n", suppressErrorsLongerThan)
		} else {
			errors.Print(w, e, cfg)
		}
	}
}
