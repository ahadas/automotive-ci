An opinionated operator that facilitates using OpenShift for Automotive-related development.

## Introduction

The purpose of this operator is to facilitate the configuration of development tasks for code repositories that are used by the Automotive SIG on public clouds, on-premise environments, or hybrid clouds, using the OpenShift platform.

Development tasks that were identified so far are:
1. Nightly builds: build container images, and optionally other artifacts, with their latest third-party dependencies every day and publish the artifacts publicly.
2. CI/CD: integrate the build-test-publish flow in Automotive-related repositories (on various platforms, like GitLab and GitHub).
3. Facilitate the creation of custom artifacts with pending changes on different platforms.

## Design

For the sake of simplicity, we implemented this with an Ansible based operator.

The properties for each code repository are defined in a custom resource [ContainerBuild](config/crd/bases/automotive.sig.centos.org_containerbuilds.yaml). The operator acts on ContainerBuild entities and create a pipeline and other entities for the code repository that is specified in the ContainerBuild instance, in the namespace that the ContainerBuild instance resides in. By removing the ContainerBuild instance, all the derived entities would be removed from the OpenShift cluster.

The definition of the aforementioned derived entities are stored in the form of [templates of the 'containerbuild' role](roles/containerbuild/templates) and are applied by [the 'main' task for the 'containerbuild' role](roles/containerbuild/tasks/main.yml). These are the resources that are most likely to change. Another resource that is likely to change is [the Role-based access control (RBAC) rules for the operator](config/rbac/role.yaml).

Examples of ContainerBuild definitions can be found in the [config/samples folder](config/samples). Further examples can be found in [the deployment folder](deployment) for ContainrBuild instances we defined for [the CentOS Automotive SIG Container Images projects](https://gitlab.com/CentOS/automotive/container-images).

## Deploying the Operator

Next we'll explain how to deploy the operator to an OpenShift cluster. The operator is not published to the Operator Hub and there is no plan to publish it there so we show how to do it from source, assuming this repository was cloned and you logged in to the OpenShift cluster already.

First, you need to choose where to publish the image of the operator and set the `IMG` environment variable accordingly, e.g., IMG=quay.io/ahadas/autoci, and then run the `docker-build` and `docker-push` targets in the Makefile to build and push it, respectively. Second, run the `deploy` target to deploy the operator to the OpenShift cluster.

## Configure a Source Code Repository

After deploying the operator, it waits for per-respository configuration to be posted to the OpenShift cluster. You look at the examples in the `deployment` folder for references of ContainerBuild configurations. We will demonstrate the process for the [automotive-osbuild repository](https://gitlab.com/CentOS/automotive/container-images/automotive-osbuild) using the [automotive-osbuild.yaml] configuration that is stored in the `deployment` folder. For further details about the supported configuration properties, see the default settings [here](roles/containerbuild/defaults/main.yml).

Deploy a ContainerBuild instance:
```bash
oc apply -f deployment/automotive-osbuild.yaml
```

This leads to having the following entities:
1. A pipeline for cloning the code repository and then build, test, and push a container image from this repository.
2. A cron job that triggers the abovementioned pipeline for nightly builds. The cron job is named after the ContainerBuild instance with a `-nightly` postfix. You can trigger this cron job for testing by creating a job from it:
```bash
oc -n autosd create job test --from=cronjob/automotive-osbuild-nightly
```
3. A route for a webhook for GitLab.


## More Resources
* [Ansible operators in the OpenShift documentation](https://docs.openshift.com/container-platform/4.16/operators/operator_sdk/ansible/osdk-ansible-tutorial.html)
