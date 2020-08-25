window.onload = function() {
	var rootDiv = document.createElement( "div" );
	var textDiv = document.createElement( "div" );
	var textArea = document.createElement( "textarea" );
	var secureDiv = document.createElement( "div" );
	var checkbox = document.createElement( "input" );
	var checkboxLabel = document.createElement( "label" );
	var phraseInput = document.createElement( "input" );
	var buttonSecureDiv = document.createElement( "div" );
	var button = document.createElement( "button" );

	rootDiv.setAttribute( "id", "rootDiv" );

	textDiv.setAttribute( "id", "textDiv" );

	secureDiv.setAttribute( "id", "secureDiv" );

	textArea.setAttribute( "id", "textArea" );

	checkbox.setAttribute( "type", "checkbox" );
	checkbox.setAttribute( "id", "checkbox" );
	checkbox.onclick = showPhraseInput;

	checkboxLabel.setAttribute( "for", "checkbox" );	
	checkboxLabel.innerHTML = "secure?";

	phraseInput.setAttribute( "type", "text" );
	phraseInput.setAttribute( "size", "20" );
	phraseInput.setAttribute( "id", "phraseInput" );
	phraseInput.style.display = "none";

	buttonSecureDiv.setAttribute( "id", "buttonSecureDiv" );

	button.textContent = "send";
	button.setAttribute( "id", "buttonSend" );
	button.onclick = send;

	textDiv.append( textArea );
	rootDiv.append( textDiv );
	secureDiv.append( checkbox );
	secureDiv.append( checkboxLabel );
	secureDiv.append( phraseInput );
	rootDiv.append( secureDiv );
	buttonSecureDiv.append( button );
	rootDiv.append( buttonSecureDiv );
	document.body.append( rootDiv );
};

function send() {
	var text = document.getElementById( "textArea" ).value;
	var buttonSend = document.getElementById( "buttonSend" );
	var xhr = new XMLHttpRequest();

	var checkbox = document.getElementById( "checkbox" );
	var phraseInput = document.getElementById( "phraseInput" );

	xhr.open( "POST", "/api/1.0/setSecret" );
	xhr.setRequestHeader( "Content-Type", "application/json;charset=UTF-8" );

	if ( checkbox.checked == true && phraseInput.value != "" ) {
		xhr.send( JSON.stringify( { "message": text, "phrase": phraseInput.value } ) );
	} else if ( checkbox.checked == true && phraseInput.value == "" ) {
    	alert( "you missed something!" );
	} else if ( checkbox.checked == false ) {
		xhr.send( JSON.stringify( { "message": text } ) );
	}

	// xhr.send( JSON.stringify( { "message": text } ) );

	xhr.onload = function() {
		if ( xhr.status != 201 ) {
			alert( "error" );
		} else {
			displayLink( xhr.response );
		}
	};
}

function displayLink( hint ) {
	var parsedHint = JSON.parse( hint );
	var rootDiv = document.getElementById( "rootDiv" );
	var secureDiv = document.getElementById( "secureDiv" );
	var buttonSecureDiv = document.getElementById( "buttonSecureDiv" );
	var buttonDiv = document.createElement( "div" );
	var button = document.createElement( "button" );

	// rootDiv.innerHTML = "";

	var linkDiv = document.createElement( "div" );

	var link = document.createElement( "input" );

	linkDiv.setAttribute( "id", "linkDiv" );
	buttonDiv.setAttribute( "id", "buttonDiv" );
	button.textContent = "copy";
	button.onclick = copy;

	link.setAttribute( "type", "text" );
	link.setAttribute( "size", parsedHint.url.length + 5 );
	link.setAttribute( "value",  parsedHint.url );
	link.setAttribute( "id", "link" );
	link.setAttribute( "readOnly", "true" );

	secureDiv.remove();
	buttonSecureDiv.remove();

	linkDiv.append( link );
	buttonDiv.append( button );
	rootDiv.append( linkDiv );
	rootDiv.append( buttonDiv );
}

function copy() {
	var link = document.getElementById( "link" );
	link.select();
	document.execCommand( "copy" );
}

function showPhraseInput() {
	var checkbox = document.getElementById( "checkbox" );
	var phraseInput = document.getElementById( "phraseInput" );

	if ( checkbox.checked == true ){
		phraseInput.style.display = "block";
	} else {
    	phraseInput.style.display = "none";
	}
}