package msg

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	NewGameType byte = iota
	GuessType
	AnswerType
	GameEndType
	RetreatType
	AcceptType
	ErrorType
)

const (
	ParseError byte = iota
	BadParametersError
	NoActiveGameError
	EverythingIsBrokenError
)

const (
	CorrectGuess int = iota
	ColorExistsGuess
	WrongGuess
)

type Game struct {
	Code       *Code
	Attributes *NewGameMsg
	Tries      int
}

// Server -> Client Messages
type AcceptMsg struct {
	Accept     bool
	Attributes *NewGameMsg
}
type GameEndMsg struct {
	Won   bool
	Tries int
	Code  *Code
}
type AnswerMsg struct {
	Answer *Answer
}
type ErrorMsg struct {
	Message   string
	ErrorType byte
}

//Client -> Server Messages
type RetreatMsg struct{}
type NewGameMsg struct {
	NbrColors  int
	MaxTries   int
	CodeLength int
}
type GuessMsg struct {
	Guess Code
}
type Code []	int[]int
type Answer []int
type MessageHandler interface {
	OnCommand(commandType byte, message interface{})
}

func ReadMessage(msg *string, handler MessageHandler) {

	var cmd interface{}
	var msgType int
	var err interface{}
	splitted := strings.Split(*msg, "|")

	//Check whether the message string might be ok
	if len(splitted) != 2 {
		cmd = ErrorMsg{ErrorType: ParseError, Message: "More than one delimiter in string"}
		msgType = int(ErrorType)
		//log.Printf("length of splitstring array  %v instead of 2:", len(splitted))
	} else {
		//Check whether the command type parameter can be converted to int
		if msgType, err = strconv.Atoi(splitted[0]); err != nil {
			cmd = ErrorMsg{ErrorType: ParseError, Message: "Cannot convert the message type to int"}
			msgType = int(ErrorType)
			//log.Printf("cannot convert  %v to int:", splitted[0])
		} else {
			switch byte(msgType) {
			case NewGameType:
				cmd, err = DecodeNewGameMessage([]byte(splitted[1]))
				if err != nil {
					//log.Printf("decode nwm error %v:", err)
					cmd = ErrorMsg{ErrorType: ParseError, Message: "Cannot unmarshall the NewGameMessage"}
					msgType = int(ErrorType)
				} else {
					//log.Printf("decoded nwm  %v:", cmd)
				}
			case GuessType:
				cmd, err = DecodeGuessMessage([]byte(splitted[1]))
				if err != nil {
					//log.Printf("decode guess error %v:", err)
					cmd = ErrorMsg{ErrorType: ParseError, Message: "Cannot unmarshall the GuessMessage"}
					msgType = int(ErrorType)
				} else {
					//log.Printf("decoded guess  %v:", cmd)
				}
			case RetreatType:
				cmd = RetreatMsg{}
			default:
				cmd = ErrorMsg{ErrorType: BadParametersError, Message: "Unsupported command type"}
				msgType = int(ErrorType)
				err = errors.New("Unsupported command type")
			}
		}

	}
	//log.Printf("going to execute oncommand with msgtype = %v, should be %v", msgType, splitted[0])
	handler.OnCommand(byte(msgType), cmd)
}

func DecodeNewGameMessage(msg []byte) (*NewGameMsg, error) {
	var cmd NewGameMsg
	err := json.Unmarshal(msg, &cmd)
	return &cmd, err
}

func DecodeGuessMessage(msg []byte) (*GuessMsg, error) {
	var cmd GuessMsg
	err := json.Unmarshal(msg, &cmd)
	return &cmd, err
}

func EncodeAnswerMsg(data *AnswerMsg) string {
	//log.Printf("going to encode answer now: %+v",data)
	header := []byte(fmt.Sprintf("%v|", AnswerType))
	p, err := json.Marshal(&data)
	//log.Printf("encoded answerdata: %+v",string(p))
	if err != nil {
		panic(err)
	}
	p = append(header, p...)
	return string(p)
}

func EncodeErrorMsg(data *ErrorMsg) *string {
	//log.Printf("going to encode error now: %+v",data)
	header := []byte(fmt.Sprintf("%v|", ErrorType))
	p, err := json.Marshal(&data)
	//log.Printf("encoded errordata: %+v",string(p))
	if err != nil {
		panic(err)
	}
	p = append(header, p...)
	r:= string(p)
	return &r
}

func EncodeAcceptMsg(data *AcceptMsg) string {
	//log.Printf("going to encode answer now: %+v",data)
	header := []byte(fmt.Sprintf("%v|", AcceptType))
	p, err := json.Marshal(&data)
	//log.Printf("encoded answerdata: %+v",string(p))
	if err != nil {
		panic(err)
	}
	p = append(header, p...)
	return string(p)
}

func EncodeGameEndMsg(data *GameEndMsg) string {
	header := []byte(fmt.Sprintf("%v|", GameEndType))
	p, err := json.Marshal(&data)
	if err != nil {
		panic(err)
	}
	p = append(header, p...)
	return string(p)
}

func EncodeNewGameMsg(data *NewGameMsg) string {
	header := []byte(fmt.Sprintf("%v|", NewGameType))
	p, err := json.Marshal(&data)
	if err != nil {
		panic(err)
	}
	p = append(header, p...)
	//binary.BigEndian.PutUint32(p[1:5], uint32(len(p)))
	return string(p)
}
func EncodeGuessMsg(data *GuessMsg) string {
	header := []byte(fmt.Sprintf("%v|", GuessType))
	p, err := json.Marshal(&data)
	if err != nil {
		panic(err)
	}
	p = append(header, p...)
	//binary.BigEndian.PutUint32(p[1:5], uint32(len(p)))
	return string(p)
}
