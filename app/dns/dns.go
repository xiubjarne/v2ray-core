package dns

import (
	"net"
	"sync"
	"time"

	"github.com/v2ray/v2ray-core/app"
)

type entry struct {
	domain     string
	ip         net.IP
	validUntil time.Time
}

func newEntry(domain string, ip net.IP) *entry {
	this := &entry{
		domain: domain,
		ip:     ip,
	}
	this.Extend()
	return this
}

func (this *entry) IsValid() bool {
	return this.validUntil.After(time.Now())
}

func (this *entry) Extend() {
	this.validUntil = time.Now().Add(time.Hour)
}

type DnsCache struct {
	sync.RWMutex
	cache  map[string]*entry
	config CacheConfig
}

func NewCache(config CacheConfig) *DnsCache {
	cache := &DnsCache{
		cache:  make(map[string]*entry),
		config: config,
	}
	go cache.cleanup()
	return cache
}

func (this *DnsCache) cleanup() {
	for range time.Tick(60 * time.Second) {
		entry2Remove := make([]*entry, 0, 128)
		this.RLock()
		for _, entry := range this.cache {
			if !entry.IsValid() {
				entry2Remove = append(entry2Remove, entry)
			}
		}
		this.RUnlock()

		for _, entry := range entry2Remove {
			if !entry.IsValid() {
				this.Lock()
				delete(this.cache, entry.domain)
				this.Unlock()
			}
		}
	}
}

func (this *DnsCache) Add(context app.Context, domain string, ip net.IP) {
	callerTag := context.CallerTag()
	if !this.config.IsTrustedSource(callerTag) {
		return
	}

	this.RLock()
	entry, found := this.cache[domain]
	this.RUnlock()
	if found {
		entry.ip = ip
		entry.Extend()
	} else {
		this.Lock()
		this.cache[domain] = newEntry(domain, ip)
		this.Unlock()
	}
}

func (this *DnsCache) Get(context app.Context, domain string) net.IP {
	this.RLock()
	entry, found := this.cache[domain]
	this.RUnlock()
	if found {
		return entry.ip
	}
	return nil
}