package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	API        API
	Kubeconfig string
	Provider   string
	CASTAI     *CASTAI
	EKS        *EKS
}

type API struct {
	Key string
	URL string
}

type CASTAI struct {
	ClusterID      string
	OrganizationID string
}

type EKS struct {
	AccountID   string
	Region      string
	ClusterName string
}

var cfg *Config

// Get configuration bound to environment variables.
func Get() Config {
	if cfg != nil {
		return *cfg
	}

	_ = viper.BindEnv("api.key", "API_KEY")
	_ = viper.BindEnv("api.url", "API_URL")

	_ = viper.BindEnv("kubeconfig")

	_ = viper.BindEnv("provider")

	_ = viper.BindEnv("castai.clusterid", "CASTAI_CLUSTER_ID")
	_ = viper.BindEnv("castai.organizationid", "CASTAI_ORGANIZATION_ID")

	_ = viper.BindEnv("eks.accountid", "EKS_ACCOUNT_ID")
	_ = viper.BindEnv("eks.region", "EKS_REGION")
	_ = viper.BindEnv("eks.clustername", "EKS_CLUSTER_NAME")

	cfg = &Config{}
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("parsing configuration: %v", err))
	}

	if cfg.API.Key == "" {
		required("API_KEY")
	}
	if cfg.API.URL == "" {
		required("API_URL")
	}

	if cfg.CASTAI != nil {
		if cfg.CASTAI.ClusterID == "" {
			requiredDiscoveryDisabled("CASTAI_CLUSTER_ID")
		}
		if cfg.CASTAI.OrganizationID == "" {
			requiredDiscoveryDisabled("CASTAI_ORGANIZATION_ID")
		}
	}

	if cfg.EKS != nil {
		if cfg.EKS.AccountID == "" {
			requiredDiscoveryDisabled("EKS_ACCOUNT_ID")
		}
		if cfg.EKS.Region == "" {
			requiredDiscoveryDisabled("EKS_REGION")
		}
		if cfg.EKS.ClusterName == "" {
			requiredDiscoveryDisabled("EKS_CLUSTER_NAME")
		}
	}

	return *cfg
}

// Reset is used only for unit testing to reset configuration and rebind variables.
func Reset() {
	cfg = nil
}

func required(variable string) {
	panic(fmt.Errorf("env variable %s is required", variable))
}

func requiredDiscoveryDisabled(variable string) {
	panic(fmt.Errorf("env variable %s is required when discovery is disabled", variable))
}