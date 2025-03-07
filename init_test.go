package acceptance_test

import (
	"fmt"
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/onsi/gomega/format"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
)

var tinyStack struct {
	BuildArchive string
	RunArchive   string
	BuildImageID string
	RunImageID   string
}

var baseStack struct {
	BuildArchive string
	RunArchive   string
	BuildImageID string
	RunImageID   string
}

var fullStack struct {
	BuildArchive string
	RunArchive   string
	BuildImageID string
	RunImageID   string
}

var localRegistryPort int

func by(_ string, f func()) { f() }

func getFreePort() (port int, err error) {
	if l, err := net.Listen("tcp", ":0"); err == nil {
		defer l.Close()
		return l.Addr().(*net.TCPAddr).Port, nil
	}
	return 0, err
}

func TestAcceptance(t *testing.T) {
	format.MaxLength = 0
	SetDefaultEventuallyTimeout(30 * time.Second)

	Expect := NewWithT(t).Expect

	root, err := filepath.Abs(".")
	Expect(err).ToNot(HaveOccurred())

	localRegistryPort, err = getFreePort()
	Expect(err).ToNot(HaveOccurred())

	tinyStack.BuildArchive = filepath.Join(root, "builds", "noble-tiny-stack", "build.oci")
	tinyStack.BuildImageID = fmt.Sprintf("localhost:%d/tiny-stack-build-%s", localRegistryPort, uuid.NewString())

	tinyStack.RunArchive = filepath.Join(root, "builds", "noble-tiny-stack", "run.oci")
	tinyStack.RunImageID = fmt.Sprintf("localhost:%d/-tiny-stack-run-%s", localRegistryPort, uuid.NewString())

	baseStack.BuildArchive = filepath.Join(root, "builds", "noble-base-stack", "build.oci")
	baseStack.BuildImageID = fmt.Sprintf("localhost:%d/base-stack-build-%s", localRegistryPort, uuid.NewString())

	baseStack.RunArchive = filepath.Join(root, "builds", "noble-base-stack", "run.oci")
	baseStack.RunImageID = fmt.Sprintf("localhost:%d/-base-stack-run-%s", localRegistryPort, uuid.NewString())

	fullStack.BuildArchive = filepath.Join(root, "builds", "noble-full-stack", "build.oci")
	fullStack.BuildImageID = fmt.Sprintf("localhost:%d/full-stack-build-%s", localRegistryPort, uuid.NewString())

	fullStack.RunArchive = filepath.Join(root, "builds", "noble-full-stack", "run.oci")
	fullStack.RunImageID = fmt.Sprintf("localhost:%d/-full-stack-run-%s", localRegistryPort, uuid.NewString())

	suite := spec.New("Acceptance", spec.Report(report.Terminal{}), spec.Parallel())
	suite("Tiny Stack Metadata", testMetadataTinyStack)
	suite("Base Stack Metadata", testMetadataBaseStack)
	suite("Full Stack Metadata", testMetadataFullStack)
	suite.Run(t)
}
