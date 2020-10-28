package kubeconfig

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os/user"
)

type Kubeconfig struct {
	ApiVersion     string      `yaml:"apiVersion,omitempty"`
	Kind           string      `yaml:"kind,omitempty"`
	Clusters       []Clusters  `yaml:"clusters,omitempty"`
	Contexts       []Contexts  `yaml:"contexts,omitempty"`
	CurrentContext string      `yaml:"current-context"`
	Preferences    interface{} `yaml:"preferences,omitempty"`
	Users          []Users     `yaml:"users,omitempty"`
}

type Clusters struct {
	Cluster Cluster `yaml:"cluster,omitempty"`
	Name    string  `yaml:"name,omitempty"`
}

type Cluster struct {
	Server string `yaml:"server,omitempty"`
}

type Contexts struct {
	Context Context `yaml:"context,omitempty"`
	Name    string  `yaml:"name,omitempty"`
}

type Context struct {
	Cluster   string `yaml:"cluster,omitempty"`
	Namespace string `yaml:"namespace,omitempty"`
	User      string `yaml:"user,omitempty"`
}

type Users struct {
	Name string `yaml:"name,omitempty"`
	User User   `yaml:"user,omitempty"`
}

type User struct {
	Token string `yaml:"token,omitempty"`
}

func (kubeconfig *Kubeconfig) SetCurrentContext(context string) *Kubeconfig {
	kubeconfig.CurrentContext = context
	return kubeconfig
}

func (kubeconfig *Kubeconfig) CreateContext(namespace string, cluster string) (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}

	context := Contexts{
		Name: fmt.Sprintf("%s/%s/%s", namespace, cluster, user.Username),
		Context: Context{
			Cluster:   cluster,
			Namespace: namespace,
			User:      fmt.Sprintf("%s/%s", user.Username, cluster),
		},
	}

	kubeconfig.Contexts = append(kubeconfig.Contexts, context)

	return context.Name, nil
}

func (kubeconfig *Kubeconfig) UserExists(cluster string) (bool, error) {
	user, err := user.Current()
	if err != nil {
		return false, err
	}

	username := fmt.Sprintf("%s/%s", user.Username, cluster)

	for _, u := range kubeconfig.Users {
		if u.Name == username {
			return true, nil
		}
	}
	return false, fmt.Errorf("Could not find user: %s", username)
}

func (kubeconfig *Kubeconfig) HasContext(cluster string, namespace string) (bool, error) {

	user, err := user.Current()
	if err != nil {
		return false, err
	}

	contextName := fmt.Sprintf("%s/%s/%s", namespace, cluster, user.Username)

	for _, ctx := range kubeconfig.Contexts {

		if ctx.Name == contextName {
			return true, nil
		}

	}
	return false, nil
}

func (kubeconfig *Kubeconfig) GetContextName(cluster string, namespace string) (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}
	contextName := fmt.Sprintf("%s/%s/%s", namespace, cluster, user.Username)

	for _, ctx := range kubeconfig.Contexts {
		if ctx.Name == contextName {
			return contextName, nil
		}
	}
	return "", nil
}

func Load(path string) (*Kubeconfig, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := Kubeconfig{}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (kubeconfig *Kubeconfig) Save(path string) error {

	data, err := yaml.Marshal(kubeconfig)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, 0644)
}
