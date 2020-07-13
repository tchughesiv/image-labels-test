package main

import (
	"context"
	"fmt"

	"github.com/containers/image/v5/image"
	"github.com/containers/image/v5/transports/alltransports"
	imagetypes "github.com/containers/image/v5/types"
	pkgerrors "github.com/pkg/errors"
)

func main() {
	println("test")
	if err := imageLookup(); err != nil {
		println(err)
	}
}

// test product type/version aggregation from pod image lookup
func imageLookup() (retErr error) {
	images := []string{"quay.io/crio/redis:alpine", "quay.io/crio/redis@sha256:1780b5a5496189974b94eb2595d86731d7a0820e4beb8ea770974298a943ed55"}
	sys := &imagetypes.SystemContext{
		OSChoice:                    "linux",
		SystemRegistriesConfDirPath: "/etc/containers",
	}
	imagesCount := uniqueCount(images)
	for img, num := range imagesCount {
		fmt.Println(img)
		fmt.Printf("image is used %d time(s)\n", num)
		imgSrc, err := parseImageSource(context.TODO(), sys, "containers-storage:"+img)
		if err != nil {
			return err
		}
		defer func() {
			if err := imgSrc.Close(); err != nil {
				retErr = pkgerrors.Wrapf(retErr, fmt.Sprintf("(could not close image: %v) ", err))
			}
		}()

		img, err := image.FromUnparsedImage(context.TODO(), sys, image.UnparsedInstance(imgSrc, nil))
		if err != nil {
			return fmt.Errorf("Error parsing manifest for image: %v", err)
		}

		config, err := img.OCIConfig(context.TODO())
		if err != nil {
			return fmt.Errorf("Error reading OCI-formatted configuration data: %v", err)
		}
		fmt.Println()
		fmt.Println("IMAGE LABELS -")
		for key, val := range config.Config.Labels {
			if key == "org.jboss.product" {
				fmt.Println(key + "=" + val)
			}
		}
	}
	fmt.Println()
	return retErr
}

func parseImageSource(ctx context.Context, sys *imagetypes.SystemContext, name string) (imagetypes.ImageSource, error) {
	ref, err := alltransports.ParseImageName(name)
	if err != nil {
		return nil, err
	}
	return ref.NewImageSource(ctx, sys)
}

func uniqueCount(intSlice []string) map[string]int {
	keys := make(map[string]bool)
	list := map[string]int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list[entry] = 1
		} else {
			list[entry] = list[entry] + 1
		}
	}
	return list
}
