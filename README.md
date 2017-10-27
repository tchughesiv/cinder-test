# cinder-test
Container for testing standalone Cinder w/ oauth-proxy on OpenShift

Assumes standalone cinder deployed to "openstack" ns and available via service at:
https://cinder.openstack.svc

```shell
$ go build
$ make
```