package geo

import (
	"time"
)

type Provider interface {
	ByAddress(address string) Location
	ByHostname(hostname string) Location
	ByIp(ip string) Location
}

type CachedProvider struct {
	store          Store
	dataUrl        string
	cacheDuration  time.Duration
	cacheTimestamp time.Time
}

func NewCachedProvider(geoDataUrl string, cacheDuration time.Duration) *CachedProvider {
	return &CachedProvider{
		cacheDuration: cacheDuration,
		dataUrl:       geoDataUrl,
		store:         nil,
	}
}

func (p *CachedProvider) validate() {
	if !p.isUpToDate() {
		p.update()
	}
}

func (p *CachedProvider) isUpToDate() bool {
	if p.store == nil {
		return true
	}

	age := time.Now().Sub(p.cacheTimestamp).Seconds()
	return age > p.cacheDuration.Seconds()
}

func (p *CachedProvider) update() {
	p.store = NewStoreFromUrl(p.dataUrl)
	p.cacheTimestamp = time.Now()
}

func (p *CachedProvider) ByAddress(address string) Location {
	p.validate()
	return p.store.ByIp(address)
}

func (p *CachedProvider) ByHostname(hostname string) Location {
	p.validate()
	return p.store.ByHostname(hostname)
}

func (p *CachedProvider) ByIp(ip string) Location {
	p.validate()
	return p.store.ByIp(ip)
}
