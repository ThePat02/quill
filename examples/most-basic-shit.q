# Heyo! This is the most basic script that I need to get to work!

LABEL start

BB: "Hello there!"
BB: "How are you doing?"
BH: "I am here too!"

CHOICE {
    "Again, again!" { GOTO start },
    "Please make it stop..." { END }
}
