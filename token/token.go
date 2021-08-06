package token

import (
	"basic/cipher"
	"bytes"
	"encoding/binary"
)

type Token struct {
	Id        int64
	Timestamp int64
}

// Encode 编码
func (t Token) Encode() string {

	//加入时间戳
	tsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(tsBytes, uint64(t.Timestamp))

	//加入id
	idBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(idBytes, uint64(t.Id))

	//合并
	tsBytes = append(tsBytes, idBytes...)

	//返回string
	return cipher.Base64EncryptBytes(tsBytes)
}

func (t Token) AccessKeyID() string {
	return cipher.Base64EncryptInt64((t.Timestamp + t.Id) * 2)
}

// Decode 解码
func (t *Token) Decode(token string) (err error) {
	bs, err := cipher.Base64DecryptBytes(token)
	if err != nil {
		return
	}

	buff := bytes.NewBuffer(bs)
	b := make([]byte, 8)

	//提取时间戳
	_, err = buff.Read(b)
	if err != nil {
		return
	}
	t.Timestamp = int64(binary.LittleEndian.Uint64(b))

	//提取时间戳
	_, err = buff.Read(b)
	if err != nil {
		return
	}
	t.Id = int64(binary.LittleEndian.Uint64(b))
	return
}
