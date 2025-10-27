package libkubernetes

import (
	snapshotv1beta1 "github.com/kubernetes-csi/external-snapshotter/client/v3/clientset/versioned/typed/volumesnapshot/v1beta1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/wangweihong/eazycloud/apis/ikubernetes"
)

type Client struct {
	client *kubernetes.Clientset
	//metricClient   *metricctl.Clientset
	//	monitorClient  *monitoringv1.MonitoringV1Client
	snapshotClient *snapshotv1beta1.SnapshotV1beta1Client
}

func NewClient(config *ikubernetes.ClusterConfig) (*Client, error) {
	k8sConfig := config.ToRestConfig()

	clientSet, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, err
	}
	//monitorclientSet, err := monitoringv1.NewForConfig(k8sConfig)
	//if err != nil {
	//	return nil, status.NewStatusDesc(scode.ScodeNewK8SClientError, err.Error())
	//}

	snapshotClinetSet, err := snapshotv1beta1.NewForConfig(k8sConfig)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: clientSet,
		//	monitorClient:  monitorclientSet,
		snapshotClient: snapshotClinetSet,
	}, nil
}

func NewClusterConfigFromKubeconfig(kubeconfig string) (*ikubernetes.ClusterConfig, error) {
	cc, err := clientcmd.NewClientConfigFromBytes([]byte(kubeconfig))
	if err != nil {
		return nil, err
	}
	clientConfig, err := cc.ClientConfig()
	if err != nil {
		return nil, err
	}

	config := &ikubernetes.ClusterConfig{
		Host: clientConfig.Host,
		TlsClientConfig: ikubernetes.TLSClientConfig{
			// CAData:   base64.StdEncoding.EncodeToString([]byte(clientConfig.TLSClientConfig.CAData)),
			// KeyData:  base64.StdEncoding.EncodeToString([]byte(clientConfig.TLSClientConfig.KeyData)),
			// CertData: base64.StdEncoding.EncodeToString([]byte(clientConfig.TLSClientConfig.CertData)),
			CAData:   string(clientConfig.TLSClientConfig.CAData),
			KeyData:  string(clientConfig.TLSClientConfig.KeyData),
			CertData: string(clientConfig.TLSClientConfig.CertData),
		},
	}
	return config, nil
}
