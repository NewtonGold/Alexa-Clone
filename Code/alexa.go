package main
import (
   "errors"
   "io"
   "net/http"
   "bytes"
   "github.com/gorilla/mux"
)


func Alexa( w http.ResponseWriter, r *http.Request ) {
   if body, err := io.ReadAll(r.Body); err == nil {
       if question_json, err := Service( body, "stt", "3002" ); err == nil {
           if answer_json, err := Service(question_json, "alpha", "3001" ); err == nil {
               if speech_json, err := Service( answer_json, "tts", "3003" ); err == nil {
                   w.WriteHeader(http.StatusOK)
                   w.Write( speech_json  )
               } else {
                  if bytes.Compare(speech_json, []byte("500")) == 0 {
                     w.WriteHeader( http.StatusInternalServerError )
                     w.Write([]byte(err.Error()))
                  } else if bytes.Compare(speech_json, []byte("404")) == 0 {
                     w.WriteHeader( http.StatusNotFound )
                     w.Write([]byte(err.Error()))
                  } else {
                     w.WriteHeader( http.StatusBadRequest )
                     w.Write([]byte(err.Error()))
                  }
               }
           } else {
               if bytes.Compare(answer_json, []byte("500")) == 0 {
                  w.WriteHeader( http.StatusInternalServerError )
                  w.Write([]byte(err.Error()))
               } else if bytes.Compare(answer_json, []byte("404")) == 0 {
                  w.WriteHeader( http.StatusNotFound )
                  w.Write([]byte(err.Error()))
               } else {
                  w.WriteHeader( http.StatusBadRequest )
                  w.Write([]byte(err.Error()))
               }
           }
       } else {
           if bytes.Compare(question_json, []byte("500")) == 0 {
              w.WriteHeader( http.StatusInternalServerError )
              w.Write([]byte(err.Error()))
           } else if bytes.Compare(question_json, []byte("404")) == 0 {
              w.WriteHeader( http.StatusNotFound )
              w.Write([]byte(err.Error()))
           } else {
              w.WriteHeader( http.StatusBadRequest )
              w.Write([]byte(err.Error()))
           }
       }
   } else {
       w.WriteHeader(http.StatusBadRequest)
       w.Write([]byte(err.Error()))
   }
}

func Service( json []byte, microservice string, port string) ( []byte, error ) {
   client := &http.Client{}
   URI := "http://localhost:"+ port + "/" + microservice
   if req, err := http.NewRequest( "POST", URI, bytes.NewBuffer( json ) ); err == nil {
      if rsp, err := client.Do( req ); err == nil {
         defer rsp.Body.Close()
         if rsp.StatusCode == http.StatusOK {
            if body, err := io.ReadAll( rsp.Body ); err == nil {
               return body, nil
            } else {
               return nil, errors.New( "cannot read response body" )
            }
         } else if rsp.StatusCode == http.StatusInternalServerError{
            return []byte("500"), errors.New( "microservice " + microservice + 
               " cannot complete request" )
         } else if rsp.StatusCode == http.StatusNotFound{
            return []byte("404"), errors.New( "microservice " + microservice + 
               " returned Not Found 404" )
         } else {
            return nil, errors.New( "invalid input parameters" )
         }
      } else {
         return nil, errors.New( "client cannot complete request" )
      }
   } else {
      return nil, errors.New( "invalid http request" )
   }
}

func main() {
   r := mux.NewRouter()
   r.HandleFunc( "/alexa",  Alexa ).Methods( "POST" )
   http.ListenAndServe( ":3000", r )
}
