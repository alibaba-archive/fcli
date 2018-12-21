package util

import (
	"github.com/stretchr/testify/suite"
	"gopkg.in/h2non/gock.v1"
	"os/exec"
	"testing"
)

func TestUtil(t *testing.T) {
	suite.Run(t, new(UtilTestSuite))
}

type UtilTestSuite struct {
	suite.Suite
}

func (s *UtilTestSuite) TestGetRegionNo() {
	testEndpoint := "1234567.fc.cn-hangzhou.aliyuncs.com:8080"
	region := GetRegionNoForSLSEndpoint(testEndpoint)
	s.Equal("cn-hangzhou", region)
}

func (s *UtilTestSuite) TestCheckImageExist() {
	output = func(cmd *exec.Cmd) (bytes []byte, e error) {
		return []byte{}, nil
	}

	isExist := CheckImageExist("aliyunfc/runtime-nodejs6", "build")
	s.False(isExist)

	output = func(cmd *exec.Cmd) (bytes []byte, e error) {
		return []byte{1}, nil
	}
	isExist = CheckImageExist("aliyunfc/runtime-nodejs6", "build")
	s.True(isExist)
}

func (s *UtilTestSuite) TestGetPublicImageDigest() {
	defer gock.Off()
	gock.New("https://auth.docker.io").
		MatchParam("service", "registry.docker.io").
		MatchParam("scope", "repository:aliyunfc/runtime-nodejs6:pull").
		Get("token").
		Reply(200).
		JSON(map[string]string{"token": "testToken"})

	mockDigest := "sha256:123456"

	gock.New("https://registry-1.docker.io").
		MatchHeader("Authorization", "Bearer testToken").
		Get("v2/aliyunfc/runtime-nodejs6/manifests/build").
		Reply(200).
		AddHeader("Docker-Content-Digest", mockDigest)

	digest, _ := GetPublicImageDigest("aliyunfc/runtime-nodejs6", "build")
	s.Equal(mockDigest, digest)
}

func (s *UtilTestSuite) TestGetLocalImageDigest() {
	mockDigest := "sha256:123456"

	output = func(cmd *exec.Cmd) (bytes []byte, e error) {
		return []byte("'aliyunfc/runtime-nodejs6@" + mockDigest + "'"), nil
	}

	digest, _ := GetLocalImageDigest("aliyunfc/runtime-nodejs6", "build")
	s.Equal(mockDigest, digest)
}
