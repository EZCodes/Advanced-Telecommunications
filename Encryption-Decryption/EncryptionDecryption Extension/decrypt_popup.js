
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

function decrypt() {
	var ciphertext = document.getElementById('decryptMsg').value;
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
			document.getElementById("output").innerHTML = "Decrypted message is: "; + json.message;
		} else {
			document.getElementById("output").innerHTML = "Decryption failed";
		}
	};
	var data = JSON.stringify({"type": "decrypt", "user": username, "password": password, "message" : ciphertext});
	xhr.send(data); // send the json to the server
	
	
}

// Listeners for the buttons
document.getElementById('addToGroup').addEventListener('click', openGroupAdder);
document.getElementById('removeFromGroup').addEventListener('click', openGroupRemover);
document.getElementById('encryptionMode').addEventListener('click', openEncryptor);
document.getElementById('decryptionMode').addEventListener('click', openDecryptor);
document.getElementById('decrypt').addEventListener('click', decrypt);