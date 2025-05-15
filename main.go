package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/mux"
)

// Variables used for command line parameters
var (
	Token      string
	AuthToken  string
	Webhookurl string
	Address    string
	Port       string
	clientHttp *resty.Client
)

type Msg struct {
	Body string `json:"body"`
}

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&AuthToken, "a", "", "Auth Token")
	flag.StringVar(&Webhookurl, "w", "", "Webhook URL")
	flag.StringVar(&Address, "l", "", "Bind IP Address")
	flag.StringVar(&Port, "p", "", "Bind Port")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Println("error creating Discord session,", err)
		return
	}

	if Address == "" {
		Address = "127.0.0.1"
	}
	if Port == "" {
		Port = "8001"
	}
	if !strings.HasPrefix(Webhookurl, "http") {
		log.Fatal("Invalid Webhook address")
	}

	clientHttp = resty.New()
	clientHttp.SetRedirectPolicy(resty.FlexibleRedirectPolicy(15))
	clientHttp.SetTimeout(5 * time.Second)
	clientHttp.SetDebug(false)

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)
	dg.LogLevel = discordgo.LogWarning

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		return
	}

	// web server?
	r := mux.NewRouter()
	r.HandleFunc("/send/{channelid}", send(dg)).Methods("POST")

	srv := &http.Server{
		Addr:    Address + ":" + Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %s\n", err)
		}
	}()

	log.Printf("Started server on %s", srv.Addr)

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Discord bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	defer close(sc)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()

}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	postmap := make(map[string]interface{})
	postmap["event"] = "Message"
	postmap["body"] = m.Content
	postmap["avatar"] = m.Author.AvatarURL("")
	postmap["author"] = m.Author
	postmap["channelid"] = m.ChannelID
	values, _ := json.Marshal(postmap)
	data := make(map[string]string)
	data["jsonData"] = string(values)
	go callHook(Webhookurl, data)

}

func callHook(myurl string, payload map[string]string) {
	log.Println("Message recieved from Discord. Call WebHook")
	clientHttp.R().SetFormData(payload).Post(myurl)
}

func send(s *discordgo.Session) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		vars := mux.Vars(request)
		channelid := vars["channelid"]

		htoken := request.Header.Get("Token")
		if htoken == "" {
			htoken = strings.Join(request.URL.Query()["token"], "")
		}
		if AuthToken != htoken {
			log.Println("Auth failure")
			return
		}

		b, err := io.ReadAll(request.Body)

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		var msg Msg

		if err := json.Unmarshal(b, &msg); err != nil {
			if jsonErr, ok := err.(*json.SyntaxError); ok {
				problemPart := b[jsonErr.Offset-10 : jsonErr.Offset+10]
				log.Printf("%w ~ error near '%s' (offset %d)", err, problemPart, jsonErr.Offset)
			}
		} else {
			log.Printf("Message received from Agent: %s", msg.Body)
		}

		s.ChannelMessageSend(channelid, msg.Body)
		err = json.NewEncoder(writer).Encode(&msg)
		if err != nil {
			log.Println("There was an error encoding the initialized struct")
		}
	}
}
