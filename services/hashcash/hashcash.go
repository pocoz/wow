package hashcash

import (
	"crypto/sha1"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/pocoz/wow/tools"
)

const (
	defaultVer  = 1
	defaultBits = 4
)

type Block struct {
	Ver      int
	Bits     int
	Date     int64
	Resource string
	Ext      string
	Rand     string
	Counter  int
	Hash     string
}

func NewBlock(resource string) *Block {
	return &Block{
		Ver:      defaultVer,
		Bits:     defaultBits,
		Date:     time.Now().Unix(),
		Resource: resource,
		Rand:     tools.GetRandomString(),
		Counter:  0,
	}
}

func (b *Block) CalculateHash() string {
	payload := fmt.Sprintf("%d:%d:%d:%s::%s:%d", b.Ver, b.Bits, b.Date, b.Resource, b.Rand, b.Counter)

	hash := sha1.New()
	hash.Write([]byte(payload))

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (b *Block) GeneratePow() {
	for b.Counter < math.MaxInt64 {
		b.Hash = b.CalculateHash()
		if b.ValidateHash() {
			return
		}

		b.Counter++
	}
}

func (b *Block) ValidateHash() bool {
	if !strings.HasPrefix(b.Hash, strings.Repeat("0", b.Bits)) {
		return false
	}

	return true
}
