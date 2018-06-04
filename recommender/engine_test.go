package recommender

import (
	"context"
	"errors"
	"testing"

	"github.com/banzaicloud/telescopes/productinfo"
	"github.com/stretchr/testify/assert"
)

var (
	dummyVmInfo = []productinfo.VmInfo{
		{OnDemandPrice: float64(10), SpotPrice: map[string]float64{"zonea": 0.021, "zoneb": 0.022, "zonec": 0.026}, Cpus: float64(10), Mem: float64(10), Gpus: float64(0)},
		{OnDemandPrice: float64(12), SpotPrice: map[string]float64{"zonea": 0.043}, Cpus: float64(10), Mem: float64(10), Gpus: float64(0)},
		{OnDemandPrice: float64(21), SpotPrice: map[string]float64{"zonea": 0.032}, Cpus: float64(12), Mem: float64(12), Gpus: float64(0)},
	}
	trueVal  = true
	falseVal = false
)

func TestNewEngine(t *testing.T) {
	tests := []struct {
		name    string
		pi      productinfo.ProductInfo
		checker func(engine *Engine, err error)
	}{
		{
			name: "engine successfully created",
			pi:   &dummyProductInfo{},
			checker: func(engine *Engine, err error) {
				assert.Nil(t, err, "should not get error ")
				assert.NotNil(t, engine, "the engine should not be nil")
			},
		},
		{
			name: "engine creation fails when registries is nil",
			pi:   nil,
			checker: func(engine *Engine, err error) {
				assert.Nil(t, engine, "the engine should be nil")
				assert.NotNil(t, err, "the error shouldn't be nil")
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.checker(NewEngine(test.pi))

		})
	}
}

// dummy network mapper
type dummyNetworkMapper struct {
}

func (dm dummyNetworkMapper) MapNetworkPerf(vm productinfo.VmInfo) (string, error) {
	return productinfo.NTW_HIGH, nil
}

// utility VmRegistry for mocking purposes
type dummyProductInfo struct {
	// test case id to drive the behaviour
	TcId string
}

func (d *dummyProductInfo) Start(ctx context.Context) {

}
func (d *dummyProductInfo) GetAttrValues(provider string, attribute string) ([]float64, error) {
	switch d.TcId {
	case "3 values between 10-20":
		return []float64{12, 13, 14}, nil
	case "2 values between 10 - 20":
		return []float64{8, 13, 14, 6}, nil
	case "all values are larger than 20, return the closest value":
		return []float64{30, 40, 50, 60}, nil
	case "all values are less than 10, return the closest value":
		return []float64{1, 2, 3, 5, 9}, nil
	case "error, min > max":
		return []float64{1}, nil
	case "error returned":
		return nil, errors.New("error")

	}
	return []float64{15, 16, 17}, nil
}

func (d *dummyProductInfo) GetVmsWithAttrValue(provider string, regionId string, attrKey string, value float64) ([]productinfo.VmInfo, error) {
	if d.TcId == "attribute value error" {
		return nil, errors.New("attribute value error")
	}
	return dummyVmInfo, nil
}
func (d *dummyProductInfo) GetZones(provider string, region string) ([]string, error) {
	if d.TcId == "zone error" {
		return nil, errors.New("no zone available")
	}
	return nil, nil
}

func (d *dummyProductInfo) GetSpotPrice(provider string, region string, instanceType string, zones []string) (float64, error) {
	return float64(6), nil
}

func (d *dummyProductInfo) GetNetworkPerfMapper(provider string) (productinfo.NetworkPerfMapper, error) {
	nm := dummyNetworkMapper{}
	return nm, nil
}

func TestEngine_RecommendAttrValues(t *testing.T) {

	tests := []struct {
		name      string
		pi        productinfo.ProductInfo
		request   ClusterRecommendationReq
		provider  string
		attribute string
		check     func([]float64, error)
	}{
		{
			name: "all attributes between limits",
			pi:   &dummyProductInfo{"3 values between 10-20"},
			request: ClusterRecommendationReq{
				MinNodes: 5,
				MaxNodes: 10,
				SumMem:   100,
				SumCpu:   100,
			},
			provider:  "dummy",
			attribute: productinfo.Cpu,

			check: func(values []float64, err error) {
				assert.Nil(t, err, "should not get error when recommending attributes")
				assert.Equal(t, 3, len(values), "recommended number of values is not as expected")
				assert.Equal(t, []float64{12, 13, 14}, values)

			},
		},
		{
			name: "attributes out of limits not recommended",
			pi:   &dummyProductInfo{"2 values between 10 - 20"},
			request: ClusterRecommendationReq{
				MinNodes: 5,
				MaxNodes: 10,
				SumMem:   100,
				SumCpu:   100,
			},
			provider:  "dummy",
			attribute: productinfo.Cpu,

			check: func(values []float64, err error) {
				assert.Nil(t, err, "should not get error when recommending attributes")
				assert.Equal(t, 2, len(values), "recommended number of values is not as expected")

			},
		},
		{
			name: "no values between limits found - smallest value returned",
			pi:   &dummyProductInfo{"all values are larger than 20, return the closest value"},
			request: ClusterRecommendationReq{
				MinNodes: 5,
				MaxNodes: 10,
				SumMem:   100,
				SumCpu:   100,
			},
			provider:  "dummy",
			attribute: productinfo.Cpu,

			check: func(values []float64, err error) {
				assert.Nil(t, err, "should not get error when recommending attributes")
				assert.Equal(t, 1, len(values), "recommended number of values is not as expected")
				assert.Equal(t, float64(30), values[0], "recommended number of values is not as expected")

			},
		},
		{
			name: "no values between limits found - largest value returned",
			pi:   &dummyProductInfo{"all values are less than 10, return the closest value"},
			request: ClusterRecommendationReq{
				MinNodes: 5,
				MaxNodes: 10,
				SumMem:   100,
				SumCpu:   100,
			},
			provider:  "dummy",
			attribute: productinfo.Cpu,

			check: func(values []float64, err error) {
				assert.Nil(t, err, "should not get error when recommending attributes")
				assert.Equal(t, 1, len(values), "recommended number of values is not as expected")
				assert.Equal(t, float64(9), values[0], "recommended number of values is not as expected")

			},
		},
		{
			name: "error - min larger than max",
			pi:   &dummyProductInfo{"error, min > max"},
			request: ClusterRecommendationReq{
				MinNodes: 10,
				MaxNodes: 5,
				SumMem:   100,
				SumCpu:   100,
			},
			provider:  "dummy",
			attribute: productinfo.Cpu,

			check: func(values []float64, err error) {
				assert.Equal(t, err.Error(), "min value cannot be larger than the max value")

			},
		},
		{
			name: "error - attribute values could not be retrieved",
			pi:   &dummyProductInfo{"error returned"},
			request: ClusterRecommendationReq{
				MinNodes: 10,
				MaxNodes: 10,
				SumMem:   100,
				SumCpu:   100,
			},
			provider:  "dummy",
			attribute: productinfo.Cpu,

			check: func(values []float64, err error) {
				assert.Nil(t, values, "returned attr values should be nils")
				assert.EqualError(t, err, "error")

			},
		},
		{
			name: "error - unsupported attribute",
			pi:   &dummyProductInfo{},
			request: ClusterRecommendationReq{
				MinNodes: 5,
				MaxNodes: 10,
				SumMem:   100,
				SumCpu:   100,
			},
			provider:  "dummy",
			attribute: "error",

			check: func(values []float64, err error) {
				assert.Nil(t, values, "the values should be nil")
				assert.EqualError(t, err, "unsupported attribute: [error]")

			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			engine, err := NewEngine(test.pi)
			assert.Nil(t, err, "the engine couldn't be created")

			test.check(engine.RecommendAttrValues("dummy", test.attribute, test.request))

		})
	}
}

func TestEngine_RecommendVms(t *testing.T) {
	tests := []struct {
		name      string
		region    string
		pi        productinfo.ProductInfo
		values    []float64
		filters   []vmFilter
		request   ClusterRecommendationReq
		provider  string
		attribute string
		check     func([]VirtualMachine, error)
	}{
		{
			name:   "error - findVmsWithAttrValues",
			region: "us-west-2",
			pi:     &dummyProductInfo{"attribute value error"},
			values: []float64{1, 2},
			check: func(vms []VirtualMachine, err error) {
				assert.EqualError(t, err, "attribute value error")
				assert.Nil(t, vms, "the vms should be nil")

			},
		},
		{
			name:   "success - recommend vms",
			region: "us-west-2",
			filters: []vmFilter{func(vm VirtualMachine, req ClusterRecommendationReq) bool {
				return true
			}},
			pi:        &dummyProductInfo{},
			values:    []float64{2},
			provider:  "dummy",
			attribute: productinfo.Cpu,

			check: func(vms []VirtualMachine, err error) {
				assert.Nil(t, err, "the error should be nil")
				assert.Equal(t, 3, len(vms))
				assert.Equal(t, []VirtualMachine{VirtualMachine{Type: "", AvgPrice: 6, OnDemandPrice: 10, Cpus: 10,
					Mem: 10, Gpus: 0, Burst: false, NetworkPerf: "high"}, VirtualMachine{Type: "", AvgPrice: 6, OnDemandPrice: 12,
					Cpus: 10, Mem: 10, Gpus: 0, Burst: false, NetworkPerf: "high"}, VirtualMachine{Type: "", AvgPrice: 6,
					OnDemandPrice: 21, Cpus: 12, Mem: 12, Gpus: 0, Burst: false, NetworkPerf: "high"}}, vms)

			},
		},
		{
			name:   "could not find any VMs to recommender",
			region: "us-west-2",
			filters: []vmFilter{func(vm VirtualMachine, req ClusterRecommendationReq) bool {
				return false
			}},
			pi:        &dummyProductInfo{},
			values:    []float64{1, 2},
			provider:  "dummy",
			attribute: productinfo.Cpu,

			check: func(vms []VirtualMachine, err error) {
				assert.EqualError(t, err, "couldn't find any VMs to recommend")
				assert.Nil(t, vms, "the vms should be nil")

			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			engine, err := NewEngine(test.pi)
			assert.Nil(t, err, "the engine couldn't be created")

			test.check(engine.RecommendVms("dummy", test.region, test.attribute, test.values, test.filters, test.request))

		})
	}
}

func TestEngine_RecommendNodePools(t *testing.T) {
	vms := []VirtualMachine{
		{OnDemandPrice: float64(10), AvgPrice: 99, Cpus: float64(10), Mem: float64(10), Gpus: float64(0)},
		{OnDemandPrice: float64(12), AvgPrice: 89, Cpus: float64(10), Mem: float64(10), Gpus: float64(0)},
	}

	tests := []struct {
		name    string
		pi      productinfo.ProductInfo
		attr    string
		vms     []VirtualMachine
		values  []float64
		request ClusterRecommendationReq
		check   func([]NodePool, error)
	}{
		{
			name:   "successful",
			pi:     &dummyProductInfo{},
			vms:    vms,
			attr:   productinfo.Cpu,
			values: []float64{4},
			request: ClusterRecommendationReq{
				MinNodes: 5,
				MaxNodes: 10,
				SumMem:   100,
				SumCpu:   100,
			},
			check: func(nps []NodePool, err error) {
				assert.Nil(t, err, "the error should be nil")
				assert.Equal(t, nps, []NodePool{NodePool{VmType: VirtualMachine{Type: "", AvgPrice: 99, OnDemandPrice: 10,
					Cpus: 10, Mem: 10, Gpus: 0, Burst: false}, SumNodes: 0, VmClass: "regular"},
					NodePool{VmType: VirtualMachine{Type: "", AvgPrice: 89, OnDemandPrice: 12, Cpus: 10, Mem: 10, Gpus: 0,
						Burst: false}, SumNodes: 5, VmClass: "spot"}, NodePool{VmType: VirtualMachine{Type: "", AvgPrice: 99,
						OnDemandPrice: 10, Cpus: 10, Mem: 10, Gpus: 0, Burst: false}, SumNodes: 5, VmClass: "spot"}})
			},
		},
		{
			name:   "attribute error",
			pi:     &dummyProductInfo{},
			vms:    vms,
			attr:   "error",
			values: []float64{4},
			request: ClusterRecommendationReq{
				MinNodes: 5,
				MaxNodes: 10,
				SumMem:   100,
				SumCpu:   100,
			},
			check: func(nps []NodePool, err error) {
				assert.EqualError(t, err, "could not get sum for attr: [error], cause: [unsupported attribute: [error]]")
				assert.Nil(t, nps, "the nps should be nil")

			},
		},
		{
			name:   "no spot instances available",
			pi:     &dummyProductInfo{},
			vms:    []VirtualMachine{{OnDemandPrice: float64(2), AvgPrice: 0, Cpus: float64(10), Mem: float64(10), Gpus: float64(0)}},
			attr:   productinfo.Cpu,
			values: []float64{4},
			request: ClusterRecommendationReq{
				MinNodes: 5,
				MaxNodes: 10,
				SumMem:   100,
				SumCpu:   100,
			},
			check: func(nps []NodePool, err error) {
				assert.EqualError(t, err, "no vm's suitable for spot pools")
				assert.Nil(t, nps, "the nps should be nil")

			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			engine, err := NewEngine(test.pi)
			assert.Nil(t, err, "the engine couldn't be created")

			test.check(engine.RecommendNodePools(test.attr, test.vms, test.values, test.request))

		})
	}
}

func TestEngine_RecommendCluster(t *testing.T) {
	tests := []struct {
		name     string
		pi       productinfo.ProductInfo
		request  ClusterRecommendationReq
		provider string
		region   string
		check    func(response *ClusterRecommendationResp, err error)
	}{
		{
			name: "cluster recommendation success",
			pi:   &dummyProductInfo{},
			request: ClusterRecommendationReq{
				MinNodes: 5,
				MaxNodes: 10,
				SumMem:   100,
				SumCpu:   100,
				Zones:    []string{"testZone1", "testZone2"},
			},
			provider: "dummy",
			region:   "us-west-2",
			check: func(response *ClusterRecommendationResp, err error) {
				assert.Nil(t, err, "should not get error when recommending")
				assert.Equal(t, &ClusterRecommendationResp{Provider: "dummy", Zones: []string{"testZone1", "testZone2"},
					NodePools: []NodePool{
						{VmType: VirtualMachine{Type: "", AvgPrice: 6, OnDemandPrice: 10, Cpus: 10, Mem: 10, Gpus: 0, Burst: false, NetworkPerf: "high"}, SumNodes: 0, VmClass: "regular"},
						{VmType: VirtualMachine{Type: "", AvgPrice: 6, OnDemandPrice: 21, Cpus: 12, Mem: 12, Gpus: 0, Burst: false, NetworkPerf: "high"}, SumNodes: 3, VmClass: "spot"},
						{VmType: VirtualMachine{Type: "", AvgPrice: 6, OnDemandPrice: 21, Cpus: 12, Mem: 12, Gpus: 0, Burst: false, NetworkPerf: "high"}, SumNodes: 2, VmClass: "spot"},
						{VmType: VirtualMachine{Type: "", AvgPrice: 6, OnDemandPrice: 21, Cpus: 12, Mem: 12, Gpus: 0, Burst: false, NetworkPerf: "high"}, SumNodes: 2, VmClass: "spot"},
						{VmType: VirtualMachine{Type: "", AvgPrice: 6, OnDemandPrice: 10, Cpus: 10, Mem: 10, Gpus: 0, Burst: false, NetworkPerf: "high"}, SumNodes: 2, VmClass: "spot"},
						{VmType: VirtualMachine{Type: "", AvgPrice: 6, OnDemandPrice: 12, Cpus: 10, Mem: 10, Gpus: 0, Burst: false, NetworkPerf: "high"}, SumNodes: 0, VmClass: "spot"},
						{VmType: VirtualMachine{Type: "", AvgPrice: 6, OnDemandPrice: 10, Cpus: 10, Mem: 10, Gpus: 0, Burst: false, NetworkPerf: "high"}, SumNodes: 0, VmClass: "spot"}}},
					response)
			},
		},
		{
			name: "error - RecommendAttrValues, min value cannot be larger than the max value",
			pi:   &dummyProductInfo{},
			request: ClusterRecommendationReq{
				MinNodes: 10,
				MaxNodes: 5,
				SumMem:   100,
				SumCpu:   100,
			},
			provider: "dummy",
			region:   "us-west-2",
			check: func(response *ClusterRecommendationResp, err error) {
				assert.EqualError(t, err, "could not get values for attr: [cpu], cause: [min value cannot be larger than the max value]")
				assert.Nil(t, response, "the response should be nil")
			},
		},
		{
			name: "error - RecommendVms, no zone available",
			pi:   &dummyProductInfo{"zone error"},
			request: ClusterRecommendationReq{
				MinNodes: 5,
				MaxNodes: 10,
				SumMem:   100,
				SumCpu:   100,
			},
			provider: "dummy",
			region:   "us-west-2",
			check: func(response *ClusterRecommendationResp, err error) {
				assert.EqualError(t, err, "could not get virtual machines for attr: [cpu], cause: [no zone available]")
				assert.Nil(t, response, "the response should be nil")
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			engine, err := NewEngine(test.pi)
			assert.Nil(t, err, "the engine couldn't be created")

			test.check(engine.RecommendCluster(test.provider, test.region, test.request))

		})
	}
}

func TestEngine_sortByAttrValue(t *testing.T) {
	vms := []VirtualMachine{
		{OnDemandPrice: float64(10), AvgPrice: 99, Cpus: float64(10), Mem: float64(10), Gpus: float64(0)},
		{OnDemandPrice: float64(12), AvgPrice: 89, Cpus: float64(10), Mem: float64(10), Gpus: float64(0)},
	}

	tests := []struct {
		name  string
		pi    productinfo.ProductInfo
		attr  string
		vms   []VirtualMachine
		check func(err error)
	}{
		{
			name: "error - unsupported attribute",
			pi:   &dummyProductInfo{},
			attr: "error",
			vms:  vms,
			check: func(err error) {
				assert.EqualError(t, err, "unsupported attribute: [error]")
			},
		},
		{
			name: "successful - memory",
			pi:   &dummyProductInfo{},
			attr: productinfo.Memory,
			vms:  vms,
			check: func(err error) {
				assert.Nil(t, err, "the error should be nil")
			},
		},
		{
			name: "successful - cpu",
			pi:   &dummyProductInfo{},
			attr: productinfo.Cpu,
			vms:  vms,
			check: func(err error) {
				assert.Nil(t, err, "the error should be nil")
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			engine, err := NewEngine(test.pi)
			assert.Nil(t, err, "the engine couldn't be created")

			test.check(engine.sortByAttrValue(test.attr, test.vms))

		})
	}
}

func TestEngine_filtersForAttr(t *testing.T) {
	tests := []struct {
		name  string
		pi    productinfo.ProductInfo
		attr  string
		check func(vms []vmFilter, err error)
	}{
		{
			name: "error - unsupported attribute",
			pi:   &dummyProductInfo{},
			attr: "error",
			check: func(vms []vmFilter, err error) {
				assert.EqualError(t, err, "unsupported attribute: [error]")
				assert.Nil(t, vms, "the vms should be nil")
			},
		},
		{
			name: "all filters added - cpu",
			pi:   &dummyProductInfo{},
			attr: productinfo.Cpu,
			check: func(vms []vmFilter, err error) {
				assert.Equal(t, 3, len(vms), "invalid filter count")
				assert.Nil(t, err, "the error should be nil")
			},
		},
		{
			name: "all filters added - memory",
			pi:   &dummyProductInfo{},
			attr: productinfo.Memory,
			check: func(vms []vmFilter, err error) {
				assert.Equal(t, 3, len(vms), "invalid filter count")
				assert.Nil(t, err, "the error should be nil")
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			engine, err := NewEngine(test.pi)
			assert.Nil(t, err, "the engine couldn't be created")

			test.check(engine.filtersForAttr(test.attr))

		})
	}
}

func TestEngine_filtersApply(t *testing.T) {
	tests := []struct {
		name   string
		engine Engine
		vm     VirtualMachine
		req    ClusterRecommendationReq
		attr   string
		check  func(filtersApply bool)
	}{
		{
			name:   "filter applies for cpu/mem and burst allowed",
			engine: Engine{},
			// minRatio = SumCpu/SumMem = 0.5
			req: ClusterRecommendationReq{SumCpu: 4, SumMem: float64(8), AllowBurst: &trueVal},
			// ratio = Cpus/Mem = 1
			vm:   VirtualMachine{Cpus: 4, Mem: float64(4), Burst: true},
			attr: productinfo.Memory,
			check: func(filtersApply bool) {
				assert.Equal(t, true, filtersApply, "vm should pass all filters")
			},
		},
		{
			name:   "filter doesn't apply for cpu/mem and burst not allowed ",
			engine: Engine{},
			// minRatio = SumCpu/SumMem = 0.5
			req: ClusterRecommendationReq{SumCpu: 4, SumMem: float64(8), AllowBurst: &falseVal},
			// ratio = Cpus/Mem = 1
			vm:   VirtualMachine{Cpus: 4, Mem: float64(4), Burst: true},
			attr: productinfo.Cpu,
			check: func(filtersApply bool) {
				assert.Equal(t, false, filtersApply, "vm should not pass all filters")
			},
		},
		{
			name:   "filter applies for mem/cpu and burst allowed",
			engine: Engine{},
			// minRatio = AumMem/SumCpu = 2
			req: ClusterRecommendationReq{SumMem: float64(8), SumCpu: 4, AllowBurst: &trueVal},
			// ratio = Mem/Cpus = 1
			vm:   VirtualMachine{Mem: float64(20), Cpus: 4, Burst: true},
			attr: productinfo.Cpu,
			check: func(filtersApply bool) {
				assert.Equal(t, true, filtersApply, "vm should pass all filters")
			},
		},
		{
			name:   "filter doesn't apply for mem/cpu and burst not allowed ",
			engine: Engine{},
			// minRatio = AumMem/SumCpu = 2
			req: ClusterRecommendationReq{SumMem: float64(8), SumCpu: 4, AllowBurst: &falseVal},
			// ratio = Mem/Cpus = 1
			vm:   VirtualMachine{Mem: float64(20), Cpus: 4, Burst: true},
			attr: productinfo.Memory,
			check: func(filtersApply bool) {
				assert.Equal(t, false, filtersApply, "vm should not pass all filters")
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filters, err := test.engine.filtersForAttr(test.attr)
			assert.Nil(t, err, "should get filters for attribute")
			test.check(test.engine.filtersApply(test.vm, filters, test.req))
		})
	}
}

func TestEngine_minCpuRatioFilter(t *testing.T) {
	tests := []struct {
		name   string
		engine Engine
		vm     VirtualMachine
		attr   string
		req    ClusterRecommendationReq
		check  func(filterApplies bool)
	}{
		{
			name:   "minCpuRatioFilter applies",
			engine: Engine{},
			// minRatio = SumCpu/SumMem = 0.5
			req: ClusterRecommendationReq{SumCpu: 4, SumMem: float64(8)},
			// ratio = Cpus/Mem = 1
			vm:   VirtualMachine{Cpus: 4, Mem: float64(4)},
			attr: productinfo.Cpu,
			check: func(filterApplies bool) {
				assert.Equal(t, true, filterApplies, "vm should pass the minCpuRatioFilter")
			},
		},
		{
			name:   "minCpuRatioFilter doesn't apply",
			engine: Engine{},
			// minRatio = SumCpu/SumMem = 1
			req: ClusterRecommendationReq{SumCpu: 4, SumMem: float64(4)},
			// ratio = Cpus/Mem = 0.5
			vm:   VirtualMachine{Cpus: 4, Mem: float64(8)},
			attr: productinfo.Cpu,
			check: func(filterApplies bool) {
				assert.Equal(t, false, filterApplies, "vm should not pass the minCpuRatioFilter")
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.check(test.engine.minCpuRatioFilter(test.vm, test.req))

		})
	}
}

func TestEngine_minMemRatioFilter(t *testing.T) {
	tests := []struct {
		name   string
		engine Engine
		req    ClusterRecommendationReq
		vm     VirtualMachine
		attr   string
		check  func(filterApplies bool)
	}{
		{
			name:   "minMemRatioFilter applies",
			engine: Engine{},
			// minRatio = SumMem/SumCpu = 2
			req: ClusterRecommendationReq{SumMem: float64(8), SumCpu: 4},
			// ratio = Mem/Cpus = 4
			vm:   VirtualMachine{Mem: float64(16), Cpus: 4},
			attr: productinfo.Cpu,
			check: func(filterApplies bool) {
				assert.Equal(t, true, filterApplies, "vm should pass the minMemRatioFilter")
			},
		},
		{
			name:   "minMemRatioFilter doesn't apply",
			engine: Engine{},
			// minRatio = SumMem/SumCpu = 2
			req: ClusterRecommendationReq{SumMem: float64(8), SumCpu: 4},
			// ratio = Mem/Cpus = 0.5
			vm:   VirtualMachine{Cpus: 4, Mem: float64(4)},
			attr: productinfo.Cpu,
			check: func(filterApplies bool) {
				assert.Equal(t, false, filterApplies, "vm should not pass the minMemRatioFilter")
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.check(test.engine.minMemRatioFilter(test.vm, test.req))

		})
	}
}

func TestEngine_burstFilter(t *testing.T) {
	tests := []struct {
		name   string
		engine Engine
		req    ClusterRecommendationReq
		vm     VirtualMachine
		check  func(filterApplies bool)
	}{
		{
			name:   "burst filter applies - burst vm, burst allowed in req",
			engine: Engine{},
			req:    ClusterRecommendationReq{AllowBurst: &trueVal},
			vm:     VirtualMachine{Burst: true},
			check: func(filterApplies bool) {
				assert.Equal(t, true, filterApplies, "vm should pass the burst filter")
			},
		},
		{
			name:   "burst filter applies - burst vm, burst not set in req",
			engine: Engine{},
			// BurstAllowed not specified
			req: ClusterRecommendationReq{},
			vm:  VirtualMachine{Burst: true},
			check: func(filterApplies bool) {
				assert.Equal(t, true, filterApplies, "vm should pass the burst filter")
			},
		},
		{
			name:   "burst filter doesn't apply - burst vm, burst not allowed",
			engine: Engine{},
			req:    ClusterRecommendationReq{AllowBurst: &falseVal},
			vm:     VirtualMachine{Burst: true},
			check: func(filterApplies bool) {
				assert.Equal(t, false, filterApplies, "vm should not pass the burst filter")
			},
		},
		{
			name:   "burst filter applies - not burst vm, burst not allowed",
			engine: Engine{},
			req:    ClusterRecommendationReq{AllowBurst: &falseVal},
			// not a burst vm!
			vm: VirtualMachine{Burst: false},
			check: func(filterApplies bool) {
				assert.Equal(t, true, filterApplies, "vm should pass the burst filter")
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.check(test.engine.burstFilter(test.vm, test.req))
		})
	}
}

func TestEngine_findValuesBetween(t *testing.T) {
	tests := []struct {
		name   string
		engine Engine
		check  func([]float64, error)
		min    float64
		max    float64
		values []float64
	}{
		{
			name:   "failure - no attribute values to filer",
			engine: Engine{},
			values: nil,
			min:    0,
			max:    100,
			check: func(res []float64, err error) {
				assert.Equal(t, "no attribute values provided", err.Error())
				assert.Nil(t, res, "the result should be nil")
			},
		},
		{
			name:   "failure - min > max",
			engine: Engine{},
			values: []float64{0},
			min:    100,
			max:    0,
			check: func(res []float64, err error) {
				assert.Equal(t, "min value cannot be larger than the max value", err.Error())
				assert.Nil(t, res, "the result should be nil")
			},
		},
		{
			name:   "success - max is lower than the smallest value, return smallest value",
			engine: Engine{},
			values: []float64{200, 100, 500, 66},
			min:    0,
			max:    50,
			check: func(res []float64, err error) {
				assert.Nil(t, err, "should not get error")
				assert.Equal(t, 1, len(res), "returned invalid number of results")
				assert.Equal(t, float64(66), res[0], "invalid value returned")
			},
		},
		{
			name:   "success - min is greater than the largest value, return the largest value",
			engine: Engine{},
			values: []float64{200, 100, 500, 66},
			min:    1000,
			max:    2000,
			check: func(res []float64, err error) {
				assert.Nil(t, err, "should not get error")
				assert.Equal(t, 1, len(res), "returned invalid number of results")
				assert.Equal(t, float64(500), res[0], "invalid value returned")
			},
		},
		{
			name:   "success - should return values between min-max",
			engine: Engine{},
			values: []float64{1, 2, 10, 5, 30, 99, 55},
			min:    5,
			max:    55,
			check: func(res []float64, err error) {
				assert.Nil(t, err, "should not get error")
				assert.Equal(t, 4, len(res), "returned invalid number of results")
				assert.Equal(t, []float64{5, 10, 30, 55}, res, "invalid value returned")
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.check(test.engine.findValuesBetween(test.values, test.min, test.max))
		})
	}
}

func TestEngine_avgNodeCount(t *testing.T) {
	tests := []struct {
		name       string
		attrValues []float64
		reqSum     float64
		check      func(avg int)
	}{
		{
			name: "calculate average value per node",
			// sum of vals = 36
			// avg = 36/8 = 4.5
			// avg/node = ceil(10/4.5)=3
			attrValues: []float64{1, 2, 3, 4, 5, 6, 7, 8},
			reqSum:     10,
			check: func(avg int) {
				assert.Equal(t, 3, avg, "the calculated avg is invalid")
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.check(avgNodeCount(test.attrValues, test.reqSum))
		})
	}
}

func TestEngine_findCheapestNodePoolSet(t *testing.T) {
	tests := []struct {
		name      string
		engine    Engine
		nodePools map[string][]NodePool
		check     func(npo []NodePool)
	}{
		{
			name:   "find cheapest node pool set",
			engine: Engine{},
			nodePools: map[string][]NodePool{
				"memory": {
					NodePool{ // price = 2*3 +2*2 = 10
						VmType: VirtualMachine{
							AvgPrice:      2,
							OnDemandPrice: 3,
						},
						SumNodes: 2,
						VmClass:  regular,
					}, NodePool{
						VmType: VirtualMachine{
							AvgPrice:      2,
							OnDemandPrice: 3,
						},
						SumNodes: 2,
						VmClass:  spot,
					},
				},
				"cpu": { // price = 2*2 +2*2 = 8
					NodePool{
						VmType: VirtualMachine{
							AvgPrice:      1,
							OnDemandPrice: 2,
						},
						SumNodes: 2,
						VmClass:  regular,
					}, NodePool{
						VmType: VirtualMachine{
							AvgPrice:      2,
							OnDemandPrice: 4,
						},
						SumNodes: 2,
						VmClass:  spot,
					}, NodePool{
						VmType: VirtualMachine{
							AvgPrice:      2,
							OnDemandPrice: 4,
						},
						SumNodes: 0,
						VmClass:  spot,
					},
				},
			},
			check: func(npo []NodePool) {
				assert.Equal(t, 3, len(npo), "wrong selection")
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.check(test.engine.findCheapestNodePoolSet(test.nodePools))
		})
	}
}

func TestEngine_filterSpots(t *testing.T) {
	tests := []struct {
		name   string
		engine Engine
		vms    []VirtualMachine
		check  func(filtered []VirtualMachine)
	}{
		{
			name:   "vm-s filtered out",
			engine: Engine{},
			vms: []VirtualMachine{
				{
					AvgPrice:      0,
					OnDemandPrice: 1,
					Type:          "t100",
				},
				{
					AvgPrice:      2,
					OnDemandPrice: 5,
					Type:          "t200",
				},
			},
			check: func(filtered []VirtualMachine) {
				assert.Equal(t, 1, len(filtered), "vm is not filtered out")
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.check(test.engine.filterSpots(test.vms))
		})
	}
}

func TestEngine_ntwPerformanceFilter(t *testing.T) {
	tests := []struct {
		name   string
		engine Engine
		req    ClusterRecommendationReq
		vm     VirtualMachine
		check  func(passed bool)
	}{
		{
			name:   "vm passes the network performance filter",
			engine: Engine{},
			req: ClusterRecommendationReq{
				NertworkPerf: &productinfo.NTW_LOW,
			},
			vm: VirtualMachine{
				NetworkPerf: productinfo.NTW_LOW,
				Type:        "instance type",
			},
			check: func(passed bool) {
				assert.True(t, passed, "vm should pass the check")
			},
		},
		{
			name:   "vm doesn't pass the network performance filter",
			engine: Engine{},
			req: ClusterRecommendationReq{
				NertworkPerf: &productinfo.NTW_LOW,
			},
			vm: VirtualMachine{
				NetworkPerf: productinfo.NTW_HIGH,
				Type:        "instance type",
			},
			check: func(passed bool) {
				assert.False(t, passed, "vm should not pass the check")
			},
		},
		{
			name:   "vm passes the network performance filter - no filter in req",
			engine: Engine{},
			req:    ClusterRecommendationReq{ // filter is missing
			},
			vm: VirtualMachine{
				NetworkPerf: productinfo.NTW_LOW,
				Type:        "instance type",
			},
			check: func(passed bool) {
				assert.True(t, passed, "vm should pass the check")
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.check(test.engine.ntwPerformanceFilter(test.vm, test.req))
		})
	}
}
