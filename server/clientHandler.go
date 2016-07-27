package server

import (
	"golang.org/x/net/websocket"
	"log"
	"fmt"
	"math/rand"
	"github.com/paradoxxl/gomastermind/msg"
)

type ClientHandler struct{
	ws *websocket.Conn
	game *msg.Game
}

func (clientHandler *ClientHandler) Listen(){
	for{
		var message string
		if err:= websocket.Message.Receive(clientHandler.ws,&message);err!= nil {
			log.Printf("cannot receive: %v",err)
			return
		}
		fmt.Printf("Receive: %s\n", message)

		msg.ReadMessage(&message,clientHandler)
	}

}

func (clientHandler *ClientHandler) Send(msg *string) (error){
	if err:= websocket.Message.Send(clientHandler.ws,*msg);err!= nil {
		log.Printf("cannot send: %v",err)
		return err
	}
	fmt.Printf("Send: %s\n", *msg)
	return nil
}

func (clientHandler *ClientHandler) OnCommand(commandType byte, message interface{}){

	switch commandType {
	case msg.NewGameType:
		data, ok := message.(*msg.NewGameMsg)
		if ok {
			log.Printf("converted ngm : %+v ",message)
			clientHandler.newGame(data.NbrColors,data.CodeLength,data.MaxTries)
		}else{
			log.Printf("cannot convert ngm: %+v ",message)
		}
	case msg.GuessType:
		data, ok := message.(*msg.GuessMsg)
		if ok {
			log.Printf("converted gm: %+v ",message)
			clientHandler.Send(clientHandler.guess(&data.Guess))
		}else{
			log.Printf("cannot convert gm: %+v ",message)
		}
	case msg.RetreatType:
		clientHandler.retreat()
	case msg.ErrorType:
		data, ok := message.(*msg.ErrorMsg)
		if ok {
			clientHandler.errorHappened(data)
		}else{
			clientHandler.errorHappened(&msg.ErrorMsg{Message:"I have no idea that happened", ErrorType:msg.EverythingIsBrokenError})
		}
	default:
		clientHandler.errorHappened(&msg.ErrorMsg{ErrorType: msg.BadParametersError, Message: "Unsupported command type"})
	}


}


func (state *ClientHandler) errorHappened(message *msg.ErrorMsg){
	log.Printf("error:: %+v",message)
	data := msg.EncodeErrorMsg(message)
	log.Printf("error:: %v",data)
	err := state.Send(data)
	if err != nil {
		//TODO: Do something
		log.Printf("cannot send Error msg: %v",err)
	}
}

func (state *ClientHandler) newGame(nbrColors,codeLen,maxTries int) {

	// Create Code
	var code = make(msg.Code, codeLen)
	for i := 0; i < codeLen; i++ {
		code[i] = rand.Intn(int(nbrColors))
	}
	state.game = &msg.Game{&code, &msg.NewGameMsg{NbrColors:nbrColors, CodeLength:codeLen, MaxTries:maxTries}, maxTries}
	data := msg.EncodeAcceptMsg(&msg.AcceptMsg{Accept:true, Attributes:state.game.Attributes})

	err := state.Send(&data)
	if err != nil {
		//TODO: Do something
	}
}

func (state *ClientHandler) guess(guess *msg.Code) *string {
	//Check if the length of the guess is the same as the codelength
	if state.game == nil {
		return msg.EncodeErrorMsg(&msg.ErrorMsg{Message:"No active Game", ErrorType:msg.NoActiveGameError})
	}
	if guess == nil || len(*guess) != len(*state.game.Code) {
		return msg.EncodeErrorMsg(&msg.ErrorMsg{Message:"Another codelength expected", ErrorType:msg.BadParametersError})
	}


	var data string
	var length = len(*guess)
	var codeCopy = make(msg.Code, length)
	copy(codeCopy, *state.game.Code)

	var j = 0
	var response = make(msg.Answer, len(*guess))
	// Check nbr of exact matches
	for p, v := range (*guess) {
		if v < 0 || v >= state.game.Attributes.NbrColors {
			return msg.EncodeErrorMsg(&msg.ErrorMsg{Message:"Guesses must be in correct range", ErrorType:msg.BadParametersError})
		}
		if v == (*state.game.Code)[p] {
			response[j] = msg.CorrectGuess
			j++
			codeCopy[p] = -1
			(*guess)[p] = -1
		}
	}

	//Check if game is finished
	if ( j == length) {
		state.game.Tries+=1
		data = msg.EncodeGameEndMsg(&msg.GameEndMsg{Won:true, Code:state.game.Code,Tries:state.game.Tries})
		state.game = nil
		return &data
	} else {
		if state.game.Tries + 1 == state.game.Attributes.MaxTries {
			state.game.Tries+=1
			data = msg.EncodeGameEndMsg(&msg.GameEndMsg{Won:false, Code:state.game.Code,Tries:state.game.Tries})
			state.game = nil
			return &data
		} else {
			// Check nbr of of color matches on wrong places
			for _, v := range (*guess) {
				if x := find(v, &codeCopy); x != -1 {
					response[j] = msg.ColorExistsGuess
					j++
					codeCopy[x] = -1
				}
			}

			//Rest has to be wrong guesses
			for ; j < length; j++ {
				response[j] = msg.WrongGuess
			}

			data = msg.EncodeAnswerMsg(&msg.AnswerMsg{Answer:&response})
		}
	}

	//fmt.Printf("Msg: %v, Guessstate: %v, Answer: %v,", data, *guess, response)
	state.game.Tries++
	return &data
}

func find(value int,slice *msg.Code) int {
	if value == -1 {return value}
	for p, v := range *slice {
		if (v == value) {
			return p
		}
	}
	return -1
}


func (state *ClientHandler) retreat(){
	data := msg.EncodeGameEndMsg(&msg.GameEndMsg{Won:false,Code:state.game.Code})
	err := state.Send(&data)
	if err != nil {
		//TODO: Do something
	}


	state.game = nil
}