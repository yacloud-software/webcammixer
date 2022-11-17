// client create: WebCamMixerClient
/* geninfo:
   filename  : protos/golang.conradwood.net/apis/webcammixer/webcammixer.proto
   gopackage : golang.conradwood.net/apis/webcammixer
   importname: ai_0
   varname   : client_WebCamMixerClient_0
   clientname: WebCamMixerClient
   servername: WebCamMixerServer
   gscvname  : webcammixer.WebCamMixer
   lockname  : lock_WebCamMixerClient_0
   activename: active_WebCamMixerClient_0
*/

package webcammixer

import (
   "sync"
   "golang.conradwood.net/go-easyops/client"
)
var (
  lock_WebCamMixerClient_0 sync.Mutex
  client_WebCamMixerClient_0 WebCamMixerClient
)

func GetWebCamMixerClient() WebCamMixerClient { 
    if client_WebCamMixerClient_0 != nil {
        return client_WebCamMixerClient_0
    }

    lock_WebCamMixerClient_0.Lock() 
    if client_WebCamMixerClient_0 != nil {
       lock_WebCamMixerClient_0.Unlock()
       return client_WebCamMixerClient_0
    }

    client_WebCamMixerClient_0 = NewWebCamMixerClient(client.Connect("webcammixer.WebCamMixer"))
    lock_WebCamMixerClient_0.Unlock()
    return client_WebCamMixerClient_0
}

