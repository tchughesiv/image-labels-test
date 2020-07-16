package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/containers/buildah"
	buildahcli "github.com/containers/buildah/pkg/cli"
	"github.com/containers/buildah/pkg/parse"
	"github.com/containers/storage"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type runInputOptions struct {
	addHistory     bool
	capAdd         []string
	capDrop        []string
	hostname       string
	isolation      string
	runtime        string
	runtimeFlag    []string
	noPivot        bool
	securityOption []string
	terminal       bool
	volumes        []string
	mounts         []string
	*buildahcli.NameSpaceResults
}

func main() {
	debug := true
	logrus.SetLevel(logrus.ErrorLevel)
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if err := runCmd(&cobra.Command{}, os.Args[1:], runInputOptions{}); err != nil {
		logrus.Error(err)
	}
}

// buildah logic
func runCmd(c *cobra.Command, args []string, iopts runInputOptions) error {
	name := args[0]

	storeOptions, err := storage.DefaultStoreOptionsAutoDetectUID()
	if err != nil {
		return err
	}
	store, err := storage.GetStore(storeOptions)
	if err != nil {
		return err
	}
	defer func() {
		if _, err := store.Shutdown(false); err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
	}()

	builder, err := openBuilder(context.TODO(), store, name)
	if err != nil {
		return errors.Wrapf(err, "error reading build container %q", name)
	}

	isolation, err := parse.IsolationOption(c.Flag("isolation").Value.String())
	if err != nil {
		return err
	}

	runtimeFlags := []string{}
	for _, arg := range iopts.runtimeFlag {
		runtimeFlags = append(runtimeFlags, "--"+arg)
	}

	noPivot := iopts.noPivot || (os.Getenv("BUILDAH_NOPIVOT") != "")

	namespaceOptions, networkPolicy, err := parse.NamespaceOptions(c)
	if err != nil {
		return errors.Wrapf(err, "error parsing namespace-related options")
	}

	options := buildah.RunOptions{
		Hostname:         iopts.hostname,
		Runtime:          iopts.runtime,
		Args:             runtimeFlags,
		NoPivot:          noPivot,
		User:             "",
		Isolation:        isolation,
		NamespaceOptions: namespaceOptions,
		ConfigureNetwork: networkPolicy,
		CNIPluginPath:    iopts.CNIPlugInPath,
		CNIConfigDir:     iopts.CNIConfigDir,
		AddCapabilities:  iopts.capAdd,
		DropCapabilities: iopts.capDrop,
	}

	mounts, err := parse.GetVolumes(iopts.volumes, iopts.mounts)
	if err != nil {
		return err
	}
	options.Mounts = mounts

	runerr := builder.Run(args, options)
	if runerr != nil {
		logrus.Debugf("error running %v in container %q: %v", args, builder.Container, runerr)
	}
	if runerr == nil {
		shell := "/bin/sh -c"
		if len(builder.Shell()) > 0 {
			shell = strings.Join(builder.Shell(), " ")
		}
		conditionallyAddHistory(builder, c, "%s %s", shell, strings.Join(args, " "))
		return builder.Save()
	}
	return runerr
}

func openBuilder(ctx context.Context, store storage.Store, name string) (builder *buildah.Builder, err error) {
	if name != "" {
		builder, err = buildah.OpenBuilder(store, name)
		if os.IsNotExist(errors.Cause(err)) {
			options := buildah.ImportOptions{
				Container: name,
			}
			builder, err = buildah.ImportBuilder(ctx, store, options)
		}
	}
	if err != nil {
		return nil, errors.Wrapf(err, "error reading build container")
	}
	if builder == nil {
		return nil, errors.Errorf("error finding build container")
	}
	return builder, nil
}

/*
func inspectCmd(args []string) error {
	var builder *buildah.Builder

	if len(args) == 0 {
		return errors.Errorf("image name must be specified")
	}
	if len(args) > 1 {
		return errors.Errorf("too many arguments specified")
	}

	//systemContext := &types.SystemContext{
	//	OSChoice: "linux",
	//}
	systemContext := &types.SystemContext{}

	name := args[0]

	storeOptions, err := storage.DefaultStoreOptionsAutoDetectUID()
	if err != nil {
		return err
	}
	store, err := storage.GetStore(storeOptions)
	if err != nil {
		return err
	}
	defer func() {
		if _, err := store.Shutdown(false); err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
	}()

	ctx := context.TODO()

	builder.Run([]string{"inspect", args[0]}, buildah.RunOptions{Isolation: buildah.IsolationChroot})
	builder, err = openImage(ctx, systemContext, store, name)
	if err != nil {
		return errors.Wrapf(err, "error reading image %q", name)
	}
	out := buildah.GetBuildInfo(builder)
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		enc.SetEscapeHTML(false)
	}
	return enc.Encode(out)
}

func openImage(ctx context.Context, sc *types.SystemContext, store storage.Store, name string) (builder *buildah.Builder, err error) {
	options := buildah.ImportFromImageOptions{
		Image:         name,
		SystemContext: sc,
	}
	builder, err = buildah.ImportBuilderFromImage(ctx, store, options)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading image")
	}
	if builder == nil {
		return nil, errors.Errorf("error mocking up build configuration")
	}
	return builder, nil
}
*/
// podman logic
// issues with overlay mounting requirement within a pod
// test product type/version aggregation from pod image lookup
/*
func imageLookup(args []string) (retErr error) {
	// ERRO[0000] error opening "/var/lib/containers/storage/storage.lock": permission denied
	// https://github.com/containers/storage/blob/ed28f2457e2f57cb3d3f2f4029a85f72b35370ab/store.go#L656-L664
	storeOptions, err := storage.DefaultStoreOptionsAutoDetectUID()
	if err != nil {
		return err
	}
	println("graphroot is " + storeOptions.GraphRoot)
	println("graphdriver is " + storeOptions.GraphDriverName)
	ir, err := image.NewImageRuntimeFromOptions(storeOptions)
	if err != nil {
		return err
	}
	for _, imgName := range args {
		img, err := ir.NewFromLocal(imgName)
		if err != nil {
			return err
		}
		println()
		println(img.InputName)
		for _, name := range img.Names() {
			println(name)
		}

		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, time.Duration(5)*time.Second)
		defer cancel()
		// get inspect image data
		imgData, err := img.Inspect(ctx)
		if err != nil {
			return err
		}
		println(imgData.ID)
		if imgData.Labels != nil {
			println("IMAGE LABELS:")
			for key, val := range imgData.Labels {
				if strings.Contains(key, "org.jboss.") {
					println(key + "=" + val)
				}
			}
			for key, val := range imgData.Labels {
				if strings.Contains(key, "com.redhat.") {
					println(key + "=" + val)
				}
			}
		}
	}
	println()
	return retErr
}
*/

func conditionallyAddHistory(builder *buildah.Builder, c *cobra.Command, createdByFmt string, args ...interface{}) {
	history := buildahcli.DefaultHistory()
	if c.Flag("add-history").Changed {
		history, _ = c.Flags().GetBool("add-history")
	}
	if history {
		now := time.Now().UTC()
		created := &now
		createdBy := fmt.Sprintf(createdByFmt, args...)
		builder.AddPrependedEmptyLayer(created, createdBy, "", "")
	}
}
