package main

import "github.com/Edwing123/udem-chat-app/pkg/codes"

var (
	// Server related.
	ErrServerInternal = codes.NewCode("server_internal")

	// Client related.
	ErrCannotDecodeJSON = codes.NewCode(("cannot_decode_json"))
	ErrAuthRequired     = codes.NewCode("auth_required")

	// Validation related.
	ErrPasswordNotValid = codes.NewCode("password_not_valid")
)
