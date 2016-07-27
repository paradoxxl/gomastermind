# Game mechanics
The server generates a code which has to be broken within a maximum amount of tries by the user. 

## Guessing
After a guess, the server answerd with hints. In the board game, different coloured pins are used for this but this server uses an array of 'AnswerTypes'. 
The Answer-Array is produced as follows. The code in the example is **"0,1,2,1,4,1"**. The guess is **"1,3,4,1,2,3"**

    1) Search for all exact matches, e.g. "x,x,x,1,x,x". As there is  one exact match, the first entry in the answer will be the constant 'CorrectGuess' 
    2) Search for all colour matches with wring positions. It is important to omit the findings from step 1. The maximum amount of findings per 'color' equals to the corresponding number of the same 'colour' which is left in the code after step 1. Example: "1,x,4,-,2,x" are the relative matches. Therefore, the next three entries in the answer are the constant ColorExistsGuess.
    3) The remaining entries of the answer are filles with the constant WrongGuess.
    
    Answer in the example: "0,1,1,1,2,2"
## Client          <->          Server
    Connect         ->
    NewGame         ->
                    <-          Accept/Error
    Guess           ->
                    <-          Answer/Error/GameEnd
    Retreat         <-          GameEnd



## Messages
Messages have to be in a well-defined format.

    Type|MarshalledJSON

### Types
	NewGameType     0
	GuessType       1
	AnswerType      2
	GameEndType     3
	RetreatType     4
	AcceptType      5
	ErrorType       6
	
### ErrorTypes
    ParseError byte             0
	BadParametersError          1
	NoActiveGameError           2
	EverythingIsBrokenError     3

### AnswerTypes
    CorrectGuess        0
	ColorExistsGuess    1
	WrongGuess          2
	
## JSON Messages
### Client -> Server
    NewGameMsg: {int NbrColors, nit MaxTries, int CodeLength}
    GuessMsg: {int[] Guess}
    RetreatMsg: {}
### Client <- Server
    AcceptMsg: {bool Accept, Attributes: {int NbrColors, nit MaxTries, int CodeLength}}   
    GameEndMsg: {bool Won, int Tries, int[] Code}
    AnswerMsg: {int[] Answer}
    ErrorMsg: {string Message, byte ErrorType}
    
## Example messages
The following message signals the winning state after a successful guess within the maximum tries. The correct Code was "0,1,2" in this example.

    3|{"Won":true,"Tries":1,"Code":[0,1,2]}
If the guess is within the maximum tries but is different to the code, the server returns an answer-message. The array is produced accorting to the game mechanics.

    2|{"Answer":[0,1,1,1,2,2]}
Errors may happen due to guessing without having an existing game or by providing bad parameters, e.g. the parameters are out of range. Errors have an additional description to determine what was wrong.

    6|{"Message":"Guesses must be in correct range","ErrorType":1}



    