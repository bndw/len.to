package main

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
	xurls "mvdan.cc/xurls/v2"
)

const (
	posts  = "../content/post"
	images = "../content/img"
)

var (
	ErrFileNotExist  = errors.New("file does not exist")
	ErrYamlMalformed = errors.New("unable to parse yaml")
)

func main() {
	// Build dgraph RDF
	fmt.Println("{")
	fmt.Println("set {")
	walk(images, imgRDF)
	walk(images, imgTagRDF)
	fmt.Println("}")
	fmt.Println("}")
}

// imgRDF prints an image as rdf triples
func imgRDF(img *Post) {
	node := nodeName(img.Title)

	// img node
	rdfObjString(node, "dgraph.type", "Image")
	rdfObjString(node, "name", img.Title)
	rdfObjString(node, "url", img.URLs()[0])
	rdfObjString(node, "location", img.Location)
	rdfObjString(node, "date", img.Date.Format(time.RFC3339))

	// edges to tags
	for _, tag := range img.Tags {
		rdfObjNode(node, "tagged", nodeName(tag))
	}
}

// imgTagRDF prints an image's tags as rdf triples
func imgTagRDF(img *Post) {
	for _, tag := range img.Tags {
		node := nodeName(tag)
		rdfObjString(node, "dgraph.type", "Tag")
		rdfObjString(node, "name", tag)
	}
}

func nodeName(s string) string {
	h := sha1.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func rdfObjString(s, p, o string) {
	fmt.Printf("_:%s <%s> \"%s\" .\n", s, p, o)
}
func rdfObjNode(s, p, o string) {
	fmt.Printf("_:%s <%s> _:%s .\n", s, p, o)
}

func walk(dir string, fn func(p *Post)) {
	err := filepath.Walk(dir, func(path string, i os.FileInfo, err error) error {
		filename := i.Name()
		if i.IsDir() {
			return nil
		}

		p, err := LoadPost(filepath.Join(dir, filename))
		if err != nil {
			return err
		}

		fn(p)
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

func (p *Post) URLs() []string {
	u := xurls.Strict()
	return u.FindAllString(string(p.Content), -1)
}

func (p *Post) String() string {
	return fmt.Sprintf(`%s
---
title: %s
date: %s
draft: %t
location: %s
tags: %q
---
%s
`,
		p.Path,
		p.Title,
		p.Date.String(),
		p.Draft,
		p.Location,
		p.Tags,
		string(p.Content),
	)
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
