package internal

import (
	"encoding/binary"
	"errors"
	"fmt"
	"image"
)

var (
	ErrImageNotSupported = errors.New("image format not supported or too small")
	ErrDataNotFound      = errors.New("no hidden data detected")
	ErrDecryptionFailed  = errors.New("decryption failed (wrong password or data corrupted)")
)

func EmbedData(src image.Image, plaintext []byte, off int, password string) (*image.NRGBA, error) {
	dst := format(src)
	op := PixOperator(dst.Pix)

	maxCapacity := op.Amount() - off
	if maxCapacity <= 0 {
		return nil, fmt.Errorf("%w: offset out of bounds", ErrImageNotSupported)
	}

	const overhead = 8 + 12 + 16 + 4
	requiredSize := len(plaintext) + overhead

	if requiredSize > maxCapacity {
		return nil, fmt.Errorf("%w: need %d bytes, have %d", ErrImageNotSupported, requiredSize, maxCapacity)
	}

	ciphertext, err := Encrypt(password, plaintext)
	if err != nil {
		return nil, fmt.Errorf("encryption internal error: %v", err)
	}

	length := uint32(len(ciphertext))
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, length)

	if err = op.Embed(header, off); err != nil {
		return nil, fmt.Errorf("%w: header embed failed: %v", ErrImageNotSupported, err)
	}

	if err = op.Embed(ciphertext, off+4); err != nil {
		return nil, fmt.Errorf("%w: body embed failed: %v", ErrImageNotSupported, err)
	}

	return dst, nil
}

func ExtractData(src image.Image, off int, password string) ([]byte, error) {
	img, ok := src.(*image.NRGBA)
	if !ok {
		return nil, fmt.Errorf("%w: input must be NRGBA", ErrImageNotSupported)
	}

	op := PixOperator(img.Pix)

	header, err := op.UnEmbed(4, off)
	if err != nil {
		return nil, fmt.Errorf("%w: cannot read header", ErrDataNotFound)
	}

	length := binary.BigEndian.Uint32(header)
	if length == 0 || int(length) > op.Amount() {
		return nil, fmt.Errorf("%w: invalid length %d", ErrDataNotFound, length)
	}

	ciphertext, err := op.UnEmbed(int(length), off+4)
	if err != nil {
		return nil, fmt.Errorf("%w: incomplete data stream", ErrDataNotFound)
	}

	plaintext, err := Decrypt(password, ciphertext)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDecryptionFailed, err)
	}

	return plaintext, nil
}
