package echo

import (
	"context"
	"math/rand"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/mkawserm/abesh/constant"
	"github.com/mkawserm/abesh/iface"
	"github.com/mkawserm/abesh/model"
	"github.com/mkawserm/abesh/registry"
)

var RandomNumber = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "random_number",
		Help: "Random number",
	},
)

type Echo struct {
	mValues model.ConfigMap
}

func (e *Echo) Name() string {
	return "abesh_example_echo"
}

func (e *Echo) Version() string {
	return "0.0.1"
}

func (e *Echo) Category() string {
	return string(constant.CategoryService)
}

func (e *Echo) ContractId() string {
	return "abesh:ex_echo"
}

func (e *Echo) GetConfigMap() model.ConfigMap {
	return e.mValues
}

func (e *Echo) Setup() error {
	return nil
}

func (e *Echo) SetConfigMap(values model.ConfigMap) error {
	e.mValues = values

	return nil
}

func (e *Echo) New() iface.ICapability {
	return &Echo{}
}

func (e *Echo) Serve(_ context.Context, input *model.Event) (*model.Event, error) {
	m := &model.Metadata{
		Headers:        map[string]string{"Content-Type": "application/text"},
		ContractIdList: []string{e.ContractId()},
		StatusCode:     200,
		Status:         "OK",
	}

	outputEvent := &model.Event{
		Metadata: m,
		TypeUrl:  "application/text",
		Value:    []byte("default echo"),
	}

	if input.TypeUrl == "application/json" {
		m.Headers["Content-Type"] = "application/json"
		outputEvent.TypeUrl = "application/json"
		outputEvent.Value = []byte("{\"message\":\"echo\"}")
	}

	if input.TypeUrl == "application/text" {
		m.Headers["Content-Type"] = "application/text"
		outputEvent.TypeUrl = "application/text"
		outputEvent.Value = []byte("echo")
	}

	min := -1000
	max := 1000
	random := rand.Intn(max-min) + min

	RandomNumber.Set(float64(random))
	return outputEvent, nil
}

func init() {
	prometheus.MustRegister(RandomNumber)
	registry.GlobalRegistry().AddCapability(&Echo{})
}
