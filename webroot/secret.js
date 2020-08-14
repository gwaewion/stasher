window.onload = function() {
	var textDiv = document.createElement( "div" );
	var buttonDiv = document.createElement( "div" );
	var button = document.createElement( "button" );

	button.textContent = "show";
	button.onclick = show;

	textDiv.innerHTML = "show message?";
	document.body.append( textDiv );
	buttonDiv.append( button );
	document.body.append( buttonDiv );
};

function show() {
	var pathArray = window.location.pathname.split('/');
	var id = pathArray[ pathArray.length - 1 ];

	var xhr = new XMLHttpRequest();

	xhr.open( "POST", "/api/1.0/getSecret" );
	xhr.setRequestHeader( "Content-Type", "application/json;charset=UTF-8" );

	xhr.send( JSON.stringify( { "id": id } ) );

	xhr.onload = function() {
		var errorMessage = JSON.parse( xhr.response );

		if ( xhr.status == 200 ) {
			displayMessage( xhr.response );
		} else if ( xhr.status == 400 ) {	
			if ( errorMessage.error == "no phrase" ) {
				askPhrase( id );
			} else if ( errorMessage.error == "wrong phrase" ) {
				alert( "wrong phrase" );
			}
		} else if ( errorMessage.error == "secret not exists" && xhr.status == 404 ) {
			alert( "secret deleted" );
		}
	};
}

function displayMessage( hint ) {
	var parsedHint = JSON.parse( hint );

	document.body.innerHTML = "";

	var messageDiv = document.createElement( "div" );

	messageDiv.innerHTML = parsedHint.message;

	document.body.append( messageDiv );
}

function askPhrase( id ) {
	document.body.innerHTML = "";

	var phraseDiv = document.createElement( "div" );
	var phraseInput = document.createElement( "input" );
	var button = document.createElement( "button" );

	phraseDiv.innerHTML = "input phrase:";

	phraseInput.setAttribute( "type", "text" );
	phraseInput.setAttribute( "size", "20" );
	phraseInput.setAttribute( "id", "phraseInput" );

	button.textContent = "send";

	button.onclick = function() {
	var xhr = new XMLHttpRequest();

	xhr.open( "POST", "/api/1.0/getSecret" );
	xhr.setRequestHeader( "Content-Type", "application/json;charset=UTF-8" );

	xhr.send( JSON.stringify( { "id": id, "phrase": phraseInput.value } ) );

	xhr.onload = function() {
		if ( xhr.status == 200 ) {
			displayMessage( xhr.response );
		} else if ( xhr.status == 400 ) {
			var errorMessage = JSON.parse( xhr.response );
			if ( errorMessage.error == "no phrase" ) {
				askPhrase();
			} else if ( errorMessage.error == "wrong phrase" ) {
				alert( "wrong phrase" );
			}
		}
	};
	};

	phraseDiv.append( phraseInput );
	phraseDiv.append( button );
	document.body.append( phraseDiv );
}