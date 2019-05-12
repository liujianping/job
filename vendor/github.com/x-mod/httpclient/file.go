package httpclient

/* TODO FILE CONTENT SUPPORT
type File struct {
	Filename  string
	Fieldname string
	Data      []byte
}

func NewFile(field string, filename string) (*File, error) {
	absFile, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	fn := filepath.Base(absFile)
	data, err := ioutil.ReadFile(absFile)
	if err != nil {
		return nil, err
	}
	return &File{
		Filename:  fn,
		Fieldname: field,
		Data:      data,
	}, nil
}

func NewFileByBytes(field string, filename string, data []byte) (*File, error) {
	fn := filepath.Base(filename)
	return &File{
		Filename:  fn,
		Fieldname: field,
		Data:      data,
	}, nil
}

func NewFileByReader(field string, filename string, rd io.Reader) (*File, error) {
	fn := filepath.Base(filename)
	data, err := ioutil.ReadAll(rd)
	if err != nil {
		return nil, err
	}
	return &File{
		Filename:  fn,
		Fieldname: field,
		Data:      data,
	}, nil
}
*/
