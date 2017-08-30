package install

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestWritePlanTemplate(t *testing.T) {
	plan := &Plan{
		Cluster: Cluster{
			AdminPassword: "password",
		},
		Etcd: NodeGroup{
			ExpectedCount: 3,
		},
		Master: MasterNodeGroup{
			ExpectedCount: 2,
		},
		Worker: NodeGroup{
			ExpectedCount: 3,
		},
		Ingress: OptionalNodeGroup{
			ExpectedCount: 2,
		},
		Storage: OptionalNodeGroup{
			ExpectedCount: 2,
		},
		NFS: NFS{
			Volumes: []NFSVolume{
				NFSVolume{Host: "", Path: "/"},
				NFSVolume{Host: "", Path: "/"},
			},
		},
	}

	tmpDir, err := ioutil.TempDir("", "test-write-plan-template")
	if err != nil {
		t.Fatalf("error creating tmp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filename := filepath.Join(tmpDir, "kismatic-cluster.yaml")
	planner := &FilePlanner{filename}
	if err := WritePlanTemplate(plan, planner); err != nil {
		t.Fatalf("error writing plan file template: %v", err)
	}

	goldenFile := "./test/plan-template-with-storage.golden.yaml"
	expected, err := ioutil.ReadFile(goldenFile)
	if err != nil {
		t.Fatalf("error reading golden file: %v", err)
	}

	got, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("error reading generated plan file template: %v", err)
	}

	if !bytes.Equal(expected, got) {
		t.Error("generated plan file template does not match expected")
		if _, err := exec.LookPath("diff"); err == nil {
			cmd := exec.Command("diff", goldenFile, filename)
			cmd.Stdout = os.Stdout
			cmd.Run()
		}
	}
}

func TestGenerateAlphaNumericPassword(t *testing.T) {
	_, err := generateAlphaNumericPassword()
	if err != nil {
		t.Error(err)
	}
}

func TestReadWithDeprecated(t *testing.T) {
	pm := &DeprecatedPackageManager{
		Enabled: true,
	}
	p := &Plan{}
	p.Features = &Features{
		PackageManager: pm,
	}
	b := false
	p.Cluster.AllowPackageInstallation = &b
	readDeprecatedFields(p)

	// features.package_manager should be set to add_ons.package_manager
	if p.AddOns.PackageManager.Disable || p.AddOns.PackageManager.Provider != "helm" {
		t.Errorf("Expected add_ons.package_manager to be read from features.package_manager")
	}
	// cluster.disable_package_installation shoule be set to cluster.allow_package_installation
	if p.Cluster.DisablePackageInstallation != true {
		t.Errorf("Expected cluster.allow_package_installation to be read from cluster.disable_package_installation")
	}
}

func TestReadWithNil(t *testing.T) {
	p := &Plan{}
	setDefaults(p)

	if p.AddOns.CNI.Provider != "calico" {
		t.Errorf("Expected add_ons.cni.provider to equal 'calico', instead got %s", p.AddOns.CNI.Provider)
	}
	if p.AddOns.CNI.Options.Calico.Mode != "overlay" {
		t.Errorf("Expected add_ons.cni.options.calico.mode to equal 'overlay', instead got %s", p.AddOns.CNI.Options.Calico.Mode)
	}

	if p.AddOns.HeapsterMonitoring.Options.Heapster.Replicas != 2 {
		t.Errorf("Expected add_ons.heapster.options.heapster.replicas to equal 2, instead got %d", p.AddOns.HeapsterMonitoring.Options.Heapster.Replicas)
	}

	if p.AddOns.HeapsterMonitoring.Options.Heapster.ServiceType != "ClusterIP" {
		t.Errorf("Expected add_ons.heapster.options.heapster.service_type to equal ClusterIP, instead got %s", p.AddOns.HeapsterMonitoring.Options.Heapster.ServiceType)
	}

	if p.AddOns.HeapsterMonitoring.Options.Heapster.Sink != "influxdb:http://heapster-influxdb.kube-system.svc:8086" {
		t.Errorf("Expected add_ons.heapster.options.heapster.service_type to equal 'influxdb:http://heapster-influxdb.kube-system.svc:8086', instead got %s", p.AddOns.HeapsterMonitoring.Options.Heapster.Sink)
	}

	if p.Cluster.Certificates.CAExpiry != defaultCAExpiry {
		t.Errorf("expected ca cert expiry to be %s, but got %s", defaultCAExpiry, p.Cluster.Certificates.CAExpiry)
	}
}

func TestReadDeprecatedDashboard(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test-read-deprecated-dashboard")
	if err != nil {
		t.Fatalf("error creating tmp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	file := filepath.Join(tmpDir, "kismatic-cluster.yaml")

	tests := []struct {
		name           string
		planStr        string
		expectDisabled bool
	}{
		{
			name:           "deprecated is set to true",
			planStr:        `{'add_ons': {'dashbard': {'disable': true}}}`,
			expectDisabled: true,
		},
		{
			name:           "deprecated is set to false",
			planStr:        `{'add_ons': {'dashbard': {'disable': false}}}`,
			expectDisabled: false,
		},
		{
			name:           "actual field is set to true",
			planStr:        `{'add_ons': {'dashboard': {'disable': true}}}`,
			expectDisabled: true,
		},
		{
			name:           "actual field is set to false",
			planStr:        `{'add_ons': {'dashboard': {'disable': false}}}`,
			expectDisabled: false,
		},
		{
			name:           "both fields are set to true",
			planStr:        `{'add_ons': {'dashboard': {'disable': true}, 'dashbard': {'disable': true}}}`,
			expectDisabled: true,
		},
		{
			name:           "both fields are set to false",
			planStr:        `{'add_ons': {'dashboard': {'disable': false}, 'dashbard': {'disable': false}}}`,
			expectDisabled: false,
		},
		{
			name:           "both are missing",
			planStr:        "",
			expectDisabled: false,
		},
		{
			name:           "deprecated is set to false, new one is set to true",
			planStr:        `{'add_ons': {'dashbard': {'disable': false}, 'dashboard': {'disable': true}}}`,
			expectDisabled: true,
		},
		{
			name:           "deprecated is set to true, new one is set to false",
			planStr:        `{'add_ons': {'dashbard': {'disable': true}, 'dashboard': {'disable': false}}}`,
			expectDisabled: false,
		},
	}

	for _, test := range tests {
		// writeFile truncates before writing, so we can reuse the file
		if err = ioutil.WriteFile(file, []byte(test.planStr), 0666); err != nil {
			t.Fatalf("error writing plan file")
		}

		planner := FilePlanner{file}
		plan, err := planner.Read()
		if err != nil {
			t.Fatalf("error reading plan file")
		}

		if plan.AddOns.Dashboard.Disable != test.expectDisabled {
			t.Errorf("name: %s: expected disabled to be %v, but got %v.", test.name, test.expectDisabled, plan.AddOns.Dashboard.Disable)
		}
	}

}
