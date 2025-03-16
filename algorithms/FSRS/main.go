package main

import (
	"cmp"
	"encoding/binary"
	"slices"
	"time"
	"unsafe"

	"github.com/open-spaced-repetition/go-fsrs/v3"
)

// based on https://github.com/open-spaced-repetition/go-fsrs/wiki/Usage
var memAddress map[uintptr][]byte = make(map[uintptr][]byte)

//go:export alloc
func Alloc(size uint32) int64 {
	a := make([]byte, size)
	memAddress[uintptr(unsafe.Pointer(&a[0]))] = a
	return int64(uintptr(unsafe.Pointer(&a[0])))
}

//go:export dealloc
func Dealloc(ptr uint64) {
	delete(memAddress, uintptr(ptr))
}

var w *fsrs.FSRS = fsrs.NewFSRS(fsrs.DefaultParam())
var cards map[uint64]fsrs.Card = make(map[uint64]fsrs.Card)
var cardsGraded map[uint64]struct{} = make(map[uint64]struct{})

//go:export add-card
func AddCard(id uint64) {
	card := fsrs.NewCard()
	cards[id] = card
}

//go:export grade-card
func GradeCard(id uint64, t_unix uint64, correct uint64) {
	t := time.Unix(int64(t_unix), 0)
	cardsGraded[id] = struct{}{}
	schedulingCards := w.Repeat(cards[id], t)

	l := fsrs.Again

	if correct == 0 {
		l = fsrs.Again
	} else if correct == 1 {
		l = fsrs.Good
	} else {
		panic("correct is incorrect")
	}
	cards[id] = schedulingCards[l].Card
}

type ThickCard struct {
	Card    uint64
	dueDate time.Time
}

func doGetCards(n_unix int64) (uint32, uint32, uint32, uint64) {
	n := time.Unix(n_unix, 0)
	dueCards := []ThickCard{}
	nonDueCards := []ThickCard{}
	newCards := []ThickCard{}
	for id, _ := range cardsGraded {
		card := cards[id]
		if card.Due.Before(n) {
			dueCards = append(dueCards, ThickCard{Card: id, dueDate: card.Due})
		} else {
			nonDueCards = append(nonDueCards, ThickCard{Card: id, dueDate: card.Due})
		}
	}
	for id := range cardsGraded {
		delete(cards, id)
	}
	for card := range cards {
		newCards = append(newCards, ThickCard{Card: card, dueDate: time.Time{}})
	}
	slices.SortFunc(newCards, func(a ThickCard, b ThickCard) int {
		return cmp.Compare(a.Card, b.Card)
	})
	slices.SortFunc(dueCards, func(a ThickCard, b ThickCard) int {
		return a.dueDate.Compare(b.dueDate)
	})
	slices.SortFunc(nonDueCards, func(a ThickCard, b ThickCard) int {
		return a.dueDate.Compare(b.dueDate)
	})
	lenDueCards := uint32(len(dueCards))
	lenNonDueCards := uint32(len(nonDueCards))
	lenNewCards := uint32(len(newCards))
	var dueCard uint64 = 0
	var nonDueCard uint64 = 0
	var newCard uint64 = 0
	if lenDueCards > 0 {
		dueCard = dueCards[0].Card
	}
	if lenNonDueCards > 0 {
		nonDueCard = nonDueCards[0].Card
	}
	if lenNewCards > 0 {
		newCard = newCards[0].Card
	}

	var nextCard uint64 = 0
	if lenDueCards > 0 {
		nextCard = dueCard
	} else if lenNewCards > 0 {
		nextCard = newCard
	} else {
		nextCard = nonDueCard
	}
	return lenDueCards, lenNonDueCards, lenNewCards, nextCard
}

//export get-cards
func GetCards(address int64, n_unix int64) {
	lenDueCards, lenNonDueCards, lenNewCards, nextCard := doGetCards(n_unix)
	b := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(b, lenDueCards)
	writeToMemory(b, address)
	binary.LittleEndian.PutUint32(b, lenNonDueCards)
	writeToMemory(b, address+4)
	binary.LittleEndian.PutUint32(b, lenNewCards)
	writeToMemory(b, address+8)
	bl := [16]byte{}
	binary.LittleEndian.PutUint64(bl[:], nextCard)
	writeToMemory(bl[:], address+16)
}

func writeToMemory(data []byte, memoryAddress int64) {
	for i, v := range data {
		*(*byte)(unsafe.Add(unsafe.Pointer(uintptr(memoryAddress)), i)) = v
	}
}
func main() {
}
