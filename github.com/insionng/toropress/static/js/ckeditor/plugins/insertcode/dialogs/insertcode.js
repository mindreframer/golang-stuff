CKEDITOR.dialog.add('insertcode', function(editor){
    var escape = function(value){
        return value;
    };
    return {
        title: '插入代码',
        resizable: CKEDITOR.DIALOG_RESIZE_BOTH,
        minWidth: 720,
        minHeight: 420,
        contents: [{  
            id: 'cb',  
            name: 'cb',  
            label: 'cb',  
            title: 'cb',  
            elements: [{  
                type: 'select',  
                label: 'Language',  
                id: 'lang',  
                required: true,  
                'default': 'go',  
                items: [['Python', 'python'],['Python profiler results', 'profile'],['Ruby', 'ruby'],['Haml', 'haml'],['Perl', 'perl'],['PHP', 'php'],['Scala', 'scala'],['Go language', 'go'],['HTML, XML', 'xml'],['Lasso', 'lasso'],['CSS', 'css'],['SCSS', 'scss'],['Markdown', 'markdown'],['AsciiDoc', 'asciidoc'],['Django', 'django'],['Handlebars', 'handlebars'],['JSON', 'json'],['JavaScript', 'javascript'],['CoffeeScript', 'coffeescript'],['ActionScript', 'actionscript'],['VBScript', 'vbscript'],['VB.Net', 'vbnet'],['HTTP', 'http'],['Lua', 'lua'],['Delphi', 'delphi'],['Java', 'java'],['C++','cpp'],['Objective C', 'objectivec'],['Vala', 'vala'],['C#', 'cs'],['F#', 'fsharp'],['D language', 'd'],['RenderMan RSL', 'rsl'],['RenderMan RIB', 'rib'],['Maya Embedded Language', 'mel'],['SQL', 'sql'],['Smalltalk', 'smalltalk'],['Lisp', 'lisp'],['Clojure', 'clojure'],['Ini', 'ini'],['Apache', 'apache'],['Nginx', 'nginx'],['Diff', 'diff'],['DOS', 'dos'],['Bash', 'bash'],['CMake', 'cmake'],['Axapta', 'axapta'],['Oracle Rules Language', 'ruleslanguage'],['1C', '1c'],['AVR assembler', 'avrasm'],['VHDL', 'vhdl'],['Parser3', 'parser3'],['TeX', 'tex'],['Haskell', 'haskell'],['Erlang', 'erlang'],['Rust', 'rust'],['Matlab', 'matlab'],['R', 'r'],['OpenGL Shading Language', 'glsl'],['AppleScript', 'applescript'],['Brainfuck', 'brainfuck'],['Mizar', 'mizar']]  
            }, {  
                type: 'textarea',  
                style: 'width:700px;height:400px',  
                label: 'Code',  
                id: 'code',  
                rows: 30,  
                'default': ''  
            }]  
        }],
	        onOk: function(){
                code = this.getValueOf('cb', 'code');
                lang = this.getValueOf('cb', 'lang');
                html = '' + "<code class='"+ lang +"'>" + CKEDITOR.tools.htmlEncode(code) + "</code>" +'';
                element = editor.document.createElement( 'pre' );
                element.setHtml(html);
                editor.insertElement(element);
                //editor.insertHtml(html);
	        }
    };
});