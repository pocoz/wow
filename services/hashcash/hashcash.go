package hashcash

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/pocoz/wow/tools"
)

const (
	defaultVer  = 1
	defaultBits = 4
	zeroByte    = 48
)

type Hashcash struct {
	Ver      int
	Bits     int
	Date     int64
	Resource string
	Ext      string
	Rand     string
	Counter  int
	Hash     string
}

func NewHashcash(resource string) *Hashcash {
	return &Hashcash{
		Ver:      defaultVer,
		Bits:     defaultBits,
		Date:     time.Now().Unix(),
		Resource: resource,
		Rand:     tools.GetRandomString(),
		Counter:  0,
	}
}

func (h *Hashcash) calculateHash() {
	header := fmt.Sprintf("%d:%d:%d:%s::%s:%d", h.Ver, h.Bits, h.Date, h.Resource, h.Rand, h.Counter)

	hash := sha1.New()
	hash.Write([]byte(header))

	h.Hash = fmt.Sprintf("%x", hash.Sum(nil))
}

func (h *Hashcash) GeneratePow() {
	for {
		h.calculateHash()
		if h.Validate() {
			return
		}

		h.Counter++
	}
}

func (h *Hashcash) Validate() bool {
	if h.Bits > len(h.Hash) {
		return false
	}

	for _, ch := range h.Hash[:h.Bits] {
		if ch != zeroByte {
			return false
		}
	}

	return true
}
