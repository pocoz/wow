package hashcash

import (
	"crypto/sha1"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/pocoz/wow/tools"
)

const (
	defaultVer  = 1
	defaultBits = 4
)

type Service struct {
	hashMap map[string]struct{}
	mu      *sync.Mutex
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
		mu:      &sync.Mutex{},
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

func (b *Block) calculateHash() string {
	payload := fmt.Sprintf("%d:%d:%d:%s::%s:%d", b.Ver, b.Bits, b.Date, b.Resource, b.Rand, b.Counter)

	hash := sha1.New()
	hash.Write([]byte(payload))

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func (b *Block) GeneratePow() {
	for b.Counter < math.MaxInt64 {
		b.Hash = b.calculateHash()
		if b.validateHash() {
			return
		}

		b.Counter++
	}
}

func (b *Block) validateHash() bool {
	if !strings.HasPrefix(b.Hash, strings.Repeat("0", b.Bits)) {
		return false
	}

	return true
}

func (s *Service) ValidateBlock(blockForClient, blockFromClient *Block) bool {
	if blockForClient.Ver != blockFromClient.Ver ||
		blockForClient.Bits != blockFromClient.Bits ||
		blockForClient.Date != blockFromClient.Date ||
		blockForClient.Rand != blockFromClient.Rand ||
		blockForClient.Resource != blockFromClient.Resource {
		return false
	}

	if blockFromClient.calculateHash() != blockFromClient.Hash {
		return false
	}

	if !blockFromClient.validateHash() {
		return false
	}

	if time.Unix(blockFromClient.Date, 0).Before(time.Now().AddDate(0, 0, -2)) {
		return false
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.hashMap[blockFromClient.Hash]
	if ok {
		return false
	}

	s.hashMap[blockFromClient.Hash] = struct{}{}

	return true
}
