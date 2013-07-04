// Copyright 2012 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Copyright 2012,2013 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package hopwatch

import (
	"io"
	"net/http"
)

func html(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w,
		`<!DOCTYPE html>
<meta charset="utf-8" />
<title>Hopwatch Debugger</title>
<head>
	<link href="hopwatch.css" rel="stylesheet" type="text/css" >	
	<script src="http://ajax.googleapis.com/ajax/libs/jquery/1.9.1/jquery.min.js" type="text/javascript"></script>
	<script type="text/javascript" src="hopwatch.js" ></script>
</head>
<body>
	<div id="header">
		<div class="container wide">
			<div id="heading">
				<a href="/hopwatch.html">Hopwatch - debugging tool</a>
			</div>		
			<div id="menu">
				<a id="resume" class="buttonDisabled" href="javascript:actionResume();">F8 - Resume</a>
				<a id="disconnect" class="buttonDisabled" href="javascript:actionDisconnect();">Disconnect</a>
				<a class="buttonEnabled" href="http://go.pkgdoc.org/github.com/emicklei/hopwatch" target="_blank">About</a>
			</div>
		</div>
	</div>
	<div id="content">
		<div id="log-pane">
			<div id="output"></div>
		</div>		
		<div id="gosource-pane">
			<div id="gofile" class="mono ">somefile.go</div>
			<div id="nrs" class="mono"></div>
			<div id="gosource" class="mono">
			</div>
		</div>
	</div>
	<div id="footer">
		&copy; 2012-2013. <a href="http://github.com/emicklei/hopwatch" target="_blank">hopwatch on github.com</a>
	</div>
</body>
</html>
`)
	return
}
