package store

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"github.com/ccbrown/keyvaluestore"
	"github.com/ccbrown/keyvaluestore/dynamodbstore"
	"github.com/ccbrown/keyvaluestore/memorystore"
	"github.com/ccbrown/keyvaluestore/redisstore"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

// Returned when a change conflicts with a simultaneous write, e.g. when you attempt to add a
// revision that already exists.
var ErrContention = fmt.Errorf("contention")

type Store struct {
	backend keyvaluestore.Backend
}

type resolverV2 struct {
	endpoint string
}

func (r *resolverV2) ResolveEndpoint(ctx context.Context, params dynamodb.EndpointParameters) (smithyendpoints.Endpoint, error) {
	if r.endpoint != "" {
		uri, err := url.Parse(r.endpoint)
		if err != nil {
			return smithyendpoints.Endpoint{}, err
		}
		return smithyendpoints.Endpoint{
			URI: *uri,
		}, nil
	}

	return dynamodb.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}

func New(ctx context.Context, cfg Config) (*Store, error) {
	if cfg.InMemory {
		return &Store{
			backend: memorystore.NewBackend(),
		}, nil
	}

	if dynamodbCfg := cfg.DynamoDB; dynamodbCfg != nil {
		awsConfig, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Fatalf("failed to load configuration, %v", err)
		}

		client := dynamodb.NewFromConfig(awsConfig, func(o *dynamodb.Options) {
			o.EndpointResolverV2 = &resolverV2{
				endpoint: dynamodbCfg.Endpoint,
			}
		})

		return &Store{
			backend: &dynamodbstore.Backend{
				Context:   ctx,
				Client:    client,
				TableName: dynamodbCfg.TableName,
			},
		}, nil
	} else if cfg.RedisAddress != "" {
		return &Store{
			backend: &redisstore.Backend{
				Client: redis.NewClient(&redis.Options{
					Addr: cfg.RedisAddress,
				}),
			},
		}, nil
	}

	return nil, fmt.Errorf("invalid store configuration")
}

func (s Store) Close() error {
	return nil
}

func (s *Store) getByIds(key string, dest interface{}, serializer Serializer, ids ...string) error {
	batch := s.backend.Batch()
	gets := make([]keyvaluestore.GetResult, 0, len(ids))
	keys := map[string]struct{}{}
	for _, id := range ids {
		key := key + ":" + id
		if _, ok := keys[key]; !ok {
			gets = append(gets, batch.Get(key))
			keys[key] = struct{}{}
		}
	}
	if err := batch.Exec(); err != nil {
		return err
	}

	objType := reflect.TypeOf(dest).Elem().Elem().Elem()
	slice := reflect.ValueOf(dest).Elem()
	for _, get := range gets {
		if v, _ := get.Result(); v != nil {
			obj := reflect.New(objType)
			if err := serializer.Deserialize(*v, obj.Interface()); err != nil {
				return err
			}
			slice = reflect.Append(slice, obj)
		}
	}
	reflect.ValueOf(dest).Elem().Set(slice)
	return nil
}

func (s *Store) getIdsByKeyAndTimeRange(key string, minTime, maxTime time.Time, limit int) ([]string, error) {
	zrange := s.backend.ZRangeByScore
	if limit < 0 {
		zrange = s.backend.ZRevRangeByScore
		limit = -limit
	}
	ids, err := zrange(key, timeMicrosecondScore(minTime), timeMicrosecondScore(maxTime), limit)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func execAtomicWrite(op keyvaluestore.AtomicWriteOperation) (bool, error) {
	ok, err := op.Exec()
	if err != nil && keyvaluestore.IsAtomicWriteConflict(err) {
		return false, ErrContention
	}
	return ok, err
}

func timeMicrosecondScore(t time.Time) float64 {
	return float64(t.Unix()*1000000 + int64(t.Nanosecond()/1000))
}

func stringsToUUIDs(s []string) []uuid.UUID {
	ret := make([]uuid.UUID, len(s))
	for i, id := range s {
		ret[i] = uuid.MustParse(id)
	}
	return ret
}
