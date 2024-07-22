// Copyright 2019 Altinity Ltd and/or its affiliates. All rights reserved.
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

package kube

import (
	"context"
	"sigs.k8s.io/controller-runtime/pkg/client"

	core "k8s.io/api/core/v1"
)

type EventKeeper struct {
	kubeClient client.Client
}

func NewEventKeeper(kubeClient client.Client) *EventKeeper {
	return &EventKeeper{
		kubeClient: kubeClient,
	}
}

func (c *EventKeeper) Create(ctx context.Context, event *core.Event) (*core.Event, error) {
	err := c.kubeClient.Create(ctx, event)
	return event, err
}