// Copyright 2012 Ernest Micklei. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package hopwatch

import (
	"io"
	"net/http"
)

func css(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	io.WriteString(w, `
	body, html {
		margin: 0;
		padding: 0;
		font-family: Helvetica, Arial, sans-serif;
		font-size: 16px;
		color: #222;		
	}
	.mono    {font-family:"Lucida Console", Monaco, monospace;font-size:13px;}
	.wide    {width:100%;}
	
	#header, #content, #footer, #log-pane, #gosource-pane {
		position:absolute;
	}
	
	/******************
	 * Heading
	 */	
	div#heading {
		float: left;
		margin: 0 0 10px 0;
		padding: 21px 0;
		font-size: 20px;
		font-weight: normal;
	}
	div#heading a {
		color: #222;
		text-decoration: none;
	}	
	div#header {
		background: #E0EBF5;
		height: 	64px;
		width:		100%;
	}	
	.container {
		padding: 	0 20px;
	}
	div#menu {
		float: left;
		min-width: 590px;
		padding: 10px 0;
		text-align: right;
		margin-top: 10px;
	}
	div#menu > a {
		margin-right: 5px;
		margin-bottom: 10px;
		padding: 10px;				
	}
	.buttonEnabled {
		color: white;
		background: #375EAB;
	}
	.buttonDisabled {
		color: #375EAB;
		background: white;
	}	
	div#menu > a,
	div#menu > input {
		padding: 10px;	
		text-decoration: none;
		font-size: 16px;	
		-webkit-border-radius: 5px;
		-moz-border-radius: 5px;
		border-radius: 5px;
	}
	
	/******************
	 * Footer
	 */	
	div#footer {
		bottom: 0;
		height: 24px;
		width: 100%;
				
		text-align: center;
		color: #666;
		font-size: 14px;
	}
	
	/******************
	 * Content
	 */
	#content {
		top: 	64px;
		bottom: 24px;
		width:	100%;		
	}
		
	/******************
	 * Log
	 */	
	#log-pane { 
		height: 100%;
		width: 60%; 
		overflow: auto; 
	}
	a { text-decoration:none; color: #375EAB; }
	a:hover { text-decoration:underline ; color:black }
                 
	.logline {}
    .srcline {}
	.toggle  {padding-left:4px;padding-right:4px;margin-left:4px;margin-right:4px;background-color:#375EAB;color:#FFF;}	
	.stack   {
		background-color:#FFD;
		padding: 4px;
		border-width: 1px;
		border-color: #ddd;
		border-style: solid;
		box-shadow: inset 0 4px 5px -5px rgba(0,0,0,0.4);
		margin: 1px 6px 0;	
	}
	.time    {color:#AAA;white-space:nowrap}
	.watch 	 {width:100%;white-space:pre}
	.goline  {color:#888;padding-left:8px;padding-right:8px;}
	.err 	 {background-color:#FF3300;width:100%;}
	.info 	 {width:100%;}
	.break   {background-color:#375EAB;color:#FFF;}
	.suspend {}

	/******************
	 * Source
	 */	
	#gosource-pane { 
		height: 100%;
		left: 60% ;
		right: 0px; 
		display: none; 
		margin: 0;
		overflow: auto;
	}
	#gosource {
		margin: 		0;
		background: 	#FFD;
		white-space:	pre; 		
	}
	#nrs { 
		width:			24px;
		float:			left;
	}
	#gofile {
		background-color:#FFF;
		color:#375EAB;
	}		
	`)
	return
}
