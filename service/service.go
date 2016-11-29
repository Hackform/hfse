package service

import (
	"github.com/Hackform/hfse/kappa"
)

type (
	ServiceSubstrate struct {
		serviceKappa *kappa.Kappa
		services     map[kappa.Const]Service
	}

	Service interface {
		SetId(k kappa.Const) kappa.Const
		GetId() kappa.Const
	}
)

func New() *ServiceSubstrate {
	return &ServiceSubstrate{
		serviceKappa: kappa.New(),
		services:     make(map[kappa.Const]Service),
	}
}

func (s *ServiceSubstrate) Set(ser Service) kappa.Const {
	k := s.serviceKappa.Get()
	ser.SetId(k)
	s.services[k] = ser
	return k
}

func (s *ServiceSubstrate) Get(k kappa.Const) Service {
	return s.services[k]
}
