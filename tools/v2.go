package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
	xurls "mvdan.cc/xurls/v2"
)

const (
	posts  = "../content/post"
	photos = "../content/img"
)

var (
	ErrFileNotExist  = errors.New("file does not exist")
	ErrYamlMalformed = errors.New("unable to parse yaml")
)

func main() {
	err := filepath.Walk(posts, func(path string, i os.FileInfo, err error) error {
		filename := i.Name()
		if i.IsDir() {
			return nil
		}

		p, err := LoadPost(filepath.Join(posts, filename))
		if err != nil {
			return err
		}

		//createContentImgForEachURLInPost(p)
		removeTagsFromEachPost(p)
		return nil
	})
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}
}

// removeTagsFromEachPost deletes the tags taxonomy from each post
func removeTagsFromEachPost(p *Post) {
	p.Tags = []string{}
	p.Write(p.Path)
}

// createContentImgForEachURLInPost creates a new file in ./content/img/ for
// each url (photo) in given post.
func createContentImgForEachURLInPost(p *Post) {
	for _, url := range p.URLs() {
		var (
			fn = parseFilenameWithNoExt(url)
			fp = filepath.Join(photos, fn+".md")
		)

		// Prompt before overwriting
		img, err := LoadPost(fp)
		if err == nil && img != nil &&
			!yesno(fmt.Sprintf("%s already exists. do you want to replace it? (y/N)", fp)) {
			fmt.Println("skipping")
			continue
		}

		fmt.Printf(`Image: %s
Tags: %q
Enter new tags ([enter] to proceed, ! to use existing tags, x to skip url)
`, url, p.Metadata.Tags)

		skip := false

		// Tag the photo
		var tags []string
		for {
			fmt.Printf("> ")
			scanner := bufio.NewScanner(os.Stdin)
			if scanner.Scan() {
				in := scanner.Text()
				if in == "" {
					// Done entering tags
					break
				}
				if in == "x" {
					// Skip this URL, it's probably not a photo
					skip = true
					continue
				}
				if in == "!" {
					// Copy existing tags
					tags = append(tags, p.Metadata.Tags...)
				}

				tags = append(tags, in)
			}
		}

		if skip {
			fmt.Println("skipping")
			continue
		}

		photo := &Post{
			Metadata: &Metadata{
				Title:    fn,
				Date:     p.Metadata.Date,
				Draft:    p.Metadata.Draft,
				Location: p.Metadata.Location,
				Tags:     tags,
			},
			Path:             fp,
			MetadataEncoding: "yaml",
			Content:          []byte(fmt.Sprintf("![](%s)", url)),
		}

		if err := photo.Write(fp); err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Printf("created %s\n", fp)
		fmt.Println()
	}
}

// parseFilenameWithNoExt takes a url for an image like https://s3...LE2329.jpg, and returns 'LE2329'
func parseFilenameWithNoExt(u string) string {
	_, urlFn := path.Split(u)
	return urlFn[:len(urlFn)-len(filepath.Ext(urlFn))]
}

func yesno(question string) bool {
	fmt.Printf("%s\n> ", question)
	var y string
	fmt.Scanln(&y)
	if strings.ToLower(y) == "y" {
		return true
	}
	return false
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
