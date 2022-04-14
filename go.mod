module github.com/bhojpur/os

go 1.17

replace github.com/knative/pkg => github.com/rancher/pkg v0.0.0-20190514055449-b30ab9de040e
replace github.com/docker/docker => github.com/moby/moby v1.4.2-0.20170731201646-1009e6a40b29

require (
	github.com/bhojpur/host v0.0.6
	github.com/docker/docker v20.10.14+incompatible
	github.com/docker/go-units v0.4.0
	github.com/ghodss/yaml v1.0.0
	github.com/mattn/go-isatty v0.0.14
	github.com/mattn/go-shellwords v1.0.12
	github.com/otiai10/copy v1.0.2
	github.com/paultag/go-modprobe v0.0.0-20180523225901-7ccf6a1b5271
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.1
	github.com/urfave/cli v1.22.5
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad
	gopkg.in/freddierice/go-losetup.v1 v1.0.0-20170407175016-fc9adea44124
	k8s.io/apimachinery v0.23.5
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	golang.org/x/term v0.0.0-20220411215600-e5f449aeb171 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	k8s.io/utils v0.0.0-20220210201930-3a6ce19ff2f9 // indirect
	pault.ag/go/modprobe v0.1.2 // indirect
	pault.ag/go/topsort v0.1.1 // indirect
)

replace (
	k8s.io/api => k8s.io/api v0.23.3
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.23.3
	k8s.io/apimachinery => k8s.io/apimachinery v0.23.3
	k8s.io/apiserver => k8s.io/apiserver v0.23.3
	k8s.io/client-go => k8s.io/client-go v0.23.3
	k8s.io/code-generator => k8s.io/code-generator v0.23.3
	k8s.io/component-base => k8s.io/component-base v0.23.3
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.23.3
	k8s.io/metrics => k8s.io/metrics v0.23.3
	k8s.io/mount-utils => k8s.io/mount-utils v0.23.3
)
