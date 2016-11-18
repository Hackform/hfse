package service

import (
  "github.com/Hackform/hfse/kappa"
)

type (
  Service interface {
    SetId(kappa.Const) kappa.Const
    GetId() kappa.Const
  }
)
