// Copyright(c) 2018 DSiSc Group. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"github.com/DSiSc/craft/types"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

var MockHash = types.Hash{
	0x1d, 0xcf, 0x7, 0xba, 0xfc, 0x42, 0xb0, 0x8d, 0xfd, 0x23, 0x9c, 0x45, 0xa4, 0xb9, 0x38, 0xd,
	0x8d, 0xfe, 0x5d, 0x6f, 0xa7, 0xdb, 0xd5, 0x50, 0xc9, 0x25, 0xb1, 0xb3, 0x4, 0xdc, 0xc5, 0x1c,
}

func MockBlock() *types.Block {
	return &types.Block{
		Header: &types.Header{
			ChainID:       1,
			PrevBlockHash: MockHash,
			StateRoot:     MockHash,
			TxRoot:        MockHash,
			ReceiptsRoot:  MockHash,
			Height:        1,
			Timestamp:     uint64(time.Date(2018, time.August, 28, 0, 0, 0, 0, time.UTC).Unix()),
		},
		Transactions: make([]*types.Transaction, 0),
	}
}

var MockBlockHash = types.Hash{
	0x9, 0x99, 0xfd, 0xff, 0x97, 0x34, 0xff, 0xa9, 0xda, 0x64, 0x69, 0xcb, 0x62, 0x6d, 0x7a, 0xec, 0x1c, 0xa1, 0xb2, 0xbf, 0x50, 0x5b, 0x71, 0x6, 0x3e, 0x20, 0x5b, 0x66, 0xb2, 0xd4, 0xbf, 0xb1,
}

func TestHeaderHash(t *testing.T) {
	block := MockBlock()
	hash := HeaderHash(block)
	assert.Equal(t, MockBlockHash, hash)
}

var emptyTx *types.Transaction

func TestNewTransaction(t *testing.T) {
	assert := assert.New(t)
	b := Address{
		0xb2, 0x6f, 0x2b, 0x34, 0x2a, 0xab, 0x24, 0xbc, 0xf6, 0x3e,
		0xa2, 0x18, 0xc6, 0xa9, 0x27, 0x4d, 0x30, 0xab, 0x9a, 0x15,
	}
	emptyTx = NewTransaction(
		0,
		&b,
		big.NewInt(0), 0, big.NewInt(0),
		nil,
		b,
	)
	notoTx := NewTransaction(
		0,
		nil,
		nil, 0, big.NewInt(0),
		nil,
		b,
	)
	assert.NotNil(emptyTx)
	assert.NotNil(notoTx)
}

func TestTxHash(t *testing.T) {
	assert := assert.New(t)
	b := Address{
		0xb2, 0x6f, 0x2b, 0x34, 0x2a, 0xab, 0x24, 0xbc, 0xf6, 0x3e,
		0xa2, 0x18, 0xc6, 0xa9, 0x27, 0x4d, 0x30, 0xab, 0x9a, 0x15,
	}
	emptyTx = NewTransaction(
		0,
		&b,
		big.NewInt(0),
		0,
		big.NewInt(0),
		b[:10],
		b,
	)
	exceptHash := TxHash(emptyTx)
	assert.Equal(exceptHash, emptyTx.Hash.Load().(types.Hash))

	exceptHash1 := TxHash(emptyTx)
	assert.Equal(exceptHash, exceptHash1)
}
