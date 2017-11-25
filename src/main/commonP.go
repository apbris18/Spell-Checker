//common_prgram.go
//Jessica Forrett and Alex Brisbois
//CSCI 324 Professor King

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// Element is an element of a linked list.
type Element struct {
	// Next and previous pointers in the doubly-linked list of elements.
	// To simplify the implementation, internally a list l is implemented
	// as a ring, such that &l.root is both the next element of the last
	// list element (l.Back()) and the previous element of the first list
	// element (l.Front()).
	next, prev *Element

	// The list to which this element belongs.
	list *List

	// The value stored with this element.
	Value interface{}
}

// List represents a doubly linked list.
// The zero value for List is an empty list ready to use.
type List struct {
	root Element // sentinel list element, only &root, root.prev, and root.next are used
	len  int     // current list length excluding (this) sentinel element
}

// Next returns the next list element or nil.
func (e *Element) Next() *Element {
	if p := e.next; e.list != nil && p != &e.list.root {
		return p
	}
	return nil
}

// Init initializes or clears list l.
func (l *List) Init() *List {
	l.root.next = &l.root
	l.root.prev = &l.root
	l.len = 0
	return l
}

// New returns an initialized list.
func New() *List { return new(List).Init() }

// Front returns the first element of list l or nil.
func (l *List) Front() *Element {
	if l.len == 0 {
		return nil
	}
	return l.root.next
}

// insert inserts e after at, increments l.len, and returns e.
func (l *List) insert(e, at *Element) *Element {
	n := at.next
	at.next = e
	e.prev = at
	e.next = n
	n.prev = e
	e.list = l
	l.len++
	return e
}

// insertValue is a convenience wrapper for insert(&Element{Value: v}, at).
func (l *List) insertValue(v interface{}, at *Element) *Element {
	return l.insert(&Element{Value: v}, at)
}

func dictionary_holder(word string) bool {
	//add dictionary words into list
	dictionary_list := New() //init a new list

	file, err := os.Open("src/main/dictionary.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	//scan and put words into list
	scanner := bufio.NewScanner(file)
	for scanner.Scan() { //LOOK HERE -- maybe do something like  for main() in setting up conditionals for reading a word til /n
		line := scanner.Text()
		if line == "\n" {
			break
		}
		list_element := Element{nil, nil, dictionary_list, line} //want to pull from dictionary_file

		dictionary_list.insert(&list_element, dictionary_list.root.prev)
	}

	for e := dictionary_list.Front(); e != nil; e = e.Next() {
		//if current Element matches "word" then return true
		if (e.Value) == word {
			return true
		}
	}
	return false
} //end dictionary_holder

func check(e error) {
	if e != nil {
		fmt.Println("An error has occurred")
		panic(e)
	}
}
func qSrt(anArray []string) []string {
	length := len(anArray)

	if length <= 1 {
		return anArray
	}

	p := anArray[length/2]

	lthan := make([]string, 0, length)
	eql := make([]string, 0, length)
	gthan := make([]string, 0, length)

	for _, word := range anArray {
		switch {
		case word < p:
			lthan = append(lthan, word)
		case word == p:
			eql = append(eql, word)
		case word > p:
			gthan = append(gthan, word)
		}
	}

	lthan, gthan = qSrt(lthan), qSrt(gthan)
	lthan = append(lthan, eql...)
	return (append(lthan, gthan...))
}

func main() {
	var word string = ""
	var endSentence bool = true
	var out bool = false
	//use array instead
	missed := make([]string, 0)
	//reads in the file location from the user
	fmt.Println("Please enter the location of the text file you would like us to spellcheck.")
	scan := bufio.NewScanner(os.Stdin)
	scan.Scan()
	var fileLoc string = scan.Text()

	//reads from the file
	fmt.Println("Opening file")
	fmt.Println("")
	fl, err := os.Open(fileLoc)
	check(err)
	defer fl.Close()

	fmt.Println("Would you like the spell check to be case sensitive? (y/n)")
	caseScan := bufio.NewScanner(os.Stdin)
	caseScan.Scan()
	var sens string = caseScan.Text()

	for sens != "y" && sens != "n" {
		fmt.Println("Sorry, I didn't catch that. Please enter y or n")
		scan := bufio.NewScanner(os.Stdin)
		scan.Scan()
		sens = scan.Text()
	}

	for {
		lttrb := make([]byte, 1)
		_, err := fl.Read(lttrb)
		if err == io.EOF {
			break
		}
		var lttrs string = string(lttrb)
		if lttrs == " " || lttrs == "?" || lttrs == "." || lttrs == "," || lttrs == "!" || lttrs == ":" || lttrs == ";" || lttrs == "\n" || lttrs == "\"" {
			//CHECK HERE TOO -- should be good condition too
			//first check if word == ""
			//if not, call function to check if it's in the dictionary
			//if it's not in the dictionary, store the word into a data structure

			if word != "" {
				if sens == "n" {
					out = (dictionary_holder(strings.ToLower(word)) || dictionary_holder(strings.ToUpper(word)))
				} else if endSentence == true {
					out = (dictionary_holder(strings.ToLower(word)) || dictionary_holder(word))
				} else {
					out = dictionary_holder(word)
				}
				endSentence = false
				// I know dictionary_holder is written well
				//call dictionary_holder.go
				//if not, check if the last word is in the dictionary. If it's not, store it in the data structure
				if out == false {
					//fmt.Println("It seems that " + word + " isn't a word.  We are adding that to a list of misspelled words.")
					missed = append(missed, word)

				} //end if out
				word = ""
			}
			if lttrs == "." || lttrs == "!" || lttrs == "?" || lttrs == "\n" {
				endSentence = true
			}
			//}
		} else {
			word = word + lttrs
		}
	}
	mLength := len(missed)
	if mLength == 0 {
		fmt.Println("Everything's right!")
	} else {
		//QUICK-SORT, then print out the misspelled words
		fmt.Println("Sorting out the misspelled words")
		fmt.Println("")
		missed = qSrt(missed)
		for i := 0; i < mLength; i++ {
			fmt.Println(missed[i])
		}
		fmt.Println("")
		fmt.Println("The misspelled words have been organized into a linked list.")
		fmt.Println("")
	}

}
