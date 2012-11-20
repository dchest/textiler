package textiler

// Those tests are based on https://github.com/ikirudennis/python-textile/blob/master/textile/tests/__init__.py

import (
	"bytes"
	"fmt"
	"testing"
)

func textileToHtml(input string) string {
	return string(ToHtml([]byte(input), false, false))
}

func textileToXhtml(input string) string {
	return string(ToXhtml([]byte(input), false, false))
}

func TestUrlRef(t *testing.T) {
	data := []string{
		"[hobix]http://hobix.com", "hobix", "http://hobix.com",
		"[]http://hobix.com", "", "http://hobix.com",
	}
	for i := 0; i < len(data)/3; i++ {
		title, url := isUrlRef([]byte(data[i*3]))
		expectedTitle := data[i*3+1]
		expectedUrl := data[i*3+2]
		if !bytes.Equal(title, []byte(expectedTitle)) {
			t.Fatalf("\nExpected[%s]\nActual  [%s]", expectedTitle, string(title))
		}
		if !bytes.Equal(url, []byte(expectedUrl)) {
			t.Fatalf("\nExpected[%s]\nActual  [%s]", expectedUrl, string(url))
		}
	}
}

func TestParseSpan(t *testing.T) {
	data := []string{
		"%{color:red}%", "color:red;", "", "",
		"%{color:red;foo:bar}foo%", "color:red;foo:bar;", "foo", "",
		"%{color:red}inside%after", "color:red;", "inside", "after",
	}
	var expected string
	for i := 0; i < len(data)/4; i++ {
		rest, inside, attrs := parseSpan([]byte(data[i*4]))
		expected = data[i*4+1]
		if !bytes.Equal(attrs.style, []byte(expected)) {
			t.Fatalf("\nExpected[%s]\nActual  [%s]", expected, string(attrs.style))
		}
		expected = data[i*4+2]
		if !bytes.Equal(inside, []byte(expected)) {
			t.Fatalf("\nExpected[%s]\nActual  [%s]", expected, string(inside))
		}
		expected = data[i*4+3]
		if !bytes.Equal(rest, []byte(expected)) {
			t.Fatalf("\nExpected[%s]\nActual  [%s]", expected, string(rest))
		}
	}
}

func TestIsHLine(t *testing.T) {
	data := []string{
		"h1. foo", "1", "foo",
		"h0. bar", "", "",
		"h3.rest", "", "",
		"h3. rest", "3", "rest",
		"h6. loh", "6", "loh",
	}
	for i := 0; i < len(data)/3; i++ {
		rest, n, _ := parseH([]byte(data[i*3]))
		expectedN := data[i*3+1]
		if n < 0 {
			if expectedN != "" {
				t.Fatalf("\nExpected[%s]\nActual  [%d]", expectedN, n)
			}
		} else {
			if expectedN != fmt.Sprintf("%d", n) {
				t.Fatalf("\nExpected[%s]\nActual  [%d]", expectedN, n)
			}
			expectedRest := data[i*3+2]
			if !bytes.Equal(rest, []byte(expectedRest)) {
				t.Fatalf("\nExpected[%s]\nActual  [%s]", expectedRest, string(rest))
			}
		}
	}
}

func TestUrl(t *testing.T) {
	data := []string{
		`"Hobix":http://hobix.com/`, "Hobix", "http://hobix.com/", "",
		`"":http://foo end`, "", "http://foo", " end",
		`"foo":Bar tender`, "foo", "Bar", " tender",
	}
	for i := 0; i < len(data)/4; i++ {
		rest, title, url := parseUrlOrRefName([]byte(data[i*4]))
		titleExpected := data[i*4+1]
		urlExpected := data[i*4+2]
		restExpected := data[i*4+3]
		if !bytes.Equal(title, []byte(titleExpected)) {
			t.Fatalf("\nExpected1[%s]\nActual   [%s]", string(titleExpected), string(title))
		}
		if !bytes.Equal(url, []byte(urlExpected)) {
			t.Fatalf("\nExpected2[%s]\nActual   [%s]", string(urlExpected), string(url))
		}
		if !bytes.Equal(rest, []byte(restExpected)) {
			t.Fatalf("\nExpected3[%s]\nActual   [%s]", string(restExpected), string(rest))
		}
	}
}

func TestParseInline(t *testing.T) {
	data := []string{
		".-in me-", ".<del>in me</del>",
		"_TR()", "_TR()",
		"__f__", "<i>f</i>",
		"__b__", "<i>b</i>",
		"__r __", "__r __",
		"__r__rest", "__r__rest",
		"__r__ rest", "<i>r</i> rest",
		"fo-not. me-", "fo-not. me-",
		"before;__ol__", "before;<i>ol</i>",
		"foo:**bold**?is here", "foo:<b>bold</b>?is here",
		`"Hobix":http://hobix.com/`, `<a href="http://hobix.com/">Hobix</a>`,
		`!http://hobix.com/sample.jpg!`, `<img src="http://hobix.com/sample.jpg" alt="">`,
		`!openwindow1.gif(Bunny.)!`, `<img src="openwindow1.gif" title="Bunny." alt="Bunny.">`,
		`!openwindow1.gif!:http://hobix.com/`, `<a href="http://hobix.com/" class="img"><img src="openwindow1.gif" alt=""></a>`,
		`@p@`, "<code>p</code>",
		`before@foo@`, "before<code>foo</code>",
		`bef@bar@after`, "bef<code>bar</code>after",
	}
	for i := 0; i < len(data)/2; i++ {
		p := NewParser(0)
		s := data[i*2]
		p.parseInline([]byte(s))
		expected := []byte(data[i*2+1])
		actual := p.out.Bytes()
		if !bytes.Equal(expected, actual) {
			ToHtml([]byte(s), false, true)
			t.Fatalf("\nSrc:[%s]\nExp:[%s]\nGot:[%s]", s, string(expected), string(actual))
		}
	}
}

func TestItalic(t *testing.T) {
	italics := []string{
		"____", "", "",
		"__f__", "f", "",
		"__foo__o", "foo", "o",
		"__a_d___lo", "a_d", "_lo",
	}
	for i := 0; i < len(italics)/3; i++ {
		rest, inside := parseItalic([]byte(italics[i*3]))
		er1, er2 := []byte(italics[i*3+1]), []byte(italics[i*3+2])
		if !bytes.Equal(inside, er1) {
			t.Fatalf("\nExpected[%#v]\nActual  [%#v]", er1, inside)
		}
		if !bytes.Equal(rest, er2) {
			t.Fatalf("\nExpected[%#v]\nActual  [%#v]", er2, rest)
		}
	}
}

func TestTextileHtml(t *testing.T) {
	lastPassingTest := 12
	for i := 0; i <= lastPassingTest; i++ {
		s := HtmlTests[i*2]
		actual := textileToHtml(s)
		expected := HtmlTests[i*2+1]
		if actual != expected {
			ToHtml([]byte(s), false, true)
			t.Fatalf("\nSrc:%#v\n\nExp:%#v\n\nGot:%#v\n", s, expected, actual)
		}
	}
}

func TestOl(t *testing.T) {
	ols := []struct {
		Str   string
		Rest  string
		Level int
	}{
		{"# ", "", 1},
	}
	for _, test := range ols {
		res, level := parseListEl([]byte(test.Str), '#')
		if test.Rest != string(res) {
			t.Fatalf("\nSrc:%#v\n\nExp:%#v\n\nGot:%#v\n", test.Str, test.Rest, res)
		}
		if test.Level != level {
			t.Fatalf("\nSrc:%#v\n\nExp:%#v\n\nGot:%#v\n", test.Str, test.Level, level)
		}
	}
}

func TestTextileXhtml(t *testing.T) {
	// TODO: for now mark tests that we expect to pass explicitly
	// 4,5,6,7,8,9,10 - smartypants for '"'
	passingTests := []int{0, 1, 2, 3, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38,
		39, 40, 41, 42, 43, 45, 46, 48, 50, 51, 52, 53, 54, 62, 64, 66, 68, 69,
		70, 74, 75, 76, 78, 79}
	// 44, 47, 65, 94 - lists
	// 49 - "foo (title)":http://my.com - parsing (title) and serializing as title="" attribute
	// 55, 83 - use CSS(Acronyms) - parsing acronyms in ()
	// 56, 57, 58, 59, 88, 89, 95, 96, 97, 98 - tables
	// 60, 72 - pre..
	// 61 - <pre>foo</pre>
	// 63 - "foo ==(bar)==":#foobar
	// 67 - #{color:blue} one - style for lists
	// 71 - *:(foo)foo bar baz* - <cite> within '*' (strong)
	// 73 - leading spaces don't induce <p>
	// 77 - H[~2~]O - is supposed to drop [] for some reason
	// 80, 81, 82 - smart quotes
	// 84 - Textpattern CMS - auto-dected all caps (CMS) and put inside <span style="caps">CMS</span>
	// 85 - url-escape urls (Ü => %C3%9Cb)
	// 86 - <-- comments
	// 87 - (c) => &#169;, (r) => #174;, (tm) => &#8482
	// 90, 91, 99, 101 - dl, dt
	// 92 - *(class) - class in *
	// 93 - #_(first#list) foo - class/id for lists
	// 100 - ###. comment
	for _, i := range passingTests {
		s := XhtmlTests[i*2]
		actual := textileToXhtml(s)
		expected := XhtmlTests[i*2+1]
		if actual != expected {
			ToXhtml([]byte(s), false, true)
			t.Fatalf("\nSrc:%#v\n\nExp:%#v\n\nGot:%#v\n", s, expected, actual)
		}
	}
}
