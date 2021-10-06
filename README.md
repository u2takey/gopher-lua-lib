# gopher-lua-lib


## http 

```bash
➜ ./luacmd
> c = http:defaultClient()
> c:get("http://baidu.com")
> resp, err = c:get("http://baidu.com")
> debug:print(err)
<nil>
> data, err = ioutil:readAll(resp.body)
> debug:print(data)
<html>
<meta http-equiv="refresh" content="0;url=http://www.baidu.com/">
</html>
```

## strings

```bash
> print(strings:contains("aab", "a"))
true
> print(strings:trimPrefix("aab", "aa"))
b
```

## regexp

```bash
> a, err = regexp:compile("gola([a-z]+)g")
> print(a:findString("golang lua example"))
golang
```

## ioutil

```bash
> file = "/tmp/test.file"
> err = ioutil:writeFile(file, "some test data", 0777)
> data, err = ioutil:readFile(file)
> debug:print(data)
some test data
> os.remove(file)
```

## time 

```bash
> print(time:now())
2021-10-07 17:37:40.1114 +0800 CST m=+183.728966880
```

## exec 

```bash
> stdout,stderr, err = exec:run("ls", "-l")
> debug:print(stdout)
total 77504
-rwxr-xr-x  1 leiwang  staff  39170616 10  7 17:41 luacmd
-rw-r--r--  1 leiwang  staff      3789 10  4 18:49 main.go
```

## crypto

```bash
> crypto = require("crypto")
> print(crypto.base64_encode("abcd"))
YWJjZA==
> print(crypto.base64_decode("YWJjZA=="))
abcd
```

## json/yaml

```bash
> json = require("json")
> debug:print(json.encode({1, 2, 3}))
[1,2,3]
```

## colly

```bash
>     c = colly:new()
>     c:onHTML("a[href]", function(e)
>>         debug:print(e)
>>     end)
>
>     c:onRequest(function(r)
>>         debug:print(r)
>>     end)
>
>     c:visit("http://go-colly.org/")
&{URL:http://go-colly.org/ Headers:0xc00096fc38 Ctx:0xc00071d210 Depth:1 Method:GET Body:<nil> ResponseCharacterEncoding: ID:1 collector:0xc00011f1e0 abort:false baseURL:<nil> ProxyURL:}
&{Name:a Text:Colly attributes:[{Namespace: Key:class Val:item} {Namespace: Key:href Val:http://go-colly.org/}] Request:0xc0005afd80 Response:0xc0009a1180 DOM:0xc000b41aa0 Index:0}
&{Name:a Text:Docs attributes:[{Namespace: Key:class Val:item} {Namespace: Key:href Val:/docs/}] Request:0xc0005afd80 Response:0xc0009a1180 DOM:0xc000b9a180 Index:1}
&{Name:a Text:Articles attributes:[{Namespace: Key:class Val:item} {Namespace: Key:href Val:/articles/}] Request:0xc0005afd80 Response:0xc0009a1180 DOM:0xc000b9a2d0 Index:2}
&{Name:a Text:Services attributes:[{Namespace: Key:class Val:item} {Namespace: Key:href Val:/services/}] Request:0xc0005afd80 Response:0xc0009a1180 DOM:0xc000b9a420 Index:3}
&{Name:a Text:Datasets attributes:[{Namespace: Key:class Val:item} {Namespace: Key:href Val:/datasets/}] Request:0xc0005afd80 Response:0xc0009a1180 DOM:0xc000b9a570 Index:4}
&{Name:a Text:GoDoc attributes:[{Namespace: Key:class Val:right item} {Namespace: Key:href Val:https://godoc.org/github.com/gocolly/colly} {Namespace: Key:target Val:_blank}] Request:0xc0005afd80 Response:0xc0009a1180 DOM:0xc000b9a6c0 Index:5}
&{Name:a Text: attributes:[{Namespace: Key:class Val:item} {Namespace: Key:href Val:https://github.com/gocolly/colly} {Namespace: Key:target Val:_blank}] Request:0xc0005afd80 Response:0xc0009a1180 DOM:0xc000b9a810 Index:6}
&{Name:a Text:Colly attributes:[{Namespace: Key:href Val:http://go-colly.org/} {Namespace: Key:class Val:brand item}] Request:0xc0005afd80 Response:0xc0009a1180 DOM:0xc000b9a960 Index:7}
...
...
```

## resty

```bash
> r = resty:new()
> l = reflect:newM()
> c, err =r:newRequest():setResult(reflect:mPtr(l)):get("https://v2.jokeapi.dev/joke/Any")
t = reflect:toLuaType(l)
> t = reflect:toLuaType(l)
> for k, v in pairs(t) do  print(k, v) end
error	false
category	Programming
setup	Why did the Python data scientist get arrested at customs?
delivery	She was caught trying to import pandas!
flags	table: 0xc0004db080
safe	true
type	twopart
id	234
lang	en
```

## more 

...
