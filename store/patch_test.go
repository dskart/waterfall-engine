package store

import (
	"testing"
	"time"

	"github.com/dskart/netword/model"
	appErrors "github.com/dskart/netword/pkg/errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestModelType struct {
	model.Model

	Foo string
	Bar int
}

func NewTestModelType(revisionUserId uuid.UUID, foo string, bar int) TestModelType {
	return TestModelType{
		Model: model.NewModel(revisionUserId),
		Foo:   foo,
		Bar:   bar,
	}
}

type TestModelTypePatch struct {
	Foo *string
	Bar *int
}

func (t TestModelType) WithPatch(revisionUserId uuid.UUID, patch TestModelTypePatch) *TestModelType {
	t.RevisionNumber++
	t.RevisionTime = time.Now()
	t.RevisionUserId = revisionUserId
	if patch.Foo != nil {
		t.Foo = *patch.Foo
	}
	if patch.Bar != nil {
		t.Bar = *patch.Bar
	}
	return &t
}

func toPointer[T any](t T) *T {
	return &t
}

func getValue(id uuid.UUID) (*TestModelType, appErrors.SanitizedError) {
	v, ok := memStore[id]
	if !ok {
		return nil, appErrors.NewResourceNotFoundError()
	}
	return &v, nil
}

func validateValue(v *TestModelType) appErrors.SanitizedError {
	return nil
}

func commitValue(current *TestModelType, previons *TestModelType) error {
	memStore[current.Id] = *current
	return nil
}

var memStore = make(map[uuid.UUID]TestModelType)

func TestPatch(t *testing.T) {
	userId := uuid.New()
	myModel := NewTestModelType(userId, "foo", 1)
	id := myModel.Id
	commitValue(&myModel, nil)

	patch := TestModelTypePatch{
		Foo: toPointer("foo2"),
		Bar: toPointer(2),
	}

	revisionUserId := uuid.New()
	new, prev, err := Patch(revisionUserId, id, patch, getValue, commitValue, WithValidate(validateValue))
	require.NoError(t, err)
	assert.Equal(t, "foo2", new.Foo)
	assert.Equal(t, 2, new.Bar)
	assert.Equal(t, revisionUserId, new.RevisionUserId)
	assert.Equal(t, 2, new.RevisionNumber)
	assert.Equal(t, "foo", prev.Foo)
	assert.Equal(t, 1, prev.Bar)
	assert.Equal(t, 1, prev.RevisionNumber)
	assert.True(t, prev.RevisionTime.Before(new.RevisionTime))
}

func TestPatchContention(t *testing.T) {
	userId := uuid.New()
	myModel := NewTestModelType(userId, "foo", 1)
	id := myModel.Id
	commitValue(&myModel, nil)

	patch := TestModelTypePatch{
		Foo: toPointer("foo2"),
		Bar: toPointer(2),
	}

	commitValue := func(current *TestModelType, prev *TestModelType) error {
		return ErrContention
	}

	_, _, err := Patch(userId, id, patch, getValue, commitValue, WithValidate(validateValue))
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrContention)
}
