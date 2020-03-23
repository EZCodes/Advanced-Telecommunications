function openLogin() {
	chrome.browserAction.setPopup({popup: 'popup_log_in.html'});
	window.location.href = "popup_log_in.html";
}

function register() {
	var username = document.getElementById('username').value;
	var password = document.getElementById('password').value
	
	//send the credentials and receive the response and save the group and username + password
	
	document.getElementById("error_output").innerHTML = "Username already exists!"
}

// Listeners for the buttons
document.getElementById('login').addEventListener('click', openLogin);
document.getElementById('register').addEventListener('click', register);
