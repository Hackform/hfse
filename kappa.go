package hfse

import (
  "sync/atomic"
)

type (
  Kappa struct {
    value uint32
  }
)

func NewKappa() *Kappa {
  return &Kappa{
    value: 0,
  }
}

func (k *Kappa) Get() uint32 {
  return atomic.AddUint32(&k.value, 1)
}
