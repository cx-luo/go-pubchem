// Package src coding=utf-8
// @Project : go-pubchem
// @Time    : 2023/10/17 14:12
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@pharmaron.com
// @File    : getelem.go
// @Software: GoLand
package src

func (c Compound) GetCid() int {
	return c.Cid
}

func (c Compound) GetInChi() string {
	return c.Inchi
}

func (c Compound) GetInChiKey() string {
	return c.Inchikey
}

func (c Compound) GetIsoSmiles() string {
	return c.Isosmiles
}

func (c Compound) GetMw() float64 {
	return c.Mw
}

func (c Compound) GetMeshHeading() string {
	return c.Meshheadings
}

func (c Compound) GetCmpdName() string {
	return c.Cmpdname
}
