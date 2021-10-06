package gopher_lua_lib

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"hash"
	"hash/crc32"

	lua "github.com/yuin/gopher-lua"
)

var cryptoExports = map[string]lua.LGFunction{
	"base64_encode": base64EncodeFn,
	"base64_decode": base64DecodeFn,
	"crc32":         crc32Fn,
	"md5":           md5Fn,
	"sha1":          sha1Fn,
	"sha256":        sha256Fn,
	"sha512":        sha512Fn,
	"hmac":          hmacFn,
	"encrypt":       encryptFn,
	"decrypt":       decryptFn,
}

func CryptoLoader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), cryptoExports)
	L.Push(mod)

	L.SetField(mod, "_DEBUG", lua.LBool(false))
	L.SetField(mod, "_VERSION", lua.LString("0.0.0"))

	// consts
	L.SetField(mod, "RAW_DATA", lua.LNumber(1))
	L.SetField(mod, "ZERO_PADDING", lua.LNumber(2))

	return 1
}

// from https://github.com/tengattack/gluacrypto
func base64EncodeFn(L *lua.LState) int {
	s := lua.LVAsString(L.Get(1))
	result := base64.StdEncoding.EncodeToString([]byte(s))
	L.Push(lua.LString(result))
	return 1
}

func base64DecodeFn(L *lua.LState) int {
	s := lua.LVAsString(L.Get(1))
	result, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	L.Push(lua.LString(result))
	return 1
}

func crc32Fn(L *lua.LState) int {
	h := crc32.NewIEEE()
	s := lua.LVAsString(L.Get(1))
	raw := lua.LVAsBool(L.Get(2))
	_, err := h.Write([]byte(s))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	var result string
	if !raw {
		result = hex.EncodeToString(h.Sum(nil))
	} else {
		result = string(h.Sum(nil))
	}
	L.Push(lua.LString(result))
	return 1
}

// PKCS5Unpadding unpad data
func PKCS5Unpadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// Decrypt data by specified method: `des-ecb`, `des-cbc`, `aes-cbc`
func Decrypt(data []byte, method string, key, iv []byte) ([]byte, error) {
	var out []byte
	switch method {
	case "des-ecb":
		block, err := des.NewCipher([]byte(key))
		if err != nil {
			return nil, err
		}

		bs := block.BlockSize()
		if len(data)%bs != 0 {
			return nil, errors.New("crypto/cipher: input not full blocks")
		}

		out = make([]byte, len(data))
		dst := out
		for len(data) > 0 {
			block.Decrypt(dst, data[:bs])
			data = data[bs:]
			dst = dst[bs:]
		}
		out = PKCS5Unpadding(out)
	case "des-cbc":
		block, err := des.NewCipher([]byte(key))
		if err != nil {
			return nil, err
		}

		// CBC mode always works in whole blocks.
		if len(data)%block.BlockSize() != 0 {
			return nil, ErrCiphertextNotMultipleBlockSize
		}

		mode := cipher.NewCBCDecrypter(block, []byte(iv))
		plaintext := make([]byte, len(data))
		mode.CryptBlocks(plaintext, data)
		out = PKCS5Unpadding(plaintext)
	case "aes-cbc":
		block, err := aes.NewCipher([]byte(key))
		if err != nil {
			return nil, err
		}

		// CBC mode always works in whole blocks.
		if len(data)%block.BlockSize() != 0 {
			return nil, ErrCiphertextNotMultipleBlockSize
		}

		mode := cipher.NewCBCDecrypter(block, []byte(iv))
		plaintext := make([]byte, len(data))
		mode.CryptBlocks(plaintext, data)
		out = PKCS5Unpadding(plaintext)
	default:
		return nil, ErrNotSupport
	}
	return out, nil
}

func decryptFn(L *lua.LState) int {
	s := lua.LVAsString(L.Get(1))
	method := lua.LVAsString(L.Get(2))
	key := lua.LVAsString(L.Get(3))
	options := L.ToInt(4)
	iv := lua.LVAsString(L.Get(5))

	var data []byte
	var err error
	if options&RawData == 0 {
		data, err = hex.DecodeString(s)
	} else {
		data = []byte(s)
	}
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	out, err := Decrypt(data, method, []byte(key), []byte(iv))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(out))
	return 1
}

// options const
const (
	RawData = 1 << iota
	// not implement
	// ZeroPadding = 1 << iota
)

// errors
var (
	ErrNotSupport                     = errors.New("unsupported encrypt method")
	ErrCiphertextNotMultipleBlockSize = errors.New("ciphertext is not a multiple of the block size")
)

// PKCS5Padding pad data
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// Encrypt data by specified method: `des-ecb`, `des-cbc`, `aes-cbc`
func Encrypt(data []byte, method string, key, iv []byte) ([]byte, error) {
	var out []byte
	switch method {
	case "des-ecb":
		block, err := des.NewCipher(key)
		if err != nil {
			return nil, err
		}

		bs := block.BlockSize()
		data := PKCS5Padding(data, bs)
		out = make([]byte, len(data))

		dst := out
		for len(data) > 0 {
			// The message is divided into blocks,
			// and each block is encrypted separately.
			block.Encrypt(dst, data[:bs])
			data = data[bs:]
			dst = dst[bs:]
		}
	case "des-cbc":
		block, err := des.NewCipher(key)
		if err != nil {
			return nil, err
		}

		data := PKCS5Padding(data, block.BlockSize())
		mode := cipher.NewCBCEncrypter(block, iv)
		out = make([]byte, len(data))
		mode.CryptBlocks(out, data)
	case "aes-cbc":
		block, err := aes.NewCipher(key)
		if err != nil {
			return nil, err
		}

		data := PKCS5Padding(data, block.BlockSize())
		mode := cipher.NewCBCEncrypter(block, iv)
		out = make([]byte, len(data))
		mode.CryptBlocks(out, data)
	default:
		return nil, ErrNotSupport
	}
	return out, nil
}

func encryptFn(L *lua.LState) int {
	s := lua.LVAsString(L.Get(1))
	method := lua.LVAsString(L.Get(2))
	key := lua.LVAsString(L.Get(3))
	options := L.ToInt(4)
	iv := lua.LVAsString(L.Get(5))

	out, err := Encrypt([]byte(s), method, []byte(key), []byte(iv))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	var result string
	if options&RawData == 0 {
		result = hex.EncodeToString(out)
	} else {
		result = string(out)
	}
	L.Push(lua.LString(result))
	return 1
}

func hmacFn(L *lua.LState) int {
	algorithm := lua.LVAsString(L.Get(1))
	s := lua.LVAsString(L.Get(2))
	key := lua.LVAsString(L.Get(3))
	raw := lua.LVAsBool(L.Get(4))

	var h hash.Hash
	switch algorithm {
	case "md5":
		h = hmac.New(md5.New, []byte(key))
	case "sha1":
		h = hmac.New(sha1.New, []byte(key))
	case "sha256":
		h = hmac.New(sha256.New, []byte(key))
	case "sha512":
		h = hmac.New(sha512.New, []byte(key))
	default:
		L.Push(lua.LNil)
		L.Push(lua.LString("unsupported algorithm"))
		return 2
	}

	_, err := h.Write([]byte(s))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	var result string
	if !raw {
		result = hex.EncodeToString(h.Sum(nil))
	} else {
		result = string(h.Sum(nil))
	}
	L.Push(lua.LString(result))
	return 1
}

func md5Fn(L *lua.LState) int {
	h := md5.New()
	s := lua.LVAsString(L.Get(1))
	raw := lua.LVAsBool(L.Get(2))
	_, err := h.Write([]byte(s))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	var result string
	if !raw {
		result = hex.EncodeToString(h.Sum(nil))
	} else {
		result = string(h.Sum(nil))
	}
	L.Push(lua.LString(result))
	return 1
}

func sha1Fn(L *lua.LState) int {
	h := sha1.New()
	s := lua.LVAsString(L.Get(1))
	raw := lua.LVAsBool(L.Get(2))
	_, err := h.Write([]byte(s))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	var result string
	if !raw {
		result = hex.EncodeToString(h.Sum(nil))
	} else {
		result = string(h.Sum(nil))
	}
	L.Push(lua.LString(result))
	return 1
}

func sha256Fn(L *lua.LState) int {
	h := sha256.New()
	s := lua.LVAsString(L.Get(1))
	raw := lua.LVAsBool(L.Get(2))
	_, err := h.Write([]byte(s))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	var result string
	if !raw {
		result = hex.EncodeToString(h.Sum(nil))
	} else {
		result = string(h.Sum(nil))
	}
	L.Push(lua.LString(result))
	return 1
}

func sha512Fn(L *lua.LState) int {
	h := sha512.New()
	s := lua.LVAsString(L.Get(1))
	raw := lua.LVAsBool(L.Get(2))
	_, err := h.Write([]byte(s))
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	var result string
	if !raw {
		result = hex.EncodeToString(h.Sum(nil))
	} else {
		result = string(h.Sum(nil))
	}
	L.Push(lua.LString(result))
	return 1
}
