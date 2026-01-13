package internal

import (
	"bytes"
	"encoding/binary"
	"image"
)

const Overhead = 8 + 12 + 16 + 4

func Capacity(src image.Image) int {
	bounds := src.Bounds()
	totalPixels := bounds.Dx() * bounds.Dy()
	
	capacity := totalPixels - Overhead
	if capacity < 0 {
		return 0
	}
	return capacity
}

func EmbedData(src image.Image, data []byte, extension string, off int, password string) (*image.NRGBA, error) {
	dst := format(src)
	op := PixOperator(dst.Pix)
	
	maxCapacity := op.Capacity() - off
	if maxCapacity <= 0 {
		return nil, ErrImageNotSupported
	}
	
	extBytes := []byte(extension)
	if len(extBytes) > 255 {
		return nil, ErrExtensionTooLong
	}
	
	payload := new(bytes.Buffer)
	payload.WriteByte(uint8(len(extBytes)))
	payload.Write(extBytes)
	payload.Write(data)
	
	plaintext := payload.Bytes()
	requiredSize := len(plaintext) + Overhead
	
	if requiredSize > maxCapacity {
		return nil, ErrImageTooSmall
	}
	
	ciphertext, err := Encrypt(password, plaintext)
	if err != nil {
		return nil, ErrInternal
	}
	
	length := uint32(len(ciphertext))
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, length)
	
	if err = op.Embed(header, off); err != nil {
		return nil, ErrInternal
	}
	
	if err = op.Embed(ciphertext, off+4); err != nil {
		return nil, ErrInternal
	}
	
	return dst, nil
}

func ExtractData(src image.Image, off int, password string) ([]byte, string, error) {
	dst := format(src)
	op := PixOperator(dst.Pix)
	
	header, err := op.UnEmbed(4, off)
	if err != nil {
		return nil, "", ErrDataNotFound
	}
	
	length := binary.BigEndian.Uint32(header)
	if length == 0 || int(length) > op.Capacity() {
		return nil, "", ErrDataNotFound
	}
	
	ciphertext, err := op.UnEmbed(int(length), off+4)
	if err != nil {
		return nil, "", ErrDataNotFound
	}
	
	plaintext, err := Decrypt(password, ciphertext)
	if err != nil {
		return nil, "", ErrDecryptionFailed
	}
	
	if len(plaintext) < 1 {
		return nil, "", ErrInternal
	}
	
	extLen := int(plaintext[0])
	if len(plaintext) < 1+extLen {
		return nil, "", ErrInternal
	}
	
	extension := string(plaintext[1 : 1+extLen])
	data := plaintext[1+extLen:]
	
	return data, extension, nil
}
