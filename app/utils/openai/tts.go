package openai

import (
	"bytes"
	"io"
	"os"
	"time"

	"github.com/golly-go/golly"
)

type TTSVoice string

const (
	VoiceAlloy   TTSVoice = "alloy"
	VoiceShimmer TTSVoice = "shimmer"
	VoiceEcho    TTSVoice = "echo"
	VoiceFable   TTSVoice = "fable"
	VoiceNova    TTSVoice = "nova"
	VoiceOnyx    TTSVoice = "onyx"
)

func (oai OpenAIClient) TTS(ctx golly.Context, filePath string, voice TTSVoice, text string) error {
	defer func(start time.Time) {
		ctx.Logger().
			Infof("TTS Request took %s", time.Since(start))
	}(time.Now())

	params := map[string]any{"model": TTSModel, "input": text, "voice": voice}

	b, err := oai.request(ctx, openAITTSURL, params)
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(b)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, buffer)

	return err
}
