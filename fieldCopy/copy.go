package fieldCopy

import (
	"github.com/jinzhu/copier"
)

// Copy 从src深度拷贝实例到dc中
func Copy(dc, src interface{}) error {
	return copier.Copy(dc, src)
}
