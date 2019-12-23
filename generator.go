package vfsgen

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	pathpkg "path"
	"sort"
	"strconv"
	"text/template"
	"time"

	"github.com/shurcooL/httpfs/vfsutil"
)

// Generate Go code that statically implements input filesystem,
// write the output to a file specified in opt.
func Generate(input http.FileSystem, opt Options) error {
	opt.fillMissing()

	// Use an in-memory buffer to generate the entire output.
	buf := new(bytes.Buffer)

	err := t.ExecuteTemplate(buf, "Header", opt)
	if err != nil {
		return err
	}

	var toc toc
	err = findAndWriteFiles(buf, input, &toc)
	if err != nil {
		return err
	}

	err = t.ExecuteTemplate(buf, "DirEntries", toc.dirs)
	if err != nil {
		return err
	}

	err = t.ExecuteTemplate(buf, "Trailer", toc)
	if err != nil {
		return err
	}

	// Write output file (all at once).
	fmt.Println("writing", opt.Filename)
	err = ioutil.WriteFile(opt.Filename, buf.Bytes(), 0644)
	return err
}

type toc struct {
	dirs []*dirInfo
}

// fileInfo is a definition of a file.
type fileInfo struct {
	Path    string
	Name    string
	ModTime time.Time
	Size    int64
}

// dirInfo is a definition of a directory.
type dirInfo struct {
	Path    string
	Name    string
	ModTime time.Time
	Entries []string
}

// findAndWriteFiles recursively finds all the file paths in the given directory tree.
// They are added to the given map as keys. Values will be safe function names
// for each file, which will be used when generating the output code.
func findAndWriteFiles(buf *bytes.Buffer, fs http.FileSystem, toc *toc) error {
	walkFn := func(path string, fi os.FileInfo, r io.ReadSeeker, err error) error {
		if err != nil {
			// Consider all errors reading the input filesystem as fatal.
			return err
		}

		switch fi.IsDir() {
		case false:
			file := &fileInfo{
				Path:    path,
				Name:    pathpkg.Base(path),
				ModTime: fi.ModTime().UTC(),
				Size:    fi.Size(),
			}

			// Write FileInfo.
			err = writeFileInfo(buf, file, r)
			if err != nil {
				return err
			}
		case true:
			entries, err := readDirPaths(fs, path)
			if err != nil {
				return err
			}

			dir := &dirInfo{
				Path:    path,
				Name:    pathpkg.Base(path),
				ModTime: fi.ModTime().UTC(),
				Entries: entries,
			}

			toc.dirs = append(toc.dirs, dir)

			// Write DirInfo.
			err = t.ExecuteTemplate(buf, "DirInfo", dir)
			if err != nil {
				return err
			}
		}

		return nil
	}

	err := vfsutil.WalkFiles(fs, "/", walkFn)
	return err
}

// readDirPaths reads the directory named by dirname and returns
// a sorted list of directory paths.
func readDirPaths(fs http.FileSystem, dirname string) ([]string, error) {
	fis, err := vfsutil.ReadDir(fs, dirname)
	if err != nil {
		return nil, err
	}
	paths := make([]string, len(fis))
	for i := range fis {
		paths[i] = pathpkg.Join(dirname, fis[i].Name())
	}
	sort.Strings(paths)
	return paths, nil
}

// Write FileInfo.
func writeFileInfo(w io.Writer, file *fileInfo, r io.Reader) error {
	err := t.ExecuteTemplate(w, "FileInfo-Before", file)
	if err != nil {
		return err
	}
	bw := &byteWriter{w: w}
	_, err = io.Copy(bw, r)
	if err != nil {
		return err
	}
	err = t.ExecuteTemplate(w, "FileInfo-After", file)
	return err
}

var t = template.Must(template.New("").Funcs(template.FuncMap{
	"quote": strconv.Quote,
	"comment": func(s string) (string, error) {
		var buf bytes.Buffer
		cw := &commentWriter{W: &buf}
		_, err := io.WriteString(cw, s)
		if err != nil {
			return "", err
		}
		err = cw.Close()
		return buf.String(), err
	},
}).Parse(`{{define "Header"}}// Code generated by vfsgen; DO NOT EDIT.

{{with .BuildTags}}// +build {{.}}

{{end}}package {{.PackageName}}

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	pathpkg "path"
	"time"
)

{{comment .VariableComment}}
var {{.VariableName}} = func() http.FileSystem {
	fs := vfsgen۰FS{
{{end}}


{{define "FileInfo-Before"}}		{{quote .Path}}: &vfsgen۰FileInfo{
			name:    {{quote .Name}},
			modTime: {{template "Time" .ModTime}},
			content: []byte("{{end}}{{define "FileInfo-After"}}"),
		},
{{end}}



{{define "DirInfo"}}		{{quote .Path}}: &vfsgen۰DirInfo{
			name:    {{quote .Name}},
			modTime: {{template "Time" .ModTime}},
		},
{{end}}



{{define "DirEntries"}}	}
{{range .}}{{if .Entries}}	fs[{{quote .Path}}].(*vfsgen۰DirInfo).entries = []os.FileInfo{{"{"}}{{range .Entries}}
		fs[{{quote .}}].(os.FileInfo),{{end}}
	}
{{end}}{{end}}
	return fs
}()
{{end}}



{{define "Trailer"}}
type vfsgen۰FS map[string]interface{}

func (fs vfsgen۰FS) Open(path string) (http.File, error) {
	path = pathpkg.Clean("/" + path)
	f, ok := fs[path]
	if !ok {
		return nil, &os.PathError{Op: "open", Path: path, Err: os.ErrNotExist}
	}

	switch f := f.(type) {{"{"}}
	case *vfsgen۰DirInfo:
		return &vfsgen۰Dir{
			vfsgen۰DirInfo: f,
		}, nil
	default:
		// This should never happen because we generate only the above types.
		panic(fmt.Sprintf("unexpected type %T", f))
	}
}

// vfsgen۰FileInfo is a static definition of a file.
type vfsgen۰FileInfo struct {
	name    string
	modTime time.Time
	content []byte
}

func (f *vfsgen۰FileInfo) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("cannot Readdir from file %s", f.name)
}
func (f *vfsgen۰FileInfo) Stat() (os.FileInfo, error) { return f, nil }

func (f *vfsgen۰FileInfo) NotWorthGzipCompressing() {}

func (f *vfsgen۰FileInfo) Name() string       { return f.name }
func (f *vfsgen۰FileInfo) Size() int64        { return int64(len(f.content)) }
func (f *vfsgen۰FileInfo) Mode() os.FileMode  { return 0444 }
func (f *vfsgen۰FileInfo) ModTime() time.Time { return f.modTime }
func (f *vfsgen۰FileInfo) IsDir() bool        { return false }
func (f *vfsgen۰FileInfo) Sys() interface{}   { return nil }

// vfsgen۰File is an opened file instance.
type vfsgen۰File struct {
	*vfsgen۰FileInfo
	*bytes.Reader
}

func (f *vfsgen۰File) Close() error {
	return nil
}

// vfsgen۰DirInfo is a static definition of a directory.
type vfsgen۰DirInfo struct {
	name    string
	modTime time.Time
	entries []os.FileInfo
}

func (d *vfsgen۰DirInfo) Read([]byte) (int, error) {
	return 0, fmt.Errorf("cannot Read from directory %s", d.name)
}
func (d *vfsgen۰DirInfo) Close() error               { return nil }
func (d *vfsgen۰DirInfo) Stat() (os.FileInfo, error) { return d, nil }

func (d *vfsgen۰DirInfo) Name() string       { return d.name }
func (d *vfsgen۰DirInfo) Size() int64        { return 0 }
func (d *vfsgen۰DirInfo) Mode() os.FileMode  { return 0755 | os.ModeDir }
func (d *vfsgen۰DirInfo) ModTime() time.Time { return d.modTime }
func (d *vfsgen۰DirInfo) IsDir() bool        { return true }
func (d *vfsgen۰DirInfo) Sys() interface{}   { return nil }

// vfsgen۰Dir is an opened dir instance.
type vfsgen۰Dir struct {
	*vfsgen۰DirInfo
	pos int // Position within entries for Seek and Readdir.
}

func (d *vfsgen۰Dir) Seek(offset int64, whence int) (int64, error) {
	if offset == 0 && whence == io.SeekStart {
		d.pos = 0
		return 0, nil
	}
	return 0, fmt.Errorf("unsupported Seek in directory %s", d.name)
}

func (d *vfsgen۰Dir) Readdir(count int) ([]os.FileInfo, error) {
	if d.pos >= len(d.entries) && count > 0 {
		return nil, io.EOF
	}
	if count <= 0 || count > len(d.entries)-d.pos {
		count = len(d.entries) - d.pos
	}
	e := d.entries[d.pos : d.pos+count]
	d.pos += count
	return e, nil
}
{{end}}



{{define "Time"}}
{{- if .IsZero -}}
	time.Time{}
{{- else -}}
	time.Date({{.Year}}, {{printf "%d" .Month}}, {{.Day}}, {{.Hour}}, {{.Minute}}, {{.Second}}, {{.Nanosecond}}, time.UTC)
{{- end -}}
{{end}}
`))
