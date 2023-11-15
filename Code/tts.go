package main
import (
   "encoding/json"
   "encoding/xml"
   "bytes"
   "errors"
   "io/ioutil"
   "net/http"
   "github.com/gorilla/mux"
)
const (
   REGION = "uksouth"
   URI = "https://" + REGION + ".tts.speech.microsoft.com/" +
      "cognitiveservices/v1"
   KEY = "d76745e51adf4408b1f29d7a4362dc39"
)
type Voice struct{
   XMLName xml.Name `xml:"voice"`
   Lang    string   `xml:"http://www.w3.org/XML/1998/namespace lang,attr,omitempty"`
   Name string `xml:"name,attr"`
   Text string `xml:",chardata"`
}
type Speak struct {
   XMLName xml.Name `xml:"speak"`
   Version string `xml:"version,attr"`
   Lang    string   `xml:"http://www.w3.org/XML/1998/namespace lang,attr,omitempty"`
   Voice Voice
}

func TextToSpeech( w http.ResponseWriter, r *http.Request ) {
   t:= map[ string ] interface{} {}
   var newmessage Speak
   if data, err := ioutil.ReadFile("text.xml"); err == nil {
      if err := xml.Unmarshal(data, &newmessage); err == nil {
         if err := json.NewDecoder( r.Body ).Decode( &t ); err == nil {
            if text, ok := t[ "text" ].( string ); ok {
               newmessage.Voice.Text = text
               if modified, err := xml.MarshalIndent(&newmessage, "", "   "); err == nil {
                  if speech, err := Service( modified ); err == nil {
                     u := map[ string ] interface{} { "speech" : speech }
                     w.WriteHeader( http.StatusOK )
                     json.NewEncoder( w ).Encode( u )
                  } else {
                     if bytes.Compare(speech, []byte("500")) == 0 {
                        w.WriteHeader( http.StatusInternalServerError )
                        w.Write([]byte(err.Error()))
                     } else if bytes.Compare(speech, []byte("404")) == 0 {
                        w.WriteHeader( http.StatusNotFound )
                        w.Write([]byte(err.Error()))
                     } else {
                        w.WriteHeader( http.StatusBadRequest )
                        w.Write([]byte(err.Error()))
                     }
                  }
               } else {
                  w.WriteHeader( http.StatusInternalServerError )
               }
            } else {
               w.WriteHeader( http.StatusBadRequest )
            }
         } else {
            w.WriteHeader( http.StatusBadRequest )
         }
      } else {
         w.WriteHeader( http.StatusInternalServerError )
      }
   } else {
      w.WriteHeader( http.StatusInternalServerError ) 
   }
}
func Service( text []byte ) ( []byte, error ) {
   client := &http.Client{}
   if req, err := http.NewRequest( "POST", URI, bytes.NewBuffer( text ) ); err == nil {
      req.Header.Set( "Content-Type", "application/ssml+xml" )
      req.Header.Set( "Ocp-Apim-Subscription-Key", KEY )
      req.Header.Set( "X-Microsoft-OutputFormat", "riff-16khz-16bit-mono-pcm" )
      if rsp, err := client.Do( req ); err == nil {
         defer rsp.Body.Close()
         if rsp.StatusCode == http.StatusOK {
            if body, err := ioutil.ReadAll( rsp.Body ); err == nil {
               return body, nil
            } else {
               return []byte("500"), errors.New( "cannot read response body" )
            }
         } else if rsp.StatusCode == http.StatusBadRequest{
            return nil, errors.New( "cannot convert speech to text" )
         } else if rsp.StatusCode == http.StatusNotFound{
            return []byte("404"), errors.New( "URL could not be found" ) 
         } else {
            return []byte("500"), errors.New( "server error when converting text to speech" )
         }
      } else {
         return nil, errors.New( "cannot process request" )
      }
   } else {
      return nil, errors.New( "cannot generate http request" )
   }
}
func main() {
   r := mux.NewRouter()
   r.HandleFunc( "/tts",  TextToSpeech ).Methods( "POST" )
     http.ListenAndServe( ":3003", r )
}
