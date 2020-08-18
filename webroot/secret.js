var id;

window.onload = function() {
	var rootDiv = document.createElement( "div" );
	var textDiv = document.createElement( "div" );
	var buttonDiv = document.createElement( "div" );
	var button = document.createElement( "button" );

	rootDiv.setAttribute( "id", "rootDiv" );

	buttonDiv.setAttribute( "id", "buttonDivSecret" );
	button.textContent = "show";
	button.onclick = show;

	textDiv.setAttribute( "id", "textDivSecret" );
	textDiv.innerHTML = "show message?";
	
	rootDiv.append( textDiv );
	buttonDiv.append( button );
	rootDiv.append( buttonDiv );
	document.body.append( rootDiv );
};

function show() {
	var pathArray = window.location.pathname.split('/');
	id = pathArray[ pathArray.length - 1 ];

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
	var rootDiv = document.getElementById( "rootDiv" );

	rootDiv.innerHTML = "";

	var messageDiv = document.createElement( "div" );

	messageDiv.innerHTML = parsedHint.message;

	rootDiv.append( messageDiv );
}

function askPhrase( id ) {
	var rootDiv = document.getElementById( "rootDiv" );

	rootDiv.innerHTML = "";

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
				askPhrase( id );
			} else if ( errorMessage.error == "wrong phrase" ) {
				alert( "wrong phrase" );
			} else if ( errorMessage.error == "secret not exists" && xhr.status == 404 ) {
			alert( "secret deleted" );
			}
		}
	};
	};

	phraseDiv.append( phraseInput );
	phraseDiv.append( button );
	rootDiv.append( phraseDiv );
}