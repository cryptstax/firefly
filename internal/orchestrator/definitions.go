// Copyright © 2021 Kaleido, Inc.
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
	"encoding/json"

	"github.com/kaleido-io/firefly/internal/i18n"
	"github.com/kaleido-io/firefly/pkg/fftypes"
)

func (or *orchestrator) broadcastDefinition(ctx context.Context, ns string, defObject interface{}, topic string) (msg *fftypes.Message, err error) {

	// Serialize it into a data object, as a piece of data we can write to a message
	data := &fftypes.Data{
		Validator: fftypes.ValidatorTypeDataDefinition,
		ID:        fftypes.NewUUID(),
		Namespace: ns,
		Created:   fftypes.Now(),
	}
	b, _ := json.Marshal(&defObject)
	_ = json.Unmarshal(b, &data.Value)
	data.Hash, _ = data.Value.Hash(ctx, "value")

	// Write as data to the local store
	if err = or.database.UpsertData(ctx, data, true, false /* we just generated the ID, so it is new */); err != nil {
		return nil, err
	}

	// Create a broadcast message referring to the data
	msg = &fftypes.Message{
		Header: fftypes.MessageHeader{
			Namespace: ns,
			Type:      fftypes.MessageTypeDefinition,
			Author:    or.nodeIDentity,
			Topic:     topic,
			Context:   fftypes.SystemContext,
			TX: fftypes.TransactionRef{
				Type: fftypes.TransactionTypePin,
			},
		},
		Data: fftypes.DataRefs{
			{ID: data.ID, Hash: data.Hash},
		},
	}

	// Broadcast the message
	if err = or.broadcast.BroadcastMessage(ctx, msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func (or *orchestrator) BroadcastDataDefinition(ctx context.Context, ns string, dataDef *fftypes.DataDefinition) (msg *fftypes.Message, err error) {

	// Validate the input data definition data
	dataDef.ID = fftypes.NewUUID()
	dataDef.Created = fftypes.Now()
	dataDef.Namespace = ns
	if dataDef.Validator == "" {
		dataDef.Validator = fftypes.ValidatorTypeJSON
	}
	if dataDef.Validator != fftypes.ValidatorTypeJSON {
		return nil, i18n.NewError(ctx, i18n.MsgUnknownFieldValue, "validator")
	}
	if err = or.verifyNamespaceExists(ctx, dataDef.Namespace); err != nil {
		return nil, err
	}
	if err = fftypes.ValidateFFNameField(ctx, dataDef.Name, "name"); err != nil {
		return nil, err
	}
	if err = fftypes.ValidateFFNameField(ctx, dataDef.Version, "version"); err != nil {
		return nil, err
	}
	if len(dataDef.Value) == 0 {
		return nil, i18n.NewError(ctx, i18n.MsgMissingRequiredField, "value")
	}
	if dataDef.Hash, err = dataDef.Value.Hash(ctx, "value"); err != nil {
		return nil, err
	}
	return or.broadcastDefinition(ctx, ns, dataDef, fftypes.DataDefinitionTopicName)
}

func (or *orchestrator) BroadcastNamespaceDefinition(ctx context.Context, ns *fftypes.Namespace) (msg *fftypes.Message, err error) {

	// Validate the input data definition data
	ns.ID = fftypes.NewUUID()
	ns.Created = fftypes.Now()
	if err = fftypes.ValidateFFNameField(ctx, ns.Name, "name"); err != nil {
		return nil, err
	}
	if err = fftypes.ValidateLength(ctx, ns.Description, "description", 4096); err != nil {
		return nil, err
	}

	return or.broadcastDefinition(ctx, ns.Name, ns, fftypes.NamespaceDefinitionTopicName)
}
