package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"gopkg.in/yaml.v2"
)

// define YAML struct as in config file
type Image struct {
	Name      string `yaml:"name"`
	Tag       string `yaml:"tag"`
	SrcSHA256 string `yaml:"srcsha,omitempty"`
	DstSHA256 string `yaml:"dstsha,omitempty"`
}

type ImagesDefinition struct {
	Image []Image `yaml:"images"`
}

func readImagesFromYaml(conf string) []Image {
	// Read the YAML file into a byte slice
	yamlFile, err := ioutil.ReadFile(conf)
	if err != nil {
		log.Fatal(err)
	}

	// Create an instance of the struct to unmarshal the YAML data into
	var imgdef ImagesDefinition
	var images []Image

	// Unmarshal the YAML data into the struct
	err = yaml.Unmarshal(yamlFile, &imgdef)
	if err != nil {
		log.Fatal(err)
	}

	// Access the data in the struct
	for _, image := range imgdef.Image {
		// fmt.Println(image.Name, image.Tag)
		images = append(images, image)
	}
	return images
}

func pullImage(imageURL string, cli *client.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	ctx := context.Background()
	reader, err := cli.ImagePull(ctx, imageURL, types.ImagePullOptions{})
	if err != nil {
		fmt.Printf("Error pulling image %s: %v\n", imageURL, err)
		return
	}
	defer reader.Close()

	_, _ = io.Copy(os.Stdout, reader)

	fmt.Printf("Image %s has been pulled successfully.\n", imageURL)
}

func getSHA256checksum(imageName string, imageTag string) {
	// Create a reference to the image.
	ref, err := name.ParseReference(fmt.Sprintf("%s:%s", imageName, imageTag))
	if err != nil {
		fmt.Printf("Failed to parse image reference: %v\n", err)
		return
	}
	// Fetch the image information.
	imgInfo, err := remote.Get(ref)
	if err != nil {
		fmt.Printf("Failed to fetch image information: %v\n", err)
		return
	}

	// Extract the image digest (SHA) from the image information.
	fmt.Printf("Image Digest (SHA): %s\n", imgInfo.Digest)
}

func main() {
	// registriesmapping := map[string]string {
	// 	"docker.io" : "",
	// 	"quay.io" : "",
	// }
	processimages := readImagesFromYaml("config.yaml")
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Printf("Error creating Docker client: %v\n", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup

	for _, image := range processimages {
		wg.Add(1)
		imageURL := image.Name + ":" + image.Tag
		fmt.Println(imageURL)
		go pullImage(imageURL, cli, &wg)
	}

	wg.Wait()
	fmt.Println("All images pulled successfully.")
}
