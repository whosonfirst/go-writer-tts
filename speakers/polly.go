package speakers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"io"
	"io/ioutil"
	"strings"
)

type PollySpeaker struct {
	Speaker
	service      *polly.Polly
	OutputFormat string
	VoiceId      string
}

func NewPollySpeaker() (*PollySpeaker, error) {

	// please fix me - this assumes shared credentials

	cfg := aws.NewConfig()
	cfg.WithRegion("us-east-1") // please fix me...

	sess, err := session.NewSession(cfg)

	if err != nil {
		return nil, err
	}

	svc := polly.New(sess)

	s := PollySpeaker{
		service:      svc,
		OutputFormat: "mp3",
		VoiceId:      "Russell",
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

	// https://docs.aws.amazon.com/polly/latest/dg/API_SynthesizeSpeech.html

	params := &polly.SynthesizeSpeechInput{
		OutputFormat: aws.String(s.OutputFormat),
		VoiceId:      aws.String(s.VoiceId),
		Text:         aws.String(string(p)),
	}

	resp, err := s.service.SynthesizeSpeech(params)

	if err != nil {
		return 0, err
	}

	// please fix me - can feed all the resp.AudioStream results
	// to a single filehandle or are there MP3 header things that
	// can't be repeated... TBD

	b, err := ioutil.ReadAll(resp.AudioStream)

	if err != nil {
		return 0, err
	}

	err = ioutil.WriteFile("test.mp3", b, 0644)

	if err != nil {
		return 0, err
	}

	// please fix me - cast *int64 as int
	// return int(resp.RequestCharacters), nil

	return 1, nil
}

func (s *PollySpeaker) Close() error {
	return nil
}
