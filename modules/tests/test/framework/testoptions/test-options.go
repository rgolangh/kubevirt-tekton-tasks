package testoptions

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"k8s.io/client-go/util/homedir"
)

var deployNamespace string
var testNamespace string
var storageClass string
var kubeConfigPath string
var scope string
var debug string
var isOKD string
var skipCreateVMFromManifestTests string
var skipExecuteInVMTests string
var skipGenerateSSHKeysTests string

type TestOptions struct {
	DeployNamespace               string
	TestNamespace                 string
	StorageClass                  string
	KubeConfigPath                string
	TestScope                     constants.TestScope
	EnvScope                      constants.EnvScope
	Debug                         bool
	SkipCreateVMFromManifestTests bool
	SkipExecuteInVMTests          bool
	SkipGenerateSSHKeysTests      bool

	CommonTemplatesVersion string

	targetNamespaces map[constants.TargetNamespace]string
}

func init() {
	flag.StringVar(&deployNamespace, "deploy-namespace", "", "Namespace where to deploy the tasks and taskrun")
	flag.StringVar(&testNamespace, "test-namespace", "", "Namespace where to create the vm/dv resources")
	flag.StringVar(&storageClass, "storage-class", "", "Storage class to be used for creating test DVs/PVCs")
	flag.StringVar(&kubeConfigPath, "kubeconfig-path", "", "Path to the kubeconfig")
	flag.StringVar(&isOKD, "is-okd", "", "Set to true if running on OKD. One of: true|false")
	flag.StringVar(&scope, "scope", "", "Scope of the tests. One of: cluster|namespace")
	flag.StringVar(&debug, "debug", "", "Debug keeps all the resources alive after the tests complete. One of: true|false")
	flag.StringVar(&skipCreateVMFromManifestTests, "skip-create-vm-from-manifests-tests", "", "Skip create vm from manifests test suite. One of: true|false")
	flag.StringVar(&skipExecuteInVMTests, "skip-execute-in-vm-tests", "", "Skip execute in vm test suite. One of: true|false")
	flag.StringVar(&skipGenerateSSHKeysTests, "skip-generate-ssh-keys-tests", "", "Skip generate ssh keys suite. One of: true|false")
}

func InitTestOptions(testOptions *TestOptions) error {
	flag.Parse()

	if deployNamespace == "" {
		return errors.New("--deploy-namespace must be specified")
	}

	if testNamespace == "" {
		return errors.New("--test-namespace must be specified")
	}

	if scope == "" {
		testOptions.TestScope = constants.NamespaceTestScope
	} else if constants.TestScope(scope) == constants.NamespaceTestScope || constants.TestScope(scope) == constants.ClusterTestScope {
		testOptions.TestScope = constants.TestScope(scope)
	} else {
		return fmt.Errorf("invalid scope, only %v or %v is allowed", constants.ClusterTestScope, constants.NamespaceTestScope)
	}

	if kubeConfigPath != "" {
		testOptions.KubeConfigPath = kubeConfigPath
	} else {
		kubeConfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
		if file, err := os.Stat(kubeConfigPath); err == nil && file.Mode().IsRegular() {
			testOptions.KubeConfigPath = kubeConfigPath
		}
	}

	testOptions.DeployNamespace = deployNamespace
	testOptions.TestNamespace = testNamespace
	testOptions.StorageClass = storageClass
	if strings.ToLower(isOKD) == "true" {
		testOptions.EnvScope = constants.OKDEnvScope
	} else {
		testOptions.EnvScope = constants.KubernetesEnvScope
	}
	testOptions.Debug = strings.ToLower(debug) == "true"

	testOptions.targetNamespaces = testOptions.resolveNamespaces()

	testOptions.SkipCreateVMFromManifestTests = strings.ToLower(skipCreateVMFromManifestTests) == "true"
	testOptions.SkipExecuteInVMTests = strings.ToLower(skipExecuteInVMTests) == "true"
	testOptions.SkipGenerateSSHKeysTests = strings.ToLower(skipGenerateSSHKeysTests) == "true"

	return nil
}

func (f *TestOptions) resolveNamespaces() map[constants.TargetNamespace]string {
	var systemNS string
	if f.DeployNamespace == "tekton-pipelines" {
		systemNS = "default"
	} else {
		systemNS = "tekton-pipelines"
	}

	return map[constants.TargetNamespace]string{
		constants.DeployTargetNS: f.DeployNamespace,
		constants.TestTargetNS:   f.TestNamespace,
		constants.SystemTargetNS: systemNS,
		constants.EmptyTargetNS:  "",
	}
}

func (f *TestOptions) ResolveNamespace(namespace constants.TargetNamespace, fallbackNamespace string) string {
	if namespace == "" {
		if fallbackNamespace != "" {
			return fallbackNamespace
		}
	} else {
		if namespace == constants.EmptyTargetNS {
			return ""
		}

		ns := f.targetNamespaces[namespace]

		if ns != "" {
			return ns
		}
	}

	return f.TestNamespace
}
