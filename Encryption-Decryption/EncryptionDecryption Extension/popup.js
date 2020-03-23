
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
	var groupMembers;
	var username;
	var password;

	
	chrome.storage.local.get(['password'], function(result) {
        password = result.password
		chrome.storage.local.get(['username'], function(result) {
			username = result.username
			chrome.storage.local.get(['group'], function(result) {
				groupMembers = result.group;
				var recipients
				if (currentOption == 'all') {
					recipients = groupMembers;
				} else {
					recipients = [currentOption];
				}
				var request = new XMLHttpRequest();
				var url = "http://localhost:420";
				request.open("POST", url, true);
				request.setRequestHeader("Content-Type", "application/json");
				request.onreadystatechange = function () { // when we receive the message, this function is a listener
					if (request.readyState === 4 && request.status === 200) { // receive json from the server
						var json = JSON.parse(request.responseText);
						document.getElementById("output").innerHTML = "Encrypted message is: " + json.message;
					} else {
						document.getElementById("output").innerHTML = "Encryption failed";
					}
				};
				var data = JSON.stringify({"type": "encrypt", "user": username, "password": password, "message" : text, "recipients" : recipients});
				request.send(data); // send the json to the server			
			});			
		});
    });
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
chrome.storage.local.get(['group'], function(result) {
	group = result.group;
	let option;
	for (let i=0;i<group.length;i++) {
		option = document.createElement('option');
		option.text = group[i];
		option.value = group[i];
		selectList.add(option);
	}
});

