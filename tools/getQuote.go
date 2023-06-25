package tools

import "math/rand"

var quotes = []string{
	"You create your own opportunities. Success doesn’t just come and find you–you have to go out and get it.",
	"Never break your promises. Keep every promise; it makes you credible.",
	"You are never as stuck as you think you are. Success is not final, and failure isn’t fatal.",
	"Happiness is a choice. For every minute you are angry, you lose 60 seconds of your own happiness.",
	"Habits develop into character. Character is the result of our mental attitude and the way we spend our time.",
}

func GetQuote() string {
	return quotes[rand.Intn(len(quotes))]
}
