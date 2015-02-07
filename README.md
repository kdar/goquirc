goquirc
=======

This is a simple wrapper around [Quirc](https://github.com/dlbeer/quirc).

### Install

1. `go get -d github.com/kdar/goquirc`
2. `go generate` and `go install`

### Notes

This package isn't meant to be exhaustive, but I'll try to add something if requested and time permitting.

### Detection rate

The tests in this package do not pass. [Quirc](https://github.com/dlbeer/quirc) cannot detect a lot of QR codes that I tested. As far as detection goes, I rate the following from best to worse:

1. [ZBar](https://github.com/ZBar/ZBar)
2. [Quirc](https://github.com/dlbeer/quirc)
3. [libdecodeqr](https://github.com/josephholsten/libdecodeqr) used in [qrcode](https://github.com/chai2010/qrcode)

ZBar is by far the best. It could decode everything I threw at it. The problem is I couldn't get it to compile under Windows 64 MinGW in a timely manner, so I moved on.

Quirc works for the purposes I needed it for, and the API is super simple.

libdecodeqr could barely decode even the basic things. I'm not too sure why it has so much trouble, even with flat out crisp generated QR codes.

### Request

Anyone know of any good QR decoding libraries that can decode at a high success rate and are in C/Go?
