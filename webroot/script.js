window.onload = function() {
	var textDiv = document.createElement( "div" );
	var textArea = document.createElement( "textarea" );
	var secureDiv = document.createElement( "div" );
	var checkbox = document.createElement( "input" );
	var checkboxLabel = document.createElement( "label" );
	var phraseInput = document.createElement( "input" );
	var buttonDiv = document.createElement( "div" );
	var button = document.createElement( "button" );

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

	button.textContent = "send";
	button.onclick = send;

	textDiv.append( textArea );
	document.body.append( textDiv );
	secureDiv.append( checkbox );
	secureDiv.append( checkboxLabel );
	secureDiv.append( phraseInput );
	document.body.append( secureDiv );
	buttonDiv.append( button );
	document.body.append( buttonDiv );
};

function send() {
	var text = document.getElementById( "textArea" ).value;
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

	document.body.innerHTML = "";

	var linkDiv = document.createElement( "div" );

	var link = document.createElement( "a" );

	link.setAttribute( "href",  parsedHint.url );
	link.innerHTML = "share this";

	linkDiv.append( link );
	document.body.append( linkDiv );
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