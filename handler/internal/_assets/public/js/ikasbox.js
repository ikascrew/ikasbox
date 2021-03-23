
var dialog = document.querySelector('#groupDialog');
var showDialogButton = document.querySelector('#show-groupDialog');
var addButton = document.querySelector('#groupAdd');
var groupSpinner = document.querySelector('#group_spinner');
var groupMenu = document.querySelectorAll('.group_selector');

if (! dialog.showModal) {
    dialogPolyfill.registerDialog(dialog);
}

showDialogButton.addEventListener('click', function(e) {
    dialog.showModal();
},false);

addButton.addEventListener('click', function(e) {

    groupSpinner.classList.add('is-active');

    xmlhttp.open("POST", "/group/add", true);
    xmlhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    xmlhttp.setRequestHeader("Cache-Control", "no-cache"); // For no cache

    var groupName = document.querySelector('#group_name').value;
    var data = {
        name : groupName
    };
    xmlhttp.onreadystatechange = function() {
        if (xmlhttp.readyState == 4 ) {
            groupSpinner.classList.remove('is-active');
            if ( xmlhttp.status != 200) {
              alert(xmlhttp.status + " : " + xmlhttp.responseText);
            } else {
              document.querySelector('#group_name').value = "";
              dialog.close();
            }
        }
    }
    xmlhttp.send(EncodeHTMLForm(data));
});

dialog.querySelector('.close').addEventListener('click', function(e) {
    dialog.close();
},false);

var xmlhttp;
if (window.XMLHttpRequest) { // code for IE7+, Firefox, Chrome, Opera, Safari
    xmlhttp=new XMLHttpRequest();
} else { // code for IE6, IE5
    xmlhttp=new ActiveXObject("Microsoft.XMLHTTP");
}

function selectGroup(e) {
    xmlhttp.open("POST", "/group/select", true);
    xmlhttp.setRequestHeader("Content-type", "application/x-www-form-urlencoded");
    xmlhttp.setRequestHeader("Cache-Control", "no-cache");

    var selectId = this.getAttribute("data-id");
    var data = {
        id : selectId
    };

    xmlhttp.onreadystatechange = function() {
        if (xmlhttp.readyState == 4 ) {
            groupSpinner.classList.remove('is-active');
            if ( xmlhttp.status != 200) {
            } else {
                location.reload()
            }
        }
    }
    xmlhttp.send(EncodeHTMLForm(data));
}

for (var i = 0; i < groupMenu.length; i++) {
    groupMenu[i].addEventListener('click', selectGroup, false);
}

function EncodeHTMLForm(data){
    var params = [];
    for(var name in data){
        var value = data[name];
        var param = encodeURIComponent(name).replace(/%20/g, '+')
            + '=' + encodeURIComponent(value).replace(/%20/g, '+');
        params.push(param);
    }
    return params.join('&');
}
