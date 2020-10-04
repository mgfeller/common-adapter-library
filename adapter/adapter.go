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

	"github.com/sirupsen/logrus"

	"k8s.io/client-go/dynamic"

	"k8s.io/client-go/kubernetes"

	"github.com/layer5io/gokit/logger"
	"github.com/mgfeller/common-adapter-library/config"
)

type Handler interface {
	GetName() string
	CreateInstance([]byte, string, *chan interface{}) error
	ApplyOperation(context.Context, OperationRequest) error
	ListOperations() (Operations, error)

	StreamErr(*Event, error)
	StreamInfo(*Event)
}

type BaseHandler struct {
	Config  config.Handler
	Log     logger.Handler
	Channel *chan interface{}

	KubeClient        *kubernetes.Clientset
	DynamicKubeClient dynamic.Interface
	KubeConfigPath    string
	SmiChart          string
}

type OperationRequest struct {
	OperationName     string
	Namespace         string
	Username          string
	CustomBody        string
	IsDeleteOperation bool
	OperationID       string
}

func (h *BaseHandler) CreateInstance(kubeconfig []byte, contextName string, ch *chan interface{}) error {
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

// creates the namespace unless it is 'default', or it is a delete operation
func (h *BaseHandler) CreateNamespace(isDelete bool, namespace string) error {
	if !isDelete && namespace != "default" {
		if err := h.createNamespace(context.TODO(), namespace); err != nil {
			logrus.Error(err)
			return err
		}
	}
	return nil
}

func (h *BaseHandler) ApplyKubernetesManifest(request OperationRequest, operation Operation, mergeData map[string]string, templatePath string) error {
	if err := h.applyK8sManifest(context.TODO(), request, operation, mergeData, templatePath); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
