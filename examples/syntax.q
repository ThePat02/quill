# Syntax Example

# Labels are used to define sections of dialogue that can be jumped to
LABEL start

# Tags can be defined with square brackets after dialogue or choices
ALEX: "Welcome to our little gathering!" [tag1]
BELLA: "Thanks for having us, Alex!" [tag1, tag2]
CHARLIE: "Hey everyone!"

ALEX: "Should we play a game?"
BELLA: "Oh, that sounds fun!"

# RANDOM lets you randomly choose between different dialogue options
RANDOM {
    { BELLA: "How about a trivia game?" } [tag1, tag2],
    { CHARLIE: "Let's play some games!" } [tag3],
    {
        GHOST: "This is the secret third string!"
    }
}

# CHOICE lets you present options to the player
CHOICE {
    "Let's play a trivia game" {
        ALEX: "Great choice! Here's a question..."
        CHARLIE: "I love trivia!"
        
        CHOICE {
            "Continue with trivia" { GOTO trivia_path } [tag1, tag2],
            "Maybe something else" { GOTO party_games }
        }
    },
    "How about party games?" {
        BELLA: "Perfect! I know some fun ones!"
        GOTO party_games
    },
    "I should probably go..." { GOTO ending }
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

LABEL ending

BELLA: "Thanks for joining us! See you next time!"

END
