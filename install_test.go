package gopher_lua_lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDebug_Http(t *testing.T) {
	L := NewLuaStateWithLib()
	err := L.DoString(`
a, err = http:defaultClient():get("http://qq.com")
assert(err == nil)
assert(a.statusCode == 200)
`)
	assert.Nil(t, err)
}

func TestDebug_Strings(t *testing.T) {
	L := NewLuaStateWithLib()
	err := L.DoString(`
assert(strings:contains("aab", "a") == true)
assert(strings:contains("aab", "c") == false)
assert(strings:trimPrefix("aab", "aa") == "b")
`)
	assert.Nil(t, err)
}

func TestDebug_Regexp(t *testing.T) {
	L := NewLuaStateWithLib()
	err := L.DoString(`
a, err = regexp:compile("gola([a-z]+)g")
assert(err == nil)

substring = a:findString("golang lua example")
assert(substring == "golang")

match, err = regexp:matchString("^golang", "golang lua example")
assert(err == nil)
assert(match == true)

match, err = regexp:matchString("^golang$", "golang lua example")
assert(err == nil)
assert(match == false)
`)
	assert.Nil(t, err)
}

func TestDebug_IoUtil(t *testing.T) {
	L := NewLuaStateWithLib()
	err := L.DoString(`
file = "/tmp/test.file"
err = ioutil:writeFile(file, "some test data", 0777)
assert(err == nil)

data, err = ioutil:readFile(file)
assert(err == nil)
assert(strings:toString(data) == "some test data")

os.remove(file)
`)
	assert.Nil(t, err)
}

func TestDebug_Json(t *testing.T) {
	L := NewLuaStateWithLib()
	err := L.DoString(`
local json = require("json")
data = json.encode({1, 2, 3})
assert(strings:toString(data) == "[1,2,3]")
`)
	assert.Nil(t, err)
}

func TestDebug_Url(t *testing.T) {
	L := NewLuaStateWithLib()
	err := L.DoString(`
u, err = url:parse("http://qq.com")
assert(err == nil)
assert(u.scheme == "http")
`)
	assert.Nil(t, err)
}

func TestColly(t *testing.T) {
	L := NewLuaStateWithLib()
	err := L.DoString(`
	c = colly:new()
	c:onHTML("a[href]", function(e) 
		debug:print(e)
	end)

	c:onRequest(function(r) 
		debug:print(r)
	end)

	c:visit("http://go-colly.org/")
`)
	assert.Nil(t, err)
}

func TestResty(t *testing.T) {
	L := NewLuaStateWithLib()
	err := L.DoString(`
	r = resty:new()
	l = reflect:newM()
	c, err =r:newRequest():setResult(reflect:mPtr(l)):get("https://v2.jokeapi.dev/joke/Any")
	t = reflect:toLuaType(l)
	assert(err == nil)
	for k, v in pairs(t) do 
    	print(k, v)
  	end

	l = reflect:newM()
	c, err = resty:newRequestWithResult(l):get("https://v2.jokeapi.dev/joke/Any")
	t = reflect:toLuaType(l)
	assert(err == nil)
	for k, v in pairs(t) do 
    	print(k, v)
  	end
`)
	assert.Nil(t, err)
}

func TestExec(t *testing.T) {
	L := NewLuaStateWithLib()
	err := L.DoString(`
	stdout, stderr, err = exec:run("ls", "-l")
	debug:print(stdout)
	assert(err == nil)
`)
	assert.Nil(t, err)
}
