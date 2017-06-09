package uuid

import (
	"crypto/sha1"
	"encoding/hex"

	"github.com/davecgh/go-spew/spew"

	gouuid "github.com/satori/go.uuid"
)

func NewV4() string {
	u4 := gouuid.NewV4()
	return u4.String()
}

func NewV5(u, name string) string {
	u5 := gouuid.NewV5(gouuid.FromStringOrNil(u), name)
	return u5.String()
}

func New(obj interface{}) string {
	u4 := gouuid.NewV4()
	h := sha1.New()
	spew.Fdump(h, obj)
	u5 := gouuid.NewV5(u4, hex.EncodeToString(h.Sum(nil)))
	return u5.String()
}
