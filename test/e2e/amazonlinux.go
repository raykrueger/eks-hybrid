//go:build e2e
// +build e2e

package e2e

import (
	"context"
	_ "embed"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

//go:embed testdata/amazonlinux/2023/cloud-init.txt
var al23CloudInit []byte

type amazonLinuxCloudInitData struct {
	NodeadmConfig     string
	NodeadmUrl        string
	KubernetesVersion string
	Provider          string
}

type AmazonLinux2023 struct {
	Architecture string
}

func NewAmazonLinux2023AMD() *AmazonLinux2023 {
	al := new(AmazonLinux2023)
	al.Architecture = amd64Arch
	return al
}

func NewAmazonLinux2023ARM() *AmazonLinux2023 {
	al := new(AmazonLinux2023)
	al.Architecture = arm64Arch
	return al
}

func (a AmazonLinux2023) Name() string {
	if a.Architecture == amd64Arch {
		return "al23-amd64"
	}
	return "al23-arm64"
}

func (a AmazonLinux2023) InstanceType() string {
	if a.Architecture == amd64Arch {
		return "m5.large"
	}
	return "t4g.large"
}

func (a AmazonLinux2023) AMIName(ctx context.Context, awsSession *session.Session) (string, error) {
	amiId, err := getAmiIDFromSSM(ctx, ssm.New(awsSession), "/aws/service/ami-amazon-linux-latest/al2023-ami-kernel-default-"+a.Architecture)
	return *amiId, err
}

func (a AmazonLinux2023) BuildUserData(nodeadmUrls NodeadmURLs, nodeadmConfigYaml, kubernetesVersion, provider string) ([]byte, error) {
	data := amazonLinuxCloudInitData{
		NodeadmConfig:     nodeadmConfigYaml,
		NodeadmUrl:        nodeadmUrls.AMD,
		KubernetesVersion: kubernetesVersion,
		Provider:          provider,
	}

	if a.Architecture == arm64Arch {
		data.NodeadmUrl = nodeadmUrls.ARM
	}

	return executeTemplate(al23CloudInit, data)
}