// Package hasher hashes a given file using
// specified hash function
package hasher

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"
)

func Hash(filePath, hashMethod string) (string, error) {
	hashMethod = strings.ToUpper(hashMethod)

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return "", err
	}

	if hashMethod == "SIZE" {
		// Convert int64 to string
		return strconv.FormatInt(info.Size(), 10), nil
	}

	if hashMethod == "MD5" {
		h := md5.New()
		if _, err := io.Copy(h, file); err != nil {
			return "", err
		}

		md5sum := h.Sum(nil)
		return hex.EncodeToString(md5sum), nil
	}

	if hashMethod == "SHA256" {
		h := sha256.New()
		if _, err := io.Copy(h, file); err != nil {
			return "", err
		}

		sha256sum := h.Sum(nil)
		return hex.EncodeToString(sha256sum), nil
	}

	return "", errors.New("invalid hash method")
}

func CompareHash(filepath, expected, hashMethod string) (bool, error) {
	hashStr, err := Hash(filepath, hashMethod)
	if err != nil {
		return false, err
	}

	return expected == hashStr, nil
}
