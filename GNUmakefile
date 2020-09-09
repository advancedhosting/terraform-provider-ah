TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
VERSION=$(shell ./scripts/git-version.sh)
PKG_NAME=ah
WEBSITE_REPO=github.com/hashicorp/terraform-website
export CGO_ENABLED:=0


testacc: 
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m