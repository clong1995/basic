package fieldCopy

import "github.com/ulule/deepcopier"

// FieldCopy 从src深度拷贝实例到dc中
func FieldCopy(src, dc interface{}) error {
	return deepcopier.Copy(src).To(dc)
}

// FieldFrom 从src深度拷贝实例到dc中
func FieldFrom(dc, src interface{}) error {
	return deepcopier.Copy(dc).From(src)
}
