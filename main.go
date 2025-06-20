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
			&cli.StringSliceFlag{
				Name:  "host",
				Usage: "Only include the specified host(s)",
			},
			&cli.StringSliceFlag{
				Name:  "group",
				Usage: "Only include hosts belonging to the specified group(s)",
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
			inv, err := parser.ParseInventory(data)
			if err != nil {
				return err
			}

			// 3) Apply filters
			hosts := c.StringSlice("host")
			groups := c.StringSlice("group")
			if len(hosts) > 0 || len(groups) > 0 {
				inv = inv.CopyFiltered(hosts, groups)
			}

			// 4) Dispatch output
			format := strings.ToLower(c.String("format"))
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
