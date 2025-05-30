package archive

type File struct {
	Path string
	Part int
}

func NewFile(path string, part int) File {
	return File{
		Path: path,
		Part: part,
	}
}
