package onnx

import (
	"github.com/yalue/onnxruntime_go"
	"go-scrfd-api/internal/config"
	"runtime"
)

type Session struct {
	Session *onnxruntime_go.Session[float32]
}

func GetONNXSession(cfg config.Config, inputTensor *onnxruntime_go.Tensor[float32], outputTensor []*onnxruntime_go.Tensor[float32]) (*Session, error) {

	inputNames := []string{"input.1"}
	outputNames := []string{"score_8", "score_16", "score_32", "bbox_8", "bbox_16", "bbox_32", "kps_8", "kps_16", "kps_32"}

	session, err := onnxruntime_go.NewSession[float32](cfg.ModelPath(), inputNames, outputNames, []*onnxruntime_go.Tensor[float32]{inputTensor}, outputTensor)
	if err != nil {
		return nil, err
	}
	return &Session{Session: session}, nil
}

func GetSharedLibPath() string {
	if runtime.GOOS == "windows" {
		if runtime.GOARCH == "amd64" {
			return "../third_party/onnxruntime.dll"
		}
	}
	if runtime.GOOS == "darwin" {
		if runtime.GOARCH == "arm64" {
			return "../third_party/onnxruntime_arm64.dylib"
		}
		if runtime.GOARCH == "amd64" {
			return "../third_party/onnxruntime_amd64.dylib"
		}

	}
	if runtime.GOOS == "linux" {
		if runtime.GOARCH == "arm64" {
			return "../third_party/onnxruntime_arm64.so"
		}
		return "../third_party/onnxruntime.so"
	}
	panic("Unable to find a version of the onnxruntime library supporting this system.")
}
