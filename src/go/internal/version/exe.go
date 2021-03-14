package version

import (
	"bytes"
	"debug/elf"
	"debug/macho"
	"debug/pe"
	"fmt"
	"io"
	"os"
)

// exe 是 OS 可执行文件（ELF、Mach-O、PE、XCOFF）的通用接口。
// 用于兼容不同平台的可执行文件
type exe interface {
	Close() error
}

func openExe(file string) (exe, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	// 读取文件前16字节内容
	data := make([]byte, 16)
	if _, err := io.ReadFull(f, data); err != nil {
		return nil, err
	}
	// 重新设置偏移量为初始
	f.Seek(0, 0)
	// \x7FELF -> linux平台
	// 如果在最前面匹配到 \x7FELF
	if bytes.HasPrefix(data, []byte("\x7FELF")) {
		e, err := elf.NewFile(f)
		if err != nil {
			f.Close()
			return nil, err
		}
		return &elfExe{f, e}, nil
	}
	// DOS下可执行文件头部
	if bytes.HasPrefix(data, []byte("MZ")) {
		e, err := pe.NewFile(f)
		if err != nil {
			f.Close()
			return nil, err
		}
		return &peExe{f, e}, nil
	}
	// \xFE\xED\xFA
	// \xFA\xED\xFE
	// Mac平台
	if bytes.HasPrefix(data, []byte("\xFE\xED\xFA")) || bytes.HasPrefix(data[1:], []byte("\xFA\xED\xFE")) {
		e, err := macho.NewFile(f)
		if err != nil {
			f.Close()
			return nil, err
		}
		return &machoExe{f, e}, nil
	}
	// 0x01 0xDF
	// 0x01 0xF7
	// aix平台
	if bytes.HasPrefix(data, []byte{0x01, 0xDF}) || bytes.HasPrefix(data, []byte{0x01, 0xF7}) {

	}
	return nil, fmt.Errorf("unrecongnized executable format")
}

// linux平台可执行文件头部
type elfExe struct {
	os *os.File
	f  *elf.File
}

func (x *elfExe) Close() error {
	return x.os.Close()
}

// dos平台可执行头部
type peExe struct {
	os *os.File
	f  *pe.File
}

func (x *peExe) Close() error {
	return x.os.Close()
}

// Mac平台
type machoExe struct {
	os *os.File
	f  *macho.File
}

func (x *machoExe) Close() error {
	return x.os.Close()
}

// aix平台
type xcoffExe struct {
}
