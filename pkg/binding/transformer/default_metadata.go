package transformer

import (
	"time"

	"github.com/google/uuid"

	"github.com/cloudevents/sdk-go/pkg/binding"
	"github.com/cloudevents/sdk-go/pkg/binding/spec"
	"github.com/cloudevents/sdk-go/pkg/event"
	"github.com/cloudevents/sdk-go/pkg/types"
)

var (
	// Sets the cloudevents id attribute to a UUID.New()
	SetUUID binding.TransformerFactory = setUUID{}
	// Add the cloudevents time attribute, if missing, to time.Now()
	AddTimeNow binding.TransformerFactory = addTimeNow{}
)

type setUUID struct{}

func (a setUUID) StructuredTransformer(binding.StructuredWriter) binding.StructuredWriter {
	return nil
}

func (a setUUID) BinaryTransformer(encoder binding.BinaryWriter) binding.BinaryWriter {
	return &setUUIDTransformer{
		BinaryWriter: encoder,
	}
}

func (a setUUID) EventTransformer() binding.EventTransformer {
	return func(event *event.Event) error {
		return event.Context.SetID(uuid.New().String())
	}
}

type setUUIDTransformer struct {
	binding.BinaryWriter
}

func (b *setUUIDTransformer) SetAttribute(attribute spec.Attribute, value interface{}) error {
	if attribute.Kind() == spec.ID {
		return b.BinaryWriter.SetAttribute(attribute.Version().AttributeFromKind(spec.ID), uuid.New().String())
	}
	return b.BinaryWriter.SetAttribute(attribute, value)
}

type addTimeNow struct{}

func (a addTimeNow) StructuredTransformer(binding.StructuredWriter) binding.StructuredWriter {
	return nil
}

func (a addTimeNow) BinaryTransformer(encoder binding.BinaryWriter) binding.BinaryWriter {
	return &addTimeNowTransformer{
		BinaryWriter: encoder,
		found:        false,
	}
}

func (a addTimeNow) EventTransformer() binding.EventTransformer {
	return func(event *event.Event) error {
		if event.Context.GetTime().IsZero() {
			return event.Context.SetTime(time.Now())
		}
		return nil
	}
}

type addTimeNowTransformer struct {
	binding.BinaryWriter
	version spec.Version
	found   bool
}

func (b *addTimeNowTransformer) SetAttribute(attribute spec.Attribute, value interface{}) error {
	if attribute.Kind() == spec.Time {
		b.found = true
	}
	b.version = attribute.Version()
	return b.BinaryWriter.SetAttribute(attribute, value)
}

func (b *addTimeNowTransformer) End() error {
	if !b.found {
		err := b.BinaryWriter.SetAttribute(b.version.AttributeFromKind(spec.Time), types.Timestamp{Time: time.Now()})
		if err != nil {
			return err
		}
	}
	return b.BinaryWriter.End()
}