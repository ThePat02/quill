MrsBennet: "My dear Mr. Bennet, have you heard that Netherfield Park is let at last?"

CHOICE {
    "Yes." { MrBennet: "I think I have!" } [lie, noncanonical],
    "No." {
        MrBennet: "I have not!"
        GOTO news
    }
}

MrsBennet: "Oh, you haven't! I am quiet sure of it!"
MrBennet: "Well, you caught me, but I cannot believe it is true."

LABEL news

MrsBennet: "But it is, for Mrs. Long has just been here, and she told me all about it."
MrsBennet: "Do you not want to know who has taken it?"
MrBennet: "You want to tell me, and I have no objection to hearing it."

END