package main

import (
	"encoding/xml"
	"os"
	"fmt"
	"strings"
	"bytes"
)

/**
 * WARNING: this code isn't very well documented or pretty, but it should be understandable
 */

var replacements map[string]string

type Interface struct {
	name string
	requests []RequestEvent
	events []RequestEvent
	enum []Enum
}

type RequestEvent struct {
	name string
	opcode int
	args []Arg
}

type Arg struct {
	name string
	t string
	iface string
}

type Enum struct {
	name string
	entries []Entry
}

type Entry struct {
	name string
	value int
}

func init() {
	replacements = map[string]string{
		"type":"typ",
		"interface":"iface",
	}
}

func main() {
	// Open protocol
	pr, err := os.Open("wayland.xml")
	if err != nil {
		fmt.Println(err)
		return
	}
	ifaces := make([]Interface,0)

	decoder := xml.NewDecoder(pr)
	for token,err := decoder.Token() ; err == nil ; token,err = decoder.Token() {
		if se,ok := token.(xml.StartElement) ; ok {
			if se.Name.Local == "interface" {
				iface := DecodeInterface(decoder, se.Attr)
				ifaces = append(ifaces, iface)
			}
		}
	}

	for _,val := range ifaces {
		content, path := PrintInterface(val)
		file, err := os.Create(path)
		if err != nil {
			fmt.Println(err)
			return
		}
		file.Write(content.Bytes())
	}
	return
}

func DecodeInterface(decoder *xml.Decoder, attr []xml.Attr) (iface Interface) {
	// Setup interface
	for _,val := range attr {
		if val.Name.Local == "name" {
			iface.name = val.Value
		}
	}

	// Parse requests, events and enums
	for token,err := decoder.Token() ; err == nil ; token,err = decoder.Token() {
		switch e := token.(type) {
		case xml.StartElement:
			switch e.Name.Local {
			case "request":
				r := DecodeRequest(decoder, e.Attr)
				iface.requests = append(iface.requests, r)
				break
			case "event":
				r := DecodeEvent(decoder, e.Attr)
				iface.events = append(iface.events, r)
				break
			}
			break
		// Handle end elements
		case xml.EndElement:
			if e.Name.Local == "interface"{
				return
			}
		}
	}

	return
}

func DecodeRequest(decoder *xml.Decoder, attr []xml.Attr) (req RequestEvent) {
	// get name
	for _,val := range attr {
		if val.Name.Local == "name" {
			req.name = val.Value
		}
	}

	// parse args
	for token,err := decoder.Token() ; err == nil ; token,err = decoder.Token() {
		switch e := token.(type) {
		case xml.StartElement:
			if e.Name.Local == "arg" {
				var arg Arg
				for _,val := range e.Attr {
					switch val.Name.Local {
					case "name":
						arg.name = val.Value
						break
					case "type":
						arg.t = val.Value
						break
					case "interface":
						arg.iface = val.Value
						break
					}
				}
				req.args = append(req.args, arg)
			}
			break
		case xml.EndElement:
			if e.Name.Local == "request" {
				return
			}
			break
		}
	}
	return
}

func DecodeEvent(decoder *xml.Decoder, attr []xml.Attr) (ev RequestEvent) {
	// get name
	for _,val := range attr {
		if val.Name.Local == "name" {
			ev.name = val.Value
		}
	}

	// parse args
	for token,err := decoder.Token() ; err == nil ; token,err = decoder.Token() {
		switch e := token.(type) {
		case xml.StartElement:
			if e.Name.Local == "arg" {
				var arg Arg
				for _,val := range e.Attr {
					switch val.Name.Local {
					case "name":
						arg.name = val.Value
						break
					case "type":
						arg.t = val.Value
						break
					case "interface":
						arg.iface = val.Value
						break
					}
				}
				ev.args = append(ev.args, arg)
			}
			break
		case xml.EndElement:
			if e.Name.Local == "event" {
				return
			}
			break
		}
	}
	return
}

func PrintInterface(iface Interface) (buf *bytes.Buffer, path string)  {
	buf = new(bytes.Buffer)

	// Make a proper interface name
	iname := strings.Replace(iface.name, "wl_", "", 1)
	tiname := strings.Title(iname)
	vname := iname[0:1]
	buf.WriteString("package gowl")
	if len(iface.events) > 0 {
		buf.WriteString(`

import (
	"bytes"
)

var _ bytes.Buffer
`)
	}
	buf.WriteString(fmt.Sprintf(`
type %s struct {
//	*WlObject
	id int32
	listeners map[int16][]chan interface{}
	events []func (%s *%s, msg []byte)
}

//// Requests`, tiname, vname, tiname))
	path = fmt.Sprintf("../%s.go", iname)

	for opcode,req := range iface.requests {
		argstrs := make([]string, 0)
		argnames := make([]string, 0)
		for _,arg := range req.args {
			argstrs = append(argstrs, fmt.Sprintf("%s %s", fixVarName(arg.name), getType(arg.t, arg.iface)))
			argnames = append(argnames, fixVarName(arg.name))
		}
		argstr := strings.Join(argstrs, ", ")
		argname := strings.Join(argnames, ", ")
		// Make func header
		buf.WriteString(fmt.Sprintf(`
func (%s *%s) %s (%s) {`, vname, tiname, strings.Title(req.name), argstr))

		// Make func body
		buf.WriteString(fmt.Sprintf(`
	msg := newMessage(%s, %d)
`, vname,opcode))
		for _,arg := range req.args {
			name := fixVarName(arg.name)
			switch arg.t {
			case "uint","int", "fixed":
				buf.WriteString(fmt.Sprintf("\twriteInteger(msg,%s)\n", name))
			case "string":
				buf.WriteString(fmt.Sprintf("\twriteString(msg,[]byte(%s))\n", name))
			case "object":
				buf.WriteString(fmt.Sprintf("\twriteInteger(msg,%s.Id())\n", name))
			case "new_id":
				buf.WriteString(fmt.Sprintf(`	appendObject(%s)
	writeInteger(msg,%s.Id())
`, name, name))
			case "fd":
				buf.WriteString(fmt.Sprintf("\twriteFd(msg,%s)\n", name))
			default:
				buf.WriteString(fmt.Sprintf("\twriteUnknown(%s)\n", arg.t))
			}
		}

		buf.WriteString(fmt.Sprintf(`
	sendmsg(msg)
	printRequest("%s", "%s", %s)
}
`, iname, req.name, argname))
	}

	//// EVENTS
	// HandleEvent
	buf.WriteString(fmt.Sprintf(`
//// Events
func (%s *%s) HandleEvent(opcode int16, msg []byte) {
	if %s.events[opcode] != nil {
		%s.events[opcode](%s, msg)
	}
}
`, vname, tiname, vname, vname, vname))

	handlers := make([]string, 0)
	for opcode, ev := range iface.events {
		handlers = append(handlers, fmt.Sprintf("%s_%s", iname, ev.name))

		//// Make listener type and adder	
		buf.WriteString(fmt.Sprintf(`
type %s%s struct {`, tiname, strings.Title(ev.name)))
		for _, arg := range ev.args {
			buf.WriteString(fmt.Sprintf(`
	%s %s`, strings.Title(fixVarName(arg.name)), getType(arg.t, arg.iface)))
		}

		buf.WriteString(fmt.Sprintf(`
}

func (%s *%s) Add%sListener(channel chan interface{}) {
	%s.listeners[%d] = append(%s.listeners[%d], channel)
}
`, vname, tiname, strings.Title(ev.name), vname, opcode, vname, opcode))

		//// Event handler
		argstrs := make([]string, 0)
		for _,arg := range ev.args {
			argstrs = append(argstrs, fixVarName(arg.name))
		}
		argstr := strings.Join(argstrs, ", ")

		buf.WriteString(fmt.Sprintf(`
func %s_%s(%s *%s, msg []byte) {
	var data %s%s
`, iname, ev.name, vname, tiname, tiname, strings.Title(ev.name)))
		if len(ev.args) > 0 {
			buf.WriteString(`	buf := bytes.NewBuffer(msg)
`)
		}
		// Func body
//		var argstr string
		for _, arg := range ev.args {
			name := fixVarName(arg.name)
			obj := "not"
//			argstr = fmt.Sprintf("%s, %s", argstr, name)
			var fname string
			switch arg.t {
			case "uint":
				fname = "readUint32"
			case "int", "fixed":
				fname = "readInt32"
			case "object":
				obj = "old"
			case "fd":
				fname = "readUintptr"
			case "new_id":
				obj = "new"
			case "string":
				fname = "readString"
				name = fmt.Sprintf("_,%s",name)
			case "array":
				fname = "readArray"
			default:
				fname = "unknownRead"
			}
			if obj != "not" {
				buf.WriteString(fmt.Sprintf(`
	%sid, err := readInt32(buf)
	if err != nil {
		// XXX Error handling
	}
	%s := new(%s)`, name, name, strings.Replace(getType(arg.t, arg.iface), "*", "", 1)))
				if obj == "old" {
					buf.WriteString(fmt.Sprintf(`
	%sobj := getObject(%sid)
	if %sobj == nil {
		return
	}
	%s = %sobj.(%s)
`, name, name, name, name, name, getType(arg.t, arg.iface)))
				} else {
					buf.WriteString(fmt.Sprintf(`
	setObject(%sid, %s)
`, name, name))
				}
			} else {
				buf.WriteString(fmt.Sprintf(`
	%s,err := %s(buf)
	if err != nil {
		// XXX Error handling
	}
`, name, fname))
			}
			buf.WriteString(fmt.Sprintf(`	data.%s = %s
`, strings.Title(fixVarName(arg.name)), fixVarName(arg.name)))
		}

		buf.WriteString(fmt.Sprintf(`
	for _,channel := range %s.listeners[%d] {
		go func () { channel <- data }()
	}
	printEvent("%s", "%s", %s)
}
`, vname, opcode, iname, ev.name, argstr))

	}

	// Constructor
	buf.WriteString(fmt.Sprintf(`
func New%s() (%s *%s) {
	%s = new(%s)
	%s.listeners = make(map[int16][]chan interface{}, 0)
`, tiname, vname, tiname, vname, tiname, vname))
	for _,handler := range handlers {
		buf.WriteString(fmt.Sprintf(`
	%s.events = append(%s.events, %s)`, vname, vname, handler))
	}
	buf.WriteString(fmt.Sprintf(`
	return
}

func (%s *%s) SetId(id int32) {
	%s.id = id
}

func (%s *%s) Id() int32 {
	return %s.id
}`, vname, tiname, vname, vname, tiname, vname))

	return
}

func getType(t string, iface string) string {
	switch t {
	case "uint":
		return "uint32"
	case "int", "fixed":
		return "int32"
	case "fd":
		return "uintptr"
	case "array":
		return "[]interface{}"
	case "object","new_id":
		typ := fmt.Sprintf("*%s", strings.Title(strings.Replace(iface, "wl_", "",1)))
		if typ == "*Object" {
			return "Object"
		}
		return typ
	}
	return t
}

func fixVarName(n string) string {
	if val,ok := replacements[n] ; ok {
		return val
	}
	return n
}
