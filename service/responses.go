package service

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	statusKO = "KO"
	statusOK = "OK"
)

// Response models standard response
type Response struct {
	Message   string        `json:"message"`
	Errors    []interface{} `json:"errors"`
	ClassName string        `json:"className"`
	Code      int           `json:"code"`
	Data      []interface{} `json:"data"`
	Name      string        `json:"name"`
	Status    string        `json:"status"`
}

func newResponse(message string) *Response {
	resp := &Response{
		Message: message,
		Errors:  make([]interface{}, 0),
		Data:    make([]interface{}, 0),
	}
	return resp
}

func (resp *Response) sendJSON(c *gin.Context) {
	c.JSON(resp.Code, resp)
}

// OK set status and code to statusOK
func (resp *Response) OK() {
	resp.Status = statusOK
	resp.Code = http.StatusOK
}

// KO set status and code to statusKO
func (resp *Response) KO(code int) {
	resp.Status = statusKO
	resp.Code = code
}

func (resp *Response) addError(errors ...interface{}) {
	if os.Getenv("ENV") == "prod" {
		return
	}
	resp.Errors = append(resp.Errors, errors...)
}

func (resp *Response) setMessage(message string) {
	resp.Message = message
}

func sendUnauthorizedAccess(c *gin.Context, errors ...error) {
	resp := newResponse("Not Allowed")
	resp.addError(errors)
	resp.KO(http.StatusForbidden)
	resp.sendJSON(c)
	c.Abort()
}

func sendBindQueryError(c *gin.Context, err error) {
	logrus.Infof("Error binding query: %v", err)
	sendStandardError(c, err)
}

func abortWithSendingStandardError(c *gin.Context, err error) {
	logrus.Warningf("Aborting with standard error: %v", err)
	sendStandardError(c, err)
}

func sendStandardError(c *gin.Context, errors ...error) {
	resp := newResponse("Bad Request, see errors details")
	resp.addError(errors)
	resp.KO(http.StatusBadRequest)
	resp.sendJSON(c)
	c.Abort()
}

func sendFirebaseError(c *gin.Context, errors ...error) {
	errorCode := inferErrorCode(errors[0])
	resp := newResponse("Firebase error")
	resp.addError(errors)
	resp.KO(errorCode)
	resp.sendJSON(c)
	c.Abort()
}
