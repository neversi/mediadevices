//go:build !cgo
// +build !cgo

package video

import (
	"image"
	"image/color"
	"unsafe"
)

const hasCGOConvert = false

func i444ToI420(img image.YCbCr, dst []uint8) image.YCbCr {
	h := img.Rect.Dy()
	addrSrc0 := 0
	addrSrc1 := img.CStride
	cLen := img.CStride * (h / 4)
	addrDst := 0
	// Divide preallocated memory to cbDst and crDst
	// and truncate cap and len to cLen
	cbDst, crDst := dst[:cLen:cLen], dst[cLen:]
	crDst = crDst[:cLen:cLen]

	cbDstAddr := uintptr(unsafe.Pointer(&cbDst[0]))
	crDstAddr := uintptr(unsafe.Pointer(&crDst[0]))
	cbSrcAddr := uintptr(unsafe.Pointer(&img.Cb[0]))
	crSrcAddr := uintptr(unsafe.Pointer(&img.Cr[0]))

	for i := 0; i < h/2; i++ {
		for j := 0; j < img.CStride/2; j++ {
			// cb := uint16(img.Cb[addrSrc0]) + uint16(img.Cb[addrSrc1]) +
			// 	uint16(img.Cb[addrSrc0+1]) + uint16(img.Cb[addrSrc1+1])

			srcCbPtr1 := (*uint8)((unsafe.Pointer)(cbSrcAddr + uintptr(addrSrc0)))
			srcCbPtr2 := (*uint8)((unsafe.Pointer)(cbSrcAddr + uintptr(addrSrc0+1)))
			srcCbPtr3 := (*uint8)((unsafe.Pointer)(cbSrcAddr + uintptr(addrSrc1)))
			srcCbPtr4 := (*uint8)((unsafe.Pointer)(cbSrcAddr + uintptr(addrSrc1+1)))
			dstCbPtr := (*uint8)((unsafe.Pointer)(cbDstAddr + uintptr(addrDst)))

			*dstCbPtr = uint8((uint16(*srcCbPtr1) + uint16(*srcCbPtr2) + uint16(*(srcCbPtr3)) + uint16(*(srcCbPtr4))) / 4)

			// cr := uint16(img.Cr[addrSrc0]) + uint16(img.Cr[addrSrc1]) +
			// 	uint16(img.Cr[addrSrc0+1]) + uint16(img.Cr[addrSrc1+1])
			srcCrPtr1 := (*uint8)((unsafe.Pointer)(crSrcAddr + uintptr(addrSrc0)))
			srcCrPtr2 := (*uint8)((unsafe.Pointer)(crSrcAddr + uintptr(addrSrc0+1)))
			srcCrPtr3 := (*uint8)((unsafe.Pointer)(crSrcAddr + uintptr(addrSrc1)))
			srcCrPtr4 := (*uint8)((unsafe.Pointer)(crSrcAddr + uintptr(addrSrc1+1)))
			dstCrPtr := (*uint8)((unsafe.Pointer)(crDstAddr + uintptr(addrDst)))

			*dstCrPtr = uint8((uint16(*srcCrPtr1) + uint16(*srcCrPtr2) + uint16(*(srcCrPtr3)) + uint16(*(srcCrPtr4))) / 4)
			// cbDst[addrDst] = uint8(cb / 4)
			// crDst[addrDst] = uint8(cr / 4)
			addrSrc0 += 2
			addrSrc1 += 2
			addrDst++
		}
		addrSrc0 += img.CStride
		addrSrc1 += img.CStride
	}
	img.CStride = img.CStride / 2
	img.Cb = cbDst
	img.Cr = crDst
	img.SubsampleRatio = image.YCbCrSubsampleRatio420
	return img
}

func i422ToI420(img image.YCbCr, dst []uint8) image.YCbCr {
	h := img.Rect.Dy()
	addrSrc := 0
	cLen := img.CStride * (h / 2)
	// Divide preallocated memory to cbDst and crDst
	// and truncate cap and len to cLen
	cbDst, crDst := dst[:cLen:cLen], dst[cLen:]
	crDst = crDst[:cLen:cLen]
	addrDst := 0

	cbDstAddr := uintptr(unsafe.Pointer(&cbDst[0]))
	crDstAddr := uintptr(unsafe.Pointer(&crDst[0]))
	cbSrcAddr := uintptr(unsafe.Pointer(&img.Cb[0]))
	crSrcAddr := uintptr(unsafe.Pointer(&img.Cr[0]))

	for i := 0; i < h/2; i++ {
		for j := 0; j < img.CStride; j++ {
			// cb := uint16(img.Cb[addrSrc]) + uint16(img.Cb[addrSrc+img.CStride])
			srcCbPtr1 := (*uint8)((unsafe.Pointer)(cbSrcAddr + uintptr(addrSrc)))
			srcCbPtr2 := (*uint8)((unsafe.Pointer)(cbSrcAddr + uintptr(addrSrc+img.CStride)))
			dstCbPtr := (*uint8)((unsafe.Pointer)(cbDstAddr + uintptr(addrDst)))

			*dstCbPtr = uint8((uint16(*srcCbPtr1) + uint16(*srcCbPtr2)) / 4)
			// cr := uint16(img.Cr[addrSrc]) + uint16(img.Cr[addrSrc+img.CStride])
			srcCrPtr1 := (*uint8)((unsafe.Pointer)(crSrcAddr + uintptr(addrSrc)))
			srcCrPtr2 := (*uint8)((unsafe.Pointer)(crSrcAddr + uintptr(addrSrc+img.CStride)))
			dstCrPtr := (*uint8)((unsafe.Pointer)(crDstAddr + uintptr(addrDst)))

			*dstCrPtr = uint8((uint16(*srcCrPtr1) + uint16(*srcCrPtr2)) / 4)
			// cbDst[addrDst] = uint8(cb / 4)
			// crDst[addrDst] = uint8(cr / 4)
			addrSrc++
			addrDst++
		}
		addrSrc += img.CStride
	}
	img.Cb = cbDst
	img.Cr = crDst
	img.SubsampleRatio = image.YCbCrSubsampleRatio420
	return img
}

func i444ToRGBA(dst *image.RGBA, src *image.YCbCr) {
	dx := src.Rect.Dx()
	dy := src.Rect.Dy()
	i := 0
	j := 0
	for yi := 0; yi < dy; yi++ {
		for xi := 0; xi < dx; xi++ {
			r, g, b := color.YCbCrToRGB(src.Y[j], src.Cb[j], src.Cr[j])
			dst.Pix[i+0] = uint8(r)
			dst.Pix[i+1] = uint8(g)
			dst.Pix[i+2] = uint8(b)
			dst.Pix[i+3] = 0xff
			i += 4
		}
	}
}

func rgbaToI444(dst *image.YCbCr, src *image.RGBA) {
	i := 0
	addr := 0
	dx := src.Rect.Dx()
	dy := src.Rect.Dy()
	// TODO implement
	for yi := 0; yi < dy; yi++ {
		for xi := 0; xi < dx; xi++ {
			dst.Y[i], dst.Cb[i], dst.Cr[i] = color.RGBToYCbCr(
				src.Pix[addr+0], src.Pix[addr+1], src.Pix[addr+2],
			)
			addr += 4
			i++
		}
	}
}
