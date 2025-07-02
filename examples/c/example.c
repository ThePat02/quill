#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "libquill.h"

// Mock tool call handler that provides realistic responses
char* get_mock_tool_result(const char* step_result) {
    // In a real application, you'd parse the JSON to extract function name and arguments
    // For demo purposes, we'll check for common function names in the JSON
    
    if (strstr(step_result, "getPlayerName") != NULL) {
        return "\"Hero\"";
    } else if (strstr(step_result, "getPlayerAge") != NULL) {
        return "25";
    } else if (strstr(step_result, "agePlusFive") != NULL) {
        return "30";
    } else if (strstr(step_result, "getData") != NULL) {
        if (strstr(step_result, "gold") != NULL) {
            return "150";
        } else if (strstr(step_result, "health") != NULL) {
            return "85";
        } else {
            return "\"Unknown\"";
        }
    } else if (strstr(step_result, "getItemPrice") != NULL) {
        return "20";
    } else {
        return "\"DefaultValue\"";
    }
}

int main() {
    char* source = 
        "# Simple tool call that retrieves the player's name with fallback\n"
        "LET player_name = <getPlayerName;>\n"
        "LET player_age = <getPlayerAge;> ?? 18\n"
        "LET age_plus_five = <agePlusFive; player_age>\n"
        "\n"
        "SYSTEM: \"Hello, {player_name}! You are {player_age} years old.\"\n"
        "IF player_age < 18 {\n"
        "    SYSTEM: \"You are quite young to be here!\"\n"
        "}\n"
        "\n"
        "IF player_age >= 18 {\n"
        "    SYSTEM: \"You are old enough to be here.\"\n"
        "}\n"
        "\n"
        "# Advanced inline tool call. Syntax: <function; argument>\n"
        "SYSTEM: \"Your current gold balance is <getData; \"gold\"> gold coins.\"\n"
        "\n"
        "# Multiple tool calls in a single line\n"
        "SYSTEM: \"Your current gold balance is <getData; \"gold\"> gold coins and your health is <getData; \"health\">.\"\n"
        "\n"
        "# Multiple tool call arguments\n"
        "# Tool function getItemPrice takes two arguments: item type and item level\n"
        "LET item_price = <getItemPrice; \"potion\", 4>\n"
        "SYSTEM: \"The price of a level 4 potion is {item_price} gold coins.\"\n"
        "END\n";

    printf("=== Quill C API Tool Call Test ===\n\n");

    // Test parse only
    printf("1. Testing parse only:\n");
    char* parse_result = quill_parse_only(source);
    printf("Parse result: %s\n\n", parse_result);
    quill_free_string(parse_result);

    // Create interpreter
    printf("2. Creating interpreter:\n");
    int interp_id = quill_new_interpreter(source);
    if (interp_id == -1) {
        printf("Failed to create interpreter\n");
        return 1;
    }
    printf("Created interpreter with ID: %d\n\n", interp_id);

    // Test interpreter methods
    printf("3. Testing interpreter methods:\n");
    
    // Get initial state
    char* state = quill_get_state(interp_id);
    printf("Initial state: %s\n", state);
    quill_free_string(state);

    // Step through execution
    for (int i = 0; i < 10; i++) {
        printf("\n--- Step %d ---\n", i + 1);
        
        char* step_result = quill_step(interp_id);
        printf("Step result: %s\n", step_result);

        // Check if waiting for choice
        char* waiting_choice = quill_is_waiting_for_choice(interp_id);
        printf("Waiting for choice: %s\n", waiting_choice);
        
        // Check if waiting for tool call
        char* waiting_tool = quill_is_waiting_for_tool_call(interp_id);
        printf("Waiting for tool call: %s\n", waiting_tool);
        
        // Handle choice if needed
        if (strstr(waiting_choice, "true") != NULL) {
            printf("Handling choice with option 0...\n");
            char* choice_result = quill_handle_choice(interp_id, 0);
            printf("Choice result: %s\n", choice_result);
            quill_free_string(choice_result);
        }
        
        // Handle tool call if needed
        if (strstr(waiting_tool, "true") != NULL) {
            printf("Handling tool call with mock result...\n");
            char* mock_result = get_mock_tool_result(step_result);
            printf("Using mock result: %s\n", mock_result);
            char* tool_result = quill_handle_tool_call_response(interp_id, mock_result);
            printf("Tool call result: %s\n", tool_result);
            quill_free_string(tool_result);
        }
        
        quill_free_string(waiting_choice);
        quill_free_string(waiting_tool);

        // Check if ended
        char* ended = quill_is_ended(interp_id);
        printf("Is ended: %s\n", ended);
        
        // Break if ended
        if (strstr(ended, "true") != NULL) {
            quill_free_string(ended);
            quill_free_string(step_result);
            break;
        }
        
        quill_free_string(ended);
        quill_free_string(step_result);
    }

    // Clean up
    quill_free_interpreter(interp_id);
    
    printf("\n=== Tool Call Test Complete ===\n");
    return 0;
}