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
		Usage:     "Extract ansible_host entries from a Terraform state JSON",
		Version:   version,
		ArgsUsage: "--input <file> [--format json|ini|txt|ansible] [--host-field path] [--ip-field path]",
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
				Value:   "json",
				Usage:   "Output format: json, ini, txt, ansible",
			},
			&cli.StringFlag{
				Name:  "host-field",
				Value: "values.name",
				Usage: "Dot-path to hostname in each object (e.g. values.name)",
			},
			&cli.StringFlag{
				Name:  "ip-field",
				Value: "values.variables.ip",
				Usage: "Dot-path to IP in each object (e.g. values.variables.ip)",
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

			// 2) Extract ansible_host objects
			objects := parser.ExtractAnsibleHosts(data)
			if len(objects) == 0 {
				// still valid: empty slice ⇒ outputs "[]" or nothing in ansible mode
			}

			// 3) Dispatch output
			format := strings.ToLower(c.String("format"))
			switch format {
			case "ansible":
				return iohandler.OutputAnsibleInventory(
					objects,
					c.String("host-field"),
					c.String("ip-field"),
				)
			default:
				return iohandler.Output(objects, format)
			}
		},
		CustomAppHelpTemplate: `{{.Name}} {{.Version}}

{{.Usage}}

USAGE:
   {{.HelpName}} {{.ArgsUsage}}

FLAGS:
{{range .VisibleFlags}}{{.}}
{{end}}
EXAMPLES:
   # JSON array of all ansible_host objects
   {{.HelpName}} --input terraform_state.json
   # INI-style output
   {{.HelpName}} -i terraform_state.json -f ini
   # One-line Ansible inventory
   {{.HelpName}} -i terraform_state.json -f ansible
`,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("ERROR: %v\n", err)
	}
}
