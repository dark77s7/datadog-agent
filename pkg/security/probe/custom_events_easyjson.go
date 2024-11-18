//go:build linux
// +build linux

// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package probe

import (
	json "encoding/json"
	events "github.com/DataDog/datadog-agent/pkg/security/events"
	containerutils "github.com/DataDog/datadog-agent/pkg/security/secl/containerutils"
	serializers "github.com/DataDog/datadog-agent/pkg/security/serializers"
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

func easyjsonF8f9ddd1DecodeGithubComDataDogDatadogAgentPkgSecurityProbe(in *jlexer.Lexer, out *EventLostWrite) {
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
		case "map":
			out.Name = string(in.String())
		case "per_event":
			if in.IsNull() {
				in.Skip()
			} else {
				in.Delim('{')
				out.Lost = make(map[string]uint64)
				for !in.IsDelim('}') {
					key := string(in.String())
					in.WantColon()
					var v1 uint64
					v1 = uint64(in.Uint64())
					(out.Lost)[key] = v1
					in.WantComma()
				}
				in.Delim('}')
			}
		case "date":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Timestamp).UnmarshalJSON(data))
			}
		case "service":
			out.Service = string(in.String())
		case "container":
			if in.IsNull() {
				in.Skip()
				out.AgentContainerContext = nil
			} else {
				if out.AgentContainerContext == nil {
					out.AgentContainerContext = new(events.AgentContainerContext)
				}
				easyjsonF8f9ddd1DecodeGithubComDataDogDatadogAgentPkgSecurityEvents(in, out.AgentContainerContext)
			}
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
func easyjsonF8f9ddd1EncodeGithubComDataDogDatadogAgentPkgSecurityProbe(out *jwriter.Writer, in EventLostWrite) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"map\":"
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"per_event\":"
		out.RawString(prefix)
		if in.Lost == nil && (out.Flags&jwriter.NilMapAsEmpty) == 0 {
			out.RawString(`null`)
		} else {
			out.RawByte('{')
			v2First := true
			for v2Name, v2Value := range in.Lost {
				if v2First {
					v2First = false
				} else {
					out.RawByte(',')
				}
				out.String(string(v2Name))
				out.RawByte(':')
				out.Uint64(uint64(v2Value))
			}
			out.RawByte('}')
		}
	}
	{
		const prefix string = ",\"date\":"
		out.RawString(prefix)
		out.Raw((in.Timestamp).MarshalJSON())
	}
	{
		const prefix string = ",\"service\":"
		out.RawString(prefix)
		out.String(string(in.Service))
	}
	{
		const prefix string = ",\"container\":"
		out.RawString(prefix)
		if in.AgentContainerContext == nil {
			out.RawString("null")
		} else {
			easyjsonF8f9ddd1EncodeGithubComDataDogDatadogAgentPkgSecurityEvents(out, *in.AgentContainerContext)
		}
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EventLostWrite) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF8f9ddd1EncodeGithubComDataDogDatadogAgentPkgSecurityProbe(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EventLostWrite) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF8f9ddd1DecodeGithubComDataDogDatadogAgentPkgSecurityProbe(l, v)
}
func easyjsonF8f9ddd1DecodeGithubComDataDogDatadogAgentPkgSecurityEvents(in *jlexer.Lexer, out *events.AgentContainerContext) {
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
		case "id":
			out.ContainerID = containerutils.ContainerID(in.String())
		case "created_at":
			out.CreatedAt = uint64(in.Uint64())
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
func easyjsonF8f9ddd1EncodeGithubComDataDogDatadogAgentPkgSecurityEvents(out *jwriter.Writer, in events.AgentContainerContext) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ContainerID != "" {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.ContainerID))
	}
	{
		const prefix string = ",\"created_at\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint64(uint64(in.CreatedAt))
	}
	out.RawByte('}')
}
func easyjsonF8f9ddd1DecodeGithubComDataDogDatadogAgentPkgSecurityProbe1(in *jlexer.Lexer, out *EventLostRead) {
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
		case "map":
			out.Name = string(in.String())
		case "lost":
			out.Lost = float64(in.Float64())
		case "date":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Timestamp).UnmarshalJSON(data))
			}
		case "service":
			out.Service = string(in.String())
		case "container":
			if in.IsNull() {
				in.Skip()
				out.AgentContainerContext = nil
			} else {
				if out.AgentContainerContext == nil {
					out.AgentContainerContext = new(events.AgentContainerContext)
				}
				easyjsonF8f9ddd1DecodeGithubComDataDogDatadogAgentPkgSecurityEvents(in, out.AgentContainerContext)
			}
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
func easyjsonF8f9ddd1EncodeGithubComDataDogDatadogAgentPkgSecurityProbe1(out *jwriter.Writer, in EventLostRead) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"map\":"
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"lost\":"
		out.RawString(prefix)
		out.Float64(float64(in.Lost))
	}
	{
		const prefix string = ",\"date\":"
		out.RawString(prefix)
		out.Raw((in.Timestamp).MarshalJSON())
	}
	{
		const prefix string = ",\"service\":"
		out.RawString(prefix)
		out.String(string(in.Service))
	}
	{
		const prefix string = ",\"container\":"
		out.RawString(prefix)
		if in.AgentContainerContext == nil {
			out.RawString("null")
		} else {
			easyjsonF8f9ddd1EncodeGithubComDataDogDatadogAgentPkgSecurityEvents(out, *in.AgentContainerContext)
		}
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EventLostRead) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF8f9ddd1EncodeGithubComDataDogDatadogAgentPkgSecurityProbe1(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EventLostRead) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF8f9ddd1DecodeGithubComDataDogDatadogAgentPkgSecurityProbe1(l, v)
}
func easyjsonF8f9ddd1DecodeGithubComDataDogDatadogAgentPkgSecurityProbe2(in *jlexer.Lexer, out *EBPFLessHelloMsgEvent) {
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
		case "nsid":
			out.NSID = uint64(in.Uint64())
		case "workload_container":
			easyjsonF8f9ddd1Decode(in, &out.Container)
		case "args":
			if in.IsNull() {
				in.Skip()
				out.EntrypointArgs = nil
			} else {
				in.Delim('[')
				if out.EntrypointArgs == nil {
					if !in.IsDelim(']') {
						out.EntrypointArgs = make([]string, 0, 4)
					} else {
						out.EntrypointArgs = []string{}
					}
				} else {
					out.EntrypointArgs = (out.EntrypointArgs)[:0]
				}
				for !in.IsDelim(']') {
					var v3 string
					v3 = string(in.String())
					out.EntrypointArgs = append(out.EntrypointArgs, v3)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "date":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Timestamp).UnmarshalJSON(data))
			}
		case "service":
			out.Service = string(in.String())
		case "container":
			if in.IsNull() {
				in.Skip()
				out.AgentContainerContext = nil
			} else {
				if out.AgentContainerContext == nil {
					out.AgentContainerContext = new(events.AgentContainerContext)
				}
				easyjsonF8f9ddd1DecodeGithubComDataDogDatadogAgentPkgSecurityEvents(in, out.AgentContainerContext)
			}
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
func easyjsonF8f9ddd1EncodeGithubComDataDogDatadogAgentPkgSecurityProbe2(out *jwriter.Writer, in EBPFLessHelloMsgEvent) {
	out.RawByte('{')
	first := true
	_ = first
	if in.NSID != 0 {
		const prefix string = ",\"nsid\":"
		first = false
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.NSID))
	}
	if true {
		const prefix string = ",\"workload_container\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		easyjsonF8f9ddd1Encode(out, in.Container)
	}
	if len(in.EntrypointArgs) != 0 {
		const prefix string = ",\"args\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		{
			out.RawByte('[')
			for v4, v5 := range in.EntrypointArgs {
				if v4 > 0 {
					out.RawByte(',')
				}
				out.String(string(v5))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"date\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Raw((in.Timestamp).MarshalJSON())
	}
	{
		const prefix string = ",\"service\":"
		out.RawString(prefix)
		out.String(string(in.Service))
	}
	{
		const prefix string = ",\"container\":"
		out.RawString(prefix)
		if in.AgentContainerContext == nil {
			out.RawString("null")
		} else {
			easyjsonF8f9ddd1EncodeGithubComDataDogDatadogAgentPkgSecurityEvents(out, *in.AgentContainerContext)
		}
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EBPFLessHelloMsgEvent) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF8f9ddd1EncodeGithubComDataDogDatadogAgentPkgSecurityProbe2(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EBPFLessHelloMsgEvent) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF8f9ddd1DecodeGithubComDataDogDatadogAgentPkgSecurityProbe2(l, v)
}
func easyjsonF8f9ddd1Decode(in *jlexer.Lexer, out *struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	ImageShortName string `json:"short_name,omitempty"`
	ImageTag       string `json:"image_tag,omitempty"`
}) {
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
		case "id":
			out.ID = string(in.String())
		case "name":
			out.Name = string(in.String())
		case "short_name":
			out.ImageShortName = string(in.String())
		case "image_tag":
			out.ImageTag = string(in.String())
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
func easyjsonF8f9ddd1Encode(out *jwriter.Writer, in struct {
	ID             string `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	ImageShortName string `json:"short_name,omitempty"`
	ImageTag       string `json:"image_tag,omitempty"`
}) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != "" {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	if in.Name != "" {
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Name))
	}
	if in.ImageShortName != "" {
		const prefix string = ",\"short_name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.ImageShortName))
	}
	if in.ImageTag != "" {
		const prefix string = ",\"image_tag\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.ImageTag))
	}
	out.RawByte('}')
}
func easyjsonF8f9ddd1DecodeGithubComDataDogDatadogAgentPkgSecurityProbe3(in *jlexer.Lexer, out *AbnormalEvent) {
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
		case "triggering_event":
			if in.IsNull() {
				in.Skip()
				out.Event = nil
			} else {
				if out.Event == nil {
					out.Event = new(serializers.EventSerializer)
				}
				(*out.Event).UnmarshalEasyJSON(in)
			}
		case "error":
			out.Error = string(in.String())
		case "date":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Timestamp).UnmarshalJSON(data))
			}
		case "service":
			out.Service = string(in.String())
		case "container":
			if in.IsNull() {
				in.Skip()
				out.AgentContainerContext = nil
			} else {
				if out.AgentContainerContext == nil {
					out.AgentContainerContext = new(events.AgentContainerContext)
				}
				easyjsonF8f9ddd1DecodeGithubComDataDogDatadogAgentPkgSecurityEvents(in, out.AgentContainerContext)
			}
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
func easyjsonF8f9ddd1EncodeGithubComDataDogDatadogAgentPkgSecurityProbe3(out *jwriter.Writer, in AbnormalEvent) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"triggering_event\":"
		out.RawString(prefix[1:])
		if in.Event == nil {
			out.RawString("null")
		} else {
			(*in.Event).MarshalEasyJSON(out)
		}
	}
	{
		const prefix string = ",\"error\":"
		out.RawString(prefix)
		out.String(string(in.Error))
	}
	{
		const prefix string = ",\"date\":"
		out.RawString(prefix)
		out.Raw((in.Timestamp).MarshalJSON())
	}
	{
		const prefix string = ",\"service\":"
		out.RawString(prefix)
		out.String(string(in.Service))
	}
	{
		const prefix string = ",\"container\":"
		out.RawString(prefix)
		if in.AgentContainerContext == nil {
			out.RawString("null")
		} else {
			easyjsonF8f9ddd1EncodeGithubComDataDogDatadogAgentPkgSecurityEvents(out, *in.AgentContainerContext)
		}
	}
	out.RawByte('}')
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v AbnormalEvent) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonF8f9ddd1EncodeGithubComDataDogDatadogAgentPkgSecurityProbe3(w, v)
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *AbnormalEvent) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonF8f9ddd1DecodeGithubComDataDogDatadogAgentPkgSecurityProbe3(l, v)
}
