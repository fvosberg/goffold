package main

import (
	"errors"
	"fmt"

	"github.com/fvosberg/goffold/internal/templates"
	"github.com/spf13/cobra"
)


var (
	newCmd = &cobra.Command{
		Use:   "new go-service",
		Short: "creates a new package",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := templates.ParseTo(args[0], args[1])
			if err != nil {
				return fmt.Errorf("Failed: %w", err)
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("missing template name")
			}
			if args[0] != "go-service" {
				return fmt.Errorf("wrong template %q, expected %q", args[0], "go-service")
			}
			moduleName, _ := cmd.PersistentFlags().GetString("name")
			if moduleName == "" {
				return errors.New("required flag name not set")
			}
			dockerImage, _ := cmd.PersistentFlags().GetString("docker")
			if dockerImage == "" {
				return errors.New("required flag docker not set")
			}
			if len(args) < 2 {
				return errors.New("missing path argument")
			}
			if len(args) > 2 {
				return errors.New("to many arguments")
			}
			return nil
		},
	}
)
