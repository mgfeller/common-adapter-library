package adapter

import (
	"context"
	"io/ioutil"

	"gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/layer5io/gokit/logger"
	"github.com/layer5io/gokit/models"
	"github.com/mgfeller/common-adapter-library/config"
)

// Handler provides the methods supported by the adapter
type Handler interface {
	GetName() string
	CreateInstance([]byte, string, *chan interface{}) error
	ApplyOperation(context.Context, string, string, bool) error
	ListOperations() (Operations, error)

	StreamErr(*Event, error)
	StreamInfo(*Event)
}

// BaseAdapter holds the dependencies for kuma-adapter
type BaseAdapter struct {
	Config  config.Handler
	Log     logger.Handler
	Channel *chan interface{}

	KubeClient     *kubernetes.Clientset
	KubeConfigPath string
	SmiChart       string
}

// newClient creates a new client
func (h *BaseAdapter) CreateInstance(kubeconfig []byte, contextName string, ch *chan interface{}) error {

	h.Channel = ch
	h.KubeConfigPath = h.Config.GetKey("kube-config-path")

	config, err := h.clientConfig(kubeconfig, contextName)
	if err != nil {
		return ErrClientConfig(err)
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return ErrClientSet(err)
	}

	h.KubeClient = clientset

	return nil
}

// configClient creates a Config client
func (h *BaseAdapter) clientConfig(kubeconfig []byte, contextName string) (*rest.Config, error) {
	if len(kubeconfig) > 0 {
		ccfg, err := clientcmd.Load(kubeconfig)
		if err != nil {
			return nil, err
		}
		if contextName != "" {
			ccfg.CurrentContext = contextName
		}
		err = writeKubeconfig(kubeconfig, contextName, h.KubeConfigPath)
		if err != nil {
			return nil, err
		}
		return clientcmd.NewDefaultClientConfig(*ccfg, &clientcmd.ConfigOverrides{}).ClientConfig()
	}
	return rest.InClusterConfig()
}

// writeKubeconfig creates kubeconfig in local container
func writeKubeconfig(kubeconfig []byte, contextName string, path string) error {

	yamlConfig := models.Kubeconfig{}
	err := yaml.Unmarshal(kubeconfig, &yamlConfig)
	if err != nil {
		return err
	}

	yamlConfig.CurrentContext = contextName

	d, err := yaml.Marshal(yamlConfig)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, d, 0600)
	if err != nil {
		return err
	}

	return nil
}
