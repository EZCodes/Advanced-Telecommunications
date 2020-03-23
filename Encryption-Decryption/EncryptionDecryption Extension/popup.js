
function openGroupAdder() {
	console.log("Entered Add To Group");
	chrome.browserAction.setPopup({popup: 'add_popup.html'});
	window.location.href = "add_popup.html";
}

function openGroupRemover() {
	console.log("Entered Remove From Group");
	chrome.browserAction.setPopup({popup: 'remove_popup.html'});
	window.location.href = "remove_popup.html";
}

function openEncryptor() {
	console.log("Entered Encryption");
	chrome.browserAction.setPopup({popup: 'popup.html'});
	window.location.href = "popup.html";
}

function openDecryptor() {
	console.log("Entered Decryption");
	chrome.browserAction.setPopup({popup: 'decrypt_popup.html'});
	window.location.href = "decrypt_popup.html";
}

function parseMessage() {
	var text = document.getElementById('encryp').value;
	var currentOption = document.getElementById('recipient').value
	var recipients;
	if (currentOption == 'all') {
		chrome.storage.sync.get(['group'], function(result) {
			recipients = result.value
		});
	} else {
		recipients = [currentOption];
	}
	
	var username;
	var password;
	chrome.storage.sync.get(['password'], function(result) {
        password = result.value
    });
	chrome.storage.sync.get(['username'], function(result) {
		username = result.value
    });
	
	var xhr = new XMLHttpRequest();
	var url = "http://localhost:420";
	xhr.open("POST", url, true);
	xhr.setRequestHeader("Content-Type", "application/json");
	xhr.onreadystatechange = function () { // when we receive the message, this function is a listener
		if (xhr.readyState === 4 && xhr.status === 200) { // receive json from the server
			var json = JSON.parse(xhr.responseText);
			document.getElementById("output").innerHTML = "Encrypted message is: " + json.message;
		} else {
			document.getElementById("output").innerHTML = "Encryption failed";
		}
	};
	var data = JSON.stringify({"type": "encrypt", "user": username, "password": password, "message" : text, "recipients" : recipients});
	xhr.send(data); // send the json to the server


	
}

// Listeners for the buttons
document.getElementById('addToGroup').addEventListener('click', openGroupAdder);
document.getElementById('removeFromGroup').addEventListener('click', openGroupRemover);
document.getElementById('encryptionMode').addEventListener('click', openEncryptor);
document.getElementById('decryptionMode').addEventListener('click', openDecryptor);
document.getElementById('send').addEventListener('click', parseMessage);

let selectList = document.getElementById('recipient');
selectList.length = 0;
let defaultOpt = document.createElement('option');
defaultOpt.text = 'All in Group';
defaultOpt.value = 'all';
selectList.add(defaultOpt);
selectList.selectedIndex = 0;

var group;
chrome.storage.sync.get(['group'], function(result) {
	group = result.value
});

let option;
for (let i=0;i<group.length;i++) {
	option = document.createElement('option');
	option.text = group[i];
	option.value = group[i];
	selectList.add(option);
}

