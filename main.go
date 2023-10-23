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
	"gopkg.in/yaml.v2"
)

var imageURLs = []string{
	"docker.io/library/alpine:latest",
	"docker.io/library/ubuntu:latest",
	"docker.io/library/debian:latest",
	// Add more image URLs as needed
}

type Image struct {
	Name string `yaml:"name"`
	Tag  string `yaml:"tag"`
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
		fmt.Println(image.Name, image.Tag)
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

func main() {
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
