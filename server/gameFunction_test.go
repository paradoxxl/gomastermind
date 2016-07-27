package server

import(
	"testing"
	"github.com/paradoxxl/gomastermind/msg"
	"github.com/stretchr/testify/assert"
)



func TestGuessFunctionErrors(t *testing.T){
	assert := assert.New(t)

	code := msg.Code{0,1,2}
	maxtries:=2
	numcolors:=3
	codelen := 3
	game := msg.Game{Code:&code,Tries:0,Attributes:&msg.NewGameMsg{NbrColors:numcolors,CodeLength:codelen,MaxTries:maxtries}}

	clihandler := ClientHandler{game:&game}


	//Test empty guess
	assert.Equal(*clihandler.guess(nil),"6|{\"Message\":\"Another codelength expected\",\"ErrorType\":1}")

	//Test guess out of range
	guess := &msg.Code{0,1,3}
	assert.Equal(*clihandler.guess(guess),"6|{\"Message\":\"Guesses must be in correct range\",\"ErrorType\":1}")

	//Test out of tries (no game anymore)
	guess = &msg.Code{0,1,1}
	clihandler.guess(guess)
	guess = &msg.Code{0,1,1}
	clihandler.guess(guess)
	guess = &msg.Code{0,1,2}
	assert.Equal(*clihandler.guess(guess),"6|{\"Message\":\"No active Game\",\"ErrorType\":2}")

	//Test no game exists
	guess = &msg.Code{0,1,2}
	assert.Equal(*clihandler.guess(guess),"6|{\"Message\":\"No active Game\",\"ErrorType\":2}")
}


func TestGuessFunctionSuccess(t *testing.T){
	assert := assert.New(t)

	code := msg.Code{0,1,2}
	maxtries:=2
	numcolors:=3
	codelen := 3
	game := msg.Game{Code:&code,Tries:0,Attributes:&msg.NewGameMsg{NbrColors:numcolors,CodeLength:codelen,MaxTries:maxtries}}

	clihandler := ClientHandler{game:&game}

	guess := &msg.Code{0,1,2}

	//Test instant success
	guess = &msg.Code{0,1,2}
	assert.Equal(*clihandler.guess(guess),"3|{\"Won\":true,\"Tries\":1,\"Code\":[0,1,2]}")
}

func TestGuessFunctionLogic(t *testing.T){
	assert := assert.New(t)

	code := msg.Code{0,1,2,1,4,1}
	maxtries:=20
	numcolors:=6
	codelen := 6
	game := msg.Game{Code:&code,Tries:0,Attributes:&msg.NewGameMsg{NbrColors:numcolors,CodeLength:codelen,MaxTries:maxtries}}
	clihandler := ClientHandler{game:&game}



	guess := &msg.Code{1,3,4,1,2,3}
	assert.Equal(*clihandler.guess(guess),"2|{\"Answer\":[0,1,1,1,2,2]}")

	guess = &msg.Code{1,1,1,1,1,1}
	assert.Equal(*clihandler.guess(guess),"2|{\"Answer\":[0,0,0,2,2,2]}")

	guess = &msg.Code{0,1,2,3,1,3}
	assert.Equal(*clihandler.guess(guess),"2|{\"Answer\":[0,0,0,1,2,2]}")

	guess = &msg.Code{1,0,4,2,1,3}
	assert.Equal(*clihandler.guess(guess),"2|{\"Answer\":[1,1,1,1,1,2]}")

	guess = &msg.Code{0,1,2,1,4,1}
	assert.Equal(*clihandler.guess(guess),"3|{\"Won\":true,\"Tries\":5,\"Code\":[0,1,2,1,4,1]}")
}






