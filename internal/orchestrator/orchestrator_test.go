// Copyright © 2022 Kaleido, Inc.
//
// SPDX-License-Identifier: Apache-2.0
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

package orchestrator

import (
	"context"
	"fmt"
	"testing"

	"github.com/hyperledger/firefly/internal/coreconfig"
	"github.com/hyperledger/firefly/internal/identity"
	"github.com/hyperledger/firefly/mocks/admineventsmocks"
	"github.com/hyperledger/firefly/mocks/assetmocks"
	"github.com/hyperledger/firefly/mocks/batchmocks"
	"github.com/hyperledger/firefly/mocks/batchpinmocks"
	"github.com/hyperledger/firefly/mocks/blockchainmocks"
	"github.com/hyperledger/firefly/mocks/broadcastmocks"
	"github.com/hyperledger/firefly/mocks/contractmocks"
	"github.com/hyperledger/firefly/mocks/databasemocks"
	"github.com/hyperledger/firefly/mocks/dataexchangemocks"
	"github.com/hyperledger/firefly/mocks/datamocks"
	"github.com/hyperledger/firefly/mocks/definitionsmocks"
	"github.com/hyperledger/firefly/mocks/eventmocks"
	"github.com/hyperledger/firefly/mocks/identitymanagermocks"
	"github.com/hyperledger/firefly/mocks/identitymocks"
	"github.com/hyperledger/firefly/mocks/metricsmocks"
	"github.com/hyperledger/firefly/mocks/networkmapmocks"
	"github.com/hyperledger/firefly/mocks/operationmocks"
	"github.com/hyperledger/firefly/mocks/privatemessagingmocks"
	"github.com/hyperledger/firefly/mocks/shareddownloadmocks"
	"github.com/hyperledger/firefly/mocks/sharedstoragemocks"
	"github.com/hyperledger/firefly/mocks/tokenmocks"
	"github.com/hyperledger/firefly/mocks/txcommonmocks"
	"github.com/hyperledger/firefly/pkg/core"
	"github.com/hyperledger/firefly/pkg/tokens"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const configDir = "../../test/data/config"

type testOrchestrator struct {
	orchestrator

	mdi *databasemocks.Plugin
	mdm *datamocks.Manager
	mbm *broadcastmocks.Manager
	mba *batchmocks.Manager
	mem *eventmocks.EventManager
	mnm *networkmapmocks.Manager
	mps *sharedstoragemocks.Plugin
	mpm *privatemessagingmocks.Manager
	mbi *blockchainmocks.Plugin
	mii *identitymocks.Plugin
	mim *identitymanagermocks.Manager
	mdx *dataexchangemocks.Plugin
	mam *assetmocks.Manager
	mti *tokenmocks.Plugin
	mcm *contractmocks.Manager
	mmi *metricsmocks.Manager
	mom *operationmocks.Manager
	mbp *batchpinmocks.Submitter
	mth *txcommonmocks.Helper
	msd *shareddownloadmocks.Manager
	mae *admineventsmocks.Manager
	mdh *definitionsmocks.DefinitionHandler
}

func (tor *testOrchestrator) cleanup(t *testing.T) {
	tor.mdi.AssertExpectations(t)
	tor.mdm.AssertExpectations(t)
	tor.mbm.AssertExpectations(t)
	tor.mba.AssertExpectations(t)
	tor.mem.AssertExpectations(t)
	tor.mnm.AssertExpectations(t)
	tor.mps.AssertExpectations(t)
	tor.mpm.AssertExpectations(t)
	tor.mbi.AssertExpectations(t)
	tor.mii.AssertExpectations(t)
	tor.mim.AssertExpectations(t)
	tor.mdx.AssertExpectations(t)
	tor.mam.AssertExpectations(t)
	tor.mti.AssertExpectations(t)
	tor.mcm.AssertExpectations(t)
	tor.mmi.AssertExpectations(t)
	tor.mom.AssertExpectations(t)
	tor.mbp.AssertExpectations(t)
	tor.mth.AssertExpectations(t)
	tor.msd.AssertExpectations(t)
	tor.mae.AssertExpectations(t)
	tor.mdh.AssertExpectations(t)
}

func newTestOrchestrator() *testOrchestrator {
	coreconfig.Reset()
	ctx, cancel := context.WithCancel(context.Background())
	tor := &testOrchestrator{
		orchestrator: orchestrator{
			ctx:       ctx,
			cancelCtx: cancel,
			namespace: "ns",
		},
		mdi: &databasemocks.Plugin{},
		mdm: &datamocks.Manager{},
		mbm: &broadcastmocks.Manager{},
		mba: &batchmocks.Manager{},
		mem: &eventmocks.EventManager{},
		mnm: &networkmapmocks.Manager{},
		mps: &sharedstoragemocks.Plugin{},
		mpm: &privatemessagingmocks.Manager{},
		mbi: &blockchainmocks.Plugin{},
		mii: &identitymocks.Plugin{},
		mim: &identitymanagermocks.Manager{},
		mdx: &dataexchangemocks.Plugin{},
		mam: &assetmocks.Manager{},
		mti: &tokenmocks.Plugin{},
		mcm: &contractmocks.Manager{},
		mmi: &metricsmocks.Manager{},
		mom: &operationmocks.Manager{},
		mbp: &batchpinmocks.Submitter{},
		mth: &txcommonmocks.Helper{},
		msd: &shareddownloadmocks.Manager{},
		mae: &admineventsmocks.Manager{},
		mdh: &definitionsmocks.DefinitionHandler{},
	}
	tor.orchestrator.database = tor.mdi
	tor.orchestrator.data = tor.mdm
	tor.orchestrator.batch = tor.mba
	tor.orchestrator.broadcast = tor.mbm
	tor.orchestrator.events = tor.mem
	tor.orchestrator.networkmap = tor.mnm
	tor.orchestrator.sharedstorage = tor.mps
	tor.orchestrator.messaging = tor.mpm
	tor.orchestrator.identity = tor.mim
	tor.orchestrator.dataexchange = tor.mdx
	tor.orchestrator.assets = tor.mam
	tor.orchestrator.contracts = tor.mcm
	tor.orchestrator.tokens = map[string]tokens.Plugin{"token": tor.mti}
	tor.orchestrator.blockchain = tor.mbi
	tor.orchestrator.metrics = tor.mmi
	tor.orchestrator.operations = tor.mom
	tor.orchestrator.batchpin = tor.mbp
	tor.orchestrator.sharedDownload = tor.msd
	tor.orchestrator.adminEvents = tor.mae
	tor.orchestrator.txHelper = tor.mth
	tor.orchestrator.definitions = tor.mdh
	tor.mdi.On("Name").Return("mock-di").Maybe()
	tor.mem.On("Name").Return("mock-ei").Maybe()
	tor.mps.On("Name").Return("mock-ps").Maybe()
	tor.mbi.On("Name").Return("mock-bi").Maybe()
	tor.mii.On("Name").Return("mock-ii").Maybe()
	tor.mdx.On("Name").Return("mock-dx").Maybe()
	tor.mam.On("Name").Return("mock-am").Maybe()
	tor.mti.On("Name").Return("mock-tk").Maybe()
	tor.mcm.On("Name").Return("mock-cm").Maybe()
	tor.mmi.On("Name").Return("mock-mm").Maybe()
	return tor
}

func TestNewOrchestrator(t *testing.T) {
	or := NewOrchestrator(
		"ns1",
		Config{
			Tokens: map[string]TokensPlugin{"token": {Plugin: &tokenmocks.Plugin{}}},
		},
		&metricsmocks.Manager{},
		&admineventsmocks.Manager{},
	)
	assert.NotNil(t, or)
}

func TestInitPluginsDatabaseFail(t *testing.T) {
	or := newTestOrchestrator()
	defer or.cleanup(t)
	or.mdi.On("Init", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("pop"))
	ctx := context.Background()
	err := or.initPlugins(ctx)
	assert.EqualError(t, err, "pop")
}

func TestInitPluginsBlockchainFail(t *testing.T) {
	or := newTestOrchestrator()
	defer or.cleanup(t)
	or.mdi.On("Init", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	or.mbi.On("Init", mock.Anything, mock.Anything, mock.Anything, or.mmi).Return(fmt.Errorf("pop"))
	ctx := context.Background()
	err := or.initPlugins(ctx)
	assert.EqualError(t, err, "pop")
}

func TestInitPluginsDataexchangeNodesFail(t *testing.T) {
	or := newTestOrchestrator()
	defer or.cleanup(t)
	or.mdi.On("Init", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	or.mbi.On("Init", mock.Anything, mock.Anything, mock.Anything, or.mmi).Return(nil)
	or.mdi.On("GetIdentities", mock.Anything, mock.Anything).Return(nil, nil, fmt.Errorf("pop"))
	ctx := context.Background()
	err := or.initPlugins(ctx)
	assert.EqualError(t, err, "pop")
}

func TestInitPluginsDataexchangeFail(t *testing.T) {
	or := newTestOrchestrator()
	defer or.cleanup(t)
	or.mdi.On("Init", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	or.mbi.On("Init", mock.Anything, mock.Anything, mock.Anything, or.mmi).Return(nil)
	or.mdi.On("GetIdentities", mock.Anything, mock.Anything).Return([]*core.Identity{{}}, nil, nil)
	or.mdx.On("Init", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("pop"))
	ctx := context.Background()
	err := or.initPlugins(ctx)
	assert.EqualError(t, err, "pop")
}

func TestInitPluginsSharedstorageFail(t *testing.T) {
	or := newTestOrchestrator()
	defer or.cleanup(t)
	or.mdi.On("Init", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	or.mbi.On("Init", mock.Anything, mock.Anything, mock.Anything, or.mmi).Return(nil)
	or.mdi.On("GetIdentities", mock.Anything, mock.Anything).Return([]*core.Identity{{}}, nil, nil)
	or.mdx.On("Init", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	or.mps.On("Init", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("pop"))
	ctx := context.Background()
	err := or.initPlugins(ctx)
	assert.EqualError(t, err, "pop")
}

func TestInitPluginsTokensFail(t *testing.T) {
	or := newTestOrchestrator()
	defer or.cleanup(t)
	or.mdi.On("Init", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	or.mbi.On("Init", mock.Anything, mock.Anything, mock.Anything, or.mmi).Return(nil)
	or.mdi.On("GetIdentities", mock.Anything, mock.Anything).Return([]*core.Identity{{}}, nil, nil)
	or.mdx.On("Init", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	or.mps.On("Init", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	or.mti.On("Init", mock.Anything, "token", mock.Anything, mock.Anything).Return(fmt.Errorf("pop"))
	ctx := context.Background()
	err := or.initPlugins(ctx)
	assert.EqualError(t, err, "pop")
}

func TestInitAllPluginsOK(t *testing.T) {
	or := newTestOrchestrator()
	defer or.cleanup(t)
	or.mdi.On("Init", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	or.mbi.On("Init", mock.Anything, mock.Anything, mock.Anything, or.mmi).Return(nil)
	or.mdi.On("GetIdentities", mock.Anything, mock.Anything).Return([]*core.Identity{{}}, nil, nil)
	or.mdx.On("Init", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	or.mps.On("Init", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	or.mti.On("Init", mock.Anything, "token", mock.Anything, mock.Anything).Return(nil)
	err := or.Init(or.ctx, or.cancelCtx)
	assert.NoError(t, err)
}

func TestInitMessagingComponentFail(t *testing.T) {
	or := newTestOrchestrator()
	defer or.cleanup(t)
	or.database = nil
	or.messaging = nil
	err := or.initComponents(context.Background())
	assert.Regexp(t, "FF10128", err)
}

func TestInitEventsComponentFail(t *testing.T) {
	or := newTestOrchestrator()
	defer or.cleanup(t)
	or.database = nil
	or.events = nil
	err := or.initComponents(context.Background())
	assert.Regexp(t, "FF10128", err)
}

func TestInitNetworkMapComponentFail(t *testing.T) {
	or := newTestOrchestrator()
	defer or.cleanup(t)
	or.database = nil
	or.networkmap = nil
	err := or.initComponents(context.Background())
	assert.Regexp(t, "FF10128", err)
}

func TestInitOperationComponentFail(t *testing.T) {
	or := newTestOrchestrator()
	defer or.cleanup(t)
	or.database = nil
	or.operations = nil
	err := or.initComponents(context.Background())
	assert.Regexp(t, "FF10128", err)
}

func TestInitSharedStorageDownloadComponentFail(t *testing.T) {
	or := newTestOrchestrator()
	defer or.cleanup(t)
	or.database = nil
	or.sharedDownload = nil
	err := or.initComponents(context.Background())
	assert.Regexp(t, "FF10128", err)
}

func TestInitBatchComponentFail(t *testing.T) {
	or := newTestOrchestrator()
	defer or.cleanup(t)
	or.database = nil
	or.batch = nil
	err := or.initComponents(context.Background())
	assert.Regexp(t, "FF10128", err)
}

func TestInitBroadcastComponentFail(t *testing.T) {
	or := newTestOrchestrator()
	defer or.cleanup(t)
	or.database = nil
	or.broadcast = nil
	err := or.initComponents(context.Background())
	assert.Regexp(t, "FF10128", err)
}

func TestInitDataComponentFail(t *testing.T) {
	or := newTestOrchestrator()
	defer or.cleanup(t)
	or.database = nil
	or.data = nil
	err := or.initComponents(context.Background())
	assert.Regexp(t, "FF10128", err)
}

func TestInitIdentityComponentFail(t *testing.T) {
	or := newTestOrchestrator()
	defer or.cleanup(t)
	or.database = nil
	or.identity = nil
	or.txHelper = nil
	err := or.initComponents(context.Background())
	assert.Regexp(t, "FF10128", err)
}

func TestInitAssetsComponentFail(t *testing.T) {
	or := newTestOrchestrator()
	defer or.cleanup(t)
	or.database = nil
	or.assets = nil
	err := or.initComponents(context.Background())
	assert.Regexp(t, "FF10128", err)
}

func TestInitContractsComponentFail(t *testing.T) {
	or := newTestOrchestrator()
	defer or.cleanup(t)
	or.database = nil
	or.contracts = nil
	err := or.initComponents(context.Background())
	assert.Regexp(t, "FF10128", err)
}

func TestInitDefinitionsComponentFail(t *testing.T) {
	or := newTestOrchestrator()
	defer or.cleanup(t)
	or.database = nil
	or.definitions = nil
	err := or.initComponents(context.Background())
	assert.Regexp(t, "FF10128", err)
}

func TestInitBatchPinComponentFail(t *testing.T) {
	or := newTestOrchestrator()
	defer or.cleanup(t)
	or.database = nil
	or.batchpin = nil
	err := or.initComponents(context.Background())
	assert.Regexp(t, "FF10128", err)
}

func TestInitOperationsComponentFail(t *testing.T) {
	or := newTestOrchestrator()
	defer or.cleanup(t)
	or.database = nil
	or.operations = nil
	err := or.initComponents(context.Background())
	assert.Regexp(t, "FF10128", err)
}

func TestStartBatchFail(t *testing.T) {
	coreconfig.Reset()
	or := newTestOrchestrator()
	defer or.cleanup(t)
	or.mdi.On("GetNamespace", mock.Anything, "ns").Return(nil, nil)
	or.mdi.On("UpsertNamespace", mock.Anything, mock.Anything, true).Return(nil)
	or.mbi.On("ConfigureContract", mock.Anything, mock.Anything).Return(nil)
	or.mbi.On("Start").Return(nil)
	or.mba.On("Start").Return(fmt.Errorf("pop"))
	err := or.Start()
	assert.EqualError(t, err, "pop")
}

func TestStartTokensFail(t *testing.T) {
	coreconfig.Reset()
	or := newTestOrchestrator()
	defer or.cleanup(t)
	or.mdi.On("GetNamespace", mock.Anything, "ns").Return(nil, nil)
	or.mdi.On("UpsertNamespace", mock.Anything, mock.Anything, true).Return(nil)
	or.mbi.On("ConfigureContract", mock.Anything, &core.FireFlyContracts{}).Return(nil)
	or.mbi.On("Start").Return(nil)
	or.mba.On("Start").Return(nil)
	or.mem.On("Start").Return(nil)
	or.mbm.On("Start").Return(nil)
	or.mpm.On("Start").Return(nil)
	or.msd.On("Start").Return(nil)
	or.mom.On("Start").Return(nil)
	or.mti.On("Start").Return(fmt.Errorf("pop"))
	err := or.Start()
	assert.EqualError(t, err, "pop")
}

func TestStartBlockchainsFail(t *testing.T) {
	coreconfig.Reset()
	or := newTestOrchestrator()
	defer or.cleanup(t)
	or.mdi.On("GetNamespace", mock.Anything, "ns").Return(nil, nil)
	or.mbi.On("ConfigureContract", mock.Anything, &core.FireFlyContracts{}).Return(nil)
	or.mbi.On("Start").Return(fmt.Errorf("pop"))
	err := or.Start()
	assert.EqualError(t, err, "pop")
}

func TestStartBlockchainsConfigureFail(t *testing.T) {
	coreconfig.Reset()
	or := newTestOrchestrator()
	defer or.cleanup(t)
	or.mdi.On("GetNamespace", mock.Anything, "ns").Return(nil, nil)
	or.mbi.On("ConfigureContract", mock.Anything, &core.FireFlyContracts{}).Return(fmt.Errorf("pop"))
	err := or.Start()
	assert.EqualError(t, err, "pop")
}

func TestStartStopOk(t *testing.T) {
	coreconfig.Reset()
	or := newTestOrchestrator()
	defer or.cleanup(t)
	or.mdi.On("GetNamespace", mock.Anything, "ns").Return(nil, nil)
	or.mdi.On("UpsertNamespace", mock.Anything, mock.Anything, true).Return(nil)
	or.mbi.On("ConfigureContract", mock.Anything, &core.FireFlyContracts{}).Return(nil)
	or.mbi.On("Start").Return(nil)
	or.mba.On("Start").Return(nil)
	or.mem.On("Start").Return(nil)
	or.mbm.On("Start").Return(nil)
	or.mpm.On("Start").Return(nil)
	or.mti.On("Start").Return(nil)
	or.msd.On("Start").Return(nil)
	or.mom.On("Start").Return(nil)
	or.mba.On("WaitStop").Return(nil)
	or.mbm.On("WaitStop").Return(nil)
	or.mdm.On("WaitStop").Return(nil)
	or.msd.On("WaitStop").Return(nil)
	or.mom.On("WaitStop").Return(nil)
	err := or.Start()
	assert.NoError(t, err)
	or.WaitStop()
	or.WaitStop() // swallows dups
}

func TestGetComponents(t *testing.T) {
	or := newTestOrchestrator()
	defer or.cleanup(t)

	or.mdi.On("Init", mock.Anything, mock.Anything, mock.Anything).Return(fmt.Errorf("pop"))

	ctx, cancelCtx := context.WithCancel(context.Background())
	err := or.Init(ctx, cancelCtx)
	assert.EqualError(t, err, "pop")

	assert.Equal(t, or.mbm, or.Broadcast())
	assert.Equal(t, or.mpm, or.PrivateMessaging())
	assert.Equal(t, or.mem, or.Events())
	assert.Equal(t, or.mba, or.BatchManager())
	assert.Equal(t, or.mnm, or.NetworkMap())
	assert.Equal(t, or.mdm, or.Data())
	assert.Equal(t, or.mam, or.Assets())
	assert.Equal(t, or.mcm, or.Contracts())
	assert.Equal(t, or.mmi, or.Metrics())
	assert.Equal(t, or.mom, or.Operations())
}

func TestNetworkAction(t *testing.T) {
	or := newTestOrchestrator()
	or.mim.On("NormalizeSigningKey", context.Background(), "ff_system", "", identity.KeyNormalizationBlockchainPlugin).Return("0x123", nil)
	or.mbi.On("SubmitNetworkAction", context.Background(), mock.Anything, "0x123", core.NetworkActionTerminate).Return(nil)
	err := or.SubmitNetworkAction(context.Background(), "ff_system", &core.NetworkAction{Type: core.NetworkActionTerminate})
	assert.NoError(t, err)
}

func TestNetworkActionBadKey(t *testing.T) {
	or := newTestOrchestrator()
	or.mim.On("NormalizeSigningKey", context.Background(), "ff_system", "", identity.KeyNormalizationBlockchainPlugin).Return("", fmt.Errorf("pop"))
	err := or.SubmitNetworkAction(context.Background(), "ff_system", &core.NetworkAction{Type: core.NetworkActionTerminate})
	assert.EqualError(t, err, "pop")
}

func TestNetworkActionBadType(t *testing.T) {
	or := newTestOrchestrator()
	or.mim.On("NormalizeSigningKey", context.Background(), "ff_system", "", identity.KeyNormalizationBlockchainPlugin).Return("0x123", nil)
	err := or.SubmitNetworkAction(context.Background(), "ff_system", &core.NetworkAction{Type: "bad"})
	assert.Regexp(t, "FF10397", err)
}

func TestNetworkActionBadNamespace(t *testing.T) {
	or := newTestOrchestrator()
	or.mim.On("NormalizeSigningKey", context.Background(), "ns", "", identity.KeyNormalizationBlockchainPlugin).Return("0x123", nil)
	err := or.SubmitNetworkAction(context.Background(), "ns", &core.NetworkAction{Type: core.NetworkActionTerminate})
	assert.Regexp(t, "FF10399", err)
}
