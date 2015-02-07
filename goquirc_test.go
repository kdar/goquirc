package goquirc

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestDecode(t *testing.T) {
	d := New()

	total := 0
	success := 0
	failed := 0

	err := filepath.Walk("testdata/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		total++

		imgdata, err := ioutil.ReadFile(path)
		if err != nil {
			t.Error(path+":", err)
			return nil
		}

		// Decode image
		m, _, err := image.Decode(bytes.NewReader(imgdata))
		if err != nil {
			t.Error(path+":", err)
			return nil
		}

		datas, err := d.Decode(m)
		if err != nil {
			failed++
			t.Error(path+":", err)
		} else if datas != nil {
			success++
			//fmt.Printf("%s: %d\n", path, len(datas))
			// for _, v := range datas {
			// 	fmt.Printf("%s[%d]: %s\n", path, v.DataType, v.Payload[:v.PayloadLen])
			// }
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	d.Destroy()

	if success != total {
		t.Errorf("out of %d images, %d failed and %d succeeded", total, failed, success)
	}
}
