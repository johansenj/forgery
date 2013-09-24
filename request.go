package f

import(
    "regexp"
    "strings"
    "github.com/ricallinson/stackr"
)

/*
    A Request represents an HTTP request received by the server.
*/
type Request struct {

    // The stackr.Request type.
    *stackr.Request

    // The application server.
    app *Server

    // Return the remote address, or when "trust proxy" is enabled - the upstream address.
    Ip string

    // When "trust proxy" is `true`, parse the "X-Forwarded-For" ip address list and return a slice, 
    // otherwise an empty slice is returned. For example if the value were "client, proxy1, proxy2" 
    // you would receive the slice {"client", "proxy1", "proxy2"} where "proxy2" is the furthest down-stream.
    Ips []string

    // This property is a slice containing properties mapped to the named route "parameters". 
    // For example if you have the route "/user/:name", then the "name" property is available 
    // to you as req.params["name"]. This object defaults to {}.
    Params map[string]string

    // The currently matched Route containing several properties such as the 
    // route's original path string, the regexp generated, and so on.
    Route interface{}
}

/*
    Returns a new Request.
*/

func createRequest(req *stackr.Request, app *Server) (*Request) {

    this := &Request{}

    this.Request = req

    this.app = app

    if t, v := this.app.Get("trust proxy"), this.Header.Get("X-Forwarded-For"); len(t) > 0 && len(v) > 0 {
        s := regexp.MustCompile(" *, *").Split(v, -1)
        this.Ip = s[0]
        this.Ips = s
    } else {
        this.Ip = this.RemoteAddr
        this.Ips = []string{}
    }

    this.Params = map[string]string{}

    // Could have been set by middleware.
    if this.Body == nil {
        this.Body = map[string]string{}
    }

    // Could have been set by middleware.
    if this.Query == nil {
        this.Query = map[string]string{}
    }

    // Could have been set by middleware.
    if this.Files == nil {
        this.Files = map[string]interface{}{}
    }

    return this
}

/*
    Contains the cookies sent by the user-agent.
*/
func (this *Request) Cookie(n string, i *interface{}) {
    //
}

/*
    Contains the signed cookies sent by the user-agent, unsigned and ready for use. 
    Signed cookies are accessed by a different function to show developer intent, otherwise a 
    malicious attack could be placed on `req.Cookie` values which are easy to spoof. 
    Note that signing a cookie does not mean it is "hidden" nor encrypted, this simply 
    prevents tampering as the secret used to sign is private.
*/
func (this *Request) SignedCookie(n string, i *interface{}) {
    //
}

/*
    Return the value of param "name" when present. Lookup is performed in the following order:

    * Params
    * Body
    * Query

    Direct access to req.body, req.params, and req.query should be favoured for clarity - 
    unless you truly accept input from each object.
*/
func (this *Request) Param(n string) (string) {
    var v string
    var ok bool
    v, ok = this.Params[n]
    if ok {
        return v
    }
    v, ok = this.Body[n]
    if ok {
        return v
    }
    v, ok = this.Query[n]
    if ok {
        return v
    }
    return ""
}

/*
    Get the case-insensitive request header field. The Referrer and Referer fields are interchangeable.
*/
func (this *Request) Get(f string) (string) {

    /*
        Possible future bug.
        http://www.w3.org/Protocols/rfc2616/rfc2616-sec4.html#sec4.2
        Message headers are case-insensitive.
    */

    return this.Header.Get(f)
}

/*
    Check if the given types are acceptable, returning the best match when true, 
    otherwise undefined - in which case you should respond with 406 "Not Acceptable".
*/
func (this *Request) Accepts(t string) (bool) {
    for _, v := range this.Accepted() {
        if strings.ToLower(t) == v {
            return true
        }
    }
    return false
}

/*
    Check if the incoming request contains the "Content-Type" header field, and it matches the give mime "type".
*/
func (this *Request) Is(t string) (bool) {
    h := this.Get("Content-Type")
    return strings.ToLower(h) == strings.ToLower(t)
}

/*
    Check if the given "charset" is acceptable.
*/
func (this *Request) AcceptsCharset(c string) {
    panic(halt)
}

/*
    Check if the given "lang" is acceptable.
*/
func (this *Request) AcceptsLanguage(l string) {
    panic(halt)
}