package gopher_lua_lib

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gocolly/colly"
	"github.com/kohkimakimoto/gluayaml"
	luaJson "github.com/layeh/gopher-json"
	mysql "github.com/tengattack/gluasql/mysql"
	sqlite3 "github.com/tengattack/gluasql/sqlite3"
	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
)

func NewLuaStateWithLib(opts ...lua.Options) *lua.LState {
	L := lua.NewState(opts...)
	InstallAll(L)
	return L
}

func InstallAll(L *lua.LState) {
	L.SetGlobal("reflect", luar.New(L, &Reflect{}))
	L.SetGlobal("debug", luar.New(L, &Debug{}))
	L.SetGlobal("http", luar.New(L, &Http{}))
	L.SetGlobal("strings", luar.New(L, &Strings{}))
	L.SetGlobal("regexp", luar.New(L, &Regexp{}))
	L.SetGlobal("ioutil", luar.New(L, &IoUtil{}))
	L.SetGlobal("url", luar.New(L, &Url{}))
	L.SetGlobal("exec", luar.New(L, &Exec{}))
	L.SetGlobal("time", luar.New(L, &Time{}))
	L.SetGlobal("resty", luar.New(L, &Resty{}))
	L.SetGlobal("colly", luar.New(L, &Colly{}))
	L.PreloadModule("json", luaJson.Loader)
	L.PreloadModule("yaml", gluayaml.Loader)
	L.PreloadModule("crypto", CryptoLoader)
	L.PreloadModule("mysql", mysql.Loader)
	L.PreloadModule("sqlite3", sqlite3.Loader)

	InstallA(L)

	// more... pls refer to: github.com/vadv/gopher-lua-libs
}

type Reflect struct{}
type M = map[string]interface{}
type E = []M

func (I Reflect) NewE() E                            { return E{} }
func (I Reflect) EPtr(a E) *E                        { return &a }
func (I Reflect) NewM() M                            { return M{} }
func (I Reflect) MPtr(a M) *M                        { return &a }
func (I Reflect) New() interface{}                   { var value interface{}; return value }
func (I Reflect) Indirect(v interface{}) interface{} { return reflect.ValueOf(v).Elem().Interface() }
func (I Reflect) ToLuaType(L *luar.LState) int {
	v := L.CheckUserData(1).Value
	if reflect.TypeOf(v).Kind() == reflect.Ptr {
		v = reflect.ValueOf(v).Elem().Interface()
	}
	ud := luaJson.DecodeValue(L.LState, v)
	L.Push(ud)
	return 1
}
func (I Reflect) NewCtx() context.Context { return context.Background() }
func (I Reflect) NewCtxTimeout(millisecond int) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(millisecond))
	return ctx
}

// Http ----------------------------------------------------------------------------------------------------------------
type Http struct{}

func (h Http) DefaultClient() *http.Client { return http.DefaultClient }

// Strings -------------------------------------------------------------------------------------------------------------
type Strings struct{}

func (s Strings) ToByte(a string) []byte                   { return []byte(a) }
func (s Strings) ToString(a []byte) string                 { return string(a) }
func (s Strings) Contains(a, b string) bool                { return strings.Contains(a, b) }
func (s Strings) ContainsAny(a, chars string) bool         { return strings.ContainsAny(a, chars) }
func (s Strings) LastIndex(a, substr string) int           { return strings.LastIndex(a, substr) }
func (s Strings) SplitN(a, sep string, n int) []string     { return strings.SplitN(a, sep, n) }
func (s Strings) Split(a, sep string) []string             { return strings.Split(a, sep) }
func (s Strings) Join(elems []string, sep string) string   { return strings.Join(elems, sep) }
func (s Strings) HasPrefix(a, prefix string) bool          { return strings.HasPrefix(a, prefix) }
func (s Strings) HasSuffix(a, suffix string) bool          { return strings.HasSuffix(a, suffix) }
func (s Strings) Repeat(a string, count int) string        { return strings.Repeat(a, count) }
func (s Strings) ToUpper(a string) string                  { return strings.ToUpper(a) }
func (s Strings) ToLower(a string) string                  { return strings.ToLower(a) }
func (s Strings) Trim(a, cutset string) string             { return strings.Trim(a, cutset) }
func (s Strings) TrimLeft(a, cutset string) string         { return strings.TrimLeft(a, cutset) }
func (s Strings) TrimRight(a, cutset string) string        { return strings.TrimRight(a, cutset) }
func (s Strings) TrimSpace(a string) string                { return strings.TrimSpace(a) }
func (s Strings) TrimPrefix(a, prefix string) string       { return strings.TrimPrefix(a, prefix) }
func (s Strings) TrimSuffix(a, suffix string) string       { return strings.TrimSuffix(a, suffix) }
func (s Strings) Replace(a, old, new string, n int) string { return strings.Replace(a, old, new, n) }
func (s Strings) ReplaceAll(a, old, new string) string     { return strings.ReplaceAll(a, old, new) }
func (s Strings) EqualFold(a, t string) bool               { return strings.EqualFold(a, t) }
func (s Strings) Index(a, substr string) int               { return strings.Index(a, substr) }

// Regex ---------------------------------------------------------------------------------------------------------------
type Regexp struct{}

func (r Regexp) Compile(a string) (*regexp.Regexp, error)     { return regexp.Compile(a) }
func (r Regexp) Match(a string, b []byte) (bool, error)       { return regexp.Match(a, b) }
func (r Regexp) MatchString(a string, b string) (bool, error) { return regexp.MatchString(a, b) }

// IoUtils -------------------------------------------------------------------------------------------------------------
type IoUtil struct{}

func (i IoUtil) ReadAll(a io.Reader) ([]byte, error)               { return ioutil.ReadAll(a) }
func (i IoUtil) ReadFile(filename string) ([]byte, error)          { return os.ReadFile(filename) }
func (i IoUtil) WriteFile(a string, b []byte, c fs.FileMode) error { return os.WriteFile(a, b, c) }
func (i IoUtil) ReadDir(dirname string) ([]fs.FileInfo, error)     { return ioutil.ReadDir(dirname) }

type Url struct{}

func (u Url) Parse(a string) (*url.URL, error) { return url.Parse(a) }

// Cmd -----------------------------------------------------------------------------------------------------------------
type Exec struct{}

func (c Exec) Cmd(a ...string) *exec.Cmd { return exec.Command(a[0], a[1:]...) }
func (c Exec) Run(a ...string) ([]byte, []byte, error) {
	out, err := bytes.NewBuffer(nil), bytes.NewBuffer(nil)
	cmd := exec.Command(a[0], a[1:]...)
	cmd.Stdout, cmd.Stderr = out, err
	err1 := cmd.Run()
	return out.Bytes(), err.Bytes(), err1
}

type Time struct{}

func (t Time) Parse(a, b string) (time.Time, error) { return time.Parse(a, b) }
func (t Time) Now() time.Time                       { return time.Now() }

// Colly ---------------------------------------------------------------------------------------------------------------
type Colly struct{}

func (c Colly) New() *colly.Collector { return colly.NewCollector() }

// Resty ---------------------------------------------------------------------------------------------------------------
type Resty struct{}

func (c Resty) New() *resty.Client { r := resty.New(); return r }
func (c Resty) NewRequestWithResult(a M) *resty.Request {
	r := resty.New().NewRequest().SetResult(&a)
	return r
}

// Debug ---------------------------------------------------------------------------------------------------------------
type Debug struct{}

func (d Debug) Print(a interface{}) {
	a = convertMapIfSlice(a)
	if b, ok := a.(interface{ String() string }); ok {
		fmt.Println(b)
	} else if b, ok := a.([]byte); ok {
		fmt.Println(string(b))
	} else {
		fmt.Printf("%+v\n", a)
	}
}

func convertMapIfSlice(a interface{}) interface{} {
	if b, ok := a.(map[interface{}]interface{}); ok {
		var c []interface{}
		for i := 1; i < len(b)+1; i++ {
			if v, ok := b[float64(i)]; ok {
				c = append(c, v)
			} else {
				return a
			}
		}
		return c
	}
	return a
}

var defaultLUseless = lua.NewState()
