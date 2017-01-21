package speakers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

type PollySpeaker struct {
	Speaker
	service *polly.Polly
}

func NewPollySpeaker() (*PollySpeaker, error) {

	sess, err := session.NewSession()

	if err != nil {
		return nil, err
	}

	svc := polly.New(sess)

	s := PollySpeaker{
		service: svc,
	}

	return &s, nil
}

func (s *PollySpeaker) Read(reader io.Reader) error {

	tee := io.TeeReader(reader, s)
	_, err := ioutil.ReadAll(tee)
	return err
}

func (s *PollySpeaker) WriteString(text string) (int64, error) {
	r := strings.NewReader(text)
	return r.WriteTo(s)
}

func (s *PollySpeaker) Write(p []byte) (int, error) {

	params := &polly.SynthesizeSpeechInput{
		OutputFormat: aws.String("OutputFormat"), // Required
		Text:         aws.String("Text"),         // Required
		VoiceId:      aws.String("VoiceId"),      // Required
		LexiconNames: []*string{
			aws.String("LexiconName"), // Required
			// More values...
		},
		SampleRate: aws.String("SampleRate"),
		TextType:   aws.String("TextType"),
	}

	resp, err := s.service.SynthesizeSpeech(params)

	log.Println("ERR")
	log.Println(err)

	if err != nil {
		return 0, err
	}

	log.Println(resp)

	return 0, nil
}

func (s *PollySpeaker) Close() error {
	return nil
}
