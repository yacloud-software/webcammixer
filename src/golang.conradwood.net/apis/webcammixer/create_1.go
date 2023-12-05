// client create: WebCamMixerClient
/*
  Created by /home/cnw/devel/go/yatools/src/golang.yacloud.eu/yatools/protoc-gen-cnw/protoc-gen-cnw.go
*/

/* geninfo:
   filename  : protos/golang.conradwood.net/apis/webcammixer/webcammixer.proto
   gopackage : golang.conradwood.net/apis/webcammixer
   importname: ai_0
   clientfunc: GetWebCamMixer
   serverfunc: NewWebCamMixer
   lookupfunc: WebCamMixerLookupID
   varname   : client_WebCamMixerClient_0
   clientname: WebCamMixerClient
   servername: WebCamMixerServer
   gsvcname  : webcammixer.WebCamMixer
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

    client_WebCamMixerClient_0 = NewWebCamMixerClient(client.Connect(WebCamMixerLookupID()))
    lock_WebCamMixerClient_0.Unlock()
    return client_WebCamMixerClient_0
}

func WebCamMixerLookupID() string { return "webcammixer.WebCamMixer" } // returns the ID suitable for lookup in the registry. treat as opaque, subject to change.

func init() {
   client.RegisterDependency("webcammixer.WebCamMixer")
   AddService("webcammixer.WebCamMixer")
}
