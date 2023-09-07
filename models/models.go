package models

import "gorm.io/gorm"

type Location struct {
	gorm.Model

	Cep            string `json:"cep"`
	TipoLogradouro string `json:"tipo_logradouro"`
	Logradouro     string `json:"logradouro"`
	Cidade         string `json:"cidade"`
	Uf             string `json:"uf"`
	CodigoIbge     string `json:"codigo_ibge"`
	Bairro         string `json:"bairro"`
}

type State struct {
	gorm.Model
	Name string `json:"name"`
}
