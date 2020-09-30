package route

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
)

type exception struct {
	Message string `json:"message"`
}

//validateMiddleware validates the JWT
// func validateMiddleware(next http.Handler) http.Handler {
// 	// return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
// 	// 	authorizationHeader := req.Header.Get("authorization")
// 	// 	if authorizationHeader != "" {
// 	// 		//fmt.Println("authorizationHeader : ", authorizationHeader)
// 	// 		bearerToken := strings.Split(authorizationHeader, " ")
// 	// 		if len(bearerToken) == 2 {
// 	// 			key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicCertContent))
// 	// 			if err != nil {
// 	// 				json.NewEncoder(w).Encode(exception{Message: err.Error()})
// 	// 				return
// 	// 			}
// 	// 			token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
// 	// 				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
// 	// 					return nil, fmt.Errorf("There was an error")
// 	// 				}
// 	// 				return key, nil
// 	// 			})
// 	// 			if error != nil {
// 	// 				json.NewEncoder(w).Encode(exception{Message: error.Error()})
// 	// 				return
// 	// 			}
// 	// 			if token.Valid {
// 	// 				next.ServeHTTP(w, req)
// 	// 			} else {
// 	// 				json.NewEncoder(w).Encode(exception{Message: "Invalid authorization token"})
// 	// 			}
// 	// 		}
// 	// 	} else {
// 	// 		json.NewEncoder(w).Encode(exception{Message: "An authorization header is required"})
// 	// 	}
// 	// })
// }

func loggingMiddleware(next http.Handler) http.Handler {
	logFile, err := os.OpenFile("access.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		log.Printf("could not create log file, err : %v\n", err)
	}
	return handlers.LoggingHandler(logFile, next)
}

func panicHandlerMiddleware(next http.Handler) http.Handler {
	logFile, _ := os.OpenFile("server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	handler := handlers.RecoveryHandler(handlers.PrintRecoveryStack(true),
		handlers.RecoveryLogger(log.New(logFile, "", log.LstdFlags)))
	return handler(next)

}
