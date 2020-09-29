// Copyright 2020 Michael Gfeller
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

package adapter

import (
	"context"
	"k8s.io/client-go/dynamic"

	"k8s.io/client-go/kubernetes"

	"github.com/layer5io/gokit/logger"
	"github.com/mgfeller/common-adapter-library/config"
)

type Handler interface {
	GetName() string
	CreateInstance([]byte, string, *chan interface{}) error
	ApplyOperation(context.Context, string, string, bool) error
	ListOperations() (Operations, error)

	StreamErr(*Event, error)
	StreamInfo(*Event)
}

type BaseAdapter struct {
	Config  config.Handler
	Log     logger.Handler
	Channel *chan interface{}

	KubeClient        *kubernetes.Clientset
	DynamicKubeClient dynamic.Interface
	KubeConfigPath    string
	SmiChart          string
}

func (h *BaseAdapter) CreateInstance(kubeconfig []byte, contextName string, ch *chan interface{}) error {
	h.Channel = ch
	h.KubeConfigPath = h.Config.GetKey("kube-config-path")

	config, err := h.k8sClientConfig(kubeconfig, contextName)
	if err != nil {
		return ErrClientConfig(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return ErrClientSet(err)
	}

	h.KubeClient = clientset

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}
	h.DynamicKubeClient = dynamicClient

	return nil
}
