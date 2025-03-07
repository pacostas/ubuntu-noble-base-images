package acceptance_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/layout"
	. "github.com/onsi/gomega"
	"github.com/paketo-buildpacks/packit/v2/vacation"
	"github.com/sclevine/spec"

	. "github.com/paketo-buildpacks/jam/integration/matchers"
	. "github.com/paketo-buildpacks/packit/v2/matchers"
)

func testMetadataFullStack(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		tmpDir string
	)

	it.Before(func() {
		var err error
		tmpDir, err = os.MkdirTemp("", "")
		Expect(err).NotTo(HaveOccurred())
	})

	it.After(func() {
		Expect(os.RemoveAll(tmpDir)).To(Succeed())
	})

	it("builds full stack", func() {
		var buildReleaseDate, runReleaseDate time.Time

		by("confirming that the build image is correct", func() {
			dir := filepath.Join(tmpDir, "build-index")
			err := os.Mkdir(dir, os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			archive, err := os.Open(fullStack.BuildArchive)
			Expect(err).NotTo(HaveOccurred())
			defer archive.Close()

			err = vacation.NewArchive(archive).Decompress(dir)
			Expect(err).NotTo(HaveOccurred())

			path, err := layout.FromPath(dir)
			Expect(err).NotTo(HaveOccurred())

			index, err := path.ImageIndex()
			Expect(err).NotTo(HaveOccurred())

			indexManifest, err := index.IndexManifest()
			Expect(err).NotTo(HaveOccurred())

			Expect(indexManifest.Manifests).To(HaveLen(2))
			Expect(indexManifest.Manifests[0].Platform).To(Equal(&v1.Platform{
				OS:           "linux",
				Architecture: "amd64",
			}))
			Expect(indexManifest.Manifests[1].Platform).To(Equal(&v1.Platform{
				OS:           "linux",
				Architecture: "arm64",
			}))

			image, err := index.Image(indexManifest.Manifests[0].Digest)
			Expect(err).NotTo(HaveOccurred())

			file, err := image.ConfigFile()
			Expect(err).NotTo(HaveOccurred())

			Expect(file.Config.Labels).To(SatisfyAll(
				HaveKeyWithValue("io.buildpacks.stack.id", "io.buildpacks.stacks.noble"),
				HaveKeyWithValue("io.buildpacks.stack.description", "ubuntu:noble with compilers and common libraries and utilities"),
				HaveKeyWithValue("io.buildpacks.stack.distro.name", "ubuntu"),
				HaveKeyWithValue("io.buildpacks.stack.distro.version", "24.04"),
				HaveKeyWithValue("io.buildpacks.stack.homepage", "https://github.com/paketo-buildpacks/noble-full-stack"),
				HaveKeyWithValue("io.buildpacks.stack.maintainer", "Paketo Buildpacks"),
				HaveKeyWithValue("io.buildpacks.stack.metadata", MatchJSON("{}")),
			))

			buildReleaseDate, err = time.Parse(time.RFC3339, file.Config.Labels["io.buildpacks.stack.released"])
			Expect(err).NotTo(HaveOccurred())
			Expect(buildReleaseDate).NotTo(BeZero())

			Expect(image).To(SatisfyAll(
				HaveFileWithContent("/etc/group", ContainSubstring("cnb:x:1000:")),
				HaveFileWithContent("/etc/passwd", ContainSubstring("cnb:x:1001:1000::/home/cnb:/bin/bash")),
				HaveDirectory("/home/cnb"),
			))

			Expect(file.Config.User).To(Equal("1001:1000"))

			Expect(file.Config.Env).To(ContainElements(
				"CNB_USER_ID=1001",
				"CNB_GROUP_ID=1000",
				"CNB_STACK_ID=io.buildpacks.stacks.noble",
			))

			Expect(image).To(HaveFileWithContent("/etc/gitconfig", ContainLines(
				"[safe]",
				"\tdirectory = /workspace",
				"\tdirectory = /workspace/source-ws",
				"\tdirectory = /workspace/source",
			)))

			Expect(image).To(HaveFileWithContent("/var/lib/dpkg/status", SatisfyAll(
				ContainSubstring("Package: automake"),
				ContainSubstring("Package: bison"),
				ContainSubstring("Package: build-essential"),
				ContainSubstring("Package: bzr"),
				ContainSubstring("Package: ca-certificates"),
				ContainSubstring("Package: cmake"),
				ContainSubstring("Package: comerr-dev"),
				ContainSubstring("Package: curl"),
				ContainSubstring("Package: dh-python"),
				ContainSubstring("Package: dnsutils"),
				ContainSubstring("Package: file"),
				ContainSubstring("Package: flex"),
				ContainSubstring("Package: gdb"),
				ContainSubstring("Package: gir1.2-gdkpixbuf-2.0"),
				ContainSubstring("Package: gir1.2-rsvg-2.0"),
				ContainSubstring("Package: git"),
				ContainSubstring("Package: gnupg"),
				ContainSubstring("Package: gnupg1"),
				ContainSubstring("Package: graphviz"),
				ContainSubstring("Package: gsfonts"),
				ContainSubstring("Package: gss-ntlmssp"),
				ContainSubstring("Package: gss-ntlmssp-dev"),
				ContainSubstring("Package: imagemagick"),
				ContainSubstring("Package: imagemagick-6-common"),
				ContainSubstring("Package: jq"),
				ContainSubstring("Package: krb5-user"),
				ContainSubstring("Package: libaio-dev"),
				ContainSubstring("Package: libaio1t64"),
				ContainSubstring("Package: libarchive-extract-perl"),
				ContainSubstring("Package: libargon2-1"),
				ContainSubstring("Package: libargon2-dev"),
				ContainSubstring("Package: libatm1"),
				ContainSubstring("Package: libatm1-dev"),
				ContainSubstring("Package: libaudiofile-dev"),
				ContainSubstring("Package: libaudiofile1"),
				ContainSubstring("Package: libavcodec-dev"),
				ContainSubstring("Package: libavcodec60"),
				ContainSubstring("Package: libbabeltrace1"),
				ContainSubstring("Package: libblas-dev"),
				ContainSubstring("Package: libblas3"),
				ContainSubstring("Package: libc6"),
				ContainSubstring("Package: libcurl4"),
				ContainSubstring("Package: libcurl4-openssl-dev"),
				ContainSubstring("Package: libdjvulibre-text"),
				ContainSubstring("Package: libdjvulibre21"),
				ContainSubstring("Package: libdw1"),
				ContainSubstring("Package: liberror-perl"),
				ContainSubstring("Package: libestr-dev"),
				ContainSubstring("Package: libestr0"),
				ContainSubstring("Package: libexif12"),
				ContainSubstring("Package: libffi-dev"),
				ContainSubstring("Package: libffi8"),
				ContainSubstring("Package: libfl-dev"),
				ContainSubstring("Package: libfl2"),
				ContainSubstring("Package: libfribidi-dev"),
				ContainSubstring("Package: libfribidi0"),
				ContainSubstring("Package: libgcrypt20"),
				ContainSubstring("Package: libgcrypt20-dev"),
				ContainSubstring("Package: libgd-dev"),
				ContainSubstring("Package: libgmp-dev"),
				ContainSubstring("Package: libgmp10"),
				ContainSubstring("Package: libgmpxx4ldbl"),
				ContainSubstring("Package: libgnutls-openssl27t64"),
				ContainSubstring("Package: libgnutls28-dev"),
				ContainSubstring("Package: libgnutls30t64"),
				ContainSubstring("Package: libgraphviz-dev"),
				ContainSubstring("Package: libharfbuzz-icu0"),
				ContainSubstring("Package: libicu-dev"),
				ContainSubstring("Package: libidn12"),
				ContainSubstring("Package: libidn11-dev"),
				ContainSubstring("Package: libjson-glib-1.0-0"),
				ContainSubstring("Package: libjson-glib-dev"),
				ContainSubstring("Package: libkrb5-dev"),
				ContainSubstring("Package: liblapack-dev"),
				ContainSubstring("Package: liblapack3"),
				ContainSubstring("Package: libldap2"),
				ContainSubstring("Package: libldap2-dev"),
				ContainSubstring("Package: liblockfile-bin"),
				ContainSubstring("Package: liblockfile-dev"),
				ContainSubstring("Package: liblockfile1"),
				ContainSubstring("Package: libmagic-dev"),
				ContainSubstring("Package: libmagic1"),
				ContainSubstring("Package: libmagickwand-dev"),
				ContainSubstring("Package: libmariadb-dev-compat"),
				ContainSubstring("Package: libmariadb3"),
				ContainSubstring("Package: libmodule-pluggable-perl"),
				ContainSubstring("Package: libncurses6"),
				ContainSubstring("Package: libncurses-dev"),
				ContainSubstring("Package: libnl-3-200"),
				ContainSubstring("Package: libnl-3-dev"),
				ContainSubstring("Package: libnl-genl-3-200"),
				ContainSubstring("Package: libnl-genl-3-dev"),
				ContainSubstring("Package: libopenblas0"),
				ContainSubstring("Package: libopenblas-dev"),
				ContainSubstring("Package: libopenexr-3-1-30"),
				ContainSubstring("Package: liborc-0.4-0"),
				ContainSubstring("Package: liborc-0.4-dev"),
				ContainSubstring("Package: libp11-kit-dev"),
				ContainSubstring("Package: libp11-kit0"),
				ContainSubstring("Package: libpam-cap"),
				ContainSubstring("Package: libpango-1.0-0"),
				ContainSubstring("Package: libpango1.0-dev"),
				ContainSubstring("Package: libpathplan4"),
				ContainSubstring("Package: libpcre32-3"),
				ContainSubstring("Package: libpq-dev"),
				ContainSubstring("Package: libpq5"),
				ContainSubstring("Package: libproxy-dev"),
				ContainSubstring("Package: libproxy1v5"),
				ContainSubstring("Package: libpython3.12-stdlib"),
				ContainSubstring("Package: libpython3.12"),
				ContainSubstring("Package: libreadline-dev"),
				ContainSubstring("Package: libreadline8"),
				ContainSubstring("Package: librtmp-dev"),
				ContainSubstring("Package: libsasl2-2"),
				ContainSubstring("Package: libsasl2-dev"),
				ContainSubstring("Package: libsasl2-modules"),
				ContainSubstring("Package: libsasl2-modules-gssapi-mit"),
				ContainSubstring("Package: libselinux1"),
				ContainSubstring("Package: libselinux1-dev"),
				ContainSubstring("Package: libsigc++-2.0-0v5"),
				ContainSubstring("Package: libsigc++-2.0-dev"),
				ContainSubstring("Package: libsigsegv2"),
				ContainSubstring("Package: libsqlite3-0"),
				ContainSubstring("Package: libsqlite3-dev"),
				ContainSubstring("Package: libssl-dev"),
				ContainSubstring("Package: libsysfs-dev"),
				ContainSubstring("Package: libsysfs2"),
				ContainSubstring("Package: libtasn1-6"),
				ContainSubstring("Package: libtasn1-6-dev"),
				ContainSubstring("Package: libterm-ui-perl"),
				ContainSubstring("Package: libtiffxx6"),
				ContainSubstring("Package: libtirpc-dev"),
				ContainSubstring("Package: libtirpc-common"),
				ContainSubstring("Package: libtool"),
				ContainSubstring("Package: libunwind8"),
				ContainSubstring("Package: libustr-1.0-1"),
				ContainSubstring("Package: libustr-dev"),
				ContainSubstring("Package: libwmf0.2-7"),
				ContainSubstring("Package: libwrap0-dev"),
				ContainSubstring("Package: libwrap0"),
				ContainSubstring("Package: libxapian-dev"),
				ContainSubstring("Package: libxapian30"),
				ContainSubstring("Package: libxdot4"),
				ContainSubstring("Package: libxslt1-dev"),
				ContainSubstring("Package: libxslt1.1"),
				ContainSubstring("Package: libyaml-0-2"),
				ContainSubstring("Package: libyaml-dev"),
				ContainSubstring("Package: libsodium23"),
				ContainSubstring("Package: libsodium-dev"),
				ContainSubstring("Package: lockfile-progs"),
				ContainSubstring("Package: lsof"),
				ContainSubstring("Package: lzma"),
				ContainSubstring("Package: mercurial"),
				ContainSubstring("Package: net-tools"),
				ContainSubstring("Package: ocaml-base-nox"),
				ContainSubstring("Package: openssh-client"),
				ContainSubstring("Package: openssl"),
				ContainSubstring("Package: psmisc"),
				ContainSubstring("Package: python3"),
				ContainSubstring("Package: rsync"),
				ContainSubstring("Package: subversion"),
				ContainSubstring("Package: sysstat"),
				ContainSubstring("Package: ubuntu-minimal"),
				ContainSubstring("Package: unixodbc"),
				ContainSubstring("Package: unixodbc-dev"),
				ContainSubstring("Package: unzip"),
				ContainSubstring("Package: uuid"),
				ContainSubstring("Package: uuid-dev"),
				ContainSubstring("Package: wget"),
				ContainSubstring("Package: zip"),
				ContainSubstring("Package: tzdata"),
			)))
		})

		by("confirming that the run image is correct", func() {
			dir := filepath.Join(tmpDir, "run-index")
			err := os.Mkdir(dir, os.ModePerm)
			Expect(err).NotTo(HaveOccurred())

			archive, err := os.Open(fullStack.RunArchive)
			Expect(err).NotTo(HaveOccurred())
			defer archive.Close()

			err = vacation.NewArchive(archive).Decompress(dir)
			Expect(err).NotTo(HaveOccurred())

			path, err := layout.FromPath(dir)
			Expect(err).NotTo(HaveOccurred())

			index, err := path.ImageIndex()
			Expect(err).NotTo(HaveOccurred())

			indexManifest, err := index.IndexManifest()
			Expect(err).NotTo(HaveOccurred())

			Expect(indexManifest.Manifests).To(HaveLen(2))
			Expect(indexManifest.Manifests[0].Platform).To(Equal(&v1.Platform{
				OS:           "linux",
				Architecture: "amd64",
			}))
			Expect(indexManifest.Manifests[1].Platform).To(Equal(&v1.Platform{
				OS:           "linux",
				Architecture: "arm64",
			}))

			image, err := index.Image(indexManifest.Manifests[0].Digest)
			Expect(err).NotTo(HaveOccurred())

			file, err := image.ConfigFile()
			Expect(err).NotTo(HaveOccurred())

			Expect(file.Config.Labels).To(SatisfyAll(
				HaveKeyWithValue("io.buildpacks.stack.id", "io.buildpacks.stacks.noble"),
				HaveKeyWithValue("io.buildpacks.stack.description", "ubuntu:noble with common libraries and utilities"),
				HaveKeyWithValue("io.buildpacks.stack.distro.name", "ubuntu"),
				HaveKeyWithValue("io.buildpacks.stack.distro.version", "24.04"),
				HaveKeyWithValue("io.buildpacks.stack.homepage", "https://github.com/paketo-buildpacks/noble-full-stack"),
				HaveKeyWithValue("io.buildpacks.stack.maintainer", "Paketo Buildpacks"),
				HaveKeyWithValue("io.buildpacks.stack.metadata", MatchJSON("{}")),
			))

			runReleaseDate, err = time.Parse(time.RFC3339, file.Config.Labels["io.buildpacks.stack.released"])
			Expect(err).NotTo(HaveOccurred())
			Expect(runReleaseDate).NotTo(BeZero())

			Expect(file.Config.User).To(Equal("1002:1000"))

			Expect(image).To(SatisfyAll(
				HaveFileWithContent("/etc/group", ContainSubstring("cnb:x:1000:")),
				HaveFileWithContent("/etc/passwd", ContainSubstring("cnb:x:1002:1000::/home/cnb:/bin/bash")),
				HaveDirectory("/home/cnb"),
			))

			Expect(image).To(SatisfyAll(
				HaveFile("/usr/share/doc/ca-certificates/copyright"),
				HaveFile("/etc/ssl/certs/ca-certificates.crt"),
				HaveDirectory("/root"),
				HaveDirectory("/tmp"),
				HaveFile("/etc/services"),
				HaveFile("/etc/nsswitch.conf"),
			))

			Expect(image).To(HaveFileWithContent("/var/lib/dpkg/status", SatisfyAll(
				ContainSubstring("Package: ca-certificates"),
				ContainSubstring("Package: curl"),
				ContainSubstring("Package: dh-python"),
				ContainSubstring("Package: dnsutils"),
				ContainSubstring("Package: file"),
				ContainSubstring("Package: gir1.2-gdkpixbuf-2.0"),
				ContainSubstring("Package: gir1.2-rsvg-2.0"),
				ContainSubstring("Package: gnupg"),
				ContainSubstring("Package: gnupg1"),
				ContainSubstring("Package: graphviz"),
				ContainSubstring("Package: gsfonts"),
				ContainSubstring("Package: gss-ntlmssp"),
				ContainSubstring("Package: imagemagick"),
				ContainSubstring("Package: imagemagick-6-common"),
				ContainSubstring("Package: jq"),
				ContainSubstring("Package: krb5-user"),
				ContainSubstring("Package: libaio1t64"),
				ContainSubstring("Package: libarchive-extract-perl"),
				ContainSubstring("Package: libargon2-1"),
				ContainSubstring("Package: libatm1"),
				ContainSubstring("Package: libaudiofile1"),
				ContainSubstring("Package: libavcodec60"),
				ContainSubstring("Package: libbabeltrace1"),
				ContainSubstring("Package: libblas3"),
				ContainSubstring("Package: libc6"),
				ContainSubstring("Package: libcurl4"),
				ContainSubstring("Package: libdjvulibre-text"),
				ContainSubstring("Package: libdjvulibre21"),
				ContainSubstring("Package: libdw1"),
				ContainSubstring("Package: liberror-perl"),
				ContainSubstring("Package: libestr0"),
				ContainSubstring("Package: libexif12"),
				ContainSubstring("Package: libffi8"),
				ContainSubstring("Package: libfl2"),
				ContainSubstring("Package: libfribidi0"),
				ContainSubstring("Package: libgcrypt20"),
				ContainSubstring("Package: libgmp10"),
				ContainSubstring("Package: libgmpxx4ldbl"),
				ContainSubstring("Package: libgnutls-openssl27t64"),
				ContainSubstring("Package: libgnutls28-dev"),
				ContainSubstring("Package: libgnutls30t64"),
				ContainSubstring("Package: libgraphviz-dev"),
				ContainSubstring("Package: libharfbuzz-icu0"),
				ContainSubstring("Package: libidn12"),
				ContainSubstring("Package: libisl23"),
				ContainSubstring("Package: libjson-glib-1.0-0"),
				ContainSubstring("Package: libjsoncpp25"),
				ContainSubstring("Package: liblapack3"),
				ContainSubstring("Package: libldap2"),
				ContainSubstring("Package: liblockfile-bin"),
				ContainSubstring("Package: liblockfile1"),
				ContainSubstring("Package: libmagic1"),
				ContainSubstring("Package: libmariadb3"),
				ContainSubstring("Package: libmodule-pluggable-perl"),
				ContainSubstring("Package: libmpc3"),
				ContainSubstring("Package: libmpfr6"),
				ContainSubstring("Package: libncurses6"),
				ContainSubstring("Package: libnl-3-200"),
				ContainSubstring("Package: libnl-genl-3-200"),
				ContainSubstring("Package: libopenblas0"),
				ContainSubstring("Package: libopenexr-3-1-30"),
				ContainSubstring("Package: liborc-0.4-0"),
				ContainSubstring("Package: libp11-kit0"),
				ContainSubstring("Package: libpam-cap"),
				ContainSubstring("Package: libpango-1.0-0"),
				ContainSubstring("Package: libpango1.0-dev"),
				ContainSubstring("Package: libpathplan4"),
				ContainSubstring("Package: libpcre32-3"),
				ContainSubstring("Package: libpq5"),
				ContainSubstring("Package: libproxy1v5"),
				ContainSubstring("Package: libpython3-stdlib"),
				ContainSubstring("Package: libpython3.12"),
				ContainSubstring("Package: libreadline8"),
				ContainSubstring("Package: librhash0"),
				ContainSubstring("Package: libsasl2-2"),
				ContainSubstring("Package: libsasl2-modules"),
				ContainSubstring("Package: libsasl2-modules-gssapi-mit"),
				ContainSubstring("Package: libselinux1"),
				ContainSubstring("Package: libsigc++-2.0-0v5"),
				ContainSubstring("Package: libsigsegv2"),
				ContainSubstring("Package: libsqlite3-0"),
				ContainSubstring("Package: libsysfs2"),
				ContainSubstring("Package: libtasn1-6"),
				ContainSubstring("Package: libterm-ui-perl"),
				ContainSubstring("Package: libtiffxx6"),
				ContainSubstring("Package: libtirpc-common"),
				ContainSubstring("Package: libunwind8"),
				ContainSubstring("Package: libustr-1.0-1"),
				ContainSubstring("Package: libuv1"),
				ContainSubstring("Package: libwmf0.2-7"),
				ContainSubstring("Package: libwrap0"),
				ContainSubstring("Package: libxapian30"),
				ContainSubstring("Package: libxdot4"),
				ContainSubstring("Package: libxslt1.1"),
				ContainSubstring("Package: libyaml-0-2"),
				ContainSubstring("Package: libsodium23"),
				ContainSubstring("Package: lockfile-progs"),
				ContainSubstring("Package: lsof"),
				ContainSubstring("Package: lzma"),
				ContainSubstring("Package: net-tools"),
				ContainSubstring("Package: ocaml-base-nox"),
				ContainSubstring("Package: openssh-client"),
				ContainSubstring("Package: openssl"),
				ContainSubstring("Package: psmisc"),
				ContainSubstring("Package: python3"),
				ContainSubstring("Package: rsync"),
				ContainSubstring("Package: subversion"),
				ContainSubstring("Package: ubuntu-minimal"),
				ContainSubstring("Package: tzdata"),
				ContainSubstring("Package: unixodbc"),
				ContainSubstring("Package: unzip"),
				ContainSubstring("Package: uuid"),
				ContainSubstring("Package: wget"),
				ContainSubstring("Package: zip"),
			)))

			Expect(image).To(HaveFileWithContent("/etc/os-release", SatisfyAll(
				ContainSubstring(`PRETTY_NAME="Paketo Buildpacks Full Noble"`),
				ContainSubstring(`HOME_URL="https://github.com/paketo-buildpacks/noble-full-stack"`),
				ContainSubstring(`SUPPORT_URL="https://github.com/paketo-buildpacks/noble-full-stack/blob/main/README.md"`),
				ContainSubstring(`BUG_REPORT_URL="https://github.com/paketo-buildpacks/noble-full-stack/issues/new"`),
			)))
		})

		Expect(runReleaseDate).To(Equal(buildReleaseDate))
	})
}
