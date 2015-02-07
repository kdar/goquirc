package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/kdar/goquirc"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	flag.Parse()

	if len(flag.Args()) == 0 {
		log.Fatalf("usage: %s <image file>", os.Args[0])
	}

	path := flag.Arg(0)

	imgdata, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(path+":", err)
	}

	// Decode image
	m, _, err := image.Decode(bytes.NewReader(imgdata))
	if err != nil {
		log.Fatal(path+":", err)
	}

	d := goquirc.New()
	defer d.Destroy()
	datas, err := d.Decode(m)
	if err != nil {
		log.Fatal(path+":", err)
	}

	for _, data := range datas {
		fmt.Printf("%s\n", data.Payload[:data.PayloadLen])
	}
}
