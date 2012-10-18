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
	rawName string // raw
    priName string
    capName string // capitalized interface name
    varName string // variable name of interface (the d in "d *Display")
	requests []RequestEvent
	events []RequestEvent
	enum []Enum
}

type RequestEvent struct {
	name string
    capName string
    listenName string
	opcode int
	args []Arg
    signature *bytes.Buffer
}

type Arg struct {
	name string
    capName string
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
        buffer := new(bytes.Buffer)
		PrintInterface(val, buffer)
        fmt.Fprintf(buffer, `
////
//// REQUESTS
////
`)
        PrintRequests(val, buffer)
        fmt.Fprintf(buffer, `
////
//// EVENTS
////
`)
        PrintEvents(val, buffer)
		file, err := os.Create(fmt.Sprintf("../%s.go", val.capName))
		if err != nil {
			fmt.Println(err)
			return
		}
		file.Write(buffer.Bytes())
        file.Close()
	}

    buffer := new(bytes.Buffer)
    PrintSignatures(ifaces, buffer)
    file, err := os.Create("../Signatures.go")
    if err != nil {
        fmt.Println(err)
        return
    }
    file.Write(buffer.Bytes())
    file.Close()
	return
}

func DecodeInterface(decoder *xml.Decoder, attr []xml.Attr) (iface Interface) {
	// Setup interface
	for _,val := range attr {
		if val.Name.Local == "name" {
			iface.rawName = val.Value
            iname := strings.Replace(iface.rawName, "wl_", "", 1)
            iface.capName = makeInterfaceName(iname)
            iface.varName = iname[0:1]
            iface.priName = fmt.Sprintf("%s%s", strings.ToLower(iface.capName[0:1]), iface.capName[1:])
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
            iname := strings.Replace(req.name, "wl_", "", 1)
            req.capName = makeInterfaceName(iname)
		}
	}

    req.signature = new(bytes.Buffer)

	// parse args
	for token,err := decoder.Token() ; err == nil ; token,err = decoder.Token() {
		switch e := token.(type) {
		case xml.StartElement:
			if e.Name.Local == "arg" {
				var arg Arg
				for _,val := range e.Attr {
					switch val.Name.Local {
					case "name":
						arg.name = fixVarName(val.Value)
						break
					case "type":
                        arg.t = val.Value
                        fmt.Fprintf(req.signature, "%s", getArgSignature(arg))
						break
					case "interface":
						arg.iface = makeInterfaceName(val.Value)
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
            ev.capName = makeInterfaceName(ev.name)
            ev.listenName = fmt.Sprintf("%s%s",strings.ToLower(ev.capName[0:1]), ev.capName[1:])
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
						arg.name = fixVarName(val.Value)
                        arg.capName = makeInterfaceName(arg.name)
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

func PrintInterface(iface Interface, buffer *bytes.Buffer) {
    // Make listener strings
    makelisteners := make([]string, 0, len(iface.events))
    listeners := make([]string, 0, len(iface.events))
    appends := make([]string, 0, len(iface.events))
    for _,event := range iface.events {
        makelisteners = append(makelisteners, fmt.Sprintf(
`%s.%sListeners = make([]chan %s%s, 0)`,
        iface.varName, event.listenName, iface.capName, event.capName))

        listeners = append(listeners, fmt.Sprintf(
`%sListeners []chan %s%s`,
        event.listenName, iface.capName, event.capName))

        appends = append(appends, fmt.Sprintf(
`%s.events = append(%s.events, %s%s)`,
        iface.varName, iface.varName, iface.priName, event.capName))
    }

    makelistener := strings.Join(makelisteners, "\n\t")
    listener := strings.Join(listeners, "\n\t")
    appendstr := strings.Join(appends, "\n\t")

	fmt.Fprintf(buffer,
`package gowl

import (
	"bytes"
)

var _ bytes.Buffer

type %s struct {
	id int32
    %s
	events []func(%s *%s, msg message) error
    name string
}

func New%s() (%s *%s) {
	%s = new(%s)
    %s.name = "%s"
    %s

    %s
	return
}

func (%s *%s) HandleEvent(msg message) {
	if %s.events[msg.opcode] != nil {
		%s.events[msg.opcode](%s, msg)
	}
}

func (%s *%s) SetId(id int32) {
	%s.id = id
}

func (%s *%s) Id() int32 {
	return %s.id
}

func (%s *%s) Name() string {
    return %s.name
}
`,
iface.capName,
listener,
iface.varName, iface.capName,
iface.capName, iface.varName, iface.capName,
iface.varName, iface.capName,
iface.varName, iface.capName,
makelistener,
appendstr,
iface.varName, iface.capName,
iface.varName,
iface.varName, iface.varName,
iface.varName, iface.capName,
iface.varName,
iface.varName, iface.capName,
iface.varName,
iface.varName, iface.capName,
iface.varName)
}

func PrintRequests(iface Interface, buffer *bytes.Buffer) {
	for _,req := range iface.requests {
        // Make parameter and argument string
        params := make([]string, 0, len(req.args))
        args := make([]string, 0, len(req.args))
        for _,arg := range req.args {
            params = append(params, fmt.Sprintf("%s %s", arg.name, getArgType(arg)))
            args = append(args, arg.name)
        }


        fmt.Fprintf(buffer,
`
func (%s *%s) %s(%s) {
    sendrequest(%s, "%s_%s", %s)
}
`,
        iface.varName, iface.capName, req.capName, strings.Join(params, ", "),
        iface.varName, iface.rawName, req.name, strings.Join(args, ", "))
        }
}

func PrintEvents(iface Interface, buffer *bytes.Buffer) {
    for _,event := range iface.events {
        // Make string to describe fields
        args := make([]string, 0, len(event.args))
        for _,arg := range event.args {
            args = append(args, fmt.Sprintf("%s %s", arg.capName, getArgType(arg)))
        }
        fields := strings.Join(args, "\n\t")

        datas := make([]string, 0, len(event.args))
        for _,arg := range event.args {
            datas = append(datas, getDataFetch(arg))
        }
        data := strings.Join(datas, "\n");

        fmt.Fprintf(buffer,
`
type %s%s struct {
    %s
}

func (%s *%s) Add%sListener(channel chan %s%s) {
    %s.%sListeners = append(%s.%sListeners, channel)
}

func %s%s(%s *%s, msg message) (err error) {
    var data %s%s
%s

    // Dispatch events
    for _,channel := range %s.%sListeners {
        go func () {
            channel <- data
        } ()
    }
    return
}
`,
        iface.capName, event.capName,
        fields,
        iface.varName, iface.capName, event.capName, iface.capName, event.capName,
        iface.varName, event.listenName, iface.varName, event.listenName,
        iface.priName, event.capName, iface.varName, iface.capName,
        iface.capName, event.capName,
        data,
        iface.varName, event.listenName)
    }
}

func getDataFetch(arg Arg) string {
    return fmt.Sprintf(`
    // Read %s
    %s,err := %s
    if err != nil {
        return
    }
    %s`,
    arg.name,
    arg.name, getFuncCall(arg),
    getPostFetch(arg))
}

func PrintSignatures(ifaces []Interface, buffer *bytes.Buffer) {
    signatures := make([]string,0)
    for _,iface := range ifaces {
        for opcode,req := range iface.requests {
            signatures = append(signatures, fmt.Sprintf(
`%ssignatures["%s_%s"] = Signature{%d, "%s"}`, "\t", iface.rawName, req.name, opcode, req.signature))
        }
        signatures = append(signatures, "")
    }
    signaturestr := strings.Join(signatures, "\n")
    fmt.Fprintf(buffer,
`package gowl

type Signature struct {
    opcode int16
    signature string
}

var signatures map[string]Signature

func init() {
    signatures = make(map[string]Signature, 0)
%s
}`, signaturestr)
    return
}

func getPostFetch(arg Arg) string {
    switch arg.t {
    case "object":
        return fmt.Sprintf(
    `%sObj := getObject(%s)
    data.%s = %sObj.(%s)`,
        arg.name, arg.name,
        arg.capName, arg.name,getArgType(arg))
    case "new_id":
        return fmt.Sprintf(
    `%sObj := new(%s)
    setObject(%s, %sObj)
    data.%s = %sObj`,
        arg.name, getArgType(arg)[1:], // Will this always work?
        arg.name, arg.name,
        arg.capName, arg.name)
    default:
        return fmt.Sprintf("data.%s = %s", arg.capName, arg.name)
    }
    return "never"
}

func getFuncCall(arg Arg) string {
    switch arg.t {
    case "int", "object", "new_id":
        return "readInt32(msg.buf)"
    case "uint":
        return "readUint32(msg.buf)"
    case "fd":
        return "msg.fd, nil"
    case "array":
        return "readArray(msg.buf)"
    case "string":
        return "readString(msg.buf)"
    case "fixed":
        return "readFixed(msg.buf)"
    }
    return arg.t
}

func getArgType(arg Arg) string {
	switch arg.t {
	case "fixed":
		return "int32"
	case "uint":
		return "uint32"
	case "int":
		return "int32"
	case "fd":
		return "uintptr"
	case "array":
		return "[]interface{}"
	case "object","new_id":
		typ := fmt.Sprintf("*%s", makeInterfaceName(arg.iface))
		if typ == "*" {
			return "Object"
		}
		return typ
	}
	return arg.t
}

func getArgSignature(arg Arg) string {
    switch arg.t {
    case "int":
        return "i"
    case "uint":
        return "u"
    case "array":
        return "a"
    case "object":
        return "o"
    case "new_id":
        return "n"
    case "string":
        return "s"
    case "fixed":
        return "f"
    case "fd":
        return "h"
    }
    return "never"
}

func fixVarName(n string) string {
	if val,ok := replacements[n] ; ok {
		return val
	}
	return n
}

func makeInterfaceName(n string) string {
    n = strings.Replace(n, "wl_", "", 1)
	parts := strings.Split(n, "_")
	for i,s := range parts {
		parts[i] = strings.Title(s)
	}
	return strings.Join(parts, "")
}
