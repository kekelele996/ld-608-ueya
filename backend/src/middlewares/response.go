package middlewares

	import "github.com/gin-gonic/gin"

// Envelope is the shared JSON shape: { ok, data } or { ok:false, error }.
type Envelope struct {
	OK    bool        `json:"ok"`
	Data  interface{} `json:"data,omitempty"`
	Error *ErrBody    `json:"error,omitempty"`
}

type ErrBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func abortEnvelope(c *gin.Context, status int, code, message string) {
	c.AbortWithStatusJSON(status, Envelope{OK: false, Error: &ErrBody{Code: code, Message: message}})
}

// AbortEnvelope is exported for controller-layer error wrapping.
func AbortEnvelope(c *gin.Context, status int, code, message string) {
	abortEnvelope(c, status, code, message)
}

// AbortBadRequest keeps controller call sites short.
func AbortBadRequest(c *gin.Context, message string) {
	abortEnvelope(c, 400, "VALIDATION_FAILED", message)
}

// OK writes a success envelope.
func OK(c *gin.Context, status int, data interface{}) {
	c.JSON(status, Envelope{OK: true, Data: data})
}

// Conflict writes the all-or-nothing rejection envelope carrying per-item
// conflict reasons for the frontend conflict panel.
func Conflict(c *gin.Context, status int, code, message string, data interface{}) {
	c.AbortWithStatusJSON(status, Envelope{
		OK:    false,
		Data:  data,
		Error: &ErrBody{Code: code, Message: message},
	})
}
