// Package images manages the local cache of guest images and the Firecracker
// runtime: streamed downloads with progress tracking, squashfs-to-ext4
// conversion, and runtime installation into the Agni data directory.
package images

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/anubhavg-icpl/agni/internal/logging"
	"github.com/anubhavg-icpl/agni/internal/templates"
)

// Job and artifact states.
const (
	StatePending      = "pending"
	StateDownloading  = "downloading"
	StateProcessing   = "processing"
	StateDone         = "done"
	StateError        = "error"
	StateRunning      = "running"
	phaseDownload     = "downloading"
	phaseConvert      = "converting rootfs"
	phaseComplete     = "complete"
	progressInterval  = 250 * time.Millisecond
	maxRecentJobsList = 20
)

// ArtifactProgress is the live progress of one artifact within a job.
type ArtifactProgress struct {
	Name        string `json:"name"`
	Filename    string `json:"filename"`
	Kind        string `json:"kind"`
	State       string `json:"state"`
	BytesDone   int64  `json:"bytes_done"`
	BytesTotal  int64  `json:"bytes_total"`
	Path        string `json:"path,omitempty"`
	Description string `json:"description,omitempty"`
}

// Job is one template preparation run: downloads plus any conversion steps.
type Job struct {
	ID         string             `json:"id"`
	TemplateID string             `json:"template_id"`
	State      string             `json:"state"`
	Phase      string             `json:"phase"`
	Artifacts  []ArtifactProgress `json:"artifacts"`
	Error      string             `json:"error,omitempty"`
	StartedAt  time.Time          `json:"started_at"`
	FinishedAt *time.Time         `json:"finished_at,omitempty"`

	// Shared pointer so value copies (snapshots, JSON responses) stay vet-clean.
	mu *sync.Mutex `json:"-"`
}

func (j *Job) snapshot() Job {
	j.mu.Lock()
	defer j.mu.Unlock()
	out := Job{
		ID:         j.ID,
		TemplateID: j.TemplateID,
		State:      j.State,
		Phase:      j.Phase,
		Error:      j.Error,
		StartedAt:  j.StartedAt,
		FinishedAt: j.FinishedAt,
		mu:         j.mu,
	}
	out.Artifacts = append([]ArtifactProgress{}, j.Artifacts...)
	return out
}

func (j *Job) setArtifact(i int, state, description string) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if i >= 0 && i < len(j.Artifacts) {
		j.Artifacts[i].State = state
		j.Artifacts[i].Description = description
	}
}

func (j *Job) setArtifactProgress(i int, done, total int64) {
	j.mu.Lock()
	defer j.mu.Unlock()
	if i >= 0 && i < len(j.Artifacts) {
		j.Artifacts[i].State = StateDownloading
		j.Artifacts[i].BytesDone = done
		j.Artifacts[i].BytesTotal = total
	}
}

func (j *Job) setPhase(phase string) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.Phase = phase
}

func (j *Job) fail(err error) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.State = StateError
	j.Error = err.Error()
	now := time.Now()
	j.FinishedAt = &now
}

func (j *Job) complete() {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.State = StateDone
	j.Phase = phaseComplete
	now := time.Now()
	j.FinishedAt = &now
}

// ArtifactState describes one template artifact on local disk.
type ArtifactState struct {
	Name       string   `json:"name"`
	Filename   string   `json:"filename"`
	Kind       string   `json:"kind"`
	URL        string   `json:"url"`
	Path       string   `json:"path"`
	Ready      bool     `json:"ready"`
	SizeBytes  int64    `json:"size_bytes"`
	SizeHintMB int      `json:"size_hint_mb"`
	ToolDeps   []string `json:"tool_deps,omitempty"`
}

// Manager owns the image cache directory and all download jobs.
type Manager struct {
	dir       string
	binDir    string
	client    *http.Client
	logger    *logging.Logger
	onRuntime func(path string)

	mu      sync.Mutex
	jobs    map[string]*Job
	order   []string
	current map[string]string // template ID -> active job ID
}

// NewManager creates an image manager rooted at <dataDir>/images with runtime
// binaries installed to <dataDir>/bin. onRuntime, when set, is called with the
// installed binary path each time a runtime install completes.
func NewManager(dataDir string, logger *logging.Logger, onRuntime func(string)) *Manager {
	if logger == nil {
		logger = logging.GetLogger().WithComponent("images")
	}
	return &Manager{
		dir:       filepath.Join(dataDir, "images"),
		binDir:    filepath.Join(dataDir, "bin"),
		client:    &http.Client{},
		logger:    logger,
		onRuntime: onRuntime,
		jobs:      make(map[string]*Job),
		current:   make(map[string]string),
	}
}

// Dir returns the image cache directory.
func (m *Manager) Dir() string { return m.dir }

func (m *Manager) artifactPath(t templates.Template, a templates.Artifact) string {
	// Kernels are shared across templates (many catalog entries reference the
	// same upstream kernel): store them in a URL-keyed cache so each kernel is
	// downloaded to disk exactly once.
	if a.Kind == templates.KindKernel {
		sum := sha256.Sum256([]byte(a.URL))
		key := hex.EncodeToString(sum[:])[:16]
		return filepath.Join(m.dir, "_cache", key, a.Filename)
	}
	return filepath.Join(m.dir, t.ID, a.Filename)
}

// State reports the on-disk state of every artifact of a template.
func (m *Manager) State(t templates.Template) []ArtifactState {
	out := make([]ArtifactState, 0, len(t.Artifacts))
	for _, a := range t.Artifacts {
		path := m.artifactPath(t, a)
		st := ArtifactState{
			Name:       a.Name,
			Filename:   a.Filename,
			Kind:       a.Kind,
			URL:        a.URL,
			Path:       path,
			SizeHintMB: a.SizeHintMB,
			ToolDeps:   a.ToolDeps,
		}
		if info, err := os.Stat(path); err == nil {
			st.Ready = true
			st.SizeBytes = info.Size()
		}
		out = append(out, st)
	}
	return out
}

// Ready reports whether every artifact of the template is on disk.
func (m *Manager) Ready(t templates.Template) bool {
	for _, a := range t.Artifacts {
		if _, err := os.Stat(m.artifactPath(t, a)); err != nil {
			return false
		}
	}
	return true
}

// RuntimePath returns the installed firecracker binary, if present.
func (m *Manager) RuntimePath() (string, bool) {
	p := filepath.Join(m.binDir, "firecracker")
	info, err := os.Stat(p)
	if err != nil || info.IsDir() || info.Mode()&0111 == 0 {
		return "", false
	}
	return p, true
}

// KVMAvailable reports whether this process can use /dev/kvm.
func (m *Manager) KVMAvailable() bool {
	info, err := os.Stat("/dev/kvm")
	return err == nil && info.Mode()&0200 != 0
}

// LatestJobFor returns the newest job for a template, if any.
func (m *Manager) LatestJobFor(templateID string) *Job {
	m.mu.Lock()
	defer m.mu.Unlock()
	var latest *Job
	for _, id := range m.order {
		if j, ok := m.jobs[id]; ok && j.TemplateID == templateID {
			latest = j
		}
	}
	if latest == nil {
		return nil
	}
	s := latest.snapshot()
	return &s
}

// Job returns a snapshot of one job.
func (m *Manager) Job(id string) (Job, bool) {
	m.mu.Lock()
	j := m.jobs[id]
	m.mu.Unlock()
	if j == nil {
		return Job{}, false
	}
	return j.snapshot(), true
}

// Jobs returns recent job snapshots, newest first.
func (m *Manager) Jobs() []Job {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]Job, 0, len(m.order))
	for i := len(m.order) - 1; i >= 0 && len(out) < maxRecentJobsList; i-- {
		if j, ok := m.jobs[m.order[i]]; ok {
			out = append(out, j.snapshot())
		}
	}
	return out
}

// Prepare starts (or joins) a download + conversion job for a template.
func (m *Manager) Prepare(t templates.Template) (*Job, error) {
	m.mu.Lock()
	if jobID, ok := m.current[t.ID]; ok {
		if j, ok := m.jobs[jobID]; ok && j.State == StateRunning {
			m.mu.Unlock()
			s := j.snapshot()
			return &s, nil
		}
	}
	m.mu.Unlock()

	// Fail fast when conversion tools are missing rather than after a download.
	for _, a := range t.Artifacts {
		for _, tool := range a.ToolDeps {
			if _, err := exec.LookPath(tool); err != nil {
				return nil, fmt.Errorf("the %q template needs the %q tool on this host. Install it and try again", t.ID, tool)
			}
		}
	}

	job := &Job{
		ID:         fmt.Sprintf("job-%d", time.Now().UnixNano()),
		TemplateID: t.ID,
		State:      StateRunning,
		Phase:      phaseDownload,
		StartedAt:  time.Now(),
		mu:         &sync.Mutex{},
	}
	for _, a := range t.Artifacts {
		job.Artifacts = append(job.Artifacts, ArtifactProgress{
			Name:     a.Name,
			Filename: a.Filename,
			Kind:     a.Kind,
			State:    StatePending,
			Path:     m.artifactPath(t, a),
		})
	}

	m.mu.Lock()
	m.jobs[job.ID] = job
	m.order = append(m.order, job.ID)
	m.current[t.ID] = job.ID
	m.mu.Unlock()

	m.logger.Info().Str("template", t.ID).Str("job", job.ID).Msg("Preparing template")
	go m.run(job, t)

	s := job.snapshot()
	return &s, nil
}

func (m *Manager) run(job *Job, t templates.Template) {
	for i, a := range t.Artifacts {
		final := m.artifactPath(t, a)
		if _, err := os.Stat(final); err == nil {
			job.setArtifact(i, StateDone, "Already on disk")
			continue
		}

		if err := os.MkdirAll(filepath.Dir(final), 0750); err != nil {
			job.fail(fmt.Errorf("creating image directory: %w", err))
			return
		}

		if a.Kind == templates.KindRuntime {
			if err := m.installRuntime(job, i, a); err != nil {
				job.fail(fmt.Errorf("installing the firecracker runtime: %w", err))
				return
			}
			continue
		}

		target := final
		switch a.PostProcess {
		case "squashfs-to-ext4":
			target = final + ".squashfs"
		case "tar-to-ext4":
			target = final + ".tar"
		case "oci-to-ext4":
			target = "" // pulled and extracted by ociToExt4 itself
		}
		if target != "" {
			if err := m.download(job, i, a, target); err != nil {
				job.fail(fmt.Errorf("downloading %s: %w", a.Name, err))
				return
			}
		}
		switch a.PostProcess {
		case "squashfs-to-ext4":
			job.setArtifact(i, StateProcessing, "Converting squashfs to ext4")
			job.setPhase(phaseConvert)
			if err := m.squashfsToExt4(target, final); err != nil {
				job.fail(fmt.Errorf("converting %s: %w", a.Name, err))
				return
			}
			_ = os.Remove(target)
		case "tar-to-ext4":
			job.setArtifact(i, StateProcessing, "Building ext4 root filesystem")
			job.setPhase(phaseConvert)
			if err := m.tarToExt4(target, final, a.PostProcessArgs); err != nil {
				job.fail(fmt.Errorf("building %s: %w", a.Name, err))
				return
			}
			_ = os.Remove(target)
		case "oci-to-ext4":
			job.setArtifact(i, StateDownloading, "Pulling container image")
			job.setPhase(phaseDownload)
			if err := m.ociToExt4(job, i, a, final); err != nil {
				job.fail(fmt.Errorf("building %s: %w", a.Name, err))
				return
			}
		}
		job.setArtifact(i, StateDone, "Ready")
	}
	job.complete()
	m.logger.Info().Str("template", t.ID).Str("job", job.ID).Msg("Template ready")
}

func (m *Manager) download(job *Job, i int, a templates.Artifact, dest string) error {
	job.setArtifact(i, StateDownloading, "Starting download")

	resp, err := m.client.Get(a.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status %s", resp.Status)
	}

	tmp := dest + ".part"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}

	total := resp.ContentLength
	pr := &progressWriter{
		total: total,
		last:  time.Now(),
		report: func(done int64) {
			job.setArtifactProgress(i, done, total)
		},
	}
	_, copyErr := io.Copy(f, io.TeeReader(resp.Body, pr))
	closeErr := f.Close()
	if copyErr != nil {
		_ = os.Remove(tmp)
		return copyErr
	}
	if closeErr != nil {
		_ = os.Remove(tmp)
		return closeErr
	}
	if err := os.Rename(tmp, dest); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	job.setArtifactProgress(i, total, total)
	return nil
}

type progressWriter struct {
	done   int64
	total  int64
	last   time.Time
	report func(int64)
}

func (p *progressWriter) Write(b []byte) (int, error) {
	p.done += int64(len(b))
	if time.Since(p.last) >= progressInterval {
		p.last = time.Now()
		p.report(p.done)
	}
	return len(b), nil
}

func (m *Manager) installRuntime(job *Job, i int, a templates.Artifact) error {
	if err := os.MkdirAll(m.binDir, 0750); err != nil {
		return err
	}
	dest := filepath.Join(m.binDir, "firecracker")
	if err := m.download(job, i, a, dest+".new"); err != nil {
		return err
	}
	if err := os.Chmod(dest+".new", 0755); err != nil {
		return err
	}
	if err := os.Rename(dest+".new", dest); err != nil {
		return err
	}
	job.setArtifact(i, StateDone, "Installed")
	if m.onRuntime != nil {
		m.onRuntime(dest)
	}
	m.logger.Info().Str("path", dest).Msg("Firecracker runtime installed")
	return nil
}

func (m *Manager) squashfsToExt4(squashPath, ext4Path string) error {
	workDir := ext4Path + ".work"
	defer os.RemoveAll(workDir)

	rootDir := filepath.Join(workDir, "root")
	if out, err := runCmd("unsquashfs", "-q", "-d", rootDir, squashPath); err != nil {
		return fmt.Errorf("unsquashfs failed: %s: %w", out, err)
	}
	return mkfsExt4FromDir(rootDir, ext4Path)
}

// ociReference identifies an image on an OCI registry.
type ociReference struct {
	registry string // e.g. registry-1.docker.io, quay.io
	name     string // e.g. library/rockylinux
	tag      string
}

// ociToExt4 builds a rootfs from an OCI container image: select the host-arch
// manifest, stream every layer (applying whiteouts), run the prep steps, and
// pack the flattened tree into an ext4 image.
func (m *Manager) ociToExt4(job *Job, i int, a templates.Artifact, ext4Path string) error {
	ref, err := parseOCIRef(a.URL)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 0}
	authHeader, err := ociToken(client, ref)
	if err != nil {
		return fmt.Errorf("authenticating to %s: %w", ref.registry, err)
	}

	manifest, err := ociFetchManifest(client, ref, authHeader)
	if err != nil {
		return fmt.Errorf("fetching manifest for %s:%s: %w", ref.name, ref.tag, err)
	}

	var total int64
	for _, l := range manifest.Layers {
		total += l.Size
	}
	job.setArtifactProgress(i, 0, total)

	workDir := ext4Path + ".work"
	defer os.RemoveAll(workDir)
	rootDir := filepath.Join(workDir, "root")
	if err := os.MkdirAll(rootDir, 0755); err != nil {
		return err
	}

	dirModes := map[string]os.FileMode{}
	var done int64
	for _, l := range manifest.Layers {
		if err := m.ociApplyLayer(client, ref, authHeader, l.Digest, rootDir, dirModes, func(n int64) {
			job.setArtifactProgress(i, done+n, total)
		}); err != nil {
			return fmt.Errorf("applying layer %s: %w", l.Digest[:19], err)
		}
		done += l.Size
		job.setArtifactProgress(i, done, total)
	}
	if err := applyDirModes(rootDir, dirModes); err != nil {
		return fmt.Errorf("applying directory modes: %w", err)
	}

	job.setArtifact(i, StateProcessing, "Configuring guest")
	job.setPhase(phaseConvert)
	// Container images ship secret files with mode 0000 (e.g. /etc/gshadow on
	// RHEL); mkfs.ext4 -d runs unprivileged and must be able to read them.
	if err := ensureOwnerReadable(rootDir); err != nil {
		return fmt.Errorf("relaxing file modes: %w", err)
	}
	for _, arg := range a.PostProcessArgs {
		switch arg {
		case "systemd-autologin":
			if err := prepareUbuntuAutologin(rootDir); err != nil {
				return err
			}
		case "shell-init":
			if err := prepareShellInit(rootDir); err != nil {
				return err
			}
		}
	}
	return mkfsExt4FromDir(rootDir, ext4Path)
}

func parseOCIRef(url string) (ociReference, error) {
	s := strings.TrimPrefix(url, "oci://")
	if s == url || s == "" {
		return ociReference{}, fmt.Errorf("not an oci:// reference: %s", url)
	}
	registry, rest, ok := strings.Cut(s, "/")
	if !ok {
		return ociReference{}, fmt.Errorf("missing repository in %s", url)
	}
	name, tag, _ := strings.Cut(rest, ":")
	if name == "" || tag == "" {
		return ociReference{}, fmt.Errorf("missing name or tag in %s", url)
	}
	if strings.HasPrefix(registry, "docker.io") {
		registry = "registry-1.docker.io"
	}
	return ociReference{registry: registry, name: name, tag: tag}, nil
}

// ociToken returns an Authorization header value; public registries that need
// no token (quay.io) yield an empty string.
func ociToken(client *http.Client, ref ociReference) (string, error) {
	if ref.registry != "registry-1.docker.io" {
		return "", nil
	}
	u := fmt.Sprintf("https://auth.docker.io/token?service=registry.docker.io&scope=repository:%s:pull", ref.name)
	resp, err := client.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var body struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}
	if body.Token == "" {
		return "", fmt.Errorf("empty token from auth.docker.io")
	}
	return "Bearer " + body.Token, nil
}

type ociLayerInfo struct {
	Digest string `json:"digest"`
	Size   int64  `json:"size"`
}

type ociManifest struct {
	MediaType string `json:"mediaType"`
	Manifests []struct {
		Digest   string `json:"digest"`
		Platform struct {
			Architecture string `json:"architecture"`
			OS           string `json:"os"`
		} `json:"platform"`
	} `json:"manifests"`
	Layers []ociLayerInfo `json:"layers"`
}

func ociFetchManifest(client *http.Client, ref ociReference, auth string) (*ociManifest, error) {
	index, err := ociGetManifestJSON(client, ref, ref.tag, auth)
	if err != nil {
		return nil, err
	}
	if len(index.Layers) > 0 {
		return index, nil
	}
	// Multi-arch index: resolve the matching platform child manifest.
	want := "amd64"
	if runtime.GOARCH == "arm64" {
		want = "arm64"
	}
	for _, child := range index.Manifests {
		if child.Platform.OS == "linux" && child.Platform.Architecture == want {
			return ociGetManifestJSON(client, ref, child.Digest, auth)
		}
	}
	return nil, fmt.Errorf("no linux/%s manifest in %s:%s", want, ref.name, ref.tag)
}

func ociGetManifestJSON(client *http.Client, ref ociReference, reference, auth string) (*ociManifest, error) {
	u := fmt.Sprintf("https://%s/v2/%s/manifests/%s", ref.registry, ref.name, reference)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}
	req.Header.Set("Accept", strings.Join([]string{
		"application/vnd.oci.image.index.v1+json",
		"application/vnd.docker.distribution.manifest.list.v2+json",
		"application/vnd.oci.image.manifest.v1+json",
		"application/vnd.docker.distribution.manifest.v2+json",
	}, ", "))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registry returned %s for %s", resp.Status, reference)
	}
	var m ociManifest
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

// ociApplyLayer downloads one layer and extracts it over the root directory,
// honoring OCI whiteouts (.wh.<name> deletes, .wh..wh..opq empties a directory).
func (m *Manager) ociApplyLayer(client *http.Client, ref ociReference, auth, digest, rootDir string, dirModes map[string]os.FileMode, progress func(int64)) error {
	u := fmt.Sprintf("https://%s/v2/%s/blobs/%s", ref.registry, ref.name, digest)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("blob fetch returned %s", resp.Status)
	}

	counting := &countingReader{r: resp.Body, fn: progress}
	gz, err := gzip.NewReader(counting)
	if err != nil {
		return fmt.Errorf("layer is not gzip: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if err := applyTarEntry(hdr, tr, rootDir, dirModes); err != nil {
			return fmt.Errorf("%s: %w", hdr.Name, err)
		}
	}
}

type countingReader struct {
	r  io.Reader
	n  int64
	fn func(int64)
}

func (c *countingReader) Read(p []byte) (int, error) {
	n, err := c.r.Read(p)
	c.n += int64(n)
	if c.fn != nil {
		c.fn(c.n)
	}
	return n, err
}

// applyTarEntry writes a single tar entry into the target root, skipping
// device nodes (devtmpfs provides them at boot). Directories are created
// owner-writable and their final mode is recorded in dirModes, because later
// layers routinely write into directories whose shipped mode is read-only.
func applyTarEntry(hdr *tar.Header, r io.Reader, rootDir string, dirModes map[string]os.FileMode) error {
	clean := path.Clean("/" + hdr.Name)
	target := filepath.Join(rootDir, clean)
	base := path.Base(clean)
	dir := filepath.Dir(target)

	switch {
	case base == ".wh..wh..opq":
		entries, err := os.ReadDir(dir)
		if err != nil {
			return err
		}
		for _, e := range entries {
			if err := os.RemoveAll(filepath.Join(dir, e.Name())); err != nil {
				return err
			}
		}
		return nil
	case strings.HasPrefix(base, ".wh."):
		return os.RemoveAll(filepath.Join(dir, strings.TrimPrefix(base, ".wh.")))
	}

	switch hdr.Typeflag {
	case tar.TypeDir:
		if err := os.MkdirAll(target, 0700); err != nil {
			return err
		}
		dirModes[clean] = hdr.FileInfo().Mode()
		return nil
	case tar.TypeReg:
		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}
		_ = os.Remove(target)
		f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, hdr.FileInfo().Mode().Perm())
		if err != nil {
			return err
		}
		if _, err := io.Copy(f, r); err != nil {
			f.Close()
			return err
		}
		if err := f.Chmod(hdr.FileInfo().Mode()); err != nil {
			f.Close()
			return err
		}
		return f.Close()
	case tar.TypeSymlink:
		_ = os.Remove(target)
		return os.Symlink(hdr.Linkname, target)
	case tar.TypeLink:
		_ = os.Remove(target)
		source := filepath.Join(rootDir, path.Clean("/"+hdr.Linkname))
		return os.Link(source, target)
	case tar.TypeBlock, tar.TypeChar, tar.TypeFifo:
		return nil // handled by devtmpfs / not needed in a VM root
	default:
		return nil
	}
}

// ensureOwnerReadable adds owner read (and, for directories, traverse) to
// every entry under root so mkfs.ext4 can populate the image without root.
func ensureOwnerReadable(rootDir string) error {
	return filepath.Walk(rootDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		mode := info.Mode().Perm()
		var want os.FileMode
		if info.IsDir() {
			want = mode | 0500
		} else if mode&0400 == 0 {
			want = mode | 0400
		} else {
			return nil
		}
		if want != mode {
			return os.Chmod(p, want)
		}
		return nil
	})
}

// applyDirModes restores the directory modes recorded during layer
// extraction, children before parents so restricted parents (e.g. /root 0700
// or /usr 0555) can still be walked while their children are chmod'ed.
func applyDirModes(rootDir string, modes map[string]os.FileMode) error {
	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() {
			if err := applyDirModes(filepath.Join(rootDir, e.Name()), modes); err != nil {
				return err
			}
		}
	}
	if mode, ok := modes["/"]; ok && rootDir != "" {
		if err := os.Chmod(rootDir, mode.Perm()); err != nil {
			return err
		}
	}
	return nil
}

// prepareShellInit gives image-only roots (no init system, e.g. the Nix image)
// a minimal /init that mounts the essentials and opens a root shell on ttyS0.
func prepareShellInit(rootDir string) error {
	script := `#!/bin/sh
# PID 1 inherits the kernel console (ttyS0), so no redirects are needed.
# Container images have no /dev nodes or mount(8); a static busybox is
# seeded at build time and provides both.
BB=/bin/busybox
[ -x "$BB" ] || BB=""
mkdir_p() { if [ -n "$BB" ]; then $BB mkdir -p "$@"; else mkdir -p "$@"; fi; }
do_mount() { if [ -n "$BB" ]; then $BB mount "$@"; else mount "$@"; fi; }
mkdir_p /dev /proc /sys /run
do_mount -t proc proc /proc
do_mount -t sysfs sys /sys
do_mount -t devtmpfs devtmpfs /dev
export PATH=/nix/var/nix/profiles/default/bin:/usr/bin:/usr/sbin:/bin:/sbin
# Restart the serial shell if it exits; the loop only ends when the VM stops.
while true; do
	/bin/sh -i
done
`
	// Seed a static busybox (from the build host, when available) so the
	// minimal container image gains mount(8) and core utilities at runtime.
	if busybox, lerr := exec.LookPath("busybox"); lerr == nil {
		if data, rerr := os.ReadFile(busybox); rerr == nil {
			_ = os.WriteFile(filepath.Join(rootDir, "bin", "busybox"), data, 0755)
		}
	}
	return os.WriteFile(filepath.Join(rootDir, "init"), []byte(script), 0755)
}

// tarToExt4 unpacks a distribution rootfs tarball (Alpine minirootfs, Ubuntu
// Base) into a fresh ext4 image, applying the requested build tweaks first.
func (m *Manager) tarToExt4(tarPath, ext4Path string, args []string) error {
	workDir := ext4Path + ".work"
	defer os.RemoveAll(workDir)

	rootDir := filepath.Join(workDir, "root")
	if err := os.MkdirAll(rootDir, 0755); err != nil {
		return fmt.Errorf("creating work directory: %w", err)
	}
	if out, err := runCmd("tar", "-xf", tarPath, "-C", rootDir, "--numeric-owner", "-p"); err != nil {
		return fmt.Errorf("tar extract failed: %s: %w", out, err)
	}

	for _, arg := range args {
		switch arg {
		case "alpine-init":
			if err := prepareAlpineInit(rootDir); err != nil {
				return err
			}
		case "ubuntu-autologin":
			if err := prepareUbuntuAutologin(rootDir); err != nil {
				return err
			}
		case "void-serial":
			if err := prepareVoidSerial(rootDir); err != nil {
				return err
			}
		}
	}

	return mkfsExt4FromDir(rootDir, ext4Path)
}

// prepareAlpineInit writes a busybox inittab that brings up the essential
// mounts and a login getty on the serial console.
func prepareAlpineInit(rootDir string) error {
	inittab := `::sysinit:/bin/mount -t proc proc /proc
::sysinit:/bin/mount -t sysfs sys /sys
::sysinit:/bin/mount -o remount,rw /
ttyS0::respawn:/sbin/getty -L 115200 ttyS0 vt100
::ctrlaltdel:/sbin/reboot
::shutdown:/bin/umount -a
`
	if err := os.WriteFile(filepath.Join(rootDir, "etc", "inittab"), []byte(inittab), 0644); err != nil {
		return fmt.Errorf("writing inittab: %w", err)
	}
	// Permit root logins on the serial console.
	f, err := os.OpenFile(filepath.Join(rootDir, "etc", "securetty"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		_, _ = f.WriteString("ttyS0\n")
		_ = f.Close()
	}
	return nil
}

// prepareUbuntuAutologin makes systemd's serial getty log in as root without
// a password (Ubuntu Base ships with the root account locked).
func prepareUbuntuAutologin(rootDir string) error {
	dir := filepath.Join(rootDir, "etc", "systemd", "system", "serial-getty@ttyS0.service.d")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating systemd override: %w", err)
	}
	conf := "[Service]\nExecStart=\nExecStart=-/sbin/agetty --autologin root --noclear %I $TERM\n"
	return os.WriteFile(filepath.Join(dir, "autologin.conf"), []byte(conf), 0644)
}

// mkfsExt4FromDir packs a root directory into a fresh ext4 image. The
// filesystem size is derived from the directory (plus headroom) because
// mkfs.ext4 cannot infer a size for a not-yet-existing image file.
func mkfsExt4FromDir(rootDir, ext4Path string) error {
	var contentBytes, files, dirs, links int64
	err := filepath.Walk(rootDir, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		switch {
		case info.Mode().IsRegular():
			contentBytes += info.Size()
			files++
		case info.IsDir():
			dirs++
		case info.Mode()&os.ModeSymlink != 0:
			links++
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("sizing root filesystem: %w", err)
	}

	// Container images such as the Nix store are huge trees of small files,
	// directories and symlinks; size the inode table from the real entry
	// count instead of mkfs' default bytes-per-inode ratio, which exhausts
	// inodes long before data blocks.
	entries := files + dirs + links
	inodes := int64(float64(entries)*1.5) + 16384

	// Content with slack, plus room for the inode table itself and the
	// ext4 journal, which the population step allocates out of the same pool.
	size := int64(float64(contentBytes)*1.45) + 64<<20 + inodes*384
	if min := int64(64 << 20); size < min {
		size = min
	}
	sizeMB := (size + (1 << 20) - 1) >> 20

	var lastErr error
	// Huge trees can still surprise the estimate; if population runs out of
	// blocks or inodes, retry with a doubled image instead of failing.
	for attempt := 0; attempt < 3; attempt++ {
		tmp := ext4Path + ".new"
		out, err := runCmd("mkfs.ext4", "-q", "-F", "-d", rootDir, "-L", "root", "-N",
			strconv.FormatInt(inodes, 10), tmp,
			strconv.FormatInt(sizeMB, 10)+"M")
		if err == nil {
			return os.Rename(tmp, ext4Path)
		}
		_ = os.Remove(tmp)
		lastErr = fmt.Errorf("mkfs.ext4 failed (size %dMB, %d inodes): %s: %w", sizeMB, inodes, out, err)
		sizeMB *= 2
		inodes = inodes*3/2 + 4096
	}
	return lastErr
}

// prepareVoidSerial wires Void Linux's runit to spawn an autologin serial
// getty on ttyS0 (Void ROOTFS ships agetty but does not enable serial).
func prepareVoidSerial(rootDir string) error {
	svDir := filepath.Join(rootDir, "etc", "sv", "agetty-ttyS0")
	if _, err := os.Stat(filepath.Join(svDir, "run")); err != nil {
		if err := os.MkdirAll(svDir, 0755); err != nil {
			return err
		}
		runScript := "#!/bin/sh\n[ -r conf ] && . ./conf\nexec agetty --noclear ${GETTY_ARGS} ttyS0 115200 vt100\n"
		if err := os.WriteFile(filepath.Join(svDir, "run"), []byte(runScript), 0755); err != nil {
			return err
		}
	}
	conf := "GETTY_ARGS=\"--autologin root\"\n"
	if err := os.WriteFile(filepath.Join(svDir, "conf"), []byte(conf), 0644); err != nil {
		return err
	}
	// Enable the service in the default runlevel.
	cur := filepath.Join(rootDir, "etc", "runit", "runsvdir", "current")
	if err := os.MkdirAll(cur, 0755); err != nil {
		return err
	}
	link := filepath.Join(cur, "agetty-ttyS0")
	if _, err := os.Lstat(link); err != nil {
		if err := os.Symlink("/etc/sv/agetty-ttyS0", link); err != nil {
			return err
		}
	}
	return nil
}

func runCmd(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}
