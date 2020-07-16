# image-labels-test

```shell
make
make build
image-labels-test inspect registry.redhat.io/rhpam-7/rhpam-kieserver-rhel8:7.7.1
```

deploy to openshift/k8s
```shell
oc adm policy add-scc-to-user hostmount-anyuid -z image-labels
oc apply -f build/ds.yaml
```