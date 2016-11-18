package kappa

import (
  "sync/atomic"
)

type (
  Kappa struct {
    value uint32
  }

  Const uint32
)

func New() *Kappa {
  return &Kappa{
    value: 0,
  }
}

func (k *Kappa) Get() Const {
  return Const(atomic.AddUint32(&k.value, 1))
}
