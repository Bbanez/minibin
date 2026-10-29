package utils

import (
	"os"
	"runtime"
	"strings"
)

type FS struct {
	isWin    bool
	slash    string
	BasePath string
	baseLen  int
}

func Cwd() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}

// @basePath If nil then use current working directory
func NewFS(basePath *[]string) FS {
	isWin := false
	slash := "/"
	if runtime.GOOS == "windows" {
		isWin = true
		slash = "\\"
	}
	bp, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if basePath != nil {
		bp = strings.Join(*basePath, slash)
	}
	fs := FS{
		isWin:    isWin,
		slash:    slash,
		BasePath: bp,
		baseLen:  len(strings.Split(bp, slash)),
	}
	fs.mkdirIfNotExists(fs.Path(""), false)
	return fs
}

func PathJoin(path ...string) string {
	slash := "/"
	if runtime.GOOS == "windows" {
		slash = "\\"
	}
	return strings.Join(path, slash)
}

func (fs *FS) Path(path ...string) string {
	p := fs.Join(path...)
	if strings.HasPrefix(p, "/") || strings.Contains(p, "\\") {
		return p
	}
	return fs.BasePath + fs.slash + p
}

func (fs FS) Join(path ...string) string {
	return strings.Join(path, fs.slash)
}

func (fs *FS) ToPath(path string) []string {
	parts := strings.Split(path, "/")
	relativePath := []string{}
	for i := range parts {
		part := parts[i]
		if part == "" {
			continue
		}
		relativePath = append(relativePath, part)
	}
	return relativePath
}

func (fs *FS) FileSize(path ...string) Result[int64] {
	fileInfo, err := os.Stat(fs.Path(path...))
	if err != nil {
		return Err[int64](err)
	}
	return Ok(fileInfo.Size())
}

func (fs *FS) Exists(path ...string) bool {
	_, err := os.Stat(fs.Path(path...))
	if err != nil {
		return false
	}
	return true
}

func (fs *FS) mkdirIfNotExists(path string, skipRoot bool) {
	parts := strings.Split(path, fs.slash)
	startAt := 0
	if skipRoot {
		startAt = fs.baseLen
	}
	for i := startAt; i < len(parts)-1; i++ {
		part := parts[i]
		if part == "" {
			continue
		}
		path = strings.Join(parts[:i+1], fs.slash)
		if !dirExists(path) {
			os.Mkdir(path, 0755)
		}
	}
}

func (fs *FS) Mkdir(path ...string) Result[bool] {
	aPath := fs.Path(path...)
	if dirExists(aPath) {
		return Ok(true)
	}
	err := os.MkdirAll(aPath, 0755)
	if err != nil {
		return Err[bool](err)
	}
	return Ok(true)
}

func (fs *FS) Write(data []byte, path ...string) Result[bool] {
	aPath := fs.Path(path...)
	fs.mkdirIfNotExists(aPath, true)
	err := os.WriteFile(aPath, data, 0644)
	if err != nil {
		return Err[bool](err)
	}
	return Ok(true)
}

func (fs *FS) Append(data []byte, path ...string) Result[bool] {
	aPath := fs.Path(path...)
	fs.mkdirIfNotExists(aPath, true)
	f, err := os.OpenFile(aPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return Err[bool](err)
	}
	defer f.Close()
	_, err = f.Write(data)
	if err != nil {
		return Err[bool](err)
	}
	return Ok(true)
}

func (fs *FS) Read(path ...string) Result[[]byte] {
	aPath := fs.Path(path...)
	data, err := os.ReadFile(aPath)
	if err != nil {
		return Err[[]byte](err)
	}
	return Ok(data)
}

func (fs *FS) ReadString(path ...string) Result[string] {
	result := fs.Read(path...)
	if result.Error != nil {
		return Err[string](result.Error)
	}
	return Ok(string(result.Value))
}

func (fs *FS) OpenFile(path ...string) Result[*os.File] {
	aPath := fs.Path(path...)
	file, err := os.Open(aPath)
	if err != nil {
		return Err[*os.File](err)
	}
	return Ok(file)
}

func (fs *FS) Delete(path ...string) Result[bool] {
	aPath := fs.Path(path...)
	err := os.Remove(aPath)
	if err != nil {
		return Err[bool](err)
	}
	return Ok(true)
}

type FileTreeNode struct {
	AbsPath string
	RelPath string
	DirPath string
}

func (fs *FS) FileTree(startLocation []string, currentLocation []string) (*Array[FileTreeNode], error) {
	startLoc := fs.Path(startLocation...)
	var currLoc string
	if len(currentLocation) == 0 {
		currLoc = ""
	} else {
		currLoc = fs.Join(currentLocation...) + fs.slash
	}
	output := Arr[FileTreeNode]()
	files, err := os.ReadDir(startLoc + fs.slash + currLoc)
	if err != nil {
		return &output, err
	}
	for _, file := range files {
		filePath := startLoc + fs.slash + currLoc + file.Name()
		if file.IsDir() {
			children, err := fs.FileTree(
				startLocation,
				append(currentLocation, file.Name()),
			)
			if err != nil {
				return &output, err
			}
			output.PushMany(children.Items...)
		} else {
			output.Push(FileTreeNode{
				AbsPath: filePath,
				RelPath: strings.TrimPrefix(
					strings.Replace(filePath, startLoc, "", 1),
					"/",
				),
				DirPath: currLoc,
			})
		}
	}
	return &output, nil
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return info.IsDir()
}

func listFiles(dir string) Result[[]string] {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return Err[[]string](err)
	}
	var files []string = []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			childFiles := listFiles(dir + "/" + entry.Name())
			if childFiles.Error != nil {
				return Err[[]string](childFiles.Error)
			}
			files = append(files, childFiles.Value...)
		} else {
			files = append(files, dir+"/"+entry.Name())
		}
	}
	return Ok(files)
}
