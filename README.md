make docker-build IMG=quay.io/ahadas/autoci
make docker-push IMG=quay.io/ahadas/autoci
make deploy IMG=quay.io/ahadas/autoci

See https://docs.openshift.com/container-platform/4.15/operators/operator_sdk/ansible/osdk-ansible-tutorial.html

After deployment, create a ContainerBuild like those we have in the examples.

Trigger a "nightly" build:
oc -n autosd create job osbuild-nightly-test --from=cronjob/automotive-osbuild-nightly
