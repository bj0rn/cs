package switcher

import (
	"github.com/bj0rn/cs/pkg/ao"
	"github.com/bj0rn/cs/pkg/kubeconfig"
	"github.com/sirupsen/logrus"
	"net/url"
	"strings"
)

type Switcher struct {
	KubeconfigPath string
	AoConfigPath   string
}

func NewSwitcher(kubeconfigPath string, aoconfigPath string) *Switcher {
	return &Switcher{
		KubeconfigPath: kubeconfigPath,
		AoConfigPath:   aoconfigPath,
	}
}

func (s *Switcher) Switch(clusterShortname string, namespace string) error {

	aoconfig, err := ao.Load(s.AoConfigPath)
	if err != nil {
		return err
	}

	kubeconfig, err := kubeconfig.Load(s.KubeconfigPath)
	if err != nil {
		return err
	}

	c, err := aoconfig.GetCluster(clusterShortname)
	if err != nil {
		return err
	}

	cluster, err := createClusterName(c.URL)
	if err != nil {
		return err
	}

	if ok, err := kubeconfig.UserExists(cluster); !ok {
		if err != nil {
			logrus.Infof("Never seen this cluster before. Please login")
			return nil
		} else {
			return err
		}
	}

	if ok, err := kubeconfig.HasContext(cluster, namespace); ok {
		if err != nil {
			return err
		}

		context, err := kubeconfig.GetContextName(cluster, namespace)
		if err != nil {
			return err
		}

		kubeconfig.SetCurrentContext(context)
		logrus.Infof("Current context is set to %s", context)

	} else {

		context, err := kubeconfig.CreateContext(namespace, cluster)
		if err != nil {
			return err
		}
		logrus.Infof("Created context: %s", context)

		kubeconfig.SetCurrentContext(context)
		logrus.Infof("Current context is set to %s", context)
	}

	err = kubeconfig.Save(s.KubeconfigPath)
	if err != nil {
		return nil
	}

	logrus.Info("Kubeconfig saved")

	return nil
}

func createClusterName(clusterUrl string) (string, error) {
	u, err := url.Parse(clusterUrl)
	if err != nil {
		return "", err
	}
	return strings.Replace(u.Host, ".", "-", -1), nil
}
