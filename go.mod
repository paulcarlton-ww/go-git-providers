module github.com/fluxcd/go-git-providers

go 1.14

require (
	github.com/drone/go-scm v1.11.0
	github.com/google/go-cmp v0.4.0
	github.com/google/go-github/v32 v32.1.0
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79
	github.com/ktrysmt/go-bitbucket v0.6.2
	github.com/nbio/st v0.0.0-20140626010706-e9e8d9816f32 // indirect
	github.com/onsi/ginkgo v1.14.0
	github.com/onsi/gomega v1.10.1
	github.com/xanzy/go-gitlab v0.33.0
	golang.org/x/oauth2 v0.0.0-20181106182150-f42d05182288
)

replace (
	github.com/fluxcd/go-git-providers => github.com/paulcarlton-ww/go-git-providers v0.0.4
)
