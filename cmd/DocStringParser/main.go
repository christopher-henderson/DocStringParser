package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/christopher-henderson/DocStringParser/compiler"
	"github.com/christopher-henderson/DocStringParser/tokenizer"
)

// const simpleDoc = `
// /*
// Passing Tests

// @column name: the name of the test which passed
// @column awesomeness: this query is 🔥🔥🔥💯💯💯
// */
// SELECT name, date
// FROM testing
// WHERE pass = true;
// `

// const hackDoc = `
// /**
//   @title "Cool Query"
//   @description "Testing out markup!!"
//   @parameter number "Number" "Which number should I output?"
// */
// /**
//   @table {
//     @title "A Number?"
//     @description "A Number!"
//     @column number "The Best Number" "The number you input"
//   }
// */
// SELECT ${number} as number;

// /**
//   @table {
//     @title "A Number?"
//     @description "A Number!"
//     @column squared_number "The Bestest Number" "The number you input SQUARED!"
//   }
// */
// SELECT ${number}^2 as squared_number;
// `
//
//

const wholeDoc = `
some things
that are
maybe queries
this is a haiku
on mount fuji

/**
@title "HEY NOW, YOU'RE A ROCKSTAR"
@description "GET YOUR SHOW ON, GET PAAAAAAAID"
@param you "the rockstar" "the person who should get their show on"
@param lol "bob" "alice"

@table {
  @title "All that glitters"
  @description "Is gold"
  @column name "rockstar name" "The name of all of the rockstars"
}
*/
SELECT 1 FROM Tests;
`

func test() {
	tok := tokenizer.NewTokenizer(strings.NewReader(wholeDoc))
	err := tok.Tokenize()
	if err != nil {
		log.Panic(err)
	}
	tree, err := compiler.Compile(tok.Tokens())
	if err != nil {
		log.Panic(err)
	}
	j, err := json.Marshal(tree)
	if err != nil {
		log.Panic(err)
	}
	log.Println(string(j))
}

func compile(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	tok := tokenizer.NewTokenizer(req.Body)
	err := tok.Tokenize()
	if err != nil {
		log.Panic(err)
	}
	tree, err := compiler.Compile(tok.Tokens())
	if err != nil {
		log.Panic(err)
	}
	j, err := json.Marshal(tree)
	if err != nil {
		log.Panic(err)
	}
	w.Write(j)
}

func main() {
	http.HandleFunc("/compile", compile)
	log.Println("Starting in server mode.")
	port := fmt.Sprintf(":%v", os.Getenv("PORT"))
	log.Printf("Listening on port %v\n", port)
	log.Fatal(http.ListenAndServe(":1337", nil))
}
