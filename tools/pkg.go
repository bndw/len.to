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

const (
	posts  = "../content/post"
	images = "../content/img"
)

var (
	ErrFileNotExist  = errors.New("file does not exist")
	ErrYamlMalformed = errors.New("unable to parse yaml")
)

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
	Path    string
	Content []byte
}

type Metadata struct {
	Title      string    `yaml:"title"`
	Date       time.Time `yaml:"date"`
	Draft      bool      `yaml:"draft"`
	Location   string    `yaml:"location"`
	ImgURL     string    `yaml:"img_url"`
	OriginalFn string    `yaml:"original_fn"`
	Tags       []string  `yaml:"tags"`
}

func (p *Post) String() string {
	return fmt.Sprintf(`%s
---
title: %s
date: %s
draft: %t
img_url: %s
original_fn: %s
location: %s
tags: %q
---
%s
`,
		p.Path,
		p.Title,
		p.Date.String(),
		p.Draft,
		p.ImgURL,
		p.OriginalFn,
		p.Location,
		p.Tags,
		string(p.Content),
	)
}
func (p *Post) Overwrite() error {
	return p.Write(p.Path)
}

func (p *Post) Write(path string) error {
	var body []byte

	m, err := yaml.Marshal(p.Metadata)
	if err != nil {
		return err
	}
	// The yaml frontmatter is surrounded by ---
	body = append(body, []byte("---\n")...)
	body = append(body, m...)
	body = append(body, []byte("\n---")...)

	if len(p.Content) > 0 {
		body = append(body, []byte("\n")...)
		body = append(body, p.Content...)
		body = append(body, []byte("\n")...)
	}

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

	parts := bytes.SplitN(b, []byte("---"), 3)
	p.Content = parts[2]

	if err := yaml.Unmarshal(parts[1], p.Metadata); err != nil {
		return nil, ErrYamlMalformed
	}

	return &p, nil
}
