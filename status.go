/*Error status methods*/
package hamster

import (
	"net/http"
)

func (s *Server) serveError(w http.ResponseWriter, err error, user_message string, base_message string, status int) {

	log_message := base_message + " : " + user_message
	s.logger.Printf("Error %v:  %v \n", log_message, err)
	http.Error(w, base_message, status)

}
func (s *Server) badRequest(r *http.Request, w http.ResponseWriter, err error, msg string) {
	s.serveError(w, err, msg, "Bad Request", http.StatusBadRequest)
}

func (s *Server) unauthorized(r *http.Request, w http.ResponseWriter, err error, msg string) {
	s.serveError(w, err, msg, "Unauthorized", http.StatusUnauthorized)
}

func (s *Server) notFound(r *http.Request, w http.ResponseWriter, err error, msg string) {
	s.serveError(w, err, msg, "Not found", http.StatusNotFound)

}

func (s *Server) internalError(r *http.Request, w http.ResponseWriter, err error, msg string) {
	s.serveError(w, err, msg, "Internal Server Error", http.StatusInternalServerError)

}
