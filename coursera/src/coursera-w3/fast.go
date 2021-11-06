package main

import (
	"bufio"
	"bytes"
	json "encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// Info struct
type Info struct {
	Browsers [][]byte
	Email    string
	Name     string
}

// FastSearch вам надо написать более быструю оптимальную этой функции
func FastSearch(out io.Writer) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	// seenBrowsers := map[string]string{}
	// seenBrowsers := make([][]uint8, 0)
	// seenBrowsers := make(map[int][]byte, 10)
	seenBrowsers := []string{}
	uniqueBrowsers := 0
	foundUsers := bytes.Buffer{}

	androidBytes := []byte("Android")
	msieBytes := []byte("MSIE")

	var user Info
	scanner := bufio.NewScanner(file)
	i := 0

	for scanner.Scan() {
		err := user.UnmarshalJSON(scanner.Bytes())
		if err != nil {
			panic(err)
		}

		isAndroid := false
		isMSIE := false

		for _, browserRaw := range user.Browsers {
			sb := string(browserRaw)
			// if bytes.Contains(browserRaw, androidBytes) {
			if strings.Contains(sb, string(androidBytes)) {
				isAndroid = true
				notSeenBefore := true

				for _, item := range seenBrowsers {
					if item == sb {
						// if bytes.Equal(item, browserRaw) {
						notSeenBefore = false
					}
				}

				if notSeenBefore {
					// seenBrowsers[uniqueBrowsers] = browserRaw
					seenBrowsers = append(seenBrowsers, sb)
					uniqueBrowsers++
				}
			}

			// if bytes.Contains(browserRaw, msieBytes) {
			if strings.Contains(sb, string(msieBytes)) {
				isMSIE = true
				notSeenBefore := true

				for _, item := range seenBrowsers {
					if item == sb {
						// if bytes.Equal(item, browserRaw) {
						notSeenBefore = false
					}
				}

				if notSeenBefore {
					// seenBrowsers[uniqueBrowsers] = browserRaw
					seenBrowsers = append(seenBrowsers, sb)
					uniqueBrowsers++
				}
			}
		}

		if !(isAndroid && isMSIE) {
			i++
			continue
		}

		email := strings.ReplaceAll(user.Email, "@", " [at] ")
		s := fmt.Sprintf("[%d] %s <%s>\n", i, user.Name, email)
		foundUsers.Write([]byte(s))
		i++
	}

	fmt.Fprintln(out, "found users:\n"+foundUsers.String())
	fmt.Fprintln(out, "Total unique browsers", len(seenBrowsers))
}

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson9f2eff5fDecodeHomeAlexanderGoSrc(in *jlexer.Lexer, out *Info) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "browsers":
			if in.IsNull() {
				in.Skip()
				out.Browsers = nil
			} else {
				in.Delim('[')
				if out.Browsers == nil {
					if !in.IsDelim(']') {
						out.Browsers = make([][]uint8, 0, 2)
					} else {
						out.Browsers = [][]uint8{}
					}
				} else {
					out.Browsers = (out.Browsers)[:0]
				}
				for !in.IsDelim(']') {
					var v1 []uint8
					if in.IsNull() {
						in.Skip()
						v1 = nil
					} else {
						v1 = in.UnsafeBytes()
					}
					out.Browsers = append(out.Browsers, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "email":
			out.Email = string(in.String())
		case "name":
			out.Name = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson9f2eff5fEncodeHomeAlexanderGoSrc(out *jwriter.Writer, in Info) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"Browsers\":"
		out.RawString(prefix[1:])
		if in.Browsers == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Browsers {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"Email\":"
		out.RawString(prefix)
		out.String(string(in.Email))
	}
	{
		const prefix string = ",\"Name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Info) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson9f2eff5fEncodeHomeAlexanderGoSrc(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Info) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson9f2eff5fEncodeHomeAlexanderGoSrc(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Info) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson9f2eff5fDecodeHomeAlexanderGoSrc(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Info) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson9f2eff5fDecodeHomeAlexanderGoSrc(l, v)
}
