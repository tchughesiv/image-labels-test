# image-labels-test

```shell
make
make build

image-labels-test inspect registry.redhat.io/amq7/amq-broker-lts-rhel7-operator@sha256:0414ea9cb57b7c018556414be501901bc4a6fc44b37fbef16d08db8ae658e20a
image-labels-test inspect ubi8
```

deploy to openshift/k8s
```shell
oc adm policy add-scc-to-user hostmount-anyuid -z image-labels
oc apply -f build/ds.yaml
```