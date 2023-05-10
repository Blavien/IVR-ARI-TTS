// Command quickstart generates an audio file with the content "Hello, World!".
package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/CyCoreSystems/ari/v6/ext/play"
	"github.com/inconshreveable/log15"
	"os"

	"github.com/CyCoreSystems/ari/v6"
	"github.com/CyCoreSystems/ari/v6/client/native"
	"io/ioutil"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
)

var log = log15.New()

func main() {

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// connect
	native.Logger = log

	/*log.Info("Connecting to ARI")
	cl, err := native.Connect(&native.Options{
		Application:  "test",
		Username:     "admin",
		Password:     "admin",
		URL:          "http://localhost:8088/ari",
		WebsocketURL: "ws://localhost:8088/ari/events",
	})
	if err != nil {
		log.Error("Failed to build ARI client", "error", err)
		return
	}

	// setup app

	log.Info("Listening for new calls")
	sub := cl.Bus().Subscribe(nil, "StasisStart")

	for {
		select {
		case e := <-sub.Events():
			v := e.(*ari.StasisStart)
			log.Info("Got stasis start", "channel", v.Channel.ID)
			go app(ctx, cl.Channel().Get(v.Key(ari.ChannelKey, v.Channel.ID)))
		case <-ctx.Done():
			return
		}
	}*/
	textToSpeech()
}

func textToSpeech() {
	// Instantiates a client.
	ctx := context.Background()
	scanner := bufio.NewScanner(os.Stdin)
	audioText := ""

	fmt.Print("Input something: ")
	if scanner.Scan() {
		audioText := scanner.Text()
		fmt.Println("Text that's going to be transformed: " + audioText)
	}

	client, err := texttospeech.NewClient(ctx)
	if err != nil {
		log.Error(err.Error())
	}
	defer client.Close()

	// Perform the text-to-speech request on the text input with the selected
	// voice parameters and audio file type.

	req := texttospeechpb.SynthesizeSpeechRequest{
		// Set the text input to be synthesized.
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{Text: audioText},
		},
		// Build the voice request, select the language code ("en-US") and the SSML
		// voice gender ("neutral").
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: "pt-PT",
			SsmlGender:   texttospeechpb.SsmlVoiceGender_FEMALE,
		},
		// Select the type of audio file you want returned.
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, &req)
	if err != nil {
		log.Error(err.Error())
	}

	// The resp's AudioContent is binary.
	filename := "output.mp3"
	err = ioutil.WriteFile(filename, resp.AudioContent, 0644)
	if err != nil {
		log.Error(err.Error())
	}
	fmt.Printf("Audio content written to file: %v\n", filename)
}

func app(ctx context.Context, h *ari.ChannelHandle) {
	defer h.Hangup()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log.Info("Running app", "channel", h.ID())

	end := h.Subscribe(ari.Events.StasisEnd)
	defer end.Cancel()

	// End the app when the channel goes away
	go func() {
		<-end.Events()
		cancel()
	}()

	if err := h.Answer(); err != nil {
		log.Error("failed to answer call", "error", err)
		return
	}

	textToSpeech()

	if err := play.Play(ctx, h, play.URI("file://output.mp3")).Err(); err != nil {
		log.Error("failed to play sound", "error", err)
		return
	}

	log.Info("completed playback")
	return
}
