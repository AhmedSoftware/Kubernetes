// This is a generated file. Do not edit directly.

module k8s.io/sample-controller

go 1.15

require (
	k8s.io/api v0.0.0
	k8s.io/apimachinery v0.0.0
	k8s.io/client-go v0.0.0
	k8s.io/code-generator v0.0.0
	k8s.io/klog/v2 v2.2.0
)

replace (
	golang.org/x/sys => golang.org/x/sys v0.0.0-20200622214017-ed371f2e16b4
	k8s.io/api => ../api
	k8s.io/apimachinery => ../apimachinery
	k8s.io/client-go => ../client-go
	k8s.io/code-generator => ../code-generator
	k8s.io/sample-controller => ../sample-controller
)
