
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

function remove() {
	var userToRemove = document.getElementById('removeUser').value;
	var username;
	var password;
	chrome.storage.sync.get(['password'], function(result) {
        password = result.value
    });
	chrome.storage.sync.get(['username'], function(result) {
		username = result.value
    });
	
	// send the username and current credentials recieve new username list and save it
	var request = new XMLHttpRequest();
	var url = "http://localhost:420";
	request.open("POST", url, true);
	request.setRequestHeader("Content-Type", "application/json");
	request.onreadystatechange = function () { // when we receive the message, this function is a listener
		if (request.readyState === 4 && request.status === 200) { // proceed accordingly when received
			var json = JSON.parse(xhr.responseText);
			chrome.storage.sync.set({"group": json.members}, function() {
				console.log('Group saved!');
			});	
			document.getElementById("output").innerHTML = "User successfully removed from the group!"
		} else {
			document.getElementById("output").innerHTML = "There was a problem removing user from the group. Try again!"
		}
	};
	var data = JSON.stringify({"type": "remove", "user": username, "password": password, "message" : userToRemove});
	request.send(data); // send the json to the server
	
}

// Listeners for the buttons
document.getElementById('addToGroup').addEventListener('click', openGroupAdder);
document.getElementById('removeFromGroup').addEventListener('click', openGroupRemover);
document.getElementById('encryptionMode').addEventListener('click', openEncryptor);
document.getElementById('decryptionMode').addEventListener('click', openDecryptor);
document.getElementById('remove').addEventListener('click', remove);