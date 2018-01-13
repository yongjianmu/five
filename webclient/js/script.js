var socket;

//First test for the browsers support for WebSockets
if (!window.WebSocket) {
    //If the user's browser does not support WebSockets, give an alert message
    alert("Your browser does not support the WebSocket API!");
    } else {
    //Set the websocket server URL
    var websocketurl = "ws://localhost:8080/echo";
    //var websocketurl = "ws://echo.websocket.org/";

    //get status element
    var connstatus = document.getElementById("connectionstatus");

    //get info div element
    var infodiv = document.getElementById("info");

    //Create the WebSocket object (web socket echo test service provided by websocket.org)
    socket = new WebSocket(websocketurl);

    //This function is called when the websocket connection is opened
    socket.onopen = function() {
        //connstatus.innerHTML = "Connected!";
        //infodiv.innerHTML += "<p>Connected to websocket server at: " + websocketurl + "</p>";
        console.log("Connection Opened")
    };

    //This function is called when the websocket connection is closed
    socket.onclose = function() {
        //connstatus.innerHTML = "Disconnected";
        //infodiv.innerHTML += "<p>Disconnected from the websocket server at: " + websocketurl + "</p>";
        console.log("Connection Closed")
    };

    var isValid = false
    //This function is called when the websocket receives a message. It is passed the message object as its only parameter
    socket.onmessage = function(message) {
        //infodiv.innerHTML += "<p>Message received from server: '" + message.data + "'</p>";
        console.log("msg sent");
    };

    var chess = document.getElementById("mycanvas");
    var context = chess.getContext('2d');
    var x_pos;
    var y_pos;
    var me = true;
    var chessBox = [];
    for(var i=0;i<15;i++){
        chessBox[i]=[];
        for(var j=0;j<15;j++){
            chessBox[i][j]=0;
        }
    }
    function drawChessBoard(){
        for(var i=0;i<15;i++){
            context.strokeStyle="#D6D1D1";
            context.moveTo(15+i*30,15);
            context.lineTo(15+i*30,435);
            context.stroke();
            context.moveTo(15,15+i*30);
            context.lineTo(435,15+i*30);
            context.stroke();
        }
    }
    drawChessBoard();
    socket.addEventListener("message", function(event) {
        var data = event.data;
        // 处理数据
        console.log("#### Receive data: ", data)
        if(data == "ok"){
            oneStep(x_pos, y_pos, me)
        }
        if(data == "win"){
            oneStep(x_pos, y_pos, me)
            closeConnection()
        }
    });

    function sendPos(pos) {
      //check to ensure that the socket variable is present i.e. the browser support tests passed
      if (socket) {
        if (pos !== "") {
          console.log("request send: " + pos)
          socket.send(pos);
          //infodiv.innerHTML += "<p>Sent message to server: '" + message + "'</p>";
        } else {
          alert("You must enter a message to be sent!");
        }
      }
    }

    function oneStep(i,j,k){
        context.beginPath();
        context.arc(15+i*30,15+j*30,13,0,2*Math.PI);//绘制棋子
        var g=context.createRadialGradient(15+i*30,15+j*30,13,15+i*30,15+j*30,0);//设置渐变
        if(k){                           //k=true是黑棋，否则是白棋
            g.addColorStop(0,'#0A0A0A');//黑棋
            g.addColorStop(1,'#636766');
        }else {
            g.addColorStop(0,'#D1D1D1');//白棋
            g.addColorStop(1,'#F9F9F9');
        }
        context.fillStyle=g;
        context.fill();
        context.closePath();

        if(me){
            chessBox[i][j]=1;
        }else{
            chessBox[i][j]=2;
        }
        me = !me
    }
    chess.onclick=function(e){
        var x = e.offsetX;//相对于棋盘左上角的x坐标
        var y = e.offsetY;//相对于棋盘左上角的y坐标
        var i = Math.floor(x/30);
        var j = Math.floor(y/30);
        if( chessBox[i][j] == 0 ) {
            var nex = 'a'.charCodeAt(0);
            pos = String.fromCharCode(nex + j) + i.toString();
            console.log("*** Current Pos: " + pos)
            x_pos = i;
            y_pos = j;
            sendPos(pos)
        }
    }
}

//function to send a message to the websocket server


function sendMessage() {
  //check to ensure that the socket variable is present i.e. the browser support tests passed
  if (socket) {
    //get the message text input element
    var message = document.getElementById("message").value;

    if (message !== "") {
      console.log("request send: " + message)
      socket.send(message);
      //infodiv.innerHTML += "<p>Sent message to server: '" + message + "'</p>";
    } else {
      alert("You must enter a message to be sent!");
    }
  }
}

function closeConnection() {
  //check to ensure that the socket variable is present i.e. the browser support tests passed
  console.log("Got close request")
  if (socket) {
    socket.close();
  }
}
