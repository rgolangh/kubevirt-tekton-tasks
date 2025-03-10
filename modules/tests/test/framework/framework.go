package framework

import (
	"context"
	"fmt"

	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/constants"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework/clients"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/framework/testoptions"
	"github.com/kubevirt/kubevirt-tekton-tasks/modules/tests/test/tekton"
	. "github.com/onsi/ginkgo/v2"
	templatev1 "github.com/openshift/api/template/v1"
	pipev1beta1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/api/core/v1"
	cdiv1beta1 "kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1"
)

var TestOptionsInstance = &testoptions.TestOptions{}
var ClientsInstance = &clients.Clients{}

type ManagedResources struct {
	taskRuns     []*pipev1beta1.TaskRun
	pipelineRuns []*pipev1beta1.PipelineRun
	pipelines    []*pipev1beta1.Pipeline
	dataVolumes  []*cdiv1beta1.DataVolume
	dataSources  []*cdiv1beta1.DataSource
	vms          []*kubevirtv1.VirtualMachine
	templates    []*templatev1.Template
	secrets      []*corev1.Secret
}
type Framework struct {
	*testoptions.TestOptions
	*clients.Clients

	managedResources  ManagedResources
	limitEnvScope     constants.EnvScope
	onBeforeTestSetup func(config TestConfig)
}

type TestConfig interface {
	GetLimitTestScope() constants.TestScope
	GetLimitEnvScope() constants.EnvScope
	Init(options *testoptions.TestOptions)
}

func NewFramework() *Framework {
	f := &Framework{
		TestOptions: TestOptionsInstance,
		Clients:     ClientsInstance,
	}

	AfterEach(f.AfterEach)
	return f
}

func (f *Framework) LimitEnvScope(limitEnvScope constants.EnvScope) *Framework {
	if f.limitEnvScope != "" {
		Fail("limitEnvScope was already set")
	}
	f.limitEnvScope = limitEnvScope

	return f
}

func (f *Framework) OnBeforeTestSetup(callback func(config TestConfig)) *Framework {
	f.onBeforeTestSetup = callback
	return f
}

func (f *Framework) TestSetup(config TestConfig) {
	limitScope := config.GetLimitTestScope()
	limitEnvScope := config.GetLimitEnvScope()

	// check global env limit first
	if f.limitEnvScope != "" && f.limitEnvScope != f.EnvScope {
		Skip(fmt.Sprintf("runs only in %v", f.limitEnvScope))
	}

	// check test case env limit
	if limitEnvScope != "" && limitEnvScope != f.EnvScope {
		Skip(fmt.Sprintf("runs only in %v", limitEnvScope))
	}

	// check test case test scope limit
	if limitScope != "" && limitScope != f.TestScope {
		Skip(fmt.Sprintf("runs only in %v scope", limitScope))
	}
	if f.onBeforeTestSetup != nil {
		f.onBeforeTestSetup(config)
	}
	config.Init(f.TestOptions)
}

func (f *Framework) AfterEach() {
	failed := CurrentSpecReport().Failed()
	taskRuns := f.managedResources.taskRuns
	pipelineRuns := f.managedResources.pipelineRuns

	if failed {
		defer func() {
			if !f.Debug {
				for _, taskRun := range taskRuns {
					defer f.TknClient.TaskRuns(taskRun.Namespace).Delete(context.TODO(), taskRun.Name, metav1.DeleteOptions{})
				}
				for _, pipelineRun := range pipelineRuns {
					defer f.TknClient.PipelineRuns(pipelineRun.Namespace).Delete(context.TODO(), pipelineRun.Name, metav1.DeleteOptions{})
				}
			}
			for _, taskRun := range taskRuns {
				tekton.PrintTaskRunDebugInfo(f.Clients, taskRun.Namespace, taskRun.Name)
			}
			for _, pipelineRun := range pipelineRuns {
				tekton.PrintPipelineRunDebugInfo(f.Clients, pipelineRun.Namespace, pipelineRun.Name)
			}
		}()
	}

	if f.Debug {
		// leave resources alive for inspection
		return
	}

	if !failed { // failed has its own cleanup
		for _, taskRun := range taskRuns {
			defer f.TknClient.TaskRuns(taskRun.Namespace).Delete(context.TODO(), taskRun.Name, metav1.DeleteOptions{})
		}
		for _, pipelineRun := range pipelineRuns {
			defer f.TknClient.PipelineRuns(pipelineRun.Namespace).Delete(context.TODO(), pipelineRun.Name, metav1.DeleteOptions{})
		}
	}
	for _, pipeline := range f.managedResources.pipelines {
		defer f.TknClient.Pipelines(pipeline.Namespace).Delete(context.TODO(), pipeline.Name, metav1.DeleteOptions{})
	}
	for _, dv := range f.managedResources.dataVolumes {
		defer f.CdiClient.DataVolumes(dv.Namespace).Delete(context.TODO(), dv.Name, metav1.DeleteOptions{})
	}
	for _, ds := range f.managedResources.dataSources {
		defer f.CdiClient.DataSources(ds.Namespace).Delete(context.TODO(), ds.Name, metav1.DeleteOptions{})
	}
	for _, vm := range f.managedResources.vms {
		defer f.KubevirtClient.VirtualMachine(vm.Namespace).Delete(vm.Name, &metav1.DeleteOptions{})
	}
	for _, t := range f.managedResources.templates {
		defer f.TemplateClient.Templates(t.Namespace).Delete(context.TODO(), t.Name, metav1.DeleteOptions{})
	}
	for _, s := range f.managedResources.secrets {
		defer f.KubevirtClient.CoreV1().Secrets(s.Namespace).Delete(context.TODO(), s.Name, metav1.DeleteOptions{})
	}
}

func (f *Framework) ManageTaskRuns(taskRuns ...*pipev1beta1.TaskRun) *Framework {
	f.managedResources.taskRuns = append(f.managedResources.taskRuns, taskRuns...)
	return f
}

func (f *Framework) ManagePipelineRuns(pipelineRuns ...*pipev1beta1.PipelineRun) *Framework {
	f.managedResources.pipelineRuns = append(f.managedResources.pipelineRuns, pipelineRuns...)
	return f
}

func (f *Framework) ManagePipelines(pipelines ...*pipev1beta1.Pipeline) *Framework {
	f.managedResources.pipelines = append(f.managedResources.pipelines, pipelines...)
	return f
}

func (f *Framework) ManageDataVolumes(dataVolumes ...*cdiv1beta1.DataVolume) *Framework {
	for _, dataVolume := range dataVolumes {
		if dataVolume != nil && dataVolume.Name != "" && dataVolume.Namespace != "" {
			f.managedResources.dataVolumes = append(f.managedResources.dataVolumes, dataVolume)
		}
	}
	return f
}

func (f *Framework) ManageDataSources(dataSources ...*cdiv1beta1.DataSource) *Framework {
	for _, dataSource := range dataSources {
		if dataSource != nil && dataSource.Name != "" && dataSource.Namespace != "" {
			f.managedResources.dataSources = append(f.managedResources.dataSources, dataSource)
		}
	}
	return f
}

func (f *Framework) ManageVMs(vms ...*kubevirtv1.VirtualMachine) *Framework {
	for _, vm := range vms {
		if vm != nil && vm.Name != "" && vm.Namespace != "" {
			f.managedResources.vms = append(f.managedResources.vms, vm)
		}
	}
	return f
}

func (f *Framework) ManageTemplates(templates ...*templatev1.Template) *Framework {
	for _, t := range templates {
		if t != nil && t.Name != "" && t.Namespace != "" {
			f.managedResources.templates = append(f.managedResources.templates, t)
		}
	}
	return f
}

func (f *Framework) ManageSecrets(secrets ...*corev1.Secret) *Framework {
	for _, s := range secrets {
		if s != nil && s.Name != "" && s.Namespace != "" {
			f.managedResources.secrets = append(f.managedResources.secrets, s)
		}
	}
	return f
}
