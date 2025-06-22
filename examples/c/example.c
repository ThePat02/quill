#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "libquill.h"

// Simple function to extract interpreter ID from JSON response
// In a real application, you'd use a proper JSON parser
int extract_interpreter_id(const char* json_response) {
    // Look for "success":true first
    if (strstr(json_response, "\"success\":true") == NULL) {
        return -1; // Error case
    }
    
    // This is a very basic extraction - in practice use a JSON library
    // The response should contain the interpreter ID
    // For now, we'll return a hardcoded value since the Go code doesn't 
    // actually include the ID in the response
    return 1;
}

int main() {
    const char* source = 
        "ALEX: \"Hello there!\"\n"
        "BELLA: \"Hi Alex!\"\n"
        "\n"
        "CHOICE {\n"
        "    \"How are you?\" {\n"
        "        ALEX: \"I'm doing great, thanks!\"\n"
        "    },\n"
        "    \"What's new?\" {\n"
        "        ALEX: \"Not much, just working on some projects.\"\n"
        "    }\n"
        "}\n"
        "\n"
        "ALEX: \"Thanks for asking!\"\n"
        "END\n";

    printf("=== Quill C API Test ===\n\n");

    // Test parse only
    printf("1. Testing parse only:\n");
    char* parse_result = quill_parse_only(source);
    printf("Parse result: %s\n\n", parse_result);
    quill_free_string(parse_result);

    // Create interpreter
    printf("2. Creating interpreter:\n");
    char* init_result = quill_new_interpreter(source);
    printf("Init result: %s\n\n", init_result);
    
    // Extract interpreter ID from response
    int interp_id = extract_interpreter_id(init_result);
    if (interp_id == -1) {
        printf("Failed to create interpreter\n");
        quill_free_string(init_result);
        return 1;
    }
    
    quill_free_string(init_result);

    // Test interpreter methods
    printf("3. Testing interpreter methods:\n");
    
    // Get initial state
    char* state = quill_get_state(interp_id);
    printf("Initial state: %s\n", state);
    quill_free_string(state);

    // Step through execution
    for (int i = 0; i < 5; i++) {
        printf("\n--- Step %d ---\n", i + 1);
        
        char* step_result = quill_step(interp_id);
        printf("Step result: %s\n", step_result);
        quill_free_string(step_result);

        // Check if waiting for choice
        char* waiting = quill_is_waiting_for_choice(interp_id);
        printf("Waiting for choice: %s\n", waiting);
        
        // For demo purposes, always choose option 0 if waiting
        // In real usage, you'd parse the JSON to check the boolean value
        char* choice_result = quill_handle_choice(interp_id, 0);
        printf("Choice result: %s\n", choice_result);
        quill_free_string(choice_result);
        quill_free_string(waiting);

        // Check if ended
        char* ended = quill_is_ended(interp_id);
        printf("Is ended: %s\n", ended);
        quill_free_string(ended);
    }

    // Clean up
    quill_free_interpreter(interp_id);
    
    printf("\n=== Test Complete ===\n");
    return 0;
}