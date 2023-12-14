package userimage

import (
	"sync"
)

var (
	extbinary_lock sync.Mutex
	extbinaries    = make(map[string]*extBinary)
)

type ExtBinary interface {
	Call_Modify(input []byte, height, width uint32) ([]byte, error)
}
type extBinary struct {
	binary_name string
}

func GetExtBinary(binary_name string) (*extBinary, error) {
	extbinary_lock.Lock()
	defer extbinary_lock.Unlock()
	eb := extbinaries[binary_name]
	if eb != nil {
		return eb, nil
	}
	eb = &extBinary{binary_name: binary_name}
	extbinaries[binary_name] = eb
	return eb, nil
}

// given a frame
func (e *extBinary) Call_Modify(input []byte, height, width uint32) ([]byte, error) {
	return input, nil
}
