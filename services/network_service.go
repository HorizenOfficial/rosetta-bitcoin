// Copyright 2020 Coinbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package services

import (
	"context"
	"github.com/HorizenOfficial/rosetta-zen/configuration"
	"github.com/HorizenOfficial/rosetta-zen/utils"
	"github.com/HorizenOfficial/rosetta-zen/zen"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

// NetworkAPIService implements the server.NetworkAPIServicer interface.
type NetworkAPIService struct {
	config *configuration.Configuration
	client Client
	i      Indexer
}

// NewNetworkAPIService creates a new instance of a NetworkAPIService.
func NewNetworkAPIService(
	config *configuration.Configuration,
	client Client,
	i Indexer,
) server.NetworkAPIServicer {
	return &NetworkAPIService{
		config: config,
		client: client,
		i:      i,
	}
}

// NetworkList implements the /network/list endpoint
func (s *NetworkAPIService) NetworkList(
	ctx context.Context,
	request *types.MetadataRequest,
) (*types.NetworkListResponse, *types.Error) {
	return &types.NetworkListResponse{
		NetworkIdentifiers: []*types.NetworkIdentifier{
			s.config.Network,
		},
	}, nil
}

// NetworkStatus implements the /network/status endpoint.
func (s *NetworkAPIService) NetworkStatus(
	ctx context.Context,
	request *types.NetworkRequest,
) (*types.NetworkStatusResponse, *types.Error) {
	if s.config.Mode != configuration.Online {
		return nil, wrapErr(ErrUnavailableOffline, nil)
	}

	peers, err := s.client.GetPeers(ctx)
	if err != nil {
		return nil, wrapErr(ErrBitcoind, err)
	}

	cachedBlockResponse, err := s.i.GetBlockLazy(ctx, nil)
	if err != nil {
		return nil, wrapErr(ErrNotReady, nil)
	}

	return &types.NetworkStatusResponse{
		CurrentBlockIdentifier: cachedBlockResponse.Block.BlockIdentifier,
		CurrentBlockTimestamp:  cachedBlockResponse.Block.Timestamp,
		GenesisBlockIdentifier: s.config.GenesisBlockIdentifier,
		Peers:                  peers,
	}, nil
}

// NetworkOptions implements the /network/options endpoint.
func (s *NetworkAPIService) NetworkOptions(
	ctx context.Context,
	request *types.NetworkRequest,
) (*types.NetworkOptionsResponse, *types.Error) {
	logger := utils.ExtractLogger(ctx, "Network Service")
	if s.config.ZendVersion == "" {
		logger.Info("No Zend version provided")
		version, err := s.client.SetZendNodeVersion(ctx)
		if err != nil {
			logger.Error("unable to retrieve network info", "error", err)
		}
		s.config.ZendVersion = version
	}
	return &types.NetworkOptionsResponse{
		Version: &types.Version{
			RosettaVersion:    types.RosettaAPIVersion,
			NodeVersion:       s.config.ZendVersion,
			MiddlewareVersion: types.String(MiddlewareVersion),
		},
		Allow: &types.Allow{
			OperationStatuses:       zen.OperationStatuses,
			OperationTypes:          zen.OperationTypes,
			Errors:                  Errors,
			HistoricalBalanceLookup: HistoricalBalanceLookup,
			MempoolCoins:            MempoolCoins,
		},
	}, nil
}
