package cmd

import (
	"errors"
	"io"

	"github.com/mmianl/api/v2/cmd/admin"
	"github.com/mmianl/api/v2/cmd/build"
	"github.com/mmianl/api/v2/cmd/user"
	"github.com/spf13/cobra"
)

func New(out io.Writer, in io.Reader, args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "Runs the api server.",
		Long:  `This is the entrypoint for the API server. There are two apis available: user and admin.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("no additional command provided")
		},
		Version: build.Version,
	}

	cmd.AddCommand(
		user.New(),
		admin.New(),
	)
	cmd.InitDefaultVersionFlag()

	return cmd
}
