package main
import (
   "encoding/json"
   "errors"
   "io/ioutil"
   "net/http"
   "net/url"
   "github.com/gorilla/mux"
)


func Alpha( w http.ResponseWriter, r *http.Request ) {
   t:= map[ string ] interface{} {}
   if err := json.NewDecoder( r.Body ).Decode( &t ); err == nil {
      if question, ok := t[ "text" ].( string ); ok { 
         if answer, err := Service( question ); err == nil {
            if answer == "" {
               answer = "Sorry I dont know the answer to that"
            }
            u := map[ string ] interface{} { "text" : answer }
            w.WriteHeader( http.StatusOK )
            json.NewEncoder( w ).Encode( u )
         } else if answer == "500" {
            w.WriteHeader( http.StatusInternalServerError )
            w.Write([]byte(err.Error()))
         } else if answer == "404" {
            w.WriteHeader( http.StatusNotFound )
            w.Write([]byte(err.Error()))
         } else {
            w.WriteHeader( http.StatusBadRequest )
            w.Write([]byte(err.Error()))
         }
      } else {
         w.WriteHeader( http.StatusBadRequest )
         w.Write([]byte(err.Error()))
      }
   } else {
      w.WriteHeader( http.StatusBadRequest )
      w.Write([]byte(err.Error()))
   }
}

func Service( question string ) ( string, error ) {
   client := &http.Client{}
   appid := "WU6GXH-5Q9724623K"
   questionEnc := url.QueryEscape( question )
   URI := "https://api.wolframalpha.com/v1/result?appid=" + appid + 
          "&i=" + questionEnc
   if req, err := http.NewRequest( "POST", URI, nil ); err == nil {
      if rsp, err := client.Do( req ); err == nil {
         defer rsp.Body.Close()
         if rsp.StatusCode == http.StatusOK {
            if body, err := ioutil.ReadAll( rsp.Body ); err == nil {
               answer := string( body )
               return answer, nil
            } else {
               return "", errors.New( "cannot read response body" )
            }
         } else if rsp.StatusCode == http.StatusNotImplemented {
            return "", nil
         } else if rsp.StatusCode == http.StatusInternalServerError {
            return "500", errors.New( "alpha api has returned 500" )
         } else if rsp.StatusCode == http.StatusNotFound {
            return "404", errors.New( "invalid URL" )
         } else {
            return "", errors.New( "cannot answer question" )
         }
      } else {
         return "", errors.New( "client cannot complete request" )
      }
   } else {
      return "", errors.New( "invalid http request" )
   }
}


func main() {
   r := mux.NewRouter()
   r.HandleFunc( "/alpha",  Alpha ).Methods( "POST" )
     http.ListenAndServe( ":3001", r )
}
