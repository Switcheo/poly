/*
 * Copyright (C) 2021 The poly network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */

package cosmos

import (
	"bytes"
	"fmt"

	"github.com/polynetwork/poly/common"
	"github.com/polynetwork/poly/common/log"
	"github.com/polynetwork/poly/native"
	"github.com/polynetwork/poly/native/service/governance/node_manager"
	hscommon "github.com/polynetwork/poly/native/service/header_sync/common"
	"github.com/polynetwork/poly/native/service/utils"
)

type CosmosHandler struct{}

func NewCosmosHandler() *CosmosHandler {
	return &CosmosHandler{}
}

func (this *CosmosHandler) SyncGenesisHeader(native *native.NativeService) error {
	param := new(hscommon.SyncGenesisHeaderParam)
	if err := param.Deserialization(common.NewZeroCopySource(native.GetInput())); err != nil {
		return fmt.Errorf("CosmosHandler SyncGenesisHeader, contract params deserialize error: %v", err)
	}
	// Get current epoch operator
	operatorAddress, err := node_manager.GetCurConOperator(native)
	if err != nil {
		return fmt.Errorf("CosmosHandler SyncGenesisHeader, get current consensus operator address error: %v", err)
	}
	//check witness
	err = utils.ValidateOwner(native, operatorAddress)
	if err != nil {
		return fmt.Errorf("CosmosHandler SyncGenesisHeader, checkWitness error: %v", err)
	}
	// get genesis header from input parameters
	var header CosmosHeader
	err = Cdc.UnmarshalBinaryBare(param.GenesisHeader, &header)
	if err != nil {
		return fmt.Errorf("CosmosHandler SyncGenesisHeader: %s", err)
	}
	// check if has genesis header
	info, err := GetEpochSwitchInfo(native, param.ChainID)
	if err == nil && info != nil {
		return fmt.Errorf("CosmosHandler SyncGenesisHeader, genesis header had been initialized")
	}
	PutEpochSwitchInfo(native, param.ChainID, &CosmosEpochSwitchInfo{
		Height:             header.Header.Height,
		NextValidatorsHash: header.Header.NextValidatorsHash,
		ChainID:            header.Header.ChainID,
		BlockHash:          HashCosmosHeader(header.Header),
	})
	return nil
}

func (this *CosmosHandler) SyncBlockHeader(native *native.NativeService) error {
	params := new(hscommon.SyncBlockHeaderParam)
	if err := params.Deserialization(common.NewZeroCopySource(native.GetInput())); err != nil {
		return fmt.Errorf("SyncBlockHeader, contract params deserialize error: %v", err)
	}
	cnt := 0
	info, err := GetEpochSwitchInfo(native, params.ChainID)
	if err != nil {
		return fmt.Errorf("SyncBlockHeader, get epoch switching height failed: %v", err)
	}
	for _, v := range params.Headers {
		var myHeader CosmosHeader
		err := Cdc.UnmarshalBinaryBare(v, &myHeader)
		if err != nil {
			return fmt.Errorf("SyncBlockHeader failed to unmarshal header: %v", err)
		}
		if bytes.Equal(myHeader.Header.NextValidatorsHash, myHeader.Header.ValidatorsHash) {
			continue
		}
		if info.Height >= myHeader.Header.Height {
			log.Debugf("SyncBlockHeader, height %d is lower or equal than epoch switching height %d",
				myHeader.Header.Height, info.Height)
			continue
		}
		if err = VerifyCosmosHeader(&myHeader, info); err != nil {
			return fmt.Errorf("SyncBlockHeader, failed to verify header: %v", err)
		}
		info.NextValidatorsHash = myHeader.Header.NextValidatorsHash
		info.Height = myHeader.Header.Height
		info.BlockHash = HashCosmosHeader(myHeader.Header)
		cnt++
	}
	if cnt == 0 {
		return fmt.Errorf("no header you commited is useful")
	}
	PutEpochSwitchInfo(native, params.ChainID, info)
	return nil
}

func (this *CosmosHandler) SyncCrossChainMsg(native *native.NativeService) error {
	return nil
}
