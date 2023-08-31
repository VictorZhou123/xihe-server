package bigmodels

import "github.com/opensourceways/xihe-server/bigmodel/domain"

type glmInfo struct {
	cfg GLM
	endpoints chan string
}

func newGLMInfo(cfg *Config) (glmInfo, error) {
	v := &cfg.GLM

	info := glmInfo{
		cfg: *v,
	}

	ce := &cfg.Endpoints
	es, _ := ce.parse(ce.GLM)

	// init endpoints
	info.endpoints = make(chan string, len(es))
	for _, e := range es {
		info.endpoints <- e
	}

	return info, nil
}

func (s *service) GLM(input domain.GLMInput) (resp string, err error) {
	// input check

	// output check
	return 
}