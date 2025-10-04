package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request) { // this fucntion is used to define a http handler as go expects
	respondWithJSON(w, 200, struct{}{})
}
