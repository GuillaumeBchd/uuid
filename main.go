package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

const description = `
Generates a UUID v4 or v5 if arguments are given.
		
Examples:
	uuid:                        generates a UUID v4 (random)
	uuid foo bar:                generates a UUID v5 using no namespace and the string "foobar"
	uuid -n <namespace> foo bar: generates a UUID v5 using the namespace "<namespace>" and the string "foobar"
`

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	app := &cli.App{
		Name:        "uuid",
		Usage:       "generates a UUID",
		UsageText:   "uuid [options] [args...]",
		Description: formatDescription(description),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "namespace",
				Usage:   `The namespace UUID to use when generating a v5 UUID.`,
				Aliases: []string{"n"},
			},
		},
		Action: run,
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal().Err(err).Send()
	}
}

func run(ctx *cli.Context) error {
	id := uuid.New()

	if ctx.Args().Len() > 0 {
		namespace := uuid.Nil
		if ns := ctx.String("namespace"); ns != "" {
			var err error
			namespace, err = uuid.Parse(ns)
			if err != nil {
				return err
			}
		}

		id = uuid.NewSHA1(namespace, []byte(strings.Join(ctx.Args().Slice(), "")))
	}

	fmt.Println(id.String())

	return nil
}

func formatDescription(s string) string {
	return strings.Trim(s, "\n \t")
}
