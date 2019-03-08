// Golang HTML5 Server Side Events Example
//
// Run this code like:
//  > go run server.go
//
// Then open up your browser to http://localhost:8000
// Your browser must support HTML5 SSE, of course.

package main

import (
    "fmt"
    "log"
    "math/rand"
    "net/http"
    "time"
)

// A single Broker will be created in this program. It is responsible
// for keeping a list of which clients (browsers) are currently attached
// and broadcasting events (messages) to those clients.
//
type Broker struct {

    // Create a map of clients, the keys of the map are the channels
    // over which we can push messages to attached clients.  (The values
    // are just booleans and are meaningless.)
    //
    clients map[chan string]bool

    // Channel into which new clients can be pushed
    //
    newClients chan chan string

    // Channel into which disconnected clients should be pushed
    //
    defunctClients chan chan string

    // Channel into which messages are pushed to be broadcast out
    // to attahed clients.
    //
    messages chan string
}

// This Broker method starts a new goroutine.  It handles
// the addition & removal of clients, as well as the broadcasting
// of messages out to clients that are currently attached.
//
func (b *Broker) Start() {

    // Start a goroutine
    //
    go func() {

        // Loop endlessly
        //
        for {

            // Block until we receive from one of the
            // three following channels.
            select {

            case s := <-b.newClients:

                // There is a new client attached and we
                // want to start sending them messages.
                b.clients[s] = true
                log.Println("Added new client")

            case s := <-b.defunctClients:

                // A client has dettached and we want to
                // stop sending them messages.
                delete(b.clients, s)
                close(s)

                log.Println("Removed client")

            case msg := <-b.messages:

                // There is a new message to send.  For each
                // attached client, push the new message
                // into the client's message channel.
                for s := range b.clients {
                    s <- msg
                }
                log.Printf("Broadcast message to %d clients", len(b.clients))
            }
        }
    }()
}

// This Broker method handles and HTTP request at the "/events/" URL.
//
func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    // Make sure that the writer supports flushing.
    //
    f, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
        return
    }

    // Create a new channel, over which the broker can
    // send this client messages.
    messageChan := make(chan string)

    // Add this client to the map of those that should
    // receive updates
    b.newClients <- messageChan

    // Listen to the closing of the http connection via the CloseNotifier
    notify := w.(http.CloseNotifier).CloseNotify()
    go func() {
        <-notify
        // Remove this client from the map of attached clients
        // when `EventHandler` exits.
        b.defunctClients <- messageChan
        log.Println("HTTP connection just closed.")
    }()

    // Set the headers related to event streaming.
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
    w.Header().Set("Transfer-Encoding", "chunked")

    // Don't close the connection, instead loop endlessly.
    for {

        // Read from our messageChan.
        msg, open := <-messageChan

        if !open {
            // If our messageChan was closed, this means that the client has
            // disconnected.
            break
        }

        // Write to the ResponseWriter, `w`.
        fmt.Fprintf(w, "data: %s\n\n", msg)

        // Flush the response.  This is only possible if
        // the repsonse supports streaming.
        f.Flush()
    }

    // Done.
    log.Println("Finished HTTP request at ", r.URL.Path)
}

// Main routine
//

func runSse() {

    // Make a new Broker instance
    b := &Broker{
        make(map[chan string]bool),
        make(chan (chan string)),
        make(chan (chan string)),
        make(chan string),
    }

    // Start processing events
    b.Start()

    // Make b the HTTP handler for "/events/".  It can do
    // this because it has a ServeHTTP method.  That method
    // is called in a separate goroutine for each
    // request to "/events/".
    http.Handle("/events/", b)

    // Generate a constant stream of events that get pushed
    // into the Broker's messages channel and are then broadcast
    // out to any clients that are attached. This will be replaced
    // with something else Kafka?
    go func() {
        var c, m string
        var t int
        for i := 0; ; i++ {
            if i > 55 {i = 0}

            // Create a little message to send to clients,
            // including the current time.
            b.messages <- fmt.Sprintf("%d,green, the time is %v", i, time.Now())
            if rand.Intn(30) == 1 {
                t = rand.Intn(55)
                switch rand.Intn(2) {
                case 0:
                    c = "red"
                    m = "PANIC! Red Alert"
                case 1:
                    c = "amber"
                    m = "OH NOES! Something's not quite right"
                case 2:
                    c = "grey"
                    m = "Meh, not sure what to do now...."
                }
                b.messages <- fmt.Sprintf("%d,%s,%s",t ,c ,m )
            }


            // Print a nice log message and sleep for 5s.
            log.Printf("Sent message")
            time.Sleep(1e9)

        }
    }()
}
