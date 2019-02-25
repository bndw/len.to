package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	images, err := Images()
	if err != nil {
		fmt.Println(err)
		return
	}

	http.HandleFunc("/", handler(images))
	log.Fatal(http.ListenAndServe(":12000", nil))
}

func handler(images []string) func(w http.ResponseWriter, r *http.Request) {
	shuffleImages := func() []string {
		r := rand.New(rand.NewSource(time.Now().Unix()))
		var _images []string
		for _, i := range r.Perm(len(images)) {
			_images = append(_images, images[i])
		}
		return _images
	}

	return func(w http.ResponseWriter, r *http.Request) {
		newTemplate().Execute(w, struct {
			Style  string
			Images []string
		}{
			Style:  randomStyle(),
			Images: shuffleImages(),
		})
	}
}
