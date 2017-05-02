package models

type Site struct {
  Id string `json:"id"`
  Name string `json:"name"`
  Host string `json:"host"`
  Port uint16 `json:"port"`
  Cpu int `json:"cpu"`
  Static string `json:"static"`
  Gzip string `json:"gzip"`
  LetsEncrypt bool `json:"letsEncrypt"`
  Http2 string `json:"http2"`
  Proxies []struct {
    Host string `json:"host"`
    Port uint16 `json:"port"`
    Path string `json:"path"`
    Replace string `json:"replace"`
  } `json:"proxies"`
}
