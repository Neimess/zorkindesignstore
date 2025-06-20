// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package dto

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson4a0f95aaDecodeGithubComNeimessZorkinStoreProjectInternalTransportDto(in *jlexer.Lexer, out *TokenResponse) {
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
		case "token":
			out.Token = string(in.String())
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
func easyjson4a0f95aaEncodeGithubComNeimessZorkinStoreProjectInternalTransportDto(out *jwriter.Writer, in TokenResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"token\":"
		out.RawString(prefix[1:])
		out.String(string(in.Token))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v TokenResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson4a0f95aaEncodeGithubComNeimessZorkinStoreProjectInternalTransportDto(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v TokenResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson4a0f95aaEncodeGithubComNeimessZorkinStoreProjectInternalTransportDto(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *TokenResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson4a0f95aaDecodeGithubComNeimessZorkinStoreProjectInternalTransportDto(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *TokenResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson4a0f95aaDecodeGithubComNeimessZorkinStoreProjectInternalTransportDto(l, v)
}
