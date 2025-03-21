package acceptance_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/layout"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/uuid"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"

	"github.com/paketo-buildpacks/occam"
	. "github.com/paketo-buildpacks/occam/matchers"
	"github.com/paketo-buildpacks/packit/vacation"
)

func testBuildpackIntegrationFullStack(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect     = NewWithT(t).Expect
		Eventually = NewWithT(t).Eventually

		buildPlanBuildpack string
		goDistBuildpack    string

		builderConfigFilepath string

		pack    occam.Pack
		docker  occam.Docker
		source  string
		name    string
		builder string

		image     occam.Image
		container occam.Container
	)

	it.Before(func() {
		pack = occam.NewPack().WithVerbose()
		docker = occam.NewDocker()

		var err error

		name, err = occam.RandomName()
		Expect(err).NotTo(HaveOccurred())

		buildPlanBuildpack = "index.docker.io/paketocommunity/build-plan"
		goDistBuildpack = "gcr.io/paketo-buildpacks/go-dist"

		source, err = occam.Source(filepath.Join("integration", "testdata", "go_simple_app"))
		Expect(err).NotTo(HaveOccurred())

		builderConfigFile, err := os.CreateTemp("", "builder.toml")
		Expect(err).NotTo(HaveOccurred())
		builderConfigFilepath = builderConfigFile.Name()

		_, err = fmt.Fprintf(builderConfigFile, `
[stack]
  build-image = "%s:latest"
  id = "io.buildpacks.stacks.noble"
  run-image = "%s:latest"
`,
			fullStack.BuildImageID,
			fullStack.RunImageID,
		)
		Expect(err).NotTo(HaveOccurred())

		Expect(archiveToDaemon(fullStack.BuildArchive, fullStack.BuildImageID)).To(Succeed())
		Expect(archiveToDaemon(fullStack.RunArchive, fullStack.RunImageID)).To(Succeed())

		builder = fmt.Sprintf("builder-%s", uuid.NewString())
		logs, err := createBuilder(builderConfigFilepath, builder)
		Expect(err).NotTo(HaveOccurred(), logs)
	})

	it.After(func() {
		Expect(docker.Container.Remove.Execute(container.ID)).To(Succeed())
		Expect(docker.Image.Remove.Execute(image.ID)).To(Succeed())
		Expect(docker.Volume.Remove.Execute(occam.CacheVolumeNames(name))).To(Succeed())

		_, err := getLifecycleVersion(builder)
		Expect(err).NotTo(HaveOccurred())

		Expect(docker.Image.Remove.Execute(builder)).To(Succeed())
		Expect(os.RemoveAll(builderConfigFilepath)).To(Succeed())

		Expect(docker.Image.Remove.Execute(fullStack.BuildImageID)).To(Succeed())
		Expect(docker.Image.Remove.Execute(fullStack.RunImageID)).To(Succeed())

		// Expect(docker.Image.Remove.Execute(fmt.Sprintf("buildpacksio/lifecycle:%s", lifecycleVersion))).To(Succeed())

		Expect(os.RemoveAll(source)).To(Succeed())
	})

	it("builds an app with a buildpack", func() {
		var err error
		var logs fmt.Stringer
		image, logs, err = pack.WithNoColor().Build.
			WithBuildpacks(
				goDistBuildpack,
				buildPlanBuildpack,
			).
			WithEnv(map[string]string{
				"BP_LOG_LEVEL": "DEBUG",
			}).
			WithPullPolicy("if-not-present").
			WithBuilder(builder).
			Execute(name, source)
		Expect(err).ToNot(HaveOccurred(), logs.String)

		container, err = docker.Container.Run.
			WithDirect().
			WithCommand("go").
			WithCommandArgs([]string{"run", "main.go"}).
			WithEnv(map[string]string{"PORT": "8080"}).
			WithPublish("8080").
			WithPublishAll().
			Execute(image.ID)
		Expect(err).NotTo(HaveOccurred())

		Eventually(container).WithTimeout(20 * time.Second).Should(BeAvailable())
		Eventually(container).Should(Serve(MatchRegexp(`go1.*`)).OnPort(8080))
	})
}

func archiveToDaemon(path, id string) error {
	tmpDir := os.TempDir()

	tarReader, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("unable to open tar: %w", err)
	}

	err = vacation.NewTarArchive(tarReader).Decompress(tmpDir)
	if err != nil {
		return fmt.Errorf("unable to extract files: %w", err)
	}

	pathLayout, err := layout.FromPath(tmpDir)
	if err != nil {
		return fmt.Errorf("unable to load image from path %s: %w", tmpDir, err)
	}

	imageIndex, err := pathLayout.ImageIndex()
	if err != nil {
		return fmt.Errorf("unable to read image index: %w", err)
	}

	ref, err := name.ParseReference(id)
	if err != nil {
		return fmt.Errorf("unable to parse reference from %s: %w", id, err)
	}

	return remote.WriteIndex(ref, imageIndex, remote.WithAuthFromKeychain(authn.DefaultKeychain))
}
