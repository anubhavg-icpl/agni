// Package images manages the local cache of guest images and the Firecracker
// runtime: streamed downloads with progress tracking, squashfs-to-ext4
// conversion, and runtime installation into the Agni data directory.
package images

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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
		if a.PostProcess == "squashfs-to-ext4" {
			target = final + ".squashfs"
		}
		if err := m.download(job, i, a, target); err != nil {
			job.fail(fmt.Errorf("downloading %s: %w", a.Name, err))
			return
		}
		if a.PostProcess == "squashfs-to-ext4" {
			job.setArtifact(i, StateProcessing, "Converting squashfs to ext4")
			job.setPhase(phaseConvert)
			if err := m.squashfsToExt4(target, final); err != nil {
				job.fail(fmt.Errorf("converting %s: %w", a.Name, err))
				return
			}
			_ = os.Remove(target)
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

	tmp := ext4Path + ".new"
	if out, err := runCmd("mkfs.ext4", "-q", "-F", "-d", rootDir, "-L", "root", tmp); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("mkfs.ext4 failed: %s: %w", out, err)
	}
	return os.Rename(tmp, ext4Path)
}

func runCmd(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	return string(out), err
}
