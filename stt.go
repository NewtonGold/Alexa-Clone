package main
import (
   "encoding/json"
   "errors"
   "net/http"
   "bytes"
   "encoding/base64"
   "github.com/gorilla/mux"
)
const (
   REGION = "uksouth"
   URI = "https://" + REGION + ".stt.speech.microsoft.com/" +
      "speech/recognition/conversation/cognitiveservices/v1?" +
      "language=en-US"
   KEY = "d76745e51adf4408b1f29d7a4362dc39"
)


func SpeechToText( w http.ResponseWriter, r *http.Request ) {
   t:= map[ string ] interface{} {}
   if err := json.NewDecoder( r.Body ).Decode( &t ); err == nil {
      if speech, ok := t[ "speech" ].( string ); ok {
         if sDec, err := base64.StdEncoding.DecodeString(speech); err == nil {  
            if text, err := Service( sDec ); err == nil {
               u := map[ string ] interface{} { "text" : text }
               w.WriteHeader( http.StatusOK )
               json.NewEncoder( w ).Encode( u )
            } else {
               if text == "500" {
                  w.WriteHeader( http.StatusInternalServerError )
                  w.Write([]byte(err.Error()))
               } else {
                  w.WriteHeader( http.StatusBadRequest )
                  w.Write([]byte(err.Error()))
               }
            }
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

func Service( speech []byte ) ( string, error ) {
   client := &http.Client{}
   if req, err := http.NewRequest( "POST", URI, bytes.NewReader( speech ) ); err == nil {
      req.Header.Set( "Content-Type",
         "audio/wav;codecs=audio/pcm;samplerate=16000" )
      req.Header.Set( "Ocp-Apim-Subscription-Key", KEY )
      if rsp, err := client.Do( req ); err == nil {
         defer rsp.Body.Close()
         if rsp.StatusCode == http.StatusOK {
            t:= map[ string ] interface{} {}
            if err := json.NewDecoder( rsp.Body ).Decode( &t ); err == nil {
               if text, ok := t[ "DisplayText" ].( string ); ok {
                  return string( text ), nil
               } else {
                  return "", errors.New( "cannot find json variable" )
               }
            } else {
               return "", errors.New( "cannot decode responce body" )
            }
         } else if rsp.StatusCode == http.StatusBadRequest{
            return "", errors.New( "cannot convert speech to text" )
         } else if rsp.StatusCode == http.StatusNotFound{
            return "404", errors.New( "URL could not be found" ) 
         } else {
            return "500", errors.New( "server error when converting speech to text" )
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
   r.HandleFunc( "/stt",  SpeechToText ).Methods( "POST" )
     http.ListenAndServe( ":3002", r )
}
