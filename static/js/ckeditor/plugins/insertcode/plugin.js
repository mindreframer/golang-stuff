CKEDITOR.plugins.add('insertcode', {  
    requires: ['dialog'],  
    init: function(a){  
        var b = a.addCommand('insertcode', new CKEDITOR.dialogCommand('insertcode'));  
        a.ui.addButton('insertcode', {  
            label:"insertcode",  
            command: 'insertcode',  
            icon: this.path + 'icons/code.png'  
        });  
        CKEDITOR.dialog.add('insertcode', this.path + 'dialogs/insertcode.js');  
    }
});  