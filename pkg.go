package lento

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

const (
	posts  = "./content/post"
	images = "./content/img"
)

var (
	ErrFileNotExist  = errors.New("file does not exist")
	ErrYamlMalformed = errors.New("unable to parse yaml")
)

func Images() []*Content {
	var res []*Content
	walk(images, func(c *Content) {
		res = append(res, c)
	})
	return res
}

func Posts() []*Content {
	var res []*Content
	walk(posts, func(c *Content) {
		res = append(res, c)
	})
	return res
}

func walk(dir string, fn func(c *Content)) {
	err := filepath.Walk(dir, func(path string, i os.FileInfo, err error) error {
		filename := i.Name()
		if i.IsDir() {
			return nil
		}

		c, err := LoadContent(filepath.Join(dir, filename))
		if err != nil {
			return err
		}

		fn(c)
		return nil
	})
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}
}

type Content struct {
	*Metadata
	Path string
	Body []byte
}

type Metadata struct {
	Title    string    `yaml:"title"`
	Date     time.Time `yaml:"date"`
	Draft    bool      `yaml:"draft"`
	ImgURL   string    `yaml:"img_url"`
	Location string    `yaml:"location"`
	Tags     []string  `yaml:"tags"`
}

func (c *Content) String() string {
	return fmt.Sprintf(`%s
---
title: %s
date: %s
draft: %t
img_url: %s
location: %s
tags: %q
---
%s
`,
		c.Path,
		c.Title,
		c.Date.String(),
		c.Draft,
		c.ImgURL,
		c.Location,
		c.Tags,
		string(c.Body),
	)
}

func (c *Content) Write(path string) error {
	var body []byte

	m, err := yaml.Marshal(c.Metadata)
	if err != nil {
		return err
	}
	// Frontmatter
	body = append(body, []byte("---\n")...)
	body = append(body, m...)
	body = append(body, []byte("\n---")...)

	// Body
	if len(c.Body) > 0 {
		body = append(body, []byte("\n")...)
		body = append(body, c.Body...)
		body = append(body, []byte("\n")...)
	}

	return ioutil.WriteFile(filepath.Clean(path), body, 0644)
}

func LoadContent(path string) (*Content, error) {
	b, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, err
	}

	c := Content{
		Metadata: &Metadata{},
		Path:     path,
	}

	parts := bytes.SplitN(b, []byte("---"), 3)
	c.Body = parts[2]

	if err := yaml.Unmarshal(parts[1], c.Metadata); err != nil {
		return nil, ErrYamlMalformed
	}

	return &c, nil
}
