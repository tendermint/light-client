package types

// ErrorResponse is returned for 4xx and 5xx errors
// type ErrorResponse struct {
// 	Success bool   `json:"success"`
// 	Error   string `json:"error"` // error message if Success is false
// 	Code    int    `json:"code"`  // error code if Success is false
// }

// GenericResponse is returned for 4xx and 5xx errors
// And the following 2xx results: BroadcastResult
type GenericResponse struct {
	Code int32  `json:"code"` // TODO: rethink this (0 = success)
	Data []byte `json:"data"` // TODO: make sure this is hex encoded
	Log  string `json:"log"`
}
