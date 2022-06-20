package token

import (
	"bytes"
	"encoding/binary"
	"github.com/clong1995/basic/cipher"
	"github.com/clong1995/basic/id"
	"time"
)

type Token struct {
	Id        int64
	session   int64
	timestamp int64
}

// Encode 编码
func (t *Token) Encode() string {
	//加入时间戳
	t.timestamp = time.Now().UnixNano()
	tsBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(tsBytes, uint64(t.timestamp))

	//加入id
	idBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(idBytes, uint64(t.Id))

	//加入session
	t.session = id.SId.Int()
	sessionBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(sessionBytes, uint64(t.session))

	//合并
	tsBytes = append(tsBytes, idBytes...)
	tsBytes = append(tsBytes, sessionBytes...)

	//返回string
	return cipher.Base64EncryptBytes(tsBytes)
}

func (t *Token) AccessKeyID() string {
	return cipher.Base64EncryptInt64((t.timestamp + t.Id) * 2)
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
	t.timestamp = int64(binary.LittleEndian.Uint64(b))

	//提取id
	_, err = buff.Read(b)
	if err != nil {
		return
	}
	t.Id = int64(binary.LittleEndian.Uint64(b))

	//提取session
	_, err = buff.Read(b)
	if err != nil {
		return
	}
	t.session = int64(binary.LittleEndian.Uint64(b))
	return
}

func (t *Token) Session() int64 {
	return t.session
}

func (t *Token) Timestamp() int64 {
	return t.timestamp
}
