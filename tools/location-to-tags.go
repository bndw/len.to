package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	ErrFileNotExist  = errors.New("file does not exist")
	ErrYamlMalformed = errors.New("unable to parse yaml")
)

// Takes the Location and adds it to the head of the Tags for each post
func main() {
	const dir = "./content/post"

	err := filepath.Walk(dir, func(path string, i os.FileInfo, err error) error {
		filename := i.Name()
		if i.IsDir() {
			fmt.Printf("%s is dir, skipping\n", filename)
			return nil
		}

		p, err := LoadPost(filepath.Join(dir, filename))
		if err != nil {
			return err
		}

		if p.MetadataEncoding == "yaml" {
			p.Metadata.Tags = append([]string{p.Metadata.Location}, p.Metadata.Tags...)
			if err := p.Write(path); err != nil {
				return err
			}

			return nil
		}

		return nil
	})
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}

}

type Post struct {
	*Metadata
	MetadataEncoding string
	Path             string
	Content          []byte
}

type Metadata struct {
	Title    string    `yaml:"title"`
	Date     time.Time `yaml:"date"`
	Draft    bool      `yaml:"draft"`
	Location string    `yaml:"location"`
	Tags     []string  `yaml:"tags"`
}

func (p *Post) Write(path string) error {
	var body []byte

	switch p.MetadataEncoding {
	case "yaml":
		m, err := yaml.Marshal(p.Metadata)
		if err != nil {
			return err
		}
		// The yaml frontmatter is surrounded by ---
		body = append(body, []byte("---\n")...)
		body = append(body, m...)
		body = append(body, []byte("\n---")...)

	default:
		return fmt.Errorf("cannot write unknown MetadataEncoding %s", p.MetadataEncoding)
	}

	body = append(body, []byte("\n")...)
	body = append(body, p.Content...)
	body = append(body, []byte("\n")...)

	return ioutil.WriteFile(filepath.Clean(path), body, 0644)
}

func LoadPost(path string) (*Post, error) {
	b, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	p := Post{
		Metadata: &Metadata{},
		Path:     path,
	}

	if bytes.Contains(b, []byte("---")) {
		p.MetadataEncoding = "yaml"

		parts := bytes.SplitN(b, []byte("---"), 3)
		p.Content = parts[2]

		if err := yaml.Unmarshal(parts[1], p.Metadata); err != nil {
			return nil, ErrYamlMalformed
		}
	}

	return &p, nil
}

func LoadYaml(path string) (*Metadata, error) {
	b, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrFileNotExist
		}
		return nil, err
	}

	var x Metadata
	if err := yaml.Unmarshal(b, &x); err != nil {
		return nil, ErrYamlMalformed
	}

	return &x, nil
}

func WriteYaml(x *Metadata, path string) error {
	b, err := yaml.Marshal(x)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filepath.Clean(path), b, 0644); err != nil {
		return err
	}

	return nil
}
