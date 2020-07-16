package main

import (
	"strings"

	"github.com/containers/buildah"
	buildahcli "github.com/containers/buildah/pkg/cli"
	"github.com/containers/buildah/pkg/parse"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type inspectResults struct {
	format string
}

func init() {
	var (
		opts               inspectResults
		inspectDescription = "\n  Inspects a build container's or built image's configuration."
	)

	inspectCommand := &cobra.Command{
		Use:   "inspect",
		Short: "Inspect the configuration of a container or image",
		Long:  inspectDescription,
		RunE: func(cmd *cobra.Command, args []string) error {
			return inspectCmd(cmd, args, opts)
		},
		Example: `buildah inspect imageID
  buildah inspect --format '{{.OCIv1.Config.Env}}' alpine`,
	}
	inspectCommand.SetUsageTemplate(UsageTemplate())

	flags := inspectCommand.Flags()
	flags.SetInterspersed(false)
	flags.StringVarP(&opts.format, "format", "f", "", "use `format` as a Go template to format the output")

	rootCmd.AddCommand(inspectCommand)
}

func inspectCmd(c *cobra.Command, args []string, iopts inspectResults) error {
	var builder *buildah.Builder

	if len(args) == 0 {
		return errors.Errorf("image name must be specified")
	}
	if err := buildahcli.VerifyFlagsArgsOrder(args); err != nil {
		return err
	}
	if len(args) > 1 {
		return errors.Errorf("too many arguments specified")
	}

	systemContext, err := parse.SystemContextFromOptions(c)
	if err != nil {
		return errors.Wrapf(err, "error building system context")
	}

	name := args[0]

	store, err := getStore(c)
	if err != nil {
		return err
	}
	println(store.GraphOptions())

	ctx := getContext()

	builder, err = openImage(ctx, systemContext, store, name)
	if err != nil {
		return errors.Wrapf(err, "error reading image %q", name)
	}
	parseFindings(builder)
	return nil
}

func parseFindings(builder *buildah.Builder) {
	println(builder.FromImage)
	println(builder.FromImageID)
	println()
	ociConfig := builder.OCIv1.Config
	if ociConfig.Labels != nil {
		println("IMAGE LABELS:")
		for key, val := range ociConfig.Labels {
			if strings.Contains(key, "org.jboss.") {
				println(key + "=" + val)
			}
		}
		for key, val := range ociConfig.Labels {
			if strings.Contains(key, "com.redhat.") {
				println(key + "=" + val)
			}
		}
	}
	println()
}
