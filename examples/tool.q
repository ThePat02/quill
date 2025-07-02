# Simple tool call that retrieves the player's name with fallback
LET player_name = <getPlayerName;>
LET player_age = <getPlayerAge;> ?? 18
LET age_plus_five = <agePlusFive; player_age>

SYSTEM: "Hello, {player_name}! You are {player_age} years old."
IF player_age < 18 {
    SYSTEM: "You are quite young to be here!"
}

IF player_age >= 18 {
    SYSTEM: "You are old enough to be here."
}

# Advanced inline tool call. Syntax: <function; argument>
SYSTEM: "Your current gold balance is <getData; "gold"> gold coins."

# Multiple tool calls in a single line
SYSTEM: "Your current gold balance is <getData; "gold"> gold coins and your health is <getData; "health">."

# Multiple tool call arguments
# Tool function getItemPrice takes two arguments: item type and item level
LET item_price = <getItemPrice; "potion", 4>
SYSTEM: "The price of a level 4 potion is {item_price} gold coins."