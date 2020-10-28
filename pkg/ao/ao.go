package ao

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// AOConfig is a structure of the configuration of ao
type AOConfig struct {
	RefName     string              `json:"refName"`
	APICluster  string              `json:"apiCluster"`
	Affiliation string              `json:"affiliation"`
	Localhost   bool                `json:"localhost"`
	Clusters    map[string]*Cluster `json:"clusters"`

	ServiceURLPatterns map[string]*ServiceURLPatterns `json:"serviceURLPatterns"`
	ClusterConfig      map[string]*ClusterConfig      `json:"clusterConfig"`

	AvailableClusters       []string `json:"availableClusters"`
	PreferredAPIClusters    []string `json:"preferredApiClusters"`
	AvailableUpdateClusters []string `json:"availableUpdateClusters"`
	ClusterURLPattern       string   `json:"clusterUrlPattern"`
	BooberURLPattern        string   `json:"booberUrlPattern"`
	UpdateURLPattern        string   `json:"updateUrlPattern"`
	GoboURLPattern          string   `json:"goboUrlPattern"`

	FileAOVersion string `json:"aoVersion"` // For detecting possible changes to saved file
}

// Cluster holds information of Openshift cluster
type Cluster struct {
	Name      string `json:"name"`
	URL       string `json:"url"`
	LoginURL  string `json:"loginUrl"`
	Token     string `json:"token"`
	Reachable bool   `json:"reachable"`
	BooberURL string `json:"booberUrl"`
	GoboURL   string `json:"goboUrl"`
}

type ClusterConfig struct {
	Type             string `json:"type"`
	IsAPICluster     bool   `json:"isApiCluster"`
	IsUpdateCluster  bool   `json:"isUpdateCluster"`
	ClusterURLPrefix string `json:"clusterUrlPrefix"`
}

// ServiceURLPatterns contains url patterns for all integrations made with AO.
// %s will be replaced with cluster name. If ClusterURLPrefix in ClusterConfig is specified
// it will be used for ClusterURLPattern and ClusterLoginURLPattern insted of cluster name.
type ServiceURLPatterns struct {
	ClusterURLPattern      string `json:"clusterUrlPattern"`
	ClusterLoginURLPattern string `json:"clusterLoginUrlPattern"`
	BooberURLPattern       string `json:"booberUrlPattern"`
	UpdateURLPattern       string `json:"updateUrlPattern"`
	GoboURLPattern         string `json:"goboUrlPattern"`
}

// ServiceURLs contains all the necessary URLs for integrations made with AO.
type ServiceURLs struct {
	BooberURL       string
	ClusterURL      string
	ClusterLoginURL string
	GoboURL         string
}

func Load(path string) (*AOConfig, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := AOConfig{}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (aoConfig *AOConfig) GetCluster(shortName string) (*Cluster, error) {
	for key, cluster := range aoConfig.Clusters {
		if key == shortName {
			return cluster, nil
		}
	}

	return nil, fmt.Errorf("Could not find cluster: %s", shortName)
}
