package config

type Config interface {
	ModelPath() string
}

type config struct {
	modelPath string
}

func (c *config) ModelPath() string {
	return c.modelPath
}

func NewConfig() Config {
	return &config{
		modelPath: "/model/scrfd_10g_bnkps_shape640x640OnnxV6.onnx",
	}
}
