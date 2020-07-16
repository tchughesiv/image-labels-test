package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/containers/buildah"
	"github.com/containers/image/v5/types"
	"github.com/containers/storage"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	debug := true
	logrus.SetLevel(logrus.ErrorLevel)
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if err := inspectCmd(os.Args[1:]); err != nil {
		logrus.Error(err)
	}
}

// buildah logic
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
