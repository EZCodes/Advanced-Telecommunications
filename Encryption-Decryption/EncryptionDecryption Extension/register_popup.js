function openLogin() {
	chrome.browserAction.setPopup({popup: 'popup_log_in.html'});
	window.location.href = "popup_log_in.html";
}

function register() {
	var username = document.getElementById('username').value;
	var password = document.getElementById('password').value
	
	//send the credentials and receive the response and save the group and username + password
	var request = new XMLHttpRequest();
	var url = "http://localhost:420";
	request.open("POST", url, true);
	request.setRequestHeader("Content-Type", "application/json");
	request.onreadystatechange = function () { // when we receive the message, this function is a listener
		if (request.readyState === 4 && request.status === 200) { // proceed accordingly when received
			chrome.storage.local.set({"username": username}, function() {
				console.log('Username saved!');
			});
			chrome.storage.local.set({"password": password}, function() {
				console.log('Password saved!');
			});
			chrome.storage.local.set({"group": []}, function() {
				console.log('Group saved!');
			});			
			chrome.browserAction.setPopup({popup: 'popup.html'});
			window.location.href = "popup.html";
		} else {
			document.getElementById("error_output").innerHTML = "Username already exists!"
		}
	};
	var data = JSON.stringify({"type": "register", "user": username, "password": password});
	request.send(data); // send the json to the server
	
}

// Listeners for the buttons
document.getElementById('login').addEventListener('click', openLogin);
document.getElementById('register').addEventListener('click', register);
