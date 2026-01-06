package internal

import (
	"bytes"
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

// Overhead includes: Salt(8) + Nonce(12) + Tag(16) + LengthHeader(4)
const Overhead = 8 + 12 + 16 + 4

func MaxCapacity(src image.Image) int {
	bounds := src.Bounds()
	totalPixels := bounds.Dx() * bounds.Dy()
	
	capacity := totalPixels - Overhead
	if capacity < 0 {
		return 0
	}
	return capacity
}

// EmbedData embeds data with an optional extension.
// If extension is empty, it's treated as raw data/text.
// Format inside encryption: [ExtLen(1 byte)][ExtBytes(...)][Data...]
func EmbedData(src image.Image, data []byte, extension string, off int, password string) (*image.NRGBA, error) {
	dst := format(src)
	op := PixOperator(dst.Pix)

	maxCapacity := op.Amount() - off
	if maxCapacity <= 0 {
		return nil, fmt.Errorf("%w: offset out of bounds", ErrImageNotSupported)
	}

	// Prepare payload: [ExtLen(1)][ExtBytes][Data]
	extBytes := []byte(extension)
	if len(extBytes) > 255 {
		return nil, errors.New("extension too long (max 255 bytes)")
	}
	
	payload := new(bytes.Buffer)
	payload.WriteByte(uint8(len(extBytes)))
	payload.Write(extBytes)
	payload.Write(data)
	
	plaintext := payload.Bytes()
	requiredSize := len(plaintext) + Overhead

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

// ExtractData extracts data and returns the content and its extension (if any).
func ExtractData(src image.Image, off int, password string) ([]byte, string, error) {
	dst := format(src)
	op := PixOperator(dst.Pix)

	header, err := op.UnEmbed(4, off)
	if err != nil {
		return nil, "", fmt.Errorf("%w: cannot read header", ErrDataNotFound)
	}

	length := binary.BigEndian.Uint32(header)
	if length == 0 || int(length) > op.Amount() {
		return nil, "", fmt.Errorf("%w: invalid length %d", ErrDataNotFound, length)
	}

	ciphertext, err := op.UnEmbed(int(length), off+4)
	if err != nil {
		return nil, "", fmt.Errorf("%w: incomplete data stream", ErrDataNotFound)
	}

	plaintext, err := Decrypt(password, ciphertext)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %v", ErrDecryptionFailed, err)
	}
	
	// Parse payload: [ExtLen(1)][ExtBytes][Data]
	if len(plaintext) < 1 {
		return nil, "", errors.New("data corrupted: missing extension length")
	}
	
	extLen := int(plaintext[0])
	if len(plaintext) < 1+extLen {
		return nil, "", errors.New("data corrupted: extension length mismatch")
	}
	
	extension := string(plaintext[1 : 1+extLen])
	data := plaintext[1+extLen:]

	return data, extension, nil
}
