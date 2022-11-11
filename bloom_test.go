package bloom_test

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/garyburd/redigo/redis"
	"github.com/jdk2588/bloom"
)

func TestRedisBloomFilter(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Error("Miniredis could not start")
	}
	defer s.Close()

	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", s.Addr()) },
	}
	conn := pool.Get()
	defer conn.Close()

	m, k := bloom.EstimateParameters(1000, .01)
	bitSet := bloom.NewRedisBitSet("test_key", m, conn)
	b := bloom.New(m, k, bitSet)
	testBloomFilter(t, b)
}

func testBloomFilter(t *testing.T, b *bloom.BloomFilter) {
	data := []byte("some key")
	existsBefore, err := b.Exists(data)
	if err != nil {
		t.Fatal("Error checking for existence in bloom filter")
	}
	if existsBefore {
		t.Fatal("Bloom filter should not contain this data")
	}
	err = b.Add(data)
	if err != nil {
		t.Fatal("Error adding to bloom filter")
	}
	existsAfter, err := b.Exists(data)
	if err != nil {
		t.Fatal("Error checking for existence in bloom filter")
	}
	if !existsAfter {
		t.Fatal("Bloom filter should contain this data")
	}
	err = b.Remove(data)
	if err != nil {
		t.Fatal("Error removing from bloom filter")
	}
	exists, err := b.Exists(data)
	if err != nil {
		t.Fatal("Error checking for existence in bloom filter")
	}
	if exists {
		t.Fatal("Bloom filter should not contain this data")
	}
}
