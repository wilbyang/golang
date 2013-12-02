package main

import (
	"bytes"
	"fmt"
	"github.com/golang/groupcache"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func renderImg(resp http.ResponseWriter, req *http.Request) {

	req.ParseForm()

	imgName := strings.Join(req.Form["img"], "")
	imgName = fmt.Sprintf("imgs/%s", imgName)

	if _, err := os.Stat(imgName); os.IsNotExist(err) {
		http.NotFound(resp, req)
	} else {
		var buf []byte
		err := cacher.Get(nil, imgName, groupcache.AllocatingByteSliceSink(&buf))
		if err != nil {
			fmt.Println("Get error when cacher.Get:", err)
		}

		img, format, err := image.Decode(bytes.NewBuffer(buf))
		if err != nil {
			fmt.Println("Get error when image.Decode:", err)
		}
		fmt.Println("format:", format)
		png.Encode(resp, img)
	}
}

func getImage(ctx groupcache.Context, key string, dest groupcache.Sink) error {

	fmt.Println("Getting image from slow backend ... ")
	time.Sleep(time.Duration(3000) * time.Millisecond)

	imgData, err := ioutil.ReadFile(key)
	if err != nil {
		fmt.Println("Get error when ioutil.ReadFile:", err)
	}

	dest.SetBytes(imgData)
	return nil
}

var cacher *groupcache.Group

func main() {

	cacher = groupcache.NewGroup("cacher", 64<<20, groupcache.GetterFunc(getImage))

	http.HandleFunc("/", renderImg)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println("Get error when http.ListenAndServe:", err)
	}
}