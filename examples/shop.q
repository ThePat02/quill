LET shop_name = "The Enchanted Emporium"
LET wallet = 50
LET has_sinister_key = FALSE

RANDOM {
    { SHOPKEEP: "Welcome! How can I help you today?" },
    { SHOPKEEP: "Hello there! Looking for something special?" },
    { SHOPKEEP: "Greetings! What brings you to my shop today?" }
}

SYSTEM: "You enter {shop_name}!"

LABEL shop
SYSTEM: "Wallet: {wallet} coins"
CHOICE {
    "Sinister Key (30)" {
        IF wallet >= 30 {
            IF has_sinister_key {
                SHOPKEEP: "You already have a Sinister Key."
                GOTO shop
            }
            wallet -= 30
            has_sinister_key = TRUE
        } ELSE { GOTO insufficient_funds }
    },
    "Mystic Potion (20)" {
        IF wallet >= 20 {
            wallet -= 20
        } ELSE { GOTO insufficient_funds }
    },
    "Healing Herb (10)" {
        IF wallet >= 10 {
            wallet -= 10
        } ELSE { GOTO insufficient_funds }
    },
    "Goodbye!" { GOTO leave_shop }
}

SHOPKEEP: "Thank you for your purchase! Is there anything else I can assist you with?"

LABEL leave_shop

RANDOM {
    { SHOPKEEP: "Thank you for visiting! Come back soon!" },
    { SHOPKEEP: "Take care! Hope to see you again!" },
    { SHOPKEEP: "Farewell! Don't forget to tell your friends about us!" }
}

END

LABEL insufficient_funds
SHOPKEEP: "Oh, it seems you don't have enough coins for that item."
GOTO shop