package main

import (
	"fmt"
	"log"
	"net/http"
    "strconv"
    "flag"
    "time"
    "math/rand"

	"golang.org/x/net/websocket"
)

/*
character a~n : row
number 0~14: col
*/

const (
    SIDE_LEN = 15
    MAX_CHESS = 225

    MSG_INVALID = "invalid"
    MSG_OK = "ok"
    MSG_WIN = "win"
    MSG_LOSE = "lose"
    MSG_DRAW = "draw"
)

const (
    _Color = iota
    EMPTY
    BLACK
    WHITE
)

var BOARD = [SIDE_LEN][SIDE_LEN]Color{}
var CUR_COLOR Color
var CUR_CHESS = 0
var FORBIDDEN_CHESS = 0
var AI_FLAG bool

func Max(x, y int) int {
    if x > y {
        return x
    }
    return y
}

func initBoard() {
    for i := 0; i < SIDE_LEN; i++ {
        for j := 0; j < SIDE_LEN; j++ {
            BOARD[i][j] = EMPTY
        }
    }
    CUR_COLOR = BLACK
    fmt.Println("Board initalized successfully !")
}

func setNextColor() {
    if CUR_COLOR != BLACK && CUR_COLOR != WHITE {
        fmt.Println("Current color: ", CUR_COLOR)
        log.Fatal("Current Color is invalid! Exit.")
    }

    if CUR_COLOR == BLACK {
        CUR_COLOR = WHITE
    } else {
        CUR_COLOR = BLACK
    }
}

func setCurrentChess(p Pos) {
    BOARD[p.x][p.y] = CUR_COLOR
    CUR_CHESS += 1
}

func getPos(msg []byte, n int) (Pos, string) {
    var p Pos
    p.x = int(msg[0] - 'a')
    p.y, _ = strconv.Atoi(string(msg[1:n]))
    return p, string(msg[0]) + string(msg[1:n])
}

func twoWaySearch(p Pos, dx int, dy int) int {
    //fmt.Println("Start two way search!")
    score := 1
    cur_x := p.x
    cur_y := p.y
    for i := 0; i < 5; i++ {
        cur_x += dx
        cur_y += dy

        if cur_x < 0 || cur_x >= SIDE_LEN || cur_y < 0 || cur_y >= SIDE_LEN  {
            break
        }
        if BOARD[cur_x][cur_y] != CUR_COLOR {
            break
        }
        score += 1
        fmt.Println(cur_x, cur_y)
    }

    cur_x = p.x
    cur_y = p.y
    for i := 0; i < 5; i++ {
        cur_x -= dx
        cur_y -= dy

        if cur_x < 0 || cur_x >= SIDE_LEN || cur_y < 0 || cur_y >= SIDE_LEN  {
            break
        }
        if BOARD[cur_x][cur_y] != CUR_COLOR {
            break
        }
        score += 1
        fmt.Println(cur_x, cur_y)
    }
    return score
}

func canWin(p Pos) bool {
    // 8 directions
    score := 1
    score = Max(score, twoWaySearch(p, 1, 0))
    score = Max(score, twoWaySearch(p, 0, 1))
    score = Max(score, twoWaySearch(p, 1, 1))
    score = Max(score, twoWaySearch(p, -1, 1))
    return score >= 5
}

func isDraw() bool {
    return CUR_CHESS + FORBIDDEN_CHESS >= MAX_CHESS
}

func isAvailable(p Pos) bool {
    return BOARD[p.x][p.y] == EMPTY
}

func sendMsg(ws *websocket.Conn, src string) {
    errMsg := []byte(src)
    errLen := len(errMsg)
    e, err := ws.Write(errMsg[:errLen])
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Send: %s\n", errMsg[:e])
}

func getRandomPos(positions []Pos) Pos {
    size := len(positions)
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    retPos := positions[r.Intn(size)]
    return retPos
}

func stupidMethod() []Pos {
    var positions []Pos
    for i := 0; i < SIDE_LEN; i++ {
        for j := 0; j < SIDE_LEN; j++ {
            if BOARD[i][j] == EMPTY {
                positions = append(positions, Pos{x:i,y:j})
            }
        }
    }
    return positions
}

func aiCalcNextStep() (Pos, string) {
    fmt.Println("AI is thinking...")
    var candidates []Pos
    candidates = stupidMethod()
    retPos := getRandomPos(candidates)
    retStr := string(byte('a' + retPos.x)) + strconv.Itoa(retPos.y)
    fmt.Println("AI Got: ", retStr)
    return retPos, retStr
}

func eventHandler(ws *websocket.Conn) {
    for {
    	msg := make([]byte, 512)
    	n, err := ws.Read(msg)
    	if err != nil {
    		log.Fatal(err)
    	}
    	fmt.Printf("Receive: %s\n", msg[:n])

        curPos, curPosStr := getPos(msg, n)
        fmt.Println("Got pos: ", curPos.x, ", ", curPos.y)
        if !isAvailable(curPos) {
            sendMsg(ws, MSG_INVALID)
            continue
        }

        setCurrentChess(curPos)
        if canWin(curPos) {
            sendMsg(ws, MSG_WIN)
            return
        }

        if isDraw() {
            sendMsg(ws, MSG_DRAW)
            return
        }

    	sendMsg(ws, MSG_OK + "," + curPosStr)
        setNextColor()

        if(AI_FLAG) {
            curPos, curPosStr = aiCalcNextStep()
            setCurrentChess(curPos)
            if canWin(curPos) {
                sendMsg(ws, MSG_WIN)
                return
            }

            if isDraw() {
                sendMsg(ws, MSG_DRAW)
                return
            }

            sendMsg(ws, MSG_OK + "," + curPosStr)
            setNextColor()
        }
    }
}

func main() {
    flag.BoolVar(&AI_FLAG, "enableai", false, "enable or disable AI") // --enableai=false
    flag.Parse()
    if(AI_FLAG) {
        fmt.Println("AI Enabled!")
    } else {
        fmt.Println("AI Disabled!")
    }

    fmt.Println("Starting Server......")
    initBoard()
	http.Handle("/echo", websocket.Handler(eventHandler))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
