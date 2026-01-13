package internal

import (
	"errors"
	"image"
	"image/draw"
)

func format(src image.Image) *image.NRGBA {
	if img, ok := src.(*image.NRGBA); ok {
		clone := *img
		clone.Pix = make([]byte, len(img.Pix))
		copy(clone.Pix, img.Pix)
		return &clone
	}
	
	bounds := src.Bounds()
	dst := image.NewNRGBA(bounds)
	draw.Draw(dst, bounds, src, bounds.Min, draw.Src)
	return dst
}

type PixOperator []uint8

func (p *PixOperator) Capacity() int {
	return len(*p) / 4
}

func (p *PixOperator) Embed(data []byte, off int) error {
	if off < 0 || off+len(data) > p.Capacity() {
		return errors.New("out of bounds")
	}
	
	for i, v := range data {
		base := (off + i) * 4
		(*p)[base+0] = ((*p)[base+0] & 0xFC) | ((v >> 6) & 0x03)
		(*p)[base+1] = ((*p)[base+1] & 0xFC) | ((v >> 4) & 0x03)
		(*p)[base+2] = ((*p)[base+2] & 0xFC) | ((v >> 2) & 0x03)
		(*p)[base+3] = ((*p)[base+3] & 0xFC) | (v & 0x03)
	}
	return nil
}

func (p *PixOperator) UnEmbed(n int, off int) ([]byte, error) {
	if off < 0 || n < 0 || off+n > p.Capacity() {
		return nil, errors.New("out of bounds")
	}
	
	out := make([]byte, n)
	for i := 0; i < n; i++ {
		base := (off + i) * 4
		var v byte
		v |= ((*p)[base+0] & 0x03) << 6
		v |= ((*p)[base+1] & 0x03) << 4
		v |= ((*p)[base+2] & 0x03) << 2
		v |= (*p)[base+3] & 0x03
		out[i] = v
	}
	return out, nil
}
