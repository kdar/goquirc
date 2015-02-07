package goquirc

/*
#cgo CFLAGS: -I./internal/quirc/lib
#cgo windows LDFLAGS: ./libquirc.a
#cgo linux LDFLAGS: ./libquirc.a

#include "internal/quirc/lib/quirc.h"
*/
import "C"

import (
	"fmt"
	"image"
	"image/color"

	// "image/png"
	"math"
	// "os"
	"unsafe"
)

// type DecodeError C.quirc_decode_error_t
// type Code C.struct_quirc_code
// type Data C.struct_quirc_data
// type Point C.struct_quirc_point

const (
	MaxBitmap  = 3917
	MaxPayload = 8896

	/* QR-code ECC types. */
	ECCLevelM = 0
	ECCLevelL = 1
	ECCLevelH = 2
	ECCLevelQ = 3

	/* QR-code data types. */
	DataTypeNumeric = 1
	DataTypeAlpha   = 2
	DataTypeByte    = 4
	DataTypeKanji   = 8
)

// Created by cgo -godefs
// cgo -godefs goquirc.go
type Data struct {
	Version    int32
	ECCLevel   int32
	Mask       int32
	DataType   int32
	Payload    [MaxPayload]uint8
	PayloadLen int32
}

type Decoder struct {
	decoder *C.struct_quirc
}

// New creates a new decoder.
func New() *Decoder {
	return &Decoder{
		decoder: C.quirc_new(),
	}
}

// Destroy frees up the memory allocated by the decoder.
// Call this when you're done decoding images.
func (d *Decoder) Destroy() {
	C.quirc_destroy(d.decoder)
}

// var tmptest = 0

// DecodeBytes decodes the data in the passed image and returns
// any bytes encoded in it.
func (d *Decoder) DecodeBytes(img image.Image) ([]byte, error) {
	datas, err := d.Decode(img)
	if err != nil {
		return nil, err
	}

	if datas != nil && len(datas) > 0 {
		return []byte(datas[0].Payload[:datas[0].PayloadLen]), nil
	}

	return nil, nil
}

// Decode decodes the passed image and returns a slice of Data
// structures of the found QR data.
func (d *Decoder) Decode(img image.Image) ([]Data, error) {
	b := img.Bounds()
	cw := C.int(b.Max.X)
	ch := C.int(b.Max.Y)

	if C.quirc_resize(d.decoder, cw, ch) < 0 {
		return nil, fmt.Errorf("failed to resize image buffer")
	}

	mem := C.quirc_begin(d.decoder, nil, nil)
	if mem == nil {
		return nil, fmt.Errorf("could not obtain image buffer")
	}

	slice := (*[1 << 30]byte)(unsafe.Pointer(mem))[:b.Max.X*b.Max.Y : b.Max.X*b.Max.Y]

	grayimg := image.NewGray(b)

	switch m := img.(type) {
	case *image.Gray:
		off := 0
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				gray := m.GrayAt(x, y)
				slice[off] = byte(gray.Y)
				grayimg.Set(x, y, gray)
				off++
			}
		}
	case *image.RGBA:
		off := 0
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				pix := toGrayLuminance(m.At(x, y))
				slice[off] = byte(pix)

				// grayimg.Set(x, y, color.Gray{toGrayLuma(m.At(x, y))})

				off++
			}
		}
	default:
		off := 0
		for y := b.Min.Y; y < b.Max.Y; y++ {
			for x := b.Min.X; x < b.Max.X; x++ {
				rgba := color.RGBAModel.Convert(m.At(x, y)).(color.RGBA)
				pix := toGrayLuminance(rgba)
				slice[off] = byte(pix)
				// grayimg.Set(x, y, color.Gray{toGrayLuma(rgba)})
				off++
			}
		}
	}

	// fp, _ := os.Create(fmt.Sprintf("out/%d.png", tmptest))
	// defer fp.Close()
	// tmptest++
	// png.Encode(fp, grayimg)

	C.quirc_end(d.decoder)

	count := int(C.quirc_count(d.decoder))
	if count == 0 {
		return nil, fmt.Errorf("no QR code found in image")
	}

	var datas []Data
	for i := 0; i < count; i++ {
		var code C.struct_quirc_code
		var data C.struct_quirc_data

		C.quirc_extract(d.decoder, C.int(i), &code)
		res := C.quirc_decode(&code, &data)
		if res == C.QUIRC_SUCCESS {
			datas = append(datas, *(*Data)(unsafe.Pointer(&data)))
		} else {
			str := C.GoString(C.quirc_strerror(res))
			return nil, fmt.Errorf("decode failed: %s", str)
		}
	}

	return datas, nil

	// var code C.struct_quirc_code
	// C.quirc_extract(d.decoder, 0, &code)

	// //fmt.Printf("%d - %d\n", (code.size-17)%4, code.size)

	// var data C.struct_quirc_data
	// res := C.quirc_decode(&code, &data)
	// if res != C.QUIRC_SUCCESS {
	// 	str := C.GoString(C.quirc_strerror(res))
	// 	return nil, fmt.Errorf("decode failed: %s", str)
	// }

	// return (*Data)(unsafe.Pointer(&data)), nil
}

func toGrayLuminance(c color.Color) uint8 {
	rr, gg, bb, _ := c.RGBA()
	r := math.Pow(float64(rr), 2.2)
	g := math.Pow(float64(gg), 2.2)
	b := math.Pow(float64(bb), 2.2)
	y := math.Pow(0.2125*r+0.7154*g+0.0721*b, 1/2.2)
	Y := uint16(y + 0.5)
	return uint8(Y >> 8)
}
