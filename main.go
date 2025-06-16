package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/HilkopterBob/terraform-ansible-inventory/internal/iohandler"
	"github.com/HilkopterBob/terraform-ansible-inventory/internal/parser"
	"github.com/urfave/cli/v2"
)

const version = "v1.0.0"

func main() {
	app := &cli.App{
		Name:      "terraform-ansible-inventory",
		Usage:     "Generate an Ansible inventory from a Terraform state produced by the ansible/ansible provider",
		Version:   version,
		ArgsUsage: "--input <file> [--format yaml|ini|json]",
		// TODO: Add CLI flags for filtering hosts/groups and for
		// selecting output formats like YAML or INI according to the
		// full capabilities of the ansible/ansible provider.
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Aliases:  []string{"i"},
				Usage:    "Path to input JSON file (or '-' for stdin)",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Value:   "yaml",
				Usage:   "Output format: yaml, ini, or json",
			},
		},
		Action: func(c *cli.Context) error {
			// 1) Read input
			var data []byte
			var err error
			path := c.String("input")
			if path == "-" {
				data, err = io.ReadAll(os.Stdin)
			} else {
				data, err = os.ReadFile(path)
			}
			if err != nil {
				return fmt.Errorf("failed to read %q: %w", path, err)
			}

			// 2) Parse inventory from Terraform state
			inv := parser.ParseInventory(data)

			// 3) Dispatch output
			format := strings.ToLower(c.String("format"))
			// TODO: support additional output styles (e.g. split
			// inventories by group) once provider parity is
			// implemented.
			return iohandler.OutputInventory(inv, format)
		},
		CustomAppHelpTemplate: `{{.Name}} {{.Version}}

{{.Usage}}

USAGE:
   {{.HelpName}} {{.ArgsUsage}}

FLAGS:
{{range .VisibleFlags}}{{.}}
{{end}}
EXAMPLES:
   # YAML inventory
   {{.HelpName}} --input terraform_state.json -f yaml
   # INI inventory
   {{.HelpName}} -i terraform_state.json -f ini
`,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("ERROR: %v\n", err)
	}
}
