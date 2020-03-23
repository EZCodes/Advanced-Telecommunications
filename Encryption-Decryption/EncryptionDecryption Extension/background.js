 chrome.runtime.onInstalled.addListener(function() {
    chrome.storage.sync.set({color: '#3aa757'}, function() {
      console.log("Entered Background File");
    });

	//chrome.storage.local.get('signed_in', function(data) {
	//  if (data.signed_in) {
    //    chrome.browserAction.setPopup({popup: 'popup.html'});
    //  } else {
    //    chrome.browserAction.setPopup({popup: 'popup_sign_in.html'});
    //  }
    //});
  });
  