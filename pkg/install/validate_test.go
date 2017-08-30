package install

import (
	"fmt"
	"testing"
)

var validPlan = Plan{
	Cluster: Cluster{
		Name:          "test",
		AdminPassword: "password",
		Networking: NetworkConfig{
			Type:             "overlay",
			PodCIDRBlock:     "172.16.0.0/16",
			ServiceCIDRBlock: "172.20.0.0/16",
		},
		Certificates: CertsConfig{
			Expiry: "17250h",
		},
		SSH: SSHConfig{
			User: "root",
			Key:  "/bin/sh",
			Port: 22,
		},
	},
	AddOns: AddOns{
		CNI: &CNI{
			Provider: "calico",
			Options: CNIOptions{
				Calico: CalicoOptions{
					Mode: "overlay",
				},
			},
		},
		HeapsterMonitoring: &HeapsterMonitoring{
			Options: HeapsterOptions{
				Heapster: Heapster{
					Replicas:    2,
					ServiceType: "ClusterIP",
				},
			},
		},
	},
	Etcd: NodeGroup{
		ExpectedCount: 1,
		Nodes: []Node{
			{
				Host: "etcd01",
				IP:   "192.168.205.10",
			},
		},
	},
	Master: MasterNodeGroup{
		ExpectedCount: 1,
		Nodes: []Node{
			{
				Host: "master01",
				IP:   "192.168.205.11",
			},
		},
		LoadBalancedFQDN:      "test",
		LoadBalancedShortName: "test",
	},
	Worker: NodeGroup{
		ExpectedCount: 1,
		Nodes: []Node{
			{
				Host: "worker01",
				IP:   "192.168.205.12",
			},
		},
	},
	Ingress: OptionalNodeGroup{
		ExpectedCount: 1,
		Nodes: []Node{
			{
				Host: "etcd01",
				IP:   "192.168.205.10",
			},
		},
	},
	NFS: NFS{
		Volumes: []NFSVolume{
			{
				Host: "10.10.2.20",
				Path: "/",
			},
		},
	},
}

func assertInvalidPlan(t *testing.T, p Plan) {
	valid, _ := ValidatePlan(&p)
	if valid {
		t.Errorf("expected invalid, but got valid")
	}
}

func TestValidateBlankPlan(t *testing.T) {
	p := Plan{}
	assertInvalidPlan(t, p)
}

func TestValidateValidPlan(t *testing.T) {
	p := validPlan
	valid, errs := ValidatePlan(&p)
	if !valid {
		t.Errorf("expected valid, but got invalid")
	}
	fmt.Println(errs)
}

func TestValidatePlanEmptyPodCIDR(t *testing.T) {
	p := validPlan
	p.Cluster.Networking.PodCIDRBlock = ""
	assertInvalidPlan(t, p)
}

func TestValidatePlanInvalidPodCIDR(t *testing.T) {
	p := validPlan
	p.Cluster.Networking.PodCIDRBlock = "foo"
	assertInvalidPlan(t, p)
}

func TestValidatePlanEmptyServicesCIDR(t *testing.T) {
	p := validPlan
	p.Cluster.Networking.ServiceCIDRBlock = ""
	assertInvalidPlan(t, p)
}

func TestValidatePlanInvalidServicesCIDR(t *testing.T) {
	p := validPlan
	p.Cluster.Networking.ServiceCIDRBlock = "foo"
	assertInvalidPlan(t, p)
}

func TestValidatePlanEmptyPassword(t *testing.T) {
	p := validPlan
	p.Cluster.AdminPassword = ""
	assertInvalidPlan(t, p)
}

func TestValidatePlanEmptyCertificatesExpiry(t *testing.T) {
	p := validPlan
	p.Cluster.Certificates.Expiry = ""
	assertInvalidPlan(t, p)
}

func TestValidatePlanInvalidCertExpiry(t *testing.T) {
	p := validPlan
	p.Cluster.Certificates.Expiry = "foo"
	assertInvalidPlan(t, p)
}

func TestValidatePlanEmptyCACertExpiryIsValid(t *testing.T) {
	p := validPlan
	p.Cluster.Certificates.CAExpiry = ""
	valid, _ := p.validate()
	if !valid {
		t.Errorf("plan was found invalid")
	}
}

func TestValidatePlanInvalidCACertificatesExpiry(t *testing.T) {
	p := validPlan
	p.Cluster.Certificates.CAExpiry = "foo"
	assertInvalidPlan(t, p)
}

func TestValidatePlanEmptySSHUser(t *testing.T) {
	p := validPlan
	p.Cluster.SSH.User = ""
	assertInvalidPlan(t, p)
}

func TestValidatePlanEmptySSHKey(t *testing.T) {
	p := validPlan
	p.Cluster.SSH.Key = ""
	assertInvalidPlan(t, p)
}

func TestValidatePlanNonExistentSSHKey(t *testing.T) {
	p := validPlan
	p.Cluster.SSH.Key = "/foo"
	assertInvalidPlan(t, p)
}

func TestValidatePlanNegativeSSHPort(t *testing.T) {
	p := validPlan
	p.Cluster.SSH.Port = -1
	assertInvalidPlan(t, p)
}

func TestValidatePlanEmptyLoadBalancedFQDN(t *testing.T) {
	p := validPlan
	p.Master.LoadBalancedFQDN = ""
	assertInvalidPlan(t, p)
}

func TestValidatePlanEmptyLoadBalancedShortName(t *testing.T) {
	p := validPlan
	p.Master.LoadBalancedShortName = ""
	assertInvalidPlan(t, p)
}

func TestValidatePlanNoEtcdNodes(t *testing.T) {
	p := validPlan
	p.Etcd.ExpectedCount = 0
	p.Etcd.Nodes = []Node{}
	assertInvalidPlan(t, p)
}

func TestValidatePlanNoMasterNodes(t *testing.T) {
	p := validPlan
	p.Master.ExpectedCount = 0
	p.Master.Nodes = []Node{}
	assertInvalidPlan(t, p)
}

func TestValidatePlanNoWorkerNodes(t *testing.T) {
	p := validPlan
	p.Worker.ExpectedCount = 0
	p.Worker.Nodes = []Node{}
	assertInvalidPlan(t, p)
}

func TestValidatePlanEtcdNodesMismatch(t *testing.T) {
	p := validPlan
	p.Etcd.ExpectedCount = 100
	assertInvalidPlan(t, p)
}

func TestValidatePlanMasterNodesMismatch(t *testing.T) {
	p := validPlan
	p.Master.ExpectedCount = 100
	assertInvalidPlan(t, p)
}

func TestValidatePlanWorkerNodesMismatch(t *testing.T) {
	p := validPlan
	p.Worker.ExpectedCount = 100
	assertInvalidPlan(t, p)
}

func TestValidatePlanUnexpectedEtcdNodes(t *testing.T) {
	p := validPlan
	p.Etcd.ExpectedCount = 1
	p.Etcd.Nodes = []Node{
		{
			Host: "etcd01",
			IP:   "192.168.205.10",
		},
		{
			Host: "etcd02",
			IP:   "192.168.205.11",
		},
	}
	assertInvalidPlan(t, p)
}

func TestValidatePlanUnexpectedMasterNodes(t *testing.T) {
	p := validPlan
	p.Master.ExpectedCount = 1
	p.Master.Nodes = []Node{
		{
			Host: "master01",
			IP:   "192.168.205.10",
		},
		{
			Host: "master02",
			IP:   "192.168.205.11",
		},
	}
	assertInvalidPlan(t, p)
}

func TestValidatePlanUnexpectedWorkerNodes(t *testing.T) {
	p := validPlan
	p.Worker.ExpectedCount = 1
	p.Worker.Nodes = []Node{
		{
			Host: "worker01",
			IP:   "192.168.205.10",
		},
		{
			Host: "worker02",
			IP:   "192.168.205.11",
		},
	}
	assertInvalidPlan(t, p)
}

func TestValidatePlanNoIngress(t *testing.T) {
	p := validPlan
	p.Ingress.ExpectedCount = 0
	p.Ingress.Nodes = []Node{}
	valid, _ := ValidatePlan(&p)
	if !valid {
		t.Errorf("expected valid, but got invalid")
	}
}

func TestValidatePlanIngressExpected(t *testing.T) {
	p := validPlan
	p.Ingress.ExpectedCount = 1
	p.Ingress.Nodes = []Node{}
	assertInvalidPlan(t, p)
}

func TestValidatePlanIngressProvidedNotExpected(t *testing.T) {
	p := validPlan
	p.Ingress.ExpectedCount = 0
	p.Ingress.Nodes = []Node{
		{
			Host: "ingress",
			IP:   "192.168.205.10",
		},
	}
	assertInvalidPlan(t, p)
}

func TestValidateStorageVolume(t *testing.T) {
	tests := []struct {
		sv    StorageVolume
		valid bool
	}{
		{
			sv: StorageVolume{
				Name:              "foo",
				SizeGB:            100,
				DistributionCount: 2,
				ReplicateCount:    2,
				ReclaimPolicy:     "Retain",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: true,
		},
		{
			sv: StorageVolume{
				Name:              "foo",
				SizeGB:            100,
				DistributionCount: 1,
				ReplicateCount:    1,
				ReclaimPolicy:     "Retain",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: true,
		},
		{
			sv: StorageVolume{
				Name:              "foo",
				SizeGB:            100,
				DistributionCount: 0,
				ReplicateCount:    1,
				ReclaimPolicy:     "Retain",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: false,
		},
		{
			sv: StorageVolume{
				Name:              "foo",
				SizeGB:            100,
				DistributionCount: 1,
				ReplicateCount:    0,
				ReclaimPolicy:     "Retain",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: false,
		},
		{
			sv: StorageVolume{
				Name:              "bad name with spaces",
				SizeGB:            100,
				DistributionCount: 2,
				ReplicateCount:    2,
				ReclaimPolicy:     "Retain",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: false,
		},
		{
			sv: StorageVolume{
				Name:              "bad:name2",
				SizeGB:            100,
				DistributionCount: 2,
				ReplicateCount:    2,
				ReclaimPolicy:     "Retain",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: false,
		},
		{
			sv: StorageVolume{
				Name:              "goodName",
				SizeGB:            0,
				DistributionCount: 2,
				ReplicateCount:    2,
				ReclaimPolicy:     "Retain",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: false,
		},
		{
			sv: StorageVolume{
				Name:              "goodName",
				SizeGB:            -1,
				DistributionCount: 2,
				ReplicateCount:    2,
				ReclaimPolicy:     "Retain",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: false,
		},
		{
			sv: StorageVolume{
				Name:              "goodName",
				SizeGB:            100,
				DistributionCount: -1,
				ReplicateCount:    2,
				ReclaimPolicy:     "Retain",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: false,
		},
		{
			sv: StorageVolume{
				Name:              "goodName",
				SizeGB:            100,
				DistributionCount: 2,
				ReplicateCount:    -1,
				ReclaimPolicy:     "Retain",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: false,
		},
		{
			sv:    StorageVolume{},
			valid: false,
		},
		{
			sv: StorageVolume{
				Name:              "goodName",
				SizeGB:            100,
				DistributionCount: 1,
				ReplicateCount:    1,
				ReclaimPolicy:     "",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: false,
		},
		{
			sv: StorageVolume{
				Name:              "goodName",
				SizeGB:            100,
				DistributionCount: 1,
				ReplicateCount:    1,
				ReclaimPolicy:     "foo",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: false,
		},
		{
			sv: StorageVolume{
				Name:              "goodName",
				SizeGB:            100,
				DistributionCount: 1,
				ReplicateCount:    1,
				ReclaimPolicy:     "Retain",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: true,
		},
		{
			sv: StorageVolume{
				Name:              "goodName",
				SizeGB:            100,
				DistributionCount: 1,
				ReplicateCount:    1,
				ReclaimPolicy:     "Recycle",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: true,
		},
		{
			sv: StorageVolume{
				Name:              "goodName",
				SizeGB:            100,
				DistributionCount: 1,
				ReplicateCount:    1,
				ReclaimPolicy:     "Delete",
				AccessModes:       []string{"ReadWriteMany"},
			},
			valid: true,
		},
		{
			sv: StorageVolume{
				Name:              "goodName",
				SizeGB:            100,
				DistributionCount: 1,
				ReplicateCount:    1,
				ReclaimPolicy:     "Delete",
				AccessModes:       []string{"ReadWriteOnce"},
			},
			valid: true,
		},
		{
			sv: StorageVolume{
				Name:              "goodName",
				SizeGB:            100,
				DistributionCount: 1,
				ReplicateCount:    1,
				ReclaimPolicy:     "Delete",
				AccessModes:       []string{"ReadOnlyMany"},
			},
			valid: true,
		},
		{
			sv: StorageVolume{
				Name:              "goodName",
				SizeGB:            100,
				DistributionCount: 1,
				ReplicateCount:    1,
				ReclaimPolicy:     "Delete",
				AccessModes:       []string{"ReadWriteMany", "ReadWriteOnce", "ReadOnlyMany"},
			},
			valid: true,
		},
		{
			sv: StorageVolume{
				Name:              "goodName",
				SizeGB:            100,
				DistributionCount: 1,
				ReplicateCount:    1,
				ReclaimPolicy:     "Delete",
				AccessModes:       []string{"someBadAccessMode"},
			},
			valid: false,
		},
	}
	for _, test := range tests {
		if valid, _ := test.sv.validate(); valid != test.valid {
			t.Errorf("expected %v with %+v, but got %v", test.valid, test.sv, !test.valid)
		}
	}
}

func TestValidateAllowAddress(t *testing.T) {
	tests := []struct {
		address string
		valid   bool
	}{
		{"192.168.205.10", true},
		{"192.168.205.*", true},
		{"192.168.*.*", true},
		{"192.*.*.*", true},
		{"*.168.205.10", true},
		{"*.*.205.10", true},
		{"*.*.*.10", true},
		{"*.*.*.*", true},
		{"-1.-1.-1.-1", false},
		{"-1.*.*.*", false},
		{"*.-1.*.*", false},
		{"*.*.-1.*", false},
		{"*.*.*.-1", false},
		{"256.256.256.256", false},
		{"256.*.*.*", false},
		{"*.256.*.*", false},
		{"*.*.256.*", false},
		{"*.*.*.256", false},
		{"a.a.a.a", false},
		{"*.*.*.a", false},
		{"*.*.a.*", false},
		{"*.a.*.*", false},
		{"a.*.*.*", false},
		{"", false},
		{"foo", false},
		{"192", false},
		{"192.168", false},
		{"192.168.205", false},
		{"...", false},
		{"192...", false},
		{"192.168..", false},
		{"192.168.205.", false},
	}
	for _, test := range tests {
		if validateAllowedAddress(test.address) != test.valid {
			t.Errorf("expected %v with address %q, but got %v", test.valid, test.address, !test.valid)
		}
	}
}

func TestValidatePlanNFSDupes(t *testing.T) {
	p := validPlan

	p.NFS.Volumes = append(p.NFS.Volumes, NFSVolume{
		Host: "10.10.2.20",
		Path: "/",
	})

	assertInvalidPlan(t, p)
}

func TestValidateNFSVolume(t *testing.T) {
	tests := []struct {
		host  string
		path  string
		valid bool
	}{
		{
			host:  "10.10.2.10",
			path:  "/foo",
			valid: true,
		},
		{
			host:  "10.10.2.10",
			path:  "",
			valid: false,
		},
		{
			host:  "10.10.2.10",
			path:  "../someRelativePath",
			valid: false,
		},
		{
			host:  "",
			path:  "/foo",
			valid: false,
		},
	}
	for _, test := range tests {
		v := NFSVolume{
			Host: test.host,
			Path: test.path,
		}
		if valid, _ := v.validate(); valid != test.valid {
			t.Errorf("Expected valid = %v, but got %v", test.valid, valid)
		}
	}
}

func TestValidatePlanCerts(t *testing.T) {
	p := &validPlan

	pki := getPKI(t)
	defer cleanup(pki.GeneratedCertsDirectory, t)
	ca, err := pki.GenerateClusterCA(p)
	if err != nil {
		t.Fatalf("error generating CA for test: %v", err)
	}
	if err := pki.GenerateClusterCertificates(p, ca); err != nil {
		t.Fatalf("failed to generate certs: %v", err)
	}

	valid, errs := ValidateCertificates(p, &pki)
	if !valid {
		t.Errorf("expected valid, but got invalid")
		fmt.Println(errs)
	}
}

func TestValidatePlanBadCerts(t *testing.T) {
	p := &validPlan

	pki := getPKI(t)
	defer cleanup(pki.GeneratedCertsDirectory, t)

	ca, err := pki.GenerateClusterCA(p)
	if err != nil {
		t.Fatalf("error generating CA for test: %v", err)
	}
	if err := pki.GenerateClusterCertificates(p, ca); err != nil {
		t.Fatalf("failed to generate certs: %v", err)
	}
	p.Master.Nodes[0] = Node{
		Host:       "master01",
		IP:         "11.12.13.14",
		InternalIP: "22.33.44.55",
	}

	valid, _ := ValidateCertificates(p, &pki)
	if valid {
		t.Errorf("expected an error, but got valid")
	}
}

func TestValidatePlanMissingCerts(t *testing.T) {
	p := validPlan

	pki := getPKI(t)
	defer cleanup(pki.GeneratedCertsDirectory, t)

	valid, errs := ValidateCertificates(&p, &pki)
	if !valid {
		t.Errorf("expected valid, but got invalid")
		fmt.Println(errs)
	}
}

func TestValidatePlanMissingSomeCerts(t *testing.T) {
	p := &validPlan

	pki := getPKI(t)
	defer cleanup(pki.GeneratedCertsDirectory, t)

	ca, err := pki.GenerateClusterCA(p)
	if err != nil {
		t.Fatalf("error generating CA for test: %v", err)
	}
	if err := pki.GenerateClusterCertificates(p, ca); err != nil {
		t.Fatalf("failed to generate certs: %v", err)
	}

	newNode := Node{
		Host:       "master2",
		IP:         "11.12.13.14",
		InternalIP: "22.33.44.55",
	}
	p.Master.Nodes = append(p.Master.Nodes, newNode)

	valid, errs := ValidateCertificates(p, &pki)
	if !valid {
		t.Errorf("expected valid, but got invalid")
		fmt.Println(errs)
	}
}

func TestValidateNodeGroupDuplicateIP(t *testing.T) {
	ng := NodeGroup{
		ExpectedCount: 2,
		Nodes: []Node{
			{
				Host: "host1",
				IP:   "10.0.0.1",
			},
			{
				Host: "host2",
				IP:   "10.0.0.1",
			},
		},
	}
	if ok, _ := ng.validate(); ok {
		t.Errorf("validation passed with duplicate IP")
	}
}

func TestValidateNodeGroupDuplicateHostname(t *testing.T) {
	ng := NodeGroup{
		ExpectedCount: 2,
		Nodes: []Node{
			{
				Host: "host1",
				IP:   "10.0.0.1",
			},
			{
				Host: "host1",
				IP:   "10.0.0.2",
			},
		},
	}
	if ok, _ := ng.validate(); ok {
		t.Errorf("validation passed with duplicate hostname")
	}
}

func TestValidateNodeGroupDuplicateInternalIPs(t *testing.T) {
	ng := NodeGroup{
		ExpectedCount: 2,
		Nodes: []Node{
			{
				Host:       "host1",
				IP:         "10.0.0.1",
				InternalIP: "192.168.205.10",
			},
			{
				Host:       "host2",
				IP:         "10.0.0.2",
				InternalIP: "192.168.205.10",
			},
		},
	}
	if ok, _ := ng.validate(); ok {
		t.Errorf("validation passed with duplicate hostname")
	}
}

func TestDisconnectedInstallationPrereq(t *testing.T) {
	tests := []struct {
		cluster  Cluster
		registry DockerRegistry
		valid    bool
	}{
		{
			cluster: Cluster{
				DisconnectedInstallation: true,
			},
			valid: false,
		},
		{
			cluster: Cluster{
				DisconnectedInstallation: true,
			},
			registry: DockerRegistry{
				SetupInternal: true,
				Address:       "",
			},
			valid: true,
		},
		{
			cluster: Cluster{
				DisconnectedInstallation: true,
			},
			registry: DockerRegistry{
				SetupInternal: false,
				Address:       "10.0.0.1",
			},
			valid: true,
		},
		{
			registry: DockerRegistry{
				SetupInternal: true,
			},
			valid: true,
		},
		{
			registry: DockerRegistry{
				Address: "10.0.0.1",
			},
			valid: true,
		},
		{
			registry: DockerRegistry{},
			valid:    true,
		},
		{
			cluster: Cluster{
				DisconnectedInstallation: true,
				DisableRegistrySeeding:   true,
			},
			registry: DockerRegistry{
				SetupInternal: true,
				Address:       "",
			},
			valid: false,
		},
		{
			cluster: Cluster{
				DisconnectedInstallation: true,
				DisableRegistrySeeding:   true,
			},
			registry: DockerRegistry{
				SetupInternal: false,
				Address:       "10.0.0.1",
			},
			valid: true,
		},
		{
			cluster: Cluster{
				DisconnectedInstallation: false,
				DisableRegistrySeeding:   true,
			},
			registry: DockerRegistry{
				SetupInternal: true,
				Address:       "",
			},
			valid: true,
		},
	}
	for i, test := range tests {
		ok, _ := disconnectedInstallation{cluster: test.cluster, registry: test.registry}.validate()
		if ok != test.valid {
			t.Errorf("test %d: expect %t, but got %t", i, test.valid, ok)
		}
	}
}

func TestValidateDockerStorageDirectLVM(t *testing.T) {
	tests := []struct {
		config DockerStorageDirectLVM
		valid  bool
	}{
		{
			config: DockerStorageDirectLVM{
				Enabled: false,
			},
			valid: true,
		},
		{
			config: DockerStorageDirectLVM{
				Enabled: true,
			},
			valid: false,
		},
		{
			config: DockerStorageDirectLVM{
				Enabled:     true,
				BlockDevice: "foo",
			},
			valid: false,
		},
		{
			config: DockerStorageDirectLVM{
				Enabled:     true,
				BlockDevice: "/dev/sdb",
			},
			valid: true,
		},
	}
	for i, test := range tests {
		ok, _ := test.config.validate()
		if ok != test.valid {
			t.Errorf("test %d: expect valid, but got invalid", i)
		}
	}
}

func TestRepository(t *testing.T) {
	tests := []struct {
		config Cluster
		valid  bool
	}{
		{
			config: Cluster{
				PackageRepoURLs: "",
			},
			valid: true,
		},
		{
			config: Cluster{
				PackageRepoURLs: "https://repo.com",
			},
			valid: true,
		},
		{
			config: Cluster{
				PackageRepoURLs: "http://repo.com",
			},
			valid: true,
		},
		{
			config: Cluster{
				PackageRepoURLs: "http://192.168.0.1",
			},
			valid: true,
		},
		{
			config: Cluster{
				PackageRepoURLs: "https://repo.com/kismatic,https://repo.com/gluster",
			},
			valid: true,
		},
		{
			config: Cluster{
				PackageRepoURLs: "repo.com",
			},
			valid: false,
		},
		{
			config: Cluster{
				PackageRepoURLs: "https://repo.com/kismatic,repo.com/gluster",
			},
			valid: false,
		},
	}
	for i, test := range tests {
		p := &validPlan
		p.Cluster.PackageRepoURLs = test.config.PackageRepoURLs
		ok, _ := p.Cluster.validate()
		if ok != test.valid {
			t.Errorf("test %d: expect %t, but got %t", i, test.valid, ok)
		}
	}
}

func TestCNIAddOn(t *testing.T) {
	tests := []struct {
		n     CNI
		valid bool
	}{
		{
			n: CNI{
				Provider: "calico",
				Options: CNIOptions{
					Calico: CalicoOptions{
						Mode: "overlay",
					},
				},
			},
			valid: true,
		},
		{
			n: CNI{
				Provider: "calico",
				Options: CNIOptions{
					Calico: CalicoOptions{
						Mode: "routed",
					},
				},
			},
			valid: true,
		},
		{
			n: CNI{
				Provider: "weave",
			},
			valid: true,
		},
		{
			n: CNI{
				Provider: "contiv",
			},
			valid: true,
		},
		{
			n: CNI{
				Provider: "foo",
				Options: CNIOptions{
					Calico: CalicoOptions{
						Mode: "overlay",
					},
				},
			},
			valid: false,
		},
		{
			n: CNI{
				Provider: "calico",
				Options: CNIOptions{
					Calico: CalicoOptions{
						Mode: "foo",
					},
				},
			},
			valid: false,
		},
		{
			n: CNI{
				Provider: "foo",
				Disable:  true,
				Options: CNIOptions{
					Calico: CalicoOptions{
						Mode: "overlay",
					},
				},
			},
			valid: true,
		},
		{
			n: CNI{
				Provider: "calico",
				Disable:  true,
				Options: CNIOptions{
					Calico: CalicoOptions{
						Mode: "foo",
					},
				},
			},
			valid: true,
		},
	}
	for i, test := range tests {
		ok, _ := test.n.validate()
		if ok != test.valid {
			t.Errorf("test %d: expect %t, but got %t", i, test.valid, ok)
		}
	}
}

func TestHeapsterAddOn(t *testing.T) {
	tests := []struct {
		h     HeapsterMonitoring
		valid bool
	}{
		{
			h: HeapsterMonitoring{
				Options: HeapsterOptions{
					Heapster: Heapster{
						Replicas:    0,
						ServiceType: "ClusterIP",
					},
				},
			},
			valid: false,
		},
		{
			h: HeapsterMonitoring{
				Options: HeapsterOptions{
					Heapster: Heapster{
						Replicas:    1,
						ServiceType: "Foo",
					},
				},
			},
			valid: false,
		},
		{
			h: HeapsterMonitoring{
				Options: HeapsterOptions{
					Heapster: Heapster{
						Replicas:    -1,
						ServiceType: "ClusterIP",
					},
				},
			},
			valid: false,
		},
		{
			h: HeapsterMonitoring{
				Options: HeapsterOptions{
					Heapster: Heapster{
						Replicas:    1,
						ServiceType: "ClusterIP",
					},
				},
			},
			valid: true,
		},
	}
	for i, test := range tests {
		ok, _ := test.h.validate()
		if ok != test.valid {
			t.Errorf("test %d: expect %t, but got %t", i, test.valid, ok)
		}
	}
}

func TestPackageManagerAddOn(t *testing.T) {
	tests := []struct {
		p     PackageManager
		valid bool
	}{
		{
			p: PackageManager{
				Disable:  false,
				Provider: "helm",
			},
			valid: true,
		},
		{
			p: PackageManager{
				Disable:  true,
				Provider: "",
			},
			valid: true,
		},
		{
			p: PackageManager{
				Disable:  true,
				Provider: "foo",
			},
			valid: true,
		},
		{
			p: PackageManager{
				Disable:  false,
				Provider: "",
			},
			valid: true,
		},
		{
			p: PackageManager{
				Disable:  false,
				Provider: "foo",
			},
			valid: false,
		},
	}
	for i, test := range tests {
		ok, _ := test.p.validate()
		if ok != test.valid {
			t.Errorf("test %d: expect %t, but got %t", i, test.valid, ok)
		}
	}
}

func TestCloudProvider(t *testing.T) {
	tests := []struct {
		c     CloudProvider
		valid bool
	}{
		{
			c: CloudProvider{
				Provider: "",
			},
			valid: true,
		},
		{
			c: CloudProvider{
				Provider: "aws",
			},
			valid: true,
		},
		{
			c: CloudProvider{
				Provider: "awss",
			},
			valid: false,
		},
		{
			c: CloudProvider{
				Provider: "gce",
				Config:   "/bin/sh",
			},
			valid: true,
		},
		{
			c: CloudProvider{
				Provider: "gce",
				Config:   "/bin/foo",
			},
			valid: false,
		},
		{
			c: CloudProvider{
				Provider: "gce",
				Config:   "foo",
			},
			valid: false,
		},
	}
	for i, test := range tests {
		ok, _ := test.c.validate()
		if ok != test.valid {
			t.Errorf("test %d: expect %t, but got %t", i, test.valid, ok)
		}
	}
}
