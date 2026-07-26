package contentimport_test

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/ikascrew/ikasbox/contentimport"
	"github.com/ikascrew/ikasbox/db"
)

func setupDB(t *testing.T) {
	t.Helper()

	f, err := sql.Open("sqlite3", filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("db open error: %+v", err)
	}
	t.Cleanup(func() { f.Close() })

	db.Use(f)
	if err := db.CreateTables(); err != nil {
		t.Fatalf("create tables error: %+v", err)
	}
}

// copySample copies the fixture jpg into dir under name and returns its path.
func copySample(t *testing.T, dir, name string) string {
	t.Helper()

	data, err := os.ReadFile(filepath.Join("testdata", "sample.jpg"))
	if err != nil {
		t.Fatalf("read fixture error: %+v", err)
	}

	dst := filepath.Join(dir, name)
	if err := os.WriteFile(dst, data, 0o644); err != nil {
		t.Fatalf("write fixture copy error: %+v", err)
	}
	return dst
}

func TestSearchFiles(t *testing.T) {
	dir := t.TempDir()

	copySample(t, dir, "a.jpg")
	copySample(t, dir, "b.jpg")
	if err := os.WriteFile(filepath.Join(dir, "c.txt"), []byte("not media"), 0o644); err != nil {
		t.Fatalf("write file error: %+v", err)
	}

	sub := filepath.Join(dir, "sub")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatalf("mkdir error: %+v", err)
	}
	copySample(t, sub, "d.jpg")

	files, err := contentimport.SearchFiles(dir, []string{"*.jpg"})
	if err != nil {
		t.Fatalf("search files error: %+v", err)
	}
	if len(files) != 3 {
		t.Fatalf("search files count: got %d, want 3 [%v]", len(files), files)
	}
}

func TestSearchFilesNoMatch(t *testing.T) {
	dir := t.TempDir()
	copySample(t, dir, "a.jpg")

	files, err := contentimport.SearchFiles(dir, []string{"*.mp4"})
	if err != nil {
		t.Fatalf("search files error: %+v", err)
	}
	if len(files) != 0 {
		t.Fatalf("search files count: got %d, want 0", len(files))
	}
}

func TestRegisterFileImage(t *testing.T) {
	setupDB(t)

	gid, err := db.RegisterGroup("g1", "")
	if err != nil {
		t.Fatalf("register group error: %+v", err)
	}

	dir := t.TempDir()
	f := copySample(t, dir, "photo.jpg")

	if err := contentimport.RegisterFile(gid, f); err != nil {
		t.Fatalf("register file error: %+v", err)
	}

	contents, err := db.SelectContent(gid)
	if err != nil {
		t.Fatalf("select content error: %+v", err)
	}
	if len(contents) != 1 {
		t.Fatalf("content count: got %d, want 1", len(contents))
	}

	c := contents[0]
	if c.Type != "img" {
		t.Errorf("content type: got %q, want img", c.Type)
	}
	if c.Width <= 0 || c.Height <= 0 {
		t.Errorf("content dimensions: got %dx%d", c.Width, c.Height)
	}
	if c.Path != f {
		t.Errorf("content path: got %q, want %q", c.Path, f)
	}

	ths, err := db.SelectContentThumbnails(c.ID)
	if err != nil {
		t.Fatalf("select thumbnails error: %+v", err)
	}
	if len(ths) != 17 {
		t.Errorf("thumbnail count: got %d, want 17", len(ths))
	}
}

func TestImportDirectory(t *testing.T) {
	setupDB(t)

	gid, err := db.RegisterGroup("g1", "")
	if err != nil {
		t.Fatalf("register group error: %+v", err)
	}

	dir := t.TempDir()
	copySample(t, dir, "a.jpg")
	copySample(t, dir, "b.jpg")

	count, err := contentimport.ImportDirectory(gid, dir, []string{"*.jpg"})
	if err != nil {
		t.Fatalf("import directory error: %+v", err)
	}
	if count != 2 {
		t.Fatalf("import count: got %d, want 2", count)
	}

	contents, err := db.SelectContent(gid)
	if err != nil {
		t.Fatalf("select content error: %+v", err)
	}
	if len(contents) != 2 {
		t.Errorf("content count: got %d, want 2", len(contents))
	}

	g, err := db.FindGroup(gid)
	if err != nil {
		t.Fatalf("find group error: %+v", err)
	}
	if g.Path != dir {
		t.Errorf("group path not updated: got %q, want %q", g.Path, dir)
	}
}

func TestImportDirectoryNoFiles(t *testing.T) {
	setupDB(t)

	gid, err := db.RegisterGroup("g1", "")
	if err != nil {
		t.Fatalf("register group error: %+v", err)
	}

	dir := t.TempDir()
	if _, err := contentimport.ImportDirectory(gid, dir, []string{"*.jpg"}); err == nil {
		t.Errorf("expected error for empty directory")
	}
}

func TestImportDirectoryGroupNotFound(t *testing.T) {
	setupDB(t)

	dir := t.TempDir()
	copySample(t, dir, "a.jpg")

	if _, err := contentimport.ImportDirectory(999, dir, []string{"*.jpg"}); err == nil {
		t.Errorf("expected error for missing group")
	}
}

func TestRegisterGenerated(t *testing.T) {
	setupDB(t)

	gid, err := db.RegisterGroup("g1", "")
	if err != nil {
		t.Fatalf("register group error: %+v", err)
	}

	c, err := contentimport.RegisterGenerated(gid, "NewYear", "cd", `{"text":"HELLO"}`)
	if err != nil {
		t.Fatalf("register generated error: %+v", err)
	}

	if c.ID <= 0 {
		t.Fatalf("content not saved: %+v", c)
	}
	if c.Type != "cd" {
		t.Errorf("content type: got %q, want cd", c.Type)
	}
	if c.Params != `{"text":"HELLO"}` {
		t.Errorf("content params: got %q", c.Params)
	}
	if c.Path != "" {
		t.Errorf("generated content should have no path, got %q", c.Path)
	}
	if c.Width <= 0 || c.Height <= 0 {
		t.Errorf("content dimensions: got %dx%d", c.Width, c.Height)
	}

	ths, err := db.SelectContentThumbnails(c.ID)
	if err != nil {
		t.Fatalf("select thumbnails error: %+v", err)
	}
	if len(ths) != 17 {
		t.Errorf("thumbnail count: got %d, want 17", len(ths))
	}
}

func TestRegisterGeneratedNormalizesLegacyType(t *testing.T) {
	setupDB(t)

	gid, err := db.RegisterGroup("g1", "")
	if err != nil {
		t.Fatalf("register group error: %+v", err)
	}

	// "countdown" is the legacy alias absorbed by video.Normalize into "cd"
	c, err := contentimport.RegisterGenerated(gid, "Legacy", "countdown", `{"text":"HI"}`)
	if err != nil {
		t.Fatalf("register generated error: %+v", err)
	}
	if c.Type != "cd" {
		t.Errorf("content type not normalized: got %q, want cd", c.Type)
	}
}

func TestRegisterGeneratedUnknownType(t *testing.T) {
	setupDB(t)

	gid, err := db.RegisterGroup("g1", "")
	if err != nil {
		t.Fatalf("register group error: %+v", err)
	}

	if _, err := contentimport.RegisterGenerated(gid, "x", "nosuchtype", ""); err == nil {
		t.Errorf("expected error for unknown type")
	}
}

func TestRegisterGeneratedInvalidTarget(t *testing.T) {
	setupDB(t)

	gid, err := db.RegisterGroup("g1", "")
	if err != nil {
		t.Fatalf("register group error: %+v", err)
	}

	// malformed target should fail plugin validation (video.Get -> cd.New)
	if _, err := contentimport.RegisterGenerated(gid, "x", "cd", `{"target":"not-a-date"}`); err == nil {
		t.Errorf("expected error for invalid target")
	}
}

func TestRegisterGeneratedGroupNotFound(t *testing.T) {
	setupDB(t)

	if _, err := contentimport.RegisterGenerated(999, "x", "terminal", `{"text":"hi"}`); err == nil {
		t.Errorf("expected error for missing group")
	}
}

func TestCheckMissing(t *testing.T) {
	setupDB(t)

	gid, err := db.RegisterGroup("g1", "")
	if err != nil {
		t.Fatalf("register group error: %+v", err)
	}

	dir := t.TempDir()
	existing := copySample(t, dir, "exists.jpg")
	missingPath := filepath.Join(dir, "gone.jpg")

	okContent := db.NewContent()
	okContent.GroupId = gid
	okContent.Name = "ok"
	okContent.Type = "img"
	okContent.Path = existing
	if _, errs := okContent.Save(); errs != nil {
		t.Fatalf("save content error: %+v", errs)
	}

	missingContent := db.NewContent()
	missingContent.GroupId = gid
	missingContent.Name = "missing"
	missingContent.Type = "img"
	missingContent.Path = missingPath
	if _, errs := missingContent.Save(); errs != nil {
		t.Fatalf("save content error: %+v", errs)
	}

	missing, err := contentimport.CheckMissing(gid)
	if err != nil {
		t.Fatalf("check missing error: %+v", err)
	}
	if len(missing) != 1 {
		t.Fatalf("missing count: got %d, want 1", len(missing))
	}
	if missing[0].ID != missingContent.ID {
		t.Errorf("missing content id: got %d, want %d", missing[0].ID, missingContent.ID)
	}
}

func TestCheckDuplicates(t *testing.T) {
	setupDB(t)

	gid, err := db.RegisterGroup("g1", "")
	if err != nil {
		t.Fatalf("register group error: %+v", err)
	}

	dupPath := "/tmp/dup.mp4"
	for i := 0; i < 2; i++ {
		c := db.NewContent()
		c.GroupId = gid
		c.Name = "dup"
		c.Type = "file"
		c.Path = dupPath
		if _, errs := c.Save(); errs != nil {
			t.Fatalf("save content error: %+v", errs)
		}
	}

	unique := db.NewContent()
	unique.GroupId = gid
	unique.Name = "unique"
	unique.Type = "file"
	unique.Path = "/tmp/unique.mp4"
	if _, errs := unique.Save(); errs != nil {
		t.Fatalf("save content error: %+v", errs)
	}

	dups, err := contentimport.CheckDuplicates(gid)
	if err != nil {
		t.Fatalf("check duplicates error: %+v", err)
	}
	if len(dups) != 1 {
		t.Fatalf("duplicate group count: got %d, want 1", len(dups))
	}
	if len(dups[0]) != 2 {
		t.Errorf("duplicate group size: got %d, want 2", len(dups[0]))
	}
}

func TestCheckDuplicatesNoneFound(t *testing.T) {
	setupDB(t)

	gid, err := db.RegisterGroup("g1", "")
	if err != nil {
		t.Fatalf("register group error: %+v", err)
	}

	c := db.NewContent()
	c.GroupId = gid
	c.Name = "unique"
	c.Type = "file"
	c.Path = "/tmp/unique.mp4"
	if _, errs := c.Save(); errs != nil {
		t.Fatalf("save content error: %+v", errs)
	}

	dups, err := contentimport.CheckDuplicates(gid)
	if err != nil {
		t.Fatalf("check duplicates error: %+v", err)
	}
	if len(dups) != 0 {
		t.Errorf("duplicate group count: got %d, want 0", len(dups))
	}
}
