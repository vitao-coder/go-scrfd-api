package recognition

import (
	"fmt"
	"github.com/yalue/onnxruntime_go"
	ort "github.com/yalue/onnxruntime_go"
	"go-scrfd-api/internal/config"
	"go-scrfd-api/internal/domain"
	"go-scrfd-api/pkg/onnx"
	"gocv.io/x/gocv"
	"image"
)

type UseCase interface {
	Detect(image []byte) ([]domain.PredictionBox, error)
}

type useCase struct {
	featStrideFpn []int
	nmsThreshold  float32
	anchorCenters map[string][]image.Point
	height        int
	width         int
	configuration config.Config
}

func NewUseCase(configuration config.Config) UseCase {
	return &useCase{
		featStrideFpn: []int{8, 16, 32},
		nmsThreshold:  0.4,
		anchorCenters: make(map[string][]image.Point),
		height:        640,
		width:         640,
		configuration: configuration,
	}
}

func (uc *useCase) Detect(image []byte) ([]domain.PredictionBox, error) {
	img, err := gocv.IMDecode(image, gocv.IMReadColor)
	if err != nil {
		return nil, err
	}
	defer img.Close()

	// Preprocess the image and convert to tensor
	height := uc.height
	width := uc.width

	inputTensorData := make([]float32, 1*3*height*width)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			vec := img.GetVecbAt(y, x)
			inputTensorData[(0*3*height*width)+(0*height*width)+(y*width)+x] = (float32(vec[2]) - 127.5) / 128.0
			inputTensorData[(0*3*height*width)+(1*height*width)+(y*width)+x] = (float32(vec[1]) - 127.5) / 128.0
			inputTensorData[(0*3*height*width)+(2*height*width)+(y*width)+x] = (float32(vec[0]) - 127.5) / 128.0
		}
	}

	//Ort initialization
	sharedLibPath := onnx.GetSharedLibPath()
	ort.SetSharedLibraryPath(sharedLibPath)
	err = ort.InitializeEnvironment()
	if err != nil {
		return nil, err
	}

	// Define input tensor
	inputTensor, err := onnxruntime_go.NewTensor([]int64{1, 3, int64(height), int64(width)}, inputTensorData)
	if err != nil {
		return nil, err
	}
	defer inputTensor.Destroy()

	// Define output tensors
	outputTensors := make([]*onnxruntime_go.Tensor[float32], 3)
	for i := range outputTensors {
		outputTensors[i], err = onnxruntime_go.NewEmptyTensor[float32]([]int64{1, 3, int64(height), int64(width)}) // Adjust dimensions as needed
		if err != nil {
			return nil, err
		}
		defer outputTensors[i].Destroy()
	}

	session, err := onnx.GetONNXSession(uc.configuration, inputTensor, outputTensors)
	if err != nil {
		return nil, err
	}
	defer session.Session.Destroy()

	// Run the model
	err = session.Session.Run()
	if err != nil {
		return nil, err
	}

	// Parse the output and create PredictionBox results
	scoresData := outputTensors[0].GetData()
	bboxsData := outputTensors[3].GetData()
	kpssData := outputTensors[6].GetData()
	results := uc.parseOutput(scoresData, bboxsData, kpssData, float32(height)/float32(img.Rows()))

	return results, nil

}

func (uc *useCase) parseOutput(scoresData, bboxsData, kpssData []float32, resizeRate float32) []domain.PredictionBox {
	preds := []domain.PredictionBox{}

	for _, stride := range uc.featStrideFpn {
		sHeight := uc.height / stride
		sWidth := uc.width / stride

		for i := 0; i < len(scoresData); i++ {
			score := scoresData[i]
			if score <= uc.nmsThreshold {
				continue
			}

			keyAnchorCenter := fmt.Sprintf("%d-%d-%d", sHeight, sWidth, stride)
			anchorCenter, found := uc.anchorCenters[keyAnchorCenter]
			if !found {
				anchorCenter = generateAnchorCenter(sHeight, sWidth, stride)
				uc.anchorCenters[keyAnchorCenter] = anchorCenter
			}

			box := scaleBox(bboxsData, stride, i)
			keypoints := scaleKeyPoints(kpssData, stride, i)

			box = distanceToBox(box, anchorCenter[i], resizeRate)
			keypoints = distanceToPoint(keypoints, anchorCenter[i], resizeRate)

			preds = append(preds, domain.PredictionBox{
				Score:     score,
				BoxLeft:   box[0],
				BoxTop:    box[1],
				BoxRight:  box[2],
				BoxBottom: box[3],
				KeyPoints: keypoints,
			})
		}
	}

	return preds
}

func generateAnchorCenter(sHeight, sWidth, stride int) []image.Point {
	anchors := make([]image.Point, sHeight*sWidth)
	for h := 0; h < sHeight; h++ {
		for w := 0; w < sWidth; w++ {
			anchors[h*sWidth+w] = image.Point{X: w * stride, Y: h * stride}
		}
	}
	return anchors
}

func scaleBox(bboxsData []float32, stride, index int) []float32 {
	box := make([]float32, 4)
	for b := 0; b < 4; b++ {
		box[b] = bboxsData[index*4+b] * float32(stride)
	}
	return box
}

func scaleKeyPoints(kpssData []float32, stride, index int) []float32 {
	keypoints := make([]float32, 10)
	for k := 0; k < 10; k++ {
		keypoints[k] = kpssData[index*10+k] * float32(stride)
	}
	return keypoints
}

func distanceToBox(box []float32, anchorCenter image.Point, rate float32) []float32 {
	box[0] = -box[0]
	box[1] = -box[1]
	return distanceToPoint(box, anchorCenter, rate)
}

func distanceToPoint(distances []float32, anchorCenter image.Point, rate float32) []float32 {
	for i := 0; i < len(distances); i += 2 {
		distances[i] = (float32(anchorCenter.X) + distances[i]) * rate
		distances[i+1] = (float32(anchorCenter.Y) + distances[i+1]) * rate
	}
	return distances
}
