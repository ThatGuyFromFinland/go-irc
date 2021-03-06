package core

import (
    "fmt"
    "strings"
    "strconv"
)

// For testability
type Handler interface {
    Read(b []byte) (n int, err error)
    Write(b []byte) (n int, err error)
    Close() error
}

// For testability
type Dialer func (AddressEntry) Handler


type AddressEntry struct {
    Id uint64
    IP string
    Port string
}

type Message struct {
    Recipient uint64
    Payload string
}

type Router struct {
    MessageQueue chan Message
    addressBook []AddressEntry
}

func (this *Router) InsertAddress(addressString string) bool {
    addressParts := strings.Split(addressString, ":")
    if (len(addressParts) != 2) {
        return false
    }
    newEntry := AddressEntry{Id: uint64(len(this.addressBook)), IP: addressParts[0], Port: addressParts[1]}
    this.addressBook = append(this.addressBook, newEntry)
    fmt.Println(len(this.addressBook))
    return true
}

func (this *Router) Init() {
    this.MessageQueue = make(chan Message)
}
// Kind of action controller, seems quite okay, nice to test also
func (this *Router) RouteRequest(request string, conn Handler) {
    requestSplit := strings.SplitN(request, ":", 2)
    if (len(requestSplit) <= 1) {
        // Invalid request, return
        return
    }
    
    action := requestSplit[0]
    body := strings.TrimSuffix(requestSplit[1], "\n")
    fmt.Printf("ACTION IS: %s \n", action)
    switch {
        case action == "JOIN":
            this.handleClientJoinRequest(body, conn)
        case action == "PEOPLE":
            this.handlePeopleRequest(body, conn)
        case action == "WHOAMI":
            this.handleWhoAmIRequest(body, conn)
        case action == "MESSAGE":
            this.handleMessageRequest(body, conn)
    }
}

func (this *Router) handleClientJoinRequest(body string, conn Handler) {
    noErrors := this.InsertAddress(body)
    if noErrors == false {
        conn.Write([]byte ("Invalid parameters, unable to add address!"))
    } else {
        newEntry := this.addressBook[len(this.addressBook) - 1]
        conn.Write([]byte ("Welcome! Your id is: " + strconv.Itoa(int(newEntry.Id)) + ", you address is: " + newEntry.IP + ":" + newEntry.Port))
    }
}

func (this *Router) handlePeopleRequest(body string, conn Handler) {
    requestId, err := strconv.ParseUint(body, 10, 64)
    if err != nil {
        fmt.Println(err.Error())
        return
    }
    // Seems a bit clumsy, but will do for now
    var resultIds []string
    for _, address := range this.addressBook {
        if (address.Id != requestId) {
            resultIds = append(resultIds, strconv.FormatUint(address.Id, 10))
        }
    }
    conn.Write([]byte (strings.Join(resultIds,",")))
}

func (this *Router) handleWhoAmIRequest(body string, conn Handler) {
    request := strings.SplitN(body, ":", 2)
    if len(request) != 2 {
        conn.Write([]byte ("Invalid WHOAMI request parameters"))
        return
    }
    for _, address := range this.addressBook {
        if (address.IP == request[0] && address.Port == request[1]) {
            conn.Write([]byte ("You are: " + strconv.Itoa(int(address.Id))))
        }
    }
}

func (this *Router) handleMessageRequest(body string, conn Handler) {
    bodySplit := strings.SplitN(body, ":", 2)
    fmt.Printf("Recipients: %s Message: %s\n", bodySplit[0], bodySplit[1])
    recipients := strings.Split(bodySplit[0], ",")
    message := bodySplit[1]
    
    // 1024 kilobyte limit
    if len(message) > 1048576 {
        conn.Write([]byte ("Error: Message too long!"))
        return
    }
    // Max 255 recipients
    if len(recipients) > 255 {
        conn.Write([]byte ("Error: Too many recipients!"))
        return
    }
    
    for _, recipient := range recipients {
        recipientId, _ := strconv.ParseUint(recipient, 10, 64)
        this.MessageQueue <- Message{recipientId, message}
    }
    conn.Write([]byte ("Sent: \"" + message + "\" to users " + strings.Join(recipients, ",")))
}

// A bit lame trick for testability; hopefully refactoring for channels will fix
func (this *Router) EstablishConnection(id uint64, contact Dialer) Handler {
    var ret Handler
    for _, address := range this.addressBook {
        if address.Id == id {
            ret = contact(address)
        }
    }
    return ret
}