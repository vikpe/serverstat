package geo

import (
	"time"
)

type Provider interface {
	ByAddress(address string) Location
	ByHostname(hostname string) Location
	ByIp(ip string) Location
}

type MemcachedProvider struct {
	store          Store
	dataUrl        string
	cacheDuration  time.Duration
	cacheTimestamp time.Time
}

func NewMemcachedProvider(geoDataUrl string, cacheDuration time.Duration) *MemcachedProvider {
	return &MemcachedProvider{
		cacheDuration: cacheDuration,
		dataUrl:       geoDataUrl,
		store:         nil,
	}
}

func (p *MemcachedProvider) validate() {
	if !p.isUpToDate() {
		p.update()
	}
}

func (p *MemcachedProvider) isUpToDate() bool {
	if p.store == nil {
		return false
	}

	age := time.Now().Sub(p.cacheTimestamp).Seconds()
	return age < p.cacheDuration.Seconds()
}

func (p *MemcachedProvider) update() {
	p.store = NewStoreFromUrl(p.dataUrl)
	p.cacheTimestamp = time.Now()
}

func (p *MemcachedProvider) ByAddress(address string) Location {
	p.validate()
	return p.store.ByAddress(address)
}

func (p *MemcachedProvider) ByHostname(hostname string) Location {
	p.validate()
	return p.store.ByHostname(hostname)
}

func (p *MemcachedProvider) ByIp(ip string) Location {
	p.validate()
	return p.store.ByIp(ip)
}
