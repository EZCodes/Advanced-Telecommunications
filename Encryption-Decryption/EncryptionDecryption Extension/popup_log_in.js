function openRegistration() {
	chrome.browserAction.setPopup({popup: 'register_popup.html'});
	window.location.href = "register_popup.html";
}

function login() {
	var username = document.getElementById('username').value;
	var password = document.getElementById('password').value
	
	//send the credentials and receive the response and save the group and username + password
	var request = new XMLHttpRequest();
	var url = "http://localhost:420";
	request.open("POST", url, true);
	request.setRequestHeader("Content-Type", "application/json");
	request.onreadystatechange = function () { // when we receive the message, this function is a listener
		if (request.readyState === 4 && request.status === 200) { // proceed accordingly when received
			var json = JSON.parse(xhr.responseText);
			chrome.storage.sync.set({"username": username}, function() {
				console.log('Username saved!');
			});
			chrome.storage.sync.set({"password": password}, function() {
				console.log('Password saved!');
			});
			chrome.storage.sync.set({"group": json.members}, function() {
				console.log('Group saved!');
			});
			chrome.browserAction.setPopup({popup: 'popup.html'});
			window.location.href = "popup.html";
		} else {	
			document.getElementById("error_output").innerHTML = "Login Failed!"
		}
	};
	var data = JSON.stringify({"type": "login", "user": username, "password": password});
	request.send(data); // send the json to the server
	
}

// Listeners for the buttons
document.getElementById('login').addEventListener('click', login);
document.getElementById('register').addEventListener('click', openRegistration);
