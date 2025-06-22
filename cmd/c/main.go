package main

import (
	"quill/internal/jsonapi"
	"sync"
	"unsafe"
)

// #include <stdlib.h>
// #include <stdio.h>
import "C"

var (
	interpreters = make(map[int]*jsonapi.QuillInterpreter)
	nextID       = 1
	mu           sync.Mutex
)

//export quill_new_interpreter
func quill_new_interpreter(source *C.char) C.int {
	goSource := C.GoString(source)
	interp, _ := jsonapi.NewQuillInterpreter(goSource)

	if interp == nil {
		return -1
	}

	mu.Lock()
	id := nextID
	interpreters[nextID] = interp
	nextID++
	mu.Unlock()

	return C.int(id)
}

//export quill_step
func quill_step(interpID C.int) *C.char {
	mu.Lock()
	interp, exists := interpreters[int(interpID)]
	mu.Unlock()

	if !exists {
		return C.CString(`{"success":false,"error":"Invalid interpreter ID"}`)
	}

	result := interp.Step()
	cResult := C.CString(result)
	if cResult == nil {
		return C.CString(`{"success":false,"error":"Failed to allocate C string"}`)
	}
	return cResult
}

//export quill_handle_choice
func quill_handle_choice(interpID C.int, choiceIndex C.int) *C.char {
	mu.Lock()
	interp, exists := interpreters[int(interpID)]
	mu.Unlock()

	if !exists {
		return C.CString(`{"success":false,"error":"Invalid interpreter ID"}`)
	}

	result := interp.HandleChoiceInput(int(choiceIndex))
	cResult := C.CString(result)
	if cResult == nil {
		return C.CString(`{"success":false,"error":"Failed to allocate C string"}`)
	}
	return cResult
}

//export quill_get_state
func quill_get_state(interpID C.int) *C.char {
	mu.Lock()
	interp, exists := interpreters[int(interpID)]
	mu.Unlock()

	if !exists {
		return C.CString(`{"success":false,"error":"Invalid interpreter ID"}`)
	}

	result := interp.GetState()
	cResult := C.CString(result)
	if cResult == nil {
		return C.CString(`{"success":false,"error":"Failed to allocate C string"}`)
	}
	return cResult
}

//export quill_is_ended
func quill_is_ended(interpID C.int) *C.char {
	mu.Lock()
	interp, exists := interpreters[int(interpID)]
	mu.Unlock()

	if !exists {
		return C.CString(`{"success":false,"error":"Invalid interpreter ID"}`)
	}

	result := interp.IsEnded()
	cResult := C.CString(result)
	if cResult == nil {
		return C.CString(`{"success":false,"error":"Failed to allocate C string"}`)
	}
	return cResult
}

//export quill_is_waiting_for_choice
func quill_is_waiting_for_choice(interpID C.int) *C.char {
	mu.Lock()
	interp, exists := interpreters[int(interpID)]
	mu.Unlock()

	if !exists {
		return C.CString(`{"success":false,"error":"Invalid interpreter ID"}`)
	}

	result := interp.IsWaitingForChoice()
	cResult := C.CString(result)
	if cResult == nil {
		return C.CString(`{"success":false,"error":"Failed to allocate C string"}`)
	}
	return cResult
}

//export quill_parse_only
func quill_parse_only(source *C.char) *C.char {
	goSource := C.GoString(source)
	result := jsonapi.ParseOnly(goSource)
	cResult := C.CString(result)
	if cResult == nil {
		return C.CString(`{"success":false,"error":"Failed to allocate C string"}`)
	}
	return cResult
}

//export quill_free_string
func quill_free_string(str *C.char) {
	C.free(unsafe.Pointer(str))
}

//export quill_free_interpreter
func quill_free_interpreter(interpID C.int) {
	mu.Lock()
	delete(interpreters, int(interpID))
	mu.Unlock()
}

func main() {}
