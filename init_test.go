package acceptance_test

import (
	"fmt"
	"os"
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

var RegistryUrl string

func by(_ string, f func()) { f() }

func TestAcceptance(t *testing.T) {
	format.MaxLength = 0
	SetDefaultEventuallyTimeout(30 * time.Second)

	Expect := NewWithT(t).Expect

	RegistryUrl = os.Getenv("REGISTRY_URL")
	Expect(RegistryUrl).NotTo(Equal(""))

	root, err := filepath.Abs(".")
	Expect(err).ToNot(HaveOccurred())

	tinyStack.BuildArchive = filepath.Join(root, "builds", "noble-tiny-stack", "build.oci")
	tinyStack.BuildImageID = fmt.Sprintf("%s/noble-tiny-stack-build-%s", RegistryUrl, uuid.NewString())

	tinyStack.RunArchive = filepath.Join(root, "builds", "noble-tiny-stack", "run.oci")
	tinyStack.RunImageID = fmt.Sprintf("%s/noble-tiny-stack-run-%s", RegistryUrl, uuid.NewString())

	baseStack.BuildArchive = filepath.Join(root, "builds", "noble-base-stack", "build.oci")
	baseStack.BuildImageID = fmt.Sprintf("%s/noble-base-stack-build-%s", RegistryUrl, uuid.NewString())

	baseStack.RunArchive = filepath.Join(root, "builds", "noble-base-stack", "run.oci")
	baseStack.RunImageID = fmt.Sprintf("%s/noble-base-stack-run-%s", RegistryUrl, uuid.NewString())

	fullStack.BuildArchive = filepath.Join(root, "builds", "noble-full-stack", "build.oci")
	fullStack.BuildImageID = fmt.Sprintf("%s/noble-full-stack-build-%s", RegistryUrl, uuid.NewString())

	fullStack.RunArchive = filepath.Join(root, "builds", "noble-full-stack", "run.oci")
	fullStack.RunImageID = fmt.Sprintf("%s/noble-full-stack-run-%s", RegistryUrl, uuid.NewString())

	suite := spec.New("Acceptance", spec.Report(report.Terminal{}), spec.Parallel())
	// suite("MetadataTinyStack", testMetadataTinyStack)
	// suite("MetadataBaseStack", testMetadataBaseStack)
	// suite("MetadataFullStack", testMetadataFullStack)
	// suite("BuildpackIntegrationTinyStack", testBuildpackIntegrationTinyStack)
	// suite("BuildpackIntegrationBaseStack", testBuildpackIntegrationBaseStack)
	suite("BuildpackIntegrationFullStack", testBuildpackIntegrationFullStack)

	suite.Run(t)
}
