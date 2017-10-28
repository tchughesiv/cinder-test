# cinder-test
Container for testing standalone Cinder w/ oauth-proxy on OpenShift

Assumes standalone cinder deployed to "openstack" ns and available via service at:
https://cinder.openstack.svc

#### OpenShift deployment
```shell
$ oc new-app tchughesiv/cinder-test
```

#### Local build
```shell
$ go build
$ make
```
