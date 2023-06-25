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

type Service struct {
	hashMap map[string]struct{}
}

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

func New() *Service {
	return &Service{
		hashMap: make(map[string]struct{}),
	}
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

func (b *Block) calculateHash() {
	header := fmt.Sprintf("%d:%d:%d:%s::%s:%d", b.Ver, b.Bits, b.Date, b.Resource, b.Rand, b.Counter)

	hash := sha1.New()
	hash.Write([]byte(header))

	b.Hash = fmt.Sprintf("%x", hash.Sum(nil))
}

func (b *Block) GeneratePow() {
	for {
		b.calculateHash()
		if b.validate() {
			return
		}

		b.Counter++
	}
}

func (b *Block) validate() bool {
	if b.Bits > len(b.Hash) {
		return false
	}

	for _, ch := range b.Hash[:b.Bits] {
		if ch != zeroByte {
			return false
		}
	}

	return true
}

func (s *Service) FullValidate(h *Block) bool {
	if !h.validate() {
		return false
	}

	if time.Unix(h.Date, 0).Before(time.Now().AddDate(0, 0, -2)) {
		return false
	}

	_, ok := s.hashMap[h.Hash]
	if ok {
		return false
	}

	s.hashMap[h.Hash] = struct{}{}

	return true
}
