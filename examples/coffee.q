# Simple Chat Example

LABEL start

ALEX: "Hey there!"
BELLA: "Hi Alex!"

CHOICE {
    "How's your day?" {
        BELLA: "Pretty good, thanks!"
        GOTO end
    },
    "Want to hang out?" {
        ALEX: "Sure! What should we do?"
        
        RANDOM {
            { BELLA: "Let's get coffee!" },
            { BELLA: "How about a walk?" }
        }
        
        GOTO end
    }
}

LABEL end

ALEX: "See you later!"

END