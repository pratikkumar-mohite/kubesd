package cli

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type client struct {
	conn *kubernetes.Clientset
}

type secretObject struct {
	Kind      string
	Type      string
	Name      string
	Namespace string
	Data      map[string][]byte
}

var (
	masterURL  string
	kubeconfig string
)

// incluster
func CreateConnection() *client {
	c := &client{}
	flag.StringVar(&kubeconfig, "kubeconfig", defaultKubeconfig(), "Path to a kubeconfig")
	flag.StringVar(&masterURL, "master", "https://127.0.0.1:57331", "The url of the Kubernetes API server.")
	flag.Parse()
	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}
	c.conn, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return c
}

func (c *client) ReadSecret(name string, namespace string) (strings.Builder, error) {
	secret, err := c.conn.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		fmt.Println(err.Error())
	}
	return secretToString(secret), nil
}

func secretToString(secret *v1.Secret) strings.Builder {
	var builder strings.Builder
	data := make(map[string][]byte)
	for k, v := range secret.Data {
		data[k] = v
	}
	s := secretObject{
		Kind:      "Secret",
		Type:      string(secret.Type),
		Name:      secret.Name,
		Namespace: secret.Namespace,
		Data:      data,
	}

	yamlData, err := yaml.Marshal(&s)

	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
	}
	builder.WriteString(string(yamlData))

	return builder
}

func defaultKubeconfig() string {
	fname := os.Getenv("KUBECONFIG")
	if fname != "" {
		return fname
	}
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("failed to get home directory: %v", err)
		return ""
	}
	return filepath.Join(home, ".kube", "config")
}
