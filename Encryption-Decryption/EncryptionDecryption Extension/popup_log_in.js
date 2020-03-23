function openRegistration() {
	chrome.browserAction.setPopup({popup: 'register_popup.html'});
	window.location.href = "register_popup.html";
}

function login() {
	var username = document.getElementById('username').value;
	var password = document.getElementById('password').value
	
	//send the credentials and receive the response and save the group and username + password
	
	document.getElementById("error_output").innerHTML = "Login Failed!"
}

// Listeners for the buttons
document.getElementById('login').addEventListener('click', login);
document.getElementById('register').addEventListener('click', openRegistration);
