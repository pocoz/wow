package tcpserver

import (
	"time"

	"github.com/pocoz/wow/services/hashcash"
)

func (s *Service) validateBlock(blockForClient, blockFromClient *hashcash.Block) bool {
	if blockForClient.Ver != blockFromClient.Ver ||
		blockForClient.Bits != blockFromClient.Bits ||
		blockForClient.Date != blockFromClient.Date ||
		blockForClient.Rand != blockFromClient.Rand ||
		blockForClient.Resource != blockFromClient.Resource {
		return false
	}

	if blockFromClient.CalculateHash() != blockFromClient.Hash {
		return false
	}

	if !blockFromClient.ValidateHash() {
		return false
	}

	if time.Unix(blockFromClient.Date, 0).Before(time.Now().AddDate(0, 0, -2)) {
		return false
	}

	if s.storage.IsHashExist(blockFromClient.Hash) {
		return false
	}

	s.storage.AddHash(blockFromClient.Hash)

	return true
}
