package main

import (
	"errors"
	"log"
	"math"
	"runtime"

	whisper "github.com/ggerganov/whisper.cpp/bindings/go"
)

type WhisperModel struct {
	ctx    *whisper.Context
	params whisper.Params
}

func NewWhisperModel() (*WhisperModel, error) {
	ctx := whisper.Whisper_init("./models/ggml-base.en.bin")
	if ctx == nil {
		return nil, errors.New("failed to initialize whisper")
	}

	params := ctx.Whisper_full_default_params(whisper.SAMPLING_GREEDY)
	params.SetPrintProgress(false)
	params.SetPrintSpecial(false)
	params.SetPrintRealtime(false)
	params.SetPrintTimestamps(false)
	params.SetSingleSegment(false)
	params.SetMaxTokensPerSegment(32)
	params.SetThreads(int(math.Min(float64(4), float64(runtime.NumCPU()))))
	params.SetSpeedup(false)
	params.SetLanguage(ctx.Whisper_lang_id("en"))

	return &WhisperModel{ctx: ctx, params: params}, nil
}

func (w *WhisperModel) Process(samples []float32) error {
	if err := w.ctx.Whisper_full(w.params, samples, nil, nil); err != nil {
		return err
	} else {
		segments := w.ctx.Whisper_full_n_segments()
		for i := 0; i < segments; i++ {
			text := w.ctx.Whisper_full_get_segment_text(i)
			log.Printf("Segment %d: %s", i, text)
		}
	}
	return nil
}
