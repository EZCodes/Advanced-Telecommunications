 chrome.runtime.onInstalled.addListener(function() {
    chrome.storage.sync.set({color: '#3aa757'}, function() {
      console.log("Entered Background File");
    });
	chrome.declarativeContent.onPageChanged.removeRules(undefined, function() {
     chrome.declarativeContent.onPageChanged.addRules([{
       conditions: [new chrome.declarativeContent.PageStateMatcher({
         pageUrl: {hostEquals: 'developer.chrome.com'},
       })
       ],
           actions: [new chrome.declarativeContent.ShowPageAction()]
     }]);
    });
	//chrome.runtime.onMessage.addListener(
	//	function(message) {
	//		if(message.type == "link") {
	//			window.location.href = message.link;
	//		}
	//	});
	//chrome.storage.local.get('signed_in', function(data) {
	//  if (data.signed_in) {
    //    chrome.browserAction.setPopup({popup: 'popup.html'});
    //  } else {
    //    chrome.browserAction.setPopup({popup: 'popup_sign_in.html'});
    //  }
    //});
  });
  