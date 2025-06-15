# Heyo! This is the most basic script that I need to get to work!

LABEL start

ALEX: "Welcome to our little gathering!"
BELLA: "Thanks for having us, Alex!"
CHARLIE: "Hey everyone!"

ALEX: "Should we play a game?"
BELLA: "Oh, that sounds fun!"

CHOICE {
    "Let's play a trivia game" {
        ALEX: "Great choice! Here's a question..."
        CHARLIE: "I love trivia!"
        
        CHOICE {
            "Continue with trivia" { GOTO trivia_path },
            "Maybe something else" { GOTO party_games }
        }
    },
    "How about party games?" {
        BELLA: "Perfect! I know some fun ones!"
        GOTO party_games
    },
    "I should probably go..." { }
}

LABEL trivia_path
ALEX: "What's the capital of France?"

CHOICE {
    "Paris" {
        ALEX: "Correct! Well done!"
        GOTO start
    },
    "London" {
        CHARLIE: "Not quite..."
        GOTO start
    }
}

LABEL party_games
BELLA: "Let's play charades!"

CHOICE {
    "Sure, I'll start!" { GOTO start },
    "Another time maybe..." { }
}

END
