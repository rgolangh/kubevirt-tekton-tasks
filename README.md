# KubeVirt Tekton Tasks

[Tekton Pipelines](https://github.com/tektoncd/pipeline) are CI/CD-style pipelines for Kubernetes.
This repository provides KubeVirt-specific Tekton tasks, which focus on:

- Creating and managing resources (VMs, DataVolumes, DataSources)
- Executing commands in VMs
- Manipulating disk images with libguestfs tools

## Deployment

### On Kubernetes

In order to install the KubeVirt Tekton tasks in the active namespace you need to apply the following manifest.
You have to repeat this for every namespace in which you'd like to run the tasks.

```bash
VERSION=$(curl -s https://api.github.com/repos/kubevirt/kubevirt-tekton-tasks/releases | \
            jq '.[] | select(.prerelease==false) | .tag_name' | sort -V | tail -n1 | tr -d '"')
kubectl apply -f "https://github.com/kubevirt/kubevirt-tekton-tasks/releases/download/${VERSION}/kubevirt-tekton-tasks-kubernetes.yaml"
```

Visit [RBAC permissions for running the tasks](docs/tasks-rbac-permissions.md) if the pipeline needs to create/access resources (VMs, PVCs, etc.) in a different namespace other than the one the pipeline runs in.

### On OKD

In order to install the KubeVirt Tekton tasks with additional OKD-specific tasks in the active namespace you need to apply the following manifest.
You have to repeat this for every namespace in which you'd like to run the tasks.

```bash
VERSION=$(curl -s https://api.github.com/repos/kubevirt/kubevirt-tekton-tasks/releases | \
            jq '.[] | select(.prerelease==false) | .tag_name' | sort -V | tail -n1 | tr -d '"')
kubectl apply -f "https://github.com/kubevirt/kubevirt-tekton-tasks/releases/download/${VERSION}/kubevirt-tekton-tasks-okd.yaml"
```

Visit [RBAC permissions for running the tasks](docs/tasks-rbac-permissions.md) if the pipeline needs to create/access resources (VMs, PVCs, etc.) in a different namespace other than the one the pipeline runs in.


## Usage

#### Create Virtual Machines

- [create-vm-from-manifest](tasks/create-vm-from-manifest)
- [create-vm-from-template](tasks/create-vm-from-template)

#### Utilize Templates

- [copy-template](tasks/copy-template)
- [modify-vm-template](tasks/modify-vm-template)

#### Modify data objects

- [modify-data-object](tasks/modify-data-object)

#### Generate SSH Keys

- [generate-ssh-keys](tasks/generate-ssh-keys)

#### Execute Commands in Virtual Machines

- [execute-in-vm](tasks/execute-in-vm): execute commands over SSH
- [cleanup-vm](tasks/cleanup-vm): execute commands and/or stop/delete VMs

#### Manipulate PVCs with libguestfs tools

- [disk-virt-customize](tasks/disk-virt-customize): execute virt-customize commands in PVCs
- [disk-virt-sysprep](tasks/disk-virt-sysprep): execute virt-sysprep commands in PVCs

#### Wait for Virtual Machine Instance Status

- [wait-for-vmi-status](tasks/wait-for-vmi-status)

## Examples

#### [Unit Tester Pipeline](examples/pipelines/unit-tester) 

Good unit tests are detached from the operating system and can run everywhere.
However, this is not always the case. Your tests may require access to entire operating system, or run as root,
or need a specific kernel.

This example shows how you can run your tests in your VM of choice.
The pipeline creates a VM, connects to it over SSH and runs tests inside it.
It also showcases the `finally` construct.


#### [Server Deployer Pipeline](examples/pipelines/server-deployer)

For complex application server deployments it might be easier to start the server as is in a VM rather than converting it to cloud-native application.

This example shows how you can initialize/modify a PVC and deploy such application in a VM.

#### [Virt-sysprep Updater Pipeline](examples/pipelines/virt-sysprep-updater)

Virt-sysprep can be used for preparing VM images which can be then used as base images for other VMs.

This example shows how you can update an operating system and seal VM's image by using virt-customize.
Then, a VM is created from such image.

#### [Windows 10 Installer Pipeline](examples/pipelines/windows10-installer)

Downloads a Windows Source ISO into a PVC and automatically installs Windows by using a custom Answer file into a new base DV.

#### [Windows 10 Customize Pipeline](examples/pipelines/windows10-customize)

Applies customizations to an existing Windows 10 installation by using a custom Answer file and creates a new base DV.

## Development Guide

See [Getting Started](docs/getting-started.md) for the environment setup and development workflow.
