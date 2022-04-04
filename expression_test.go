package httpfinger

import (
	"fmt"
	"testing"
	"time"
)

func TestExpression(t *testing.T) {
	startTime := time.Now()
	var expr = `(((response= "django" || response= "python") || header="django") || header="python")`
	e, err := NewExpression(expr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(e.expr)
	elapsed := time.Since(startTime)
	fmt.Printf("程序执行总时长为：[%s]", elapsed.String())
}

func TestBoolParse(t *testing.T) {
	var expr = `false || true`
	fmt.Println(ParseBoolFromString(expr))
}

func TestFindCoupleBracketIndex(t *testing.T) {
	expr := "(adf(   ())as()dfs)()()()()()()"

	index, err := findCoupleBracketIndex(expr)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(expr[:index+1])

}

func TestAbcdefg(t *testing.T) {
	banner := &Banner{
		Title:    "abc",
		Header:   "def",
		Body:     "gjo",
		Response: "django",
	}

	expr, err := NewExpression(`(((response= "django" || response= "python") || header="django") || header="python")`)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(expr.expr)
	fmt.Println(expr.value)
	fmt.Println(expr.Match(banner))
}

func TestTemp(t *testing.T) {
	Data, err := New("")
	if err != nil {
		fmt.Println(err)
	}

	banner := &Banner{
		Title: "xfinity ner324系统管理",
		Body: `<HTML>
<HEAD><TITLE>DVR WebClient</TITLE></HEAD>
<SCRIPT language="javascript">generator
<!--
  var a_server_name = "255.255.255.255";
  var a_server_port = 5191;
  var a_server_tcpmaxdelay = 30000;
  function window_onload()
  {
    a_server_name = location.hostname;
    rMedia1.ServerPort = a_server_port;
    rMedia1.TcpMaxDelay = a_server_tcpmaxdelay;
    rMedia1.ServerName = a_server_name;
  }
  function window_onunload()
  {
  }
//-->
</SCRIPT>
<BODY bgColor=#757e8d language=javascript onload="return window_onload()"
onunload="return window_onunload()" topMargin=0 marginheight="0" leftMargin=0
marginwidth="0">
<center><p>
<table width="600" border="0">
<tr><td>
<OBJECT
  classid="clsid:259F9FDF-97EA-4C59-B957-5160CAB6884E"
  codebase="ShareIE.cab#version=3,6,0,0" id=rMedia1 name=rMedia1
  width=896
  height=642
  align=center
  hspace=0
  vspace=0
>
</OBJECT>
</td></tr>
<tr><td>
<script language = "javascript">
  if(navigator.systemLanguage == "zh-cn")
  else
    document.write("If download ActiveX controls fail, please click <a href=IEClient.EXE>here</a> to install")
</script>
</td></tr>
</table>
</p></center>
</BODY></HTML>`,
		Header:   "asdb",
		Response: "adsfasdf",
	}

	products := Data.Search(banner)
	for _, productName := range products {
		fmt.Println(productName)
	}
}
