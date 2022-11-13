package bloom_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/jdk2588/bloom"
)

func TestRedisBitSet_New_Set_Test(t *testing.T) {

	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", "127.0.0.1:6379") },
	}
	conn := pool.Get()
	defer conn.Close()

	m, k := bloom.EstimateParameters(1000, .01)
	bitSet := bloom.NewRedisBitSet("test_key", m, conn, 1<<16)
	b := bloom.New(m, k, bitSet)
	b.Add([]byte("a"))
	b.Add([]byte("b"))
	exists, err := b.Exists([]byte("a"))
	fmt.Println(exists)
	if err != nil {
		t.Error("Could not test existence")
	}
	b.Remove([]byte("a"))
	exists1, _ := b.Exists([]byte("a"))
	fmt.Println(exists1)
	exists2, _ := b.Exists([]byte("b"))
	fmt.Println(exists2)
	//	isSetBefore, err := bitSet.Test([]uint{0})
	//	if err != nil {
	//		t.Error("Could not test bitset in redis")
	//	}
	//	if isSetBefore {
	//		t.Error("Bit should not be set")
	//	}
	//	err = bitSet.Set([]uint{512})
	//	if err != nil {
	//		t.Error("Could not set bitset in redis")
	//	}
	//	isSetAfter, err := bitSet.Test([]uint{512})
	//	if err != nil {
	//		t.Error("Could not test bitset in redis")
	//	}
	//	if !isSetAfter {
	//		t.Error("Bit should be set")
	//	}
	//	err = bitSet.Expire(3600)
	//	if err != nil {
	//		t.Errorf("Error adding expiration to bitset: %v", err)
	//	}
	//	err = bitSet.Delete()
	//	if err != nil {
	//		t.Errorf("Error cleaning up bitset: %v", err)
	//	}
}
