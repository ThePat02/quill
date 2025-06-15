# Example Syntax

Greeta: "Oh, hello there!"
Greeta: "Look who is here!"

LABEL options
Greeta: [
    "So, what do you want to know?",
    "How can I help you?",
    "Anything on your mind?",
    "Shoot."
]

CHOICE {
    "Who are you?" { GOTO introduction },
    "What is this?" {
        Greeta: "This is a dialog language test!"
        GOTO options
    },
    "Leave me alone!" {
        Greeta: "Oh, no worries then!"
        END
    }
}

LABEL introduction
Greeta: "I am Greeta, the greeter!"
GOTO options