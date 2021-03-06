package main

import (
	"testing"
	"github.com/stretchr/testify/suite"
	"github.com/stretchr/testify/mock"
	"os"
	"fmt"
	"os/exec"
)

type DockerComposeTestSuite struct {
	suite.Suite
	dockerComposePath	string
	serviceName 		string
	target            	string
	sideTargets			[]string
	color             	string
	blueGreen         	bool
	host 			  	string
	certPath		  	string
	project 		  	string
}

func (s *DockerComposeTestSuite) SetupTest() {
	s.dockerComposePath = "test-docker-compose.yml"
	s.serviceName = "myService"
	s.target = "my-target"
	s.sideTargets = []string{"my-side-target-1", "my-side-target-2"}
	s.color = "red"
	s.blueGreen = false
	s.host = "tcp://1.2.3.4:1234"
	s.certPath = "/path/to/docker/cert"
	s.project = "my-project"
	readFile = func(fileName string) ([]byte, error) {
		return []byte(""), nil
	}
	writeFile = func(fileName string, data []byte, perm os.FileMode) error {
		return nil
	}
	removeFile = func(name string) error {
		return nil
	}
	execCmd = func(name string, arg ...string) *exec.Cmd {
		return &exec.Cmd{}
	}
}

// CreateFlow

func (s DockerComposeTestSuite) Test_CreateFlowFile_ReturnsNil() {
	actual := DockerCompose{}.CreateFlowFile(s.dockerComposePath, s.serviceName, s.target, s.sideTargets, s.color, s.blueGreen)

	s.Nil(actual)
}

func (s DockerComposeTestSuite) Test_CreateFlowFile_ReturnsError_WhenReadFile() {
	readFile = func(fileName string) ([]byte, error) {
		return []byte(""), fmt.Errorf("Some error")
	}

	err := DockerCompose{}.CreateFlowFile(s.dockerComposePath, s.serviceName, s.target, s.sideTargets, s.color, s.blueGreen)

	s.Error(err)
}

func (s DockerComposeTestSuite) Test_CreateFlowFile_CreatesTheFile() {
	var actual string
	writeFile = func(filename string, data []byte, perm os.FileMode) error {
		actual = filename
		return nil
	}

	DockerCompose{}.CreateFlowFile(s.dockerComposePath, s.serviceName, s.target, s.sideTargets, s.color, s.blueGreen)

	s.Equal(dockerComposeFlowPath, actual)
}

func (s DockerComposeTestSuite) Test_CreateFlowFile_CreatesDockerComposeReplica() {
	var actual string
	readFile = func(filename string) ([]byte, error) {
		actual = filename
		return []byte(""), nil
	}

	DockerCompose{}.CreateFlowFile(s.dockerComposePath, s.serviceName, s.target, s.sideTargets, s.color, s.blueGreen)

	s.Equal(s.dockerComposePath, actual)
}

func (s DockerComposeTestSuite) Test_CreateFlowFile_CreatesNewTarget_WhenBlueGreen() {
	color := "orange"
	var actual string
	var dcContent = fmt.Sprintf(`
%s:
  image: vfarcic/books-ms`,
		s.target,
	)
	newTarget := fmt.Sprintf("%s-%s", s.target, color)
	expected := fmt.Sprintf(`%s:
  extends:
    file: %s
    service: %s
  environment:
    - SERVICE_NAME=%s-%s
%s:
  extends:
    file: %s
    service: %s
%s:
  extends:
    file: %s
    service: %s`,
		newTarget,
		s.dockerComposePath,
		s.target,
		s.serviceName,
		color,
		s.sideTargets[0],
		s.dockerComposePath,
		s.sideTargets[0],
		s.sideTargets[1],
		s.dockerComposePath,
		s.sideTargets[1],
	)
	readFile = func(filename string) ([]byte, error) {
		return []byte(dcContent), nil
	}
	writeFile = func(filename string, data []byte, perm os.FileMode) error {
		actual = string(data)
		return nil
	}

	DockerCompose{}.CreateFlowFile(s.dockerComposePath, s.serviceName, s.target, s.sideTargets, color, true)

	s.Equal(expected, actual)
}

func (s DockerComposeTestSuite) Test_CreateFlowFile_UsesV2_WhenBlueGreen() {
	color := "orange"
	var actual string
	newTarget := fmt.Sprintf("%s-%s", s.target, color)
	var dcContent = fmt.Sprintf(`version: '2'

services:
  %s:
    image: vfarcic/books-ms`,
		s.target,
	)
	expected := fmt.Sprintf(`version: '2'

services:
  %s:
    extends:
      file: %s
      service: %s
    environment:
      - SERVICE_NAME=%s-%s
  %s:
    extends:
      file: %s
      service: %s
  %s:
    extends:
      file: %s
      service: %s`,
		newTarget,
		s.dockerComposePath,
		s.target,
		s.serviceName,
		color,
		s.sideTargets[0],
		s.dockerComposePath,
		s.sideTargets[0],
		s.sideTargets[1],
		s.dockerComposePath,
		s.sideTargets[1],

	)
	readFile = func(filename string) ([]byte, error) {
		return []byte(dcContent), nil
	}
	writeFile = func(filename string, data []byte, perm os.FileMode) error {
		actual = string(data)
		return nil
	}

	DockerCompose{}.CreateFlowFile(s.dockerComposePath, s.serviceName, s.target, s.sideTargets, color, true)

	s.Equal(expected, actual)
}

func (s DockerComposeTestSuite) Test_CreateFlowFile_ReturnsError_WhenWriteFile() {
	writeFile = func(filename string, data []byte, perm os.FileMode) error {
		return fmt.Errorf("Some error")
	}

	err := DockerCompose{}.CreateFlowFile(s.dockerComposePath, s.serviceName, s.target, s.sideTargets, s.color, s.blueGreen)

	s.Error(err)
}

// RemoveFlow

func (s DockerComposeTestSuite) Test_RemoveFlow_RemovesTheFile() {
	var actual string
	removeFile = func(name string) error {
		actual = name
		return nil
	}

	DockerCompose{}.RemoveFlow()

	s.Equal(dockerComposeFlowPath, actual)
}

func (s DockerComposeTestSuite) Test_RemoveFlow_ReturnsError() {
	removeFile = func(name string) error {
		return fmt.Errorf("Some error")
	}

	err := DockerCompose{}.RemoveFlow()

	s.Error(err)
}

// PullTargets

func (s DockerComposeTestSuite) Test_PullTargets_ReturnsNil_WhenTargetsAreEmpty() {
	actual := DockerCompose{}.PullTargets(s.host, s.certPath, s.project, []string{})

	s.Nil(actual)
}

func (s DockerComposeTestSuite) Test_PullTargets() {
	s.testCmd(DockerCompose{}.PullTargets, "pull", s.target)
}

// UpTargets

func (s DockerComposeTestSuite) Test_UpTargets() {
	s.testCmd(DockerCompose{}.UpTargets, "up", "-d", s.target)
}

// ScaleTargets

func (s DockerComposeTestSuite) Test_ScaleTargets_ReturnsNil_WhenTargetIsEmpty() {
	actual := DockerCompose{}.ScaleTargets(s.host, s.certPath, s.project, "", 8)

	s.Nil(actual)
}

func (s DockerComposeTestSuite) Test_ScaleTargets_CreatesTheCommand() {
	var scale = 7
	expected := []string{"docker-compose", "-f", dockerComposeFlowPath, "-p", s.project, "scale", fmt.Sprintf("%s=%d", s.target, scale)}
	actual := s.mockExecCmd()

	DockerCompose{}.ScaleTargets(s.host, s.certPath, s.project, s.target, scale)

	s.Equal(expected, *actual)
}

// RmTargets

func (s DockerComposeTestSuite) Test_RmTargets() {
	s.testCmd(DockerCompose{}.RmTargets, "rm", "-f", s.target)
}

// StopTargets

func (s DockerComposeTestSuite) Test_StopTargets() {
	s.testCmd(DockerCompose{}.StopTargets, "stop", s.target)
}

// Suite

func TestDockerComposeTestSuite(t *testing.T) {
	dockerHost := os.Getenv("DOCKER_HOST")
	dockerCertPath := os.Getenv("DOCKER_CERT_PATH")
	defer func() {
		os.Setenv("DOCKER_HOST", dockerHost)
		os.Setenv("DOCKER_CERT_PATH", dockerCertPath)
	}()
	suite.Run(t, new(DockerComposeTestSuite))
}

// Helper

func (s DockerComposeTestSuite) mockExecCmd() *[]string {
	var actualCommand []string
	execCmd = func(name string, arg ...string) *exec.Cmd {
		actualCommand = append([]string{name}, arg...)
		cmd := &exec.Cmd{}
		return cmd
	}
	return &actualCommand
}

type testCmdType func(host, certPath, project string, targets []string) error

func (s DockerComposeTestSuite) testCmd(f testCmdType, args ...string) {
	var expected []string
	var actual *[]string

	// Returns nil when targets are empty
	s.Nil(f(s.host, s.certPath, s.project, []string{}))

	// Creates command
	expected = append([]string{"docker-compose", "-f", dockerComposeFlowPath, "-p", s.project}, args...)
	actual = s.mockExecCmd()
	f(s.host, s.certPath, s.project, []string{s.target})
	s.Equal(expected, *actual)

	// Does not add project when empty
	expected = append([]string{"docker-compose", "-f", dockerComposeFlowPath}, args...)
	actual = s.mockExecCmd()
	f(s.host, s.certPath, "", []string{s.target})
	s.Equal(expected, *actual)

	// Adds DOCKER_HOST variable
	f(s.host, s.certPath, s.project, []string{s.target})
	host := s.host
	s.Equal(host, s.host)

	// Does not add DOCKER_HOST variable when empty
	f("", s.certPath, s.project, []string{s.target})
	s.NotEqual(os.Getenv("DOCKER_HOST"), s.host)

	// Adds DOCKER_CERT_PATH variable
	f(s.host, s.certPath, s.project, []string{s.target})
	s.Equal(os.Getenv("DOCKER_CERT_PATH"), s.certPath)

}


// Mock

type DockerComposeMock struct{
	mock.Mock
}

func (m *DockerComposeMock) CreateFlowFile(
		dcPath,
		serviceName,
		target string,
		sideTargets []string,
		color string,
		blueGreen bool,
	) error {
	args := m.Called(dcPath, serviceName, target, sideTargets, color, blueGreen)
	return args.Error(0)
}

func (m *DockerComposeMock) RemoveFlow() error {
	args := m.Called()
	return args.Error(0)
}

func (m *DockerComposeMock) PullTargets(host, certPath, project string, targets []string) error {
	args := m.Called(host, certPath, project, targets)
	return args.Error(0)
}

func (m *DockerComposeMock) UpTargets(host, certPath, project string, targets []string) error {
	args := m.Called(host, certPath, project, targets)
	return args.Error(0)
}

func (m *DockerComposeMock) ScaleTargets(host, certPath, project, target string, scale int) error {
	args := m.Called(host, certPath, project, target, scale)
	return args.Error(0)
}

func (m *DockerComposeMock) RmTargets(host, certPath, project string, targets []string) error {
	args := m.Called(host, certPath, project, targets)
	return args.Error(0)
}

func (m *DockerComposeMock) StopTargets(host, certPath, project string, targets []string) error {
	args := m.Called(host, certPath, project, targets)
	return args.Error(0)
}

func getDockerComposeMock(opts Opts, skipMethod string) *DockerComposeMock {
	mockObj := new(DockerComposeMock)
	if skipMethod != "PullTargets" {
		mockObj.On("PullTargets", opts.Host, opts.CertPath, opts.Project, Flow{}.GetPullTargets(opts)).Return(nil)
	}
	if skipMethod != "UpTargets" {
		mockObj.On("UpTargets", opts.Host, opts.CertPath, opts.Project, append(opts.SideTargets, opts.NextTarget)).Return(nil)
	}
	if skipMethod != "RmTargets" {
		mockObj.On("RmTargets", opts.Host, opts.CertPath, opts.Project, []string{opts.NextTarget}).Return(nil)
	}
	if skipMethod != "ScaleTargets" {
		mockObj.On("ScaleTargets", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	}
	if skipMethod != "CreateFlowFile" {
		mockObj.On(
			"CreateFlowFile",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(nil)
	}
	if skipMethod != "StopTargets" {
		mockObj.On("StopTargets", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	}
	if skipMethod != "RemoveFlow" {
		mockObj.On("RemoveFlow").Return(nil)
	}
	return mockObj
}
