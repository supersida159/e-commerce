package hasher

import (
	"crypto/md5"
	"encoding/hex"
)

type md5hash struct {
}

func NewMd5() *md5hash {
	return &md5hash{}
}

func (m *md5hash) Hash(data string) string {
	haser := md5.New()
	haser.Write([]byte(data))
	return hex.EncodeToString(haser.Sum(nil))
}
