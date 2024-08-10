package evmctypes

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type syncing struct {
	StartingBlock          uint64 `json:"startingBlock" validate:"-"`
	CurrentBlock           uint64 `json:"currentBlock" validate:"-"`
	HighestBlock           uint64 `json:"highestBlock" validate:"-"`
	SyncedAccounts         uint64 `json:"syncedAccounts" validate:"-"`
	SyncedAccountBytes     uint64 `json:"syncedAccountBytes" validate:"-"`
	SyncedBytecodes        uint64 `json:"syncedBytecodes" validate:"-"`
	SyncedBytecodeBytes    uint64 `json:"syncedBytecodeBytes" validate:"-"`
	SyncedStorage          uint64 `json:"syncedStorage" validate:"-"`
	SyncedStorageBytes     uint64 `json:"syncedStorageBytes" validate:"-"`
	HealedTrienodes        uint64 `json:"healedTrienodes" validate:"-"`
	HealedTrienodeBytes    uint64 `json:"healedTrienodeBytes" validate:"-"`
	HealedBytecodes        uint64 `json:"healedBytecodes" validate:"-"`
	HealedBytecodeBytes    uint64 `json:"healedBytecodeBytes" validate:"-"`
	HealingTrienodes       uint64 `json:"healingTrienodes" validate:"-"`
	HealingBytecode        uint64 `json:"healingBytecode" validate:"-"`
	TxIndexFinishedBlocks  uint64 `json:"txIndexFinishedBlocks" validate:"-"`
	TxIndexRemainingBlocks uint64 `json:"txIndexRemainingBlocks" validate:"-"`
}

type _syncing struct {
	StartingBlock          *string `json:"startingBlock"`
	CurrentBlock           *string `json:"currentBlock"`
	HighestBlock           *string `json:"highestBlock"`
	SyncedAccounts         *string `json:"syncedAccounts"`
	SyncedAccountBytes     *string `json:"syncedAccountBytes"`
	SyncedBytecodes        *string `json:"syncedBytecodes"`
	SyncedBytecodeBytes    *string `json:"syncedBytecodeBytes"`
	SyncedStorage          *string `json:"syncedStorage"`
	SyncedStorageBytes     *string `json:"syncedStorageBytes"`
	HealedTrienodes        *string `json:"healedTrienodes"`
	HealedTrienodeBytes    *string `json:"healedTrienodeBytes"`
	HealedBytecodes        *string `json:"healedBytecodes"`
	HealedBytecodeBytes    *string `json:"healedBytecodeBytes"`
	HealingTrienodes       *string `json:"healingTrienodes"`
	HealingBytecode        *string `json:"healingBytecode"`
	TxIndexFinishedBlocks  *string `json:"txIndexFinishedBlocks"`
	TxIndexRemainingBlocks *string `json:"txIndexRemainingBlocks"`
}

func (_s *_syncing) unmarshal(s *syncing) error {
	if _s.StartingBlock != nil {
		startingBlock, err := hexutil.DecodeUint64(*_s.StartingBlock)
		if err != nil {
			return err
		}
		s.StartingBlock = startingBlock
	}
	if _s.CurrentBlock != nil {
		currentBlock, err := hexutil.DecodeUint64(*_s.CurrentBlock)
		if err != nil {
			return err
		}
		s.CurrentBlock = currentBlock
	}
	if _s.HighestBlock != nil {
		highestBlock, err := hexutil.DecodeUint64(*_s.HighestBlock)
		if err != nil {
			return err
		}
		s.HighestBlock = highestBlock
	}
	if _s.SyncedAccounts != nil {
		syncedAccounts, err := hexutil.DecodeUint64(*_s.SyncedAccounts)
		if err != nil {
			return err
		}
		s.SyncedAccounts = syncedAccounts
	}
	if _s.SyncedAccountBytes != nil {
		syncedAccountBytes, err := hexutil.DecodeUint64(*_s.SyncedAccountBytes)
		if err != nil {
			return err
		}
		s.SyncedAccountBytes = syncedAccountBytes
	}
	if _s.SyncedBytecodeBytes != nil {
		syncedBytecodeBytes, err := hexutil.DecodeUint64(*_s.SyncedBytecodeBytes)
		if err != nil {
			return err
		}
		s.SyncedBytecodeBytes = syncedBytecodeBytes
	}
	if _s.SyncedBytecodes != nil {
		syncedBytecodes, err := hexutil.DecodeUint64(*_s.SyncedBytecodes)
		if err != nil {
			return err
		}
		s.SyncedBytecodes = syncedBytecodes
	}
	if _s.SyncedStorage != nil {
		syncedStorage, err := hexutil.DecodeUint64(*_s.SyncedStorage)
		if err != nil {
			return err
		}
		s.SyncedStorage = syncedStorage
	}
	if _s.SyncedStorageBytes != nil {
		syncedStorageBytes, err := hexutil.DecodeUint64(*_s.SyncedStorageBytes)
		if err != nil {
			return err
		}
		s.SyncedStorageBytes = syncedStorageBytes
	}
	if _s.HealedTrienodes != nil {
		healedTrienodes, err := hexutil.DecodeUint64(*_s.HealedTrienodes)
		if err != nil {
			return err
		}
		s.HealedTrienodes = healedTrienodes
	}
	if _s.HealedTrienodeBytes != nil {
		healedTrienodeBytes, err := hexutil.DecodeUint64(*_s.HealedTrienodeBytes)
		if err != nil {
			return err
		}
		s.HealedTrienodeBytes = healedTrienodeBytes
	}
	if _s.HealedBytecodes != nil {
		healedBytecodes, err := hexutil.DecodeUint64(*_s.HealedBytecodes)
		if err != nil {
			return err
		}
		s.HealedBytecodes = healedBytecodes
	}
	if _s.HealedBytecodeBytes != nil {
		healedBytecodeBytes, err := hexutil.DecodeUint64(*_s.HealedBytecodeBytes)
		if err != nil {
			return err
		}
		s.HealedBytecodeBytes = healedBytecodeBytes
	}
	if _s.HealingTrienodes != nil {
		healingTrienodes, err := hexutil.DecodeUint64(*_s.HealingTrienodes)
		if err != nil {
			return err
		}
		s.HealingTrienodes = healingTrienodes
	}
	if _s.HealingBytecode != nil {
		healingBytecode, err := hexutil.DecodeUint64(*_s.HealingBytecode)
		if err != nil {
			return err
		}
		s.HealingBytecode = healingBytecode
	}
	if _s.TxIndexFinishedBlocks != nil {
		txIndexFinishedBlocks, err := hexutil.DecodeUint64(*_s.TxIndexFinishedBlocks)
		if err != nil {
			return err
		}
		s.TxIndexFinishedBlocks = txIndexFinishedBlocks
	}
	if _s.TxIndexRemainingBlocks != nil {
		txIndexRemainingBlocks, err := hexutil.DecodeUint64(*_s.TxIndexRemainingBlocks)
		if err != nil {
			return err
		}
		s.TxIndexRemainingBlocks = txIndexRemainingBlocks
	}
	return nil
}

func (s *Syncing) UnmarshalJSON(input []byte) error {
	fmt.Println(input)
	type __syncing struct {
		_syncing
	}
	var dec __syncing
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	return dec.unmarshal(&s.syncing)
}
