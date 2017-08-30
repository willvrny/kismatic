package ansible

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

type ClusterCatalog struct {
	ClusterName               string `yaml:"kubernetes_cluster_name"`
	AdminPassword             string `yaml:"kubernetes_admin_password"`
	TLSDirectory              string `yaml:"tls_directory"`
	ServicesCIDR              string `yaml:"kubernetes_services_cidr"`
	PodCIDR                   string `yaml:"kubernetes_pods_cidr"`
	DNSServiceIP              string `yaml:"kubernetes_dns_service_ip"`
	EnableModifyHosts         bool   `yaml:"modify_hosts_file"`
	EnablePackageInstallation bool   `yaml:"allow_package_installation"`
	PackageRepoURLs           string `yaml:"package_repository_urls"`
	DisconnectedInstallation  bool   `yaml:"disconnected_installation"`
	SeedRegistry              bool   `yaml:"seed_registry"`
	KuberangPath              string `yaml:"kuberang_path"`
	LoadBalancedFQDN          string `yaml:"kubernetes_load_balanced_fqdn"`

	APIServerOptions map[string]string `yaml:"kubernetes_api_server_option_overrides"`

	ConfigureDockerWithPrivateRegistry bool   `yaml:"configure_docker_with_private_registry"`
	DeployInternalDockerRegistry       bool   `yaml:"deploy_internal_docker_registry"`
	DockerCAPath                       string `yaml:"docker_certificates_ca_path"`
	DockerRegistryAddress              string `yaml:"docker_registry_address"`
	DockerRegistryPort                 string `yaml:"docker_registry_port"`

	ForceEtcdRestart              bool `yaml:"force_etcd_restart"`
	ForceAPIServerRestart         bool `yaml:"force_apiserver_restart"`
	ForceControllerManagerRestart bool `yaml:"force_controller_manager_restart"`
	ForceSchedulerRestart         bool `yaml:"force_scheduler_restart"`
	ForceProxyRestart             bool `yaml:"force_proxy_restart"`
	ForceKubeletRestart           bool `yaml:"force_kubelet_restart"`
	ForceCalicoNodeRestart        bool `yaml:"force_calico_node_restart"`
	ForceDockerRestart            bool `yaml:"force_docker_restart"`

	EnableConfigureIngress bool `yaml:"configure_ingress"`

	KismaticPreflightCheckerLinux string `yaml:"kismatic_preflight_checker"`
	KismaticPreflightCheckerLocal string `yaml:"kismatic_preflight_checker_local"`

	WorkerNode string `yaml:"worker_node"`

	NFSVolumes []NFSVolume `yaml:"nfs_volumes"`

	EnableGluster bool `yaml:"configure_storage"`

	// volume add vars
	VolumeName              string   `yaml:"volume_name"`
	VolumeReplicaCount      int      `yaml:"volume_replica_count"`
	VolumeDistributionCount int      `yaml:"volume_distribution_count"`
	VolumeStorageClass      string   `yaml:"volume_storage_class"`
	VolumeQuotaGB           int      `yaml:"volume_quota_gb"`
	VolumeQuotaBytes        int      `yaml:"volume_quota_bytes"`
	VolumeMount             string   `yaml:"volume_mount"`
	VolumeAllowedIPs        string   `yaml:"volume_allow_ips"`
	VolumeReclaimPolicy     string   `yaml:"volume_reclaim_policy"`
	VolumeAccessModes       []string `yaml:"volume_access_modes"`

	TargetVersion string `yaml:"kismatic_short_version"`

	OnlineUpgrade bool `yaml:"online_upgrade"`

	DiagnosticsDirectory string `yaml:"diagnostics_dir"`
	DiagnosticsDateTime  string `yaml:"diagnostics_date_time"`

	DockerDirectLVMEnabled                 bool   `yaml:"docker_direct_lvm_enabled"`
	DockerDirectLVMBlockDevicePath         string `yaml:"docker_direct_lvm_block_device_path"`
	DockerDirectLVMDeferredDeletionEnabled bool   `yaml:"docker_direct_lvm_deferred_deletion_enabled"`

	LocalKubeconfigDirectory string `yaml:"local_kubeconfig_directory"`

	CloudProvider string `yaml:"cloud_provider"`
	CloudConfig   string `yaml:"cloud_config_local"`

	DNS struct {
		Enabled bool
	}

	RunPodValidation bool `yaml:"run_pod_validation"`

	CNI struct {
		Enabled  bool
		Provider string
		Options  struct {
			Calico struct {
				Mode string
			}
		}
	}

	Heapster struct {
		Enabled bool
		Options struct {
			Heapster struct {
				Replicas    int    `yaml:"replicas"`
				Sink        string `yaml:"sink"`
				ServiceType string `yaml:"service_type"`
			}
			InfluxDB struct {
				PVCName string `yaml:"pvc_name"`
			}
		}
	}

	Dashboard struct {
		Enabled bool
	}

	Helm struct {
		Enabled bool
	}

	InsecureNetworkingEtcd bool `yaml:"insecure_networking_etcd"`

	HTTPProxy  string `yaml:"http_proxy"`
	HTTPSProxy string `yaml:"https_proxy"`
	NoProxy    string `yaml:"no_proxy"`
}

type NFSVolume struct {
	Host string
	Path string
}

func (c *ClusterCatalog) EnableRestart() {
	c.ForceEtcdRestart = true
	c.ForceAPIServerRestart = true
	c.ForceControllerManagerRestart = true
	c.ForceSchedulerRestart = true
	c.ForceProxyRestart = true
	c.ForceKubeletRestart = true
	c.ForceCalicoNodeRestart = true
	c.ForceDockerRestart = true
}

func (c *ClusterCatalog) ToYAML() ([]byte, error) {
	bytez, marshalErr := yaml.Marshal(c)
	if marshalErr != nil {
		return []byte{}, fmt.Errorf("error marshalling plan to yaml: %v", marshalErr)
	}

	return bytez, nil
}
