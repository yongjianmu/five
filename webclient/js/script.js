//First test for the browsers support for WebSockets
if (!window.WebSocket) {
  //If the user's browser does not support WebSockets, give an alert message
  alert("Your browser does not support the WebSocket API!");
} else {
  //Set the websocket server URL
  var websocketurl = "ws://localhost:8080/echo";

  //get status element
  var connstatus = document.getElementById("connectionstatus");

  //get info div element
  var infodiv = document.getElementById("info");

  //Create the WebSocket object (web socket echo test service provided by websocket.org)
  socket = new WebSocket(websocketurl);

  //This function is called when the websocket connection is opened
  socket.onopen = function() {
    connstatus.innerHTML = "Connected!";
    infodiv.innerHTML += "<p>Connected to websocket server at: " + websocketurl + "</p>";
  };

  //This function is called when the websocket connection is closed
  socket.onclose = function() {
    connstatus.innerHTML = "Disconnected";
    infodiv.innerHTML += "<p>Disconnected from the websocket server at: " + websocketurl + "</p>";
  };

  //This function is called when the websocket receives a message. It is passed the message object as its only parameter
  socket.onmessage = function(message) {
    infodiv.innerHTML += "<p>Message received from server: '" + message.data + "'</p>";
  };
}

//function to send a message to the websocket server


function sendMessage() {
  //check to ensure that the socket variable is present i.e. the browser support tests passed
  if (socket) {
    //get the message text input element
    var message = document.getElementById("message").value;

    if (message !== "") {
      socket.send(message);
      infodiv.innerHTML += "<p>Sent message to server: '" + message + "'</p>";
    } else {
      alert("You must enter a message to be sent!");
    }
  }
}

function closeConnection() {
  //check to ensure that the socket variable is present i.e. the browser support tests passed
  if (socket) {
    socket.close();
  }
}
