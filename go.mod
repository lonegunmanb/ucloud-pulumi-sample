module github.com/lonegunmanb/ucloud-pulumi-sample

go 1.13

replace (
	github.com/pulumi/pulumi-ucloud v0.0.3 => github.com/lonegunmanb/pulumi-ucloud v0.0.3
	github.com/uber/jaeger-lib v2.1.1+incompatible => github.com/uber/jaeger-lib v1.5.0
)

require (
	github.com/pulumi/pulumi v1.9.1
	github.com/pulumi/pulumi-ucloud v0.0.3
)
