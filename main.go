package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
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

type Person struct {
	Name string `yaml:"name"`
	Age  int    `yaml:"age"`
	City string `yaml:"city"`
}

func readConfigYaml() {
	// Read the YAML file into a byte slice
	yamlFile, err := ioutil.ReadFile("example.yaml")
	if err != nil {
		fmt.Printf("Error reading YAML file: %v\n", err)
		return
	}

	// Create an instance of the struct to unmarshal the YAML data into
	var person Person

	// Unmarshal the YAML data into the struct
	err = yaml.Unmarshal(yamlFile, &person)
	if err != nil {
		fmt.Printf("Error unmarshaling YAML: %v\n", err)
		return
	}

	// Access the data in the struct
	fmt.Printf("Name: %s\n", person.Name)
	fmt.Printf("Age: %d\n", person.Age)
	fmt.Printf("City: %s\n", person.City)
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
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Printf("Error creating Docker client: %v\n", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup

	for _, imageURL := range imageURLs {
		wg.Add(1)
		fmt.Println(imageURL)
		go pullImage(imageURL, cli, &wg)
	}

	wg.Wait()
	fmt.Println("All images pulled successfully.")
}



 server:
#  address: "127.0.0.1"
#  port: 8080
#
#database:
#  host: "localhost"
#  port: 3306
#  username: "user"
#  password: "password"
#And you want to unmarshal only the server section into a Go struct:
#
#go
#Copy code
#package main
#
#import (
#	"fmt"
#	"gopkg.in/yaml.v2"
#	"io/ioutil"
#)
#
#type ServerConfig struct {
#	Address string `yaml:"address"`
#	Port    int    `yaml:"port"`
#}
#
#func main() {
#	// Read the YAML file into a byte slice
#	yamlFile, err := ioutil.ReadFile("config.yaml")
#	if err != nil {
#		fmt.Printf("Error reading YAML file: %v\n", err)
#		return
#	}
#
#	// Create an instance of the struct to unmarshal the YAML data into
#	var serverConfig ServerConfig
#
#	// Unmarshal the YAML data into the struct
#	err = yaml.Unmarshal(yamlFile, &serverConfig)
#	if err != nil {
#		fmt.Printf("Error unmarshaling YAML: %v\n", err)
#		return
#	}
#
#	// Access the data in the struct
#	fmt.Printf("Server Address: %s\n", serverConfig.Address)
#	fmt.Printf("Server Port: %d\n", serverConfig.Port)
#}
#In this code, we unmarshal only the server section of the YAML file by defining a ServerConfig struct that corresponds to that section. We use the same yaml.Unmarshal function, but it will only populate the fields specified in the ServerConfig struct. The result will be that the serverConfig struct contains only the data from the server section of the YAML file.
#
#Make sure to replace "config.yaml" with the path to your YAML file.
