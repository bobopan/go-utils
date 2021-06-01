package baseutils

import (
	"errors"
	"strings"
	"time"

	"gopkg.in/redis.v5"
)

type RedisConfig interface {
	GetPassword() string
	GetAddr() string
	GetPoolNum() int
	GetReadTimeout() time.Duration
	GetWriteTimeout() time.Duration
	GetPoolTimeout() time.Duration
	GetDialTimeout() time.Duration
}

type Redis struct {
	client redis.Cmdable
	Err error
}

func NewRedis(redisConfig RedisConfig) *Redis {
	client, err := initRedis(redisConfig)
	return &Redis{client: client, Err: err}
}

func initRedisNormal(redisConfig RedisConfig) (redis.Cmdable, error) {
	if len(redisConfig.GetAddr()) == 0 {
		return nil, errors.New("null redis addr")
	}
	client := redis.NewClient(&redis.Options{
		Addr:         redisConfig.GetAddr(),
		PoolSize:     redisConfig.GetPoolNum(),
		ReadTimeout:  redisConfig.GetReadTimeout(),
		WriteTimeout: redisConfig.GetWriteTimeout(),
		PoolTimeout:  redisConfig.GetPoolTimeout(),
		DialTimeout:  redisConfig.GetDialTimeout(),
		Password:     redisConfig.GetPassword(),
	})
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, err
}

func initRedisCluster(redisConfig RedisConfig) (redis.Cmdable, error) {
	if len(redisConfig.GetAddr()) == 0 {
		return nil, errors.New("链接地址不能为空")
	}
	addrSegs := strings.Split(redisConfig.GetAddr(), ",")
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        addrSegs,
		PoolSize:     redisConfig.GetPoolNum(),
		ReadTimeout:  redisConfig.GetReadTimeout(),
		WriteTimeout: redisConfig.GetWriteTimeout(),
		PoolTimeout:  redisConfig.GetPoolTimeout(),
		DialTimeout:  redisConfig.GetDialTimeout(),
		Password:     redisConfig.GetPassword(),
	})
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, err
}

func initRedis(redisConfig RedisConfig) (redis.Cmdable, error) {
	var client redis.Cmdable
	var err error
	if strings.Contains(redisConfig.GetAddr(), ",") {
		client, err = initRedisCluster(redisConfig)
	} else {
		client, err = initRedisNormal(redisConfig)
	}
	return client, err
}

//get string key just for freq condition
func (r *Redis) Get(key string) (string, error) {
	if r.client == nil {
		return "", errors.New("not inited")
	}
	cmd := r.client.Get(key)
	err := cmd.Err()
	if err == redis.Nil {
		//Log.Debug("key:(%s) is not exists", key)
		return "", nil
	} else if err != nil {
		return "", err
	}
	value := cmd.Val()
	return value, err
}

func (r *Redis) GetByte(key string) ([]byte, error) {
	value := make([]byte, 0)
	if r.client == nil {
		return value, errors.New("not inited")
	}
	cmd := r.client.Get(key)
	err := cmd.Err()
	if err == redis.Nil {
		//Log.Debug("key:(%s) is not exists", key)
		return value, nil
	} else if err != nil {
		return value, err
	}
	value, err = cmd.Bytes()
	return value, err
}

func (r *Redis) Set(key string, value interface{}, expiration time.Duration) error {
	if r.client == nil {
		return errors.New("not inited")
	}
	cmd := r.client.Set(key, value, expiration)
	return cmd.Err()
}

func (r *Redis) SetNx(key string, value interface{}, expiration time.Duration) (bool, error) {
	if r.client == nil {
		return false, errors.New("not inited")
	}
	cmd := r.client.SetNX(key, value, expiration)
	return cmd.Result()
}

func (r *Redis) Mget(keys ...string) ([]interface{}, error) {
	if r.client == nil {
		return nil, errors.New("not inited")
	}
	result := r.client.MGet(keys...)
	return result.Result()
}

func (r *Redis) HMSet(key string, fields map[string]string) error {
	if r.client == nil {
		return errors.New("not initied")
	}
	cmd := r.client.HMSet(key, fields)
	return cmd.Err()
}

func (r *Redis) HGet(key, field string) (string, error) {
	if r.client == nil {
		return "", errors.New("not initied")
	}
	cmd := r.client.HGet(key, field)
	err := cmd.Err()
	if err == redis.Nil {
		//Log.Debug("key:(%s) is not exists", key)
		return "", nil
	} else if err != nil {
		return "", err
	}
	value := cmd.Val()
	return value, err
}

func (r *Redis) HDel(key, field string) error {
	if r.client == nil {
		return errors.New("not initied")
	}
	return r.client.HDel(key, field).Err()
}

func (r *Redis) HGetAllMap(key string) (map[string]string, error) {
	if r.client == nil {
		return nil, errors.New("not inited")
	}
	cmd := r.client.HGetAll(key)
	err := cmd.Err()
	if err == redis.Nil {
		//Log.Debug("key:(%s) is not exists", key)
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	value := cmd.Val()
	return value, nil
}

func (r *Redis) HSet(key, feild, value string) error {
	if r.client == nil {
		return  errors.New("not inited")
	}
	cmd := r.client.HSet(key, feild, value)
	return cmd.Err()
}

func (r *Redis) SetExpire(key string, expiration time.Duration) error {
	if r.client == nil {
		return  errors.New("not inited")
	}
	cmd := r.client.Expire(key, expiration)
	return cmd.Err()
}

func (r *Redis) Exists(key string) (bool, error) {
	if r.client == nil {
		return false, errors.New("not inited")
	}
	cmd := r.client.Exists(key)
	return cmd.Val(), cmd.Err()
}

//IncrBy(key string, value int64) *IntCmd
func (r *Redis) Incby(key string, value int64) (int64, error) {
	if r.client == nil {
		return 0, errors.New("not initied")
	}
	return r.client.IncrBy(key, value).Result()
}

func (r *Redis) TTL(key string) (time.Duration, error) {
	if r.client == nil {
		return time.Nanosecond, errors.New("not initied")
	}
	return r.client.TTL(key).Result()
}
func (r *Redis) Del(key ...string) error {
	if r.client == nil {
		return errors.New("not initied")
	}
	return r.client.Del(key...).Err()
}

func (r *Redis) Keys(pattern string) ([]string, error) {
	if r.client == nil {
		return nil, errors.New("not inited")
	}
	stringSliceCmd := r.client.Keys(pattern)
	err := stringSliceCmd.Err()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return stringSliceCmd.Val(), nil
}

//Publish 由于publish 不是基本命令，需要特殊实现
func (r *Redis) Publish(channel, message string) error {
	if v, ok := (r.client).(*redis.Client); ok {
		return v.Publish(channel, message).Err()
	}
	if v, ok := (r.client).(*redis.ClusterClient); ok {
		return v.Publish(channel, message).Err()
	}
	return redis.NewIntCmd("bad instance").Err()
}

func (r *Redis) PSubscribe(channels ...string) (*redis.PubSub, error) {
	if v, ok := (r.client).(*redis.Client); ok {
		return v.PSubscribe(channels...)
	}
	return nil, redis.NewIntCmd("bad instance").Err()
}

func (r *Redis) Subscribe(channel string) (*redis.PubSub, error) {
	if v, ok := (r.client).(*redis.Client); ok {
		return v.Subscribe(channel)
	}
	return nil, redis.NewIntCmd("bad instance").Err()
}

func (r *Redis) Client() redis.Cmdable {
	return r.client
}
