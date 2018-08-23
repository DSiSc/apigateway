package query_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/DSiSc/apigateway/core/libs/pubsub"
	"github.com/DSiSc/apigateway/core/libs/pubsub/query"
)

func TestEmptyQueryMatchesAnything(t *testing.T) {
	q := query.Empty{}
	assert.True(t, q.Matches(pubsub.NewTagMap(map[string]string{})))
	assert.True(t, q.Matches(pubsub.NewTagMap(map[string]string{"Asher": "Roth"})))
	assert.True(t, q.Matches(pubsub.NewTagMap(map[string]string{"Route": "66"})))
	assert.True(t, q.Matches(pubsub.NewTagMap(map[string]string{"Route": "66", "Billy": "Blue"})))
}
