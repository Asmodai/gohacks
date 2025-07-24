// -*- Mode: Go -*-
//
// channel.go --- AMQP channel mockery.
//
// Copyright (c) 2025 Paul Ward <paul@lisphacker.uk>
//
// Author:     Paul Ward <paul@lisphacker.uk>
// Maintainer: Paul Ward <paul@lisphacker.uk>
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation files
// (the "Software"), to deal in the Software without restriction,
// including without limitation the rights to use, copy, modify, merge,
// publish, distribute, sublicense, and/or sell copies of the Software,
// and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
// BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
// ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// * Comments:
//
//

// * Package:

package amqpmock

// * Imports:

import (
	"context"

	"github.com/Asmodai/gohacks/mocks/amqpshim"
	goamqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/mock/gomock"
)

// * Code:

// ** Types:

type MockChannel struct {
	ack                                   AckFn
	cancel                                CancelFn
	close                                 SimpleErrorFn
	confirm                               ConfirmFn
	consume                               ConsumeFn
	consumeWithContext                    ConsumeContextFn
	exchangeBind                          ExchangeBindFn
	exchangeDeclare                       ExchangeDeclareFn
	exchangeDeclarePassive                ExchangeDeclareFn
	exchangeDelete                        ExchangeDeleteFn
	exchangeUnbind                        ExchangeUnbindFn
	flow                                  FlowFn
	get                                   GetFn
	getNextPublishSeqNo                   SimpleUInt64Fn
	isClosed                              SimpleBoolFn
	nack                                  NackFn
	notifyCancel                          NotifyCancelFn
	notifyClose                           NotifyCloseFn
	notifyConfirm                         NotifyConfirmFn
	notifyFlow                            NotifyFlowFn
	notifyPublish                         NotifyPublishFn
	notifyReturn                          NotifyReturnFn
	publish                               PublishFn
	publishWithContext                    PublishContextFn
	publishWithDeferredConfirm            PublishDeferredFn
	publishWithDeferredConfirmWithContext PublishDeferredContextFn
	qos                                   QosFn
	queueBind                             QueueBindFn
	queueDeclare                          QueueDeclareFn
	queueDeclarePassive                   QueueDeclareFn
	queueDelete                           QueueDeleteFn
	queuePurge                            QueuePurgeFn
	queueUnbind                           QueueUnbindFn
	reject                                RejectFn
	tx                                    SimpleErrorFn
	txCommit                              SimpleErrorFn
	txRollback                            SimpleErrorFn

	CallLog CallLog
}

// ** Methods:

// *** Initialisation:

func (obj *MockChannel) Init() {
	obj.BuildAck(ErrorResults{})
	obj.BuildCancel(ErrorResults{})
	obj.BuildClose(ErrorResults{})
	obj.BuildConfirm(ErrorResults{})
	obj.BuildConsume(ConsumeResults{})
	obj.BuildConsumeWithContext(ConsumeResults{})
	obj.BuildExchangeBind(ErrorResults{})
	obj.BuildExchangeDeclare(ErrorResults{})
	obj.BuildExchangeDeclarePassive(ErrorResults{})
	obj.BuildExchangeDelete(ErrorResults{})
	obj.BuildExchangeUnbind(ErrorResults{})
	obj.BuildFlow(ErrorResults{})
	obj.BuildGet(GetResults{})
	obj.BuildGetNextPublishSeqNo(UInt64Results{})
	obj.BuildIsClosed(BoolResults{})
	obj.BuildNack(ErrorResults{})
	obj.BuildNotifyCancel(NotifyCancelResults{})
	obj.BuildNotifyClose(NotifyCloseResults{})
	obj.BuildNotifyConfirm(NotifyConfirmResults{})
	obj.BuildNotifyFlow(NotifyFlowResults{})
	obj.BuildNotifyPublish(NotifyPublishResults{})
	obj.BuildNotifyReturn(NotifyReturnResults{})
	obj.BuildPublish(ErrorResults{})
	obj.BuildPublishWithContext(ErrorResults{})
	obj.BuildPublishWithDeferredConfirm(PublishDeferredResults{})
	obj.BuildPublishWithDeferredConfirmWithContext(PublishDeferredResults{})
	obj.BuildQos(ErrorResults{})
	obj.BuildQueueBind(ErrorResults{})
	obj.BuildQueueDeclare(QueueDeclareResults{})
	obj.BuildQueueDeclarePassive(QueueDeclareResults{})
	obj.BuildQueueDelete(QueueDeleteResults{})
	obj.BuildQueuePurge(QueueDeleteResults{})
	obj.BuildQueueUnbind(ErrorResults{})
	obj.BuildReject(ErrorResults{})
	obj.BuildTx(ErrorResults{})
	obj.BuildTxCommit(ErrorResults{})
	obj.BuildTxRollback(ErrorResults{})
}

func (obj *MockChannel) AddCallLog(fname string, values ...any) {
	if obj.CallLog == nil {
		obj.CallLog = make(CallLog)
	}

	obj.CallLog[fname] = append(obj.CallLog[fname], CallLogList{values})
}

// *** `Ack`:

func (obj *MockChannel) SetAck(fn AckFn) {
	obj.ack = fn
}

func (obj *MockChannel) BuildAck(results ErrorResults) {
	obj.SetAck(func(_ uint64, _ bool) error {
		return results.Error
	})
}

func (obj *MockChannel) MockAck(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		Ack(gomock.Any(), gomock.Any()).
		DoAndReturn(func(tag uint64, multiple bool) error {
			obj.AddCallLog("Ack", tag, multiple)

			return obj.ack(tag, multiple)
		})
}

// *** `Cancel`:

func (obj *MockChannel) SetCancel(fn CancelFn) {
	obj.cancel = fn
}

func (obj *MockChannel) BuildCancel(results ErrorResults) {
	obj.SetCancel(func(_ string, _ bool) error {
		return results.Error
	})
}

func (obj *MockChannel) MockCancel(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		Cancel(gomock.Any(), gomock.Any()).
		DoAndReturn(func(consumer string, noWait bool) error {
			obj.AddCallLog("Cancel", consumer, noWait)

			return obj.cancel(consumer, noWait)
		})
}

// *** `Close`:

func (obj *MockChannel) SetClose(fn SimpleErrorFn) {
	obj.close = fn
}

func (obj *MockChannel) BuildClose(results ErrorResults) {
	obj.SetClose(func() error {
		return results.Error
	})
}

func (obj *MockChannel) MockClose(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		Close().
		DoAndReturn(func() error {
			obj.AddCallLog("Close")

			return obj.close()
		})
}

// *** `Confirm`:

func (obj *MockChannel) SetConfirm(fn ConfirmFn) {
	obj.confirm = fn
}

func (obj *MockChannel) BuildConfirm(results ErrorResults) {
	obj.SetConfirm(func(_ bool) error {
		return results.Error
	})
}

func (obj *MockChannel) MockConfirm(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		Confirm(gomock.Any()).
		DoAndReturn(func(noWait bool) error {
			obj.AddCallLog("Confirm", noWait)

			return obj.confirm(noWait)
		})
}

// *** `Consume`:

func (obj *MockChannel) SetConsume(fn ConsumeFn) {
	obj.consume = fn
}

func (obj *MockChannel) BuildConsume(results ConsumeResults) {
	obj.SetConsume(
		func(
			_, _ string,
			_, _, _, _ bool,
			_ goamqp.Table,
		) (<-chan goamqp.Delivery, error) {
			return results.Channel, results.Error
		})
}

func (obj *MockChannel) MockConsume(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		Consume(
			gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(),
		).
		DoAndReturn(
			func(
				queue, consumer string,
				autoAck, exclusive, noLocal, noWait bool,
				args goamqp.Table,
			) (<-chan goamqp.Delivery, error) {
				obj.AddCallLog(
					"Consume",
					queue, consumer,
					autoAck, exclusive, noLocal, noWait,
					args,
				)

				return obj.consume(
					queue, consumer,
					autoAck, exclusive, noLocal, noWait,
					args,
				)
			})
}

// *** `ConsumeWithContext`:

func (obj *MockChannel) SetConsumeWithContext(fn ConsumeContextFn) {
	obj.consumeWithContext = fn
}

func (obj *MockChannel) BuildConsumeWithContext(results ConsumeResults) {
	obj.SetConsumeWithContext(
		func(
			_ context.Context,
			_, _ string,
			_, _, _, _ bool,
			_ goamqp.Table,
		) (<-chan goamqp.Delivery, error) {
			return results.Channel, results.Error
		})
}

func (obj *MockChannel) MockConsumeWithContext(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		ConsumeWithContext(
			gomock.Any(),
			gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(),
		).
		DoAndReturn(
			func(
				ctx context.Context,
				queue, consumer string,
				autoAck, exclusive, noLocal, noWait bool,
				args goamqp.Table,
			) (<-chan goamqp.Delivery, error) {
				obj.AddCallLog(
					"ConsumeWithContext",
					ctx,
					queue, consumer,
					autoAck, exclusive, noLocal, noWait,
					args,
				)

				return obj.consumeWithContext(
					ctx,
					queue, consumer,
					autoAck, exclusive, noLocal, noWait,
					args,
				)
			})
}

// *** `ExchangeBind`:

func (obj *MockChannel) SetExchangeBind(fn ExchangeBindFn) {
	obj.exchangeBind = fn
}

func (obj *MockChannel) BuildExchangeBind(results ErrorResults) {
	obj.SetExchangeBind(
		func(
			_, _, _ string,
			_ bool,
			_ goamqp.Table,
		) error {
			return results.Error
		})
}

func (obj *MockChannel) MockExchangeBind(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		ExchangeBind(
			gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).
		DoAndReturn(
			func(
				destination, key, source string,
				noWait bool,
				args goamqp.Table,
			) error {
				obj.AddCallLog(
					"ExchangeBind",
					destination, key, source,
					noWait,
					args,
				)

				return obj.exchangeBind(
					destination, key, source,
					noWait,
					args,
				)
			})
}

// *** `ExchangeDeclare`:

func (obj *MockChannel) SetExchangeDeclare(fn ExchangeDeclareFn) {
	obj.exchangeDeclare = fn
}

func (obj *MockChannel) BuildExchangeDeclare(results ErrorResults) {
	obj.SetExchangeDeclare(
		func(
			_, _ string,
			_, _, _, _ bool,
			_ goamqp.Table,
		) error {
			return results.Error
		})
}

func (obj *MockChannel) MockExchangeDeclare(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		ExchangeDeclare(
			gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(),
		).
		DoAndReturn(func(
			name, kind string,
			durable, autoDelete, internal, noWait bool,
			args goamqp.Table,
		) error {
			obj.AddCallLog(
				"ExchangeDeclare",
				name, kind,
				durable, autoDelete, internal, noWait,
				args,
			)

			return obj.exchangeDeclare(
				name, kind,
				durable, autoDelete, internal, noWait,
				args,
			)
		})
}

// *** `ExchangeDeclarePassive`:

func (obj *MockChannel) SetExchangeDeclarePassive(fn ExchangeDeclareFn) {
	obj.exchangeDeclarePassive = fn
}

func (obj *MockChannel) BuildExchangeDeclarePassive(results ErrorResults) {
	obj.SetExchangeDeclarePassive(
		func(
			_, _ string,
			_, _, _, _ bool,
			_ goamqp.Table,
		) error {
			return results.Error
		})
}

func (obj *MockChannel) MockExchangeDeclarePassive(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		ExchangeDeclarePassive(
			gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(),
		).
		DoAndReturn(func(
			name, kind string,
			durable, autoDelete, internal, noWait bool,
			args goamqp.Table,
		) error {
			obj.AddCallLog(
				"ExchangeDeclarePassive",
				name, kind,
				durable, autoDelete, internal, noWait,
				args,
			)

			return obj.exchangeDeclarePassive(
				name, kind,
				durable, autoDelete, internal, noWait,
				args,
			)
		})
}

// *** `ExchangeDelete`:

func (obj *MockChannel) SetExchangeDelete(fn ExchangeDeleteFn) {
	obj.exchangeDelete = fn
}

func (obj *MockChannel) BuildExchangeDelete(results ErrorResults) {
	obj.SetExchangeDelete(func(_ string, _, _ bool) error {
		return results.Error
	})
}

func (obj *MockChannel) MockExchangeDelete(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		ExchangeDelete(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(name string, ifUnused, noWait bool) error {
			obj.AddCallLog(
				"ExchangeDelete",
				name,
				ifUnused, noWait,
			)

			return obj.exchangeDelete(name, ifUnused, noWait)
		})
}

// *** `ExchangeUnbind`:

func (obj *MockChannel) SetExchangeUnbind(fn ExchangeUnbindFn) {
	obj.exchangeUnbind = fn
}

func (obj *MockChannel) BuildExchangeUnbind(results ErrorResults) {
	obj.SetExchangeUnbind(
		func(
			_, _, _ string,
			_ bool,
			_ goamqp.Table,
		) error {
			return results.Error
		})
}

func (obj *MockChannel) MockExchangeUnbind(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		ExchangeUnbind(
			gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).
		DoAndReturn(
			func(
				destination, key, source string,
				noWait bool,
				args goamqp.Table,
			) error {
				obj.AddCallLog(
					"ExchangeUnbind",
					destination, key, source,
					noWait,
					args,
				)

				return obj.exchangeUnbind(
					destination, key, source,
					noWait,
					args,
				)
			})
}

// *** `Flow`:

func (obj *MockChannel) SetFlow(fn FlowFn) {
	obj.flow = fn
}

func (obj *MockChannel) BuildFlow(results ErrorResults) {
	obj.SetFlow(func(_ bool) error {
		return results.Error
	})
}

func (obj *MockChannel) MockFlow(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		Flow(gomock.Any()).
		DoAndReturn(func(active bool) error {
			obj.AddCallLog("Flow", active)

			return obj.flow(active)
		})
}

// *** `Get`:

func (obj *MockChannel) SetGet(fn GetFn) {
	obj.get = fn
}

func (obj *MockChannel) BuildGet(results GetResults) {
	obj.SetGet(func(_ string, _ bool) (goamqp.Delivery, bool, error) {
		return results.Message, results.Ok, results.Error
	})
}

func (obj *MockChannel) MockGet(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		Get(gomock.Any(), gomock.Any()).
		DoAndReturn(func(
			queue string,
			autoAck bool,
		) (goamqp.Delivery, bool, error) {
			obj.AddCallLog("Get", queue, autoAck)

			return obj.get(queue, autoAck)
		})
}

// *** `GetNextPublishSeqNo`:

func (obj *MockChannel) SetGetNextPublishSeqNo(fn SimpleUInt64Fn) {
	obj.getNextPublishSeqNo = fn
}

func (obj *MockChannel) BuildGetNextPublishSeqNo(results UInt64Results) {
	obj.SetGetNextPublishSeqNo(func() uint64 {
		return results.Value
	})
}

func (obj *MockChannel) MockGetNextPublishSeqNo(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		GetNextPublishSeqNo().
		DoAndReturn(func() uint64 {
			obj.AddCallLog("GetNextPublishSeqNo")

			return obj.getNextPublishSeqNo()
		})
}

// *** `IsClosed`:

func (obj *MockChannel) SetIsClosed(fn SimpleBoolFn) {
	obj.isClosed = fn
}

func (obj *MockChannel) BuildIsClosed(results BoolResults) {
	obj.SetIsClosed(func() bool {
		return results.Value
	})
}

func (obj *MockChannel) MockIsClosed(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		IsClosed().
		DoAndReturn(func() bool {
			obj.AddCallLog("IsClosed")

			return obj.isClosed()
		})
}

// *** `Nack`:

func (obj *MockChannel) SetNack(fn NackFn) {
	obj.nack = fn
}

func (obj *MockChannel) BuildNack(results ErrorResults) {
	obj.SetNack(func(_ uint64, _, _ bool) error {
		return results.Error
	})
}

func (obj *MockChannel) MockNack(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		Nack(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(tag uint64, multiple, requeue bool) error {
			obj.AddCallLog("Nack", tag, multiple, requeue)

			return obj.nack(tag, multiple, requeue)
		})
}

// *** `NotifyCancel`:

func (obj *MockChannel) SetNotifyCancel(fn NotifyCancelFn) {
	obj.notifyCancel = fn
}

func (obj *MockChannel) BuildNotifyCancel(results NotifyCancelResults) {
	obj.SetNotifyCancel(func(_ chan string) chan string {
		return results.StringChan
	})
}

func (obj *MockChannel) MockNotifyCancel(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		NotifyCancel(gomock.Any()).
		DoAndReturn(func(c chan string) chan string {
			obj.AddCallLog("NotifyCancel", c)

			return obj.notifyCancel(c)
		})
}

// *** `NotifyClose`:

func (obj *MockChannel) SetNotifyClose(fn NotifyCloseFn) {
	obj.notifyClose = fn
}

func (obj *MockChannel) BuildNotifyClose(results NotifyCloseResults) {
	obj.SetNotifyClose(func(_ chan *goamqp.Error) chan *goamqp.Error {
		return results.ErrorChan
	})
}

func (obj *MockChannel) MockNotifyClose(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		NotifyClose(gomock.Any()).
		DoAndReturn(func(c chan *goamqp.Error) chan *goamqp.Error {
			obj.AddCallLog("NotifyClose", c)

			return obj.notifyClose(c)
		})
}

// *** `NotifyConfirm`:

func (obj *MockChannel) SetNotifyConfirm(fn NotifyConfirmFn) {
	obj.notifyConfirm = fn
}

func (obj *MockChannel) BuildNotifyConfirm(results NotifyConfirmResults) {
	obj.SetNotifyConfirm(func(_, _ chan uint64) (chan uint64, chan uint64) {
		return results.AckChan, results.NackChan
	})
}

func (obj *MockChannel) MockNotifyConfirm(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		NotifyConfirm(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ack, nack chan uint64) (chan uint64, chan uint64) {
			obj.AddCallLog("NotifyConfirm", ack, nack)

			return obj.notifyConfirm(ack, nack)
		})
}

// *** `NotifyFlow`:

func (obj *MockChannel) SetNotifyFlow(fn NotifyFlowFn) {
	obj.notifyFlow = fn
}

func (obj *MockChannel) BuildNotifyFlow(results NotifyFlowResults) {
	obj.SetNotifyFlow(func(_ chan bool) chan bool {
		return results.FlowChan
	})
}

func (obj *MockChannel) MockNotifyFlow(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		NotifyFlow(gomock.Any()).
		DoAndReturn(func(c chan bool) chan bool {
			obj.AddCallLog("NotifyFlow", c)

			return obj.notifyFlow(c)
		})
}

// *** `NotifyPublish`:

func (obj *MockChannel) SetNotifyPublish(fn NotifyPublishFn) {
	obj.notifyPublish = fn
}

func (obj *MockChannel) BuildNotifyPublish(results NotifyPublishResults) {
	obj.SetNotifyPublish(func(_ chan goamqp.Confirmation) chan goamqp.Confirmation {
		return results.ConfirmChan
	})
}

func (obj *MockChannel) MockNotifyPublish(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		NotifyPublish(gomock.Any()).
		DoAndReturn(func(c chan goamqp.Confirmation) chan goamqp.Confirmation {
			obj.AddCallLog("NotifyPublish", c)

			return obj.notifyPublish(c)
		})
}

// *** `NotifyReturn`:

func (obj *MockChannel) SetNotifyReturn(fn NotifyReturnFn) {
	obj.notifyReturn = fn
}

func (obj *MockChannel) BuildNotifyReturn(results NotifyReturnResults) {
	obj.SetNotifyReturn(func(_ chan goamqp.Return) chan goamqp.Return {
		return results.ReturnChan
	})
}

func (obj *MockChannel) MockNotifyReturn(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		NotifyReturn(gomock.Any()).
		DoAndReturn(func(c chan goamqp.Return) chan goamqp.Return {
			obj.AddCallLog("NotifyReturn", c)

			return obj.notifyReturn(c)
		})
}

// *** `Publish`:

func (obj *MockChannel) SetPublish(fn PublishFn) {
	obj.publish = fn
}

func (obj *MockChannel) BuildPublish(results ErrorResults) {
	obj.SetPublish(
		func(
			_, _ string,
			_, _ bool,
			_ goamqp.Publishing,
		) error {
			return results.Error
		})
}

func (obj *MockChannel) MockPublish(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		Publish(
			gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(),
			gomock.Any(),
		).
		DoAndReturn(
			func(
				exchange, key string,
				mandatory, immediate bool,
				msg goamqp.Publishing,
			) error {
				obj.AddCallLog(
					"Publish",
					exchange, key,
					mandatory, immediate,
					msg,
				)

				return obj.publish(
					exchange, key,
					mandatory, immediate,
					msg,
				)
			})
}

// *** `PublishWithContext`:

func (obj *MockChannel) SetPublishWithContext(fn PublishContextFn) {
	obj.publishWithContext = fn
}

func (obj *MockChannel) BuildPublishWithContext(results ErrorResults) {
	obj.SetPublishWithContext(
		func(
			_ context.Context,
			_, _ string,
			_, _ bool,
			_ goamqp.Publishing,
		) error {
			return results.Error
		})
}

func (obj *MockChannel) MockPublishWithContext(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		PublishWithContext(
			gomock.Any(),
			gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(),
			gomock.Any(),
		).
		DoAndReturn(
			func(
				ctx context.Context,
				exchange, key string,
				mandatory, immediate bool,
				msg goamqp.Publishing,
			) error {
				obj.AddCallLog(
					"PublishWithContext",
					ctx,
					exchange, key,
					mandatory, immediate,
					msg,
				)

				return obj.publishWithContext(
					ctx,
					exchange, key,
					mandatory, immediate,
					msg,
				)
			})
}

// *** `PublishWithDeferredConfirm`:

func (obj *MockChannel) SetPublishWithDeferredConfirm(fn PublishDeferredFn) {
	obj.publishWithDeferredConfirm = fn
}

func (obj *MockChannel) BuildPublishWithDeferredConfirm(results PublishDeferredResults) {
	obj.SetPublishWithDeferredConfirm(
		func(
			_, _ string,
			_, _ bool,
			_ goamqp.Publishing,
		) (*goamqp.DeferredConfirmation, error) {
			return results.Confirmation, results.Error
		})
}

func (obj *MockChannel) MockPublishWithDeferredConfirm(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		PublishWithDeferredConfirm(
			gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(),
			gomock.Any(),
		).
		DoAndReturn(
			func(
				exchange, key string,
				mandatory, immediate bool,
				msg goamqp.Publishing,
			) (*goamqp.DeferredConfirmation, error) {
				obj.AddCallLog(
					"PublishWithDeferredConfirm",
					exchange, key,
					mandatory, immediate,
					msg,
				)

				return obj.publishWithDeferredConfirm(
					exchange, key,
					mandatory, immediate,
					msg,
				)
			})
}

// *** `PublishWithDeferredConfirmWithContext`:

func (obj *MockChannel) SetPublishWithDeferredConfirmWithContext(fn PublishDeferredContextFn) {
	obj.publishWithDeferredConfirmWithContext = fn
}

func (obj *MockChannel) BuildPublishWithDeferredConfirmWithContext(results PublishDeferredResults) {
	obj.SetPublishWithDeferredConfirmWithContext(
		func(
			_ context.Context,
			_, _ string,
			_, _ bool,
			_ goamqp.Publishing,
		) (*goamqp.DeferredConfirmation, error) {
			return results.Confirmation, results.Error
		})
}

func (obj *MockChannel) MockPublishWithDeferredConfirmWithContext(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		PublishWithDeferredConfirmWithContext(
			gomock.Any(),
			gomock.Any(), gomock.Any(),
			gomock.Any(), gomock.Any(),
			gomock.Any(),
		).
		DoAndReturn(
			func(
				ctx context.Context,
				exchange, key string,
				mandatory, immediate bool,
				msg goamqp.Publishing,
			) (*goamqp.DeferredConfirmation, error) {
				obj.AddCallLog(
					"PublishWithDeferredConfirmWithContext",
					ctx,
					exchange, key,
					mandatory, immediate,
					msg,
				)

				return obj.publishWithDeferredConfirmWithContext(
					ctx,
					exchange, key,
					mandatory, immediate,
					msg,
				)
			})
}

// *** `Qos`:

func (obj *MockChannel) SetQos(fn QosFn) {
	obj.qos = fn
}

func (obj *MockChannel) BuildQos(results ErrorResults) {
	obj.SetQos(func(_, _ int, _ bool) error {
		return results.Error
	})
}

func (obj *MockChannel) MockQos(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		Qos(gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(prefetch, size int, global bool) error {
			obj.AddCallLog("Qos", prefetch, size, global)

			return obj.qos(prefetch, size, global)
		})
}

// *** `QueueBind`:

func (obj *MockChannel) SetQueueBind(fn QueueBindFn) {
	obj.queueBind = fn
}

func (obj *MockChannel) BuildQueueBind(results ErrorResults) {
	obj.SetQueueBind(
		func(
			_, _, _ string,
			_ bool,
			_ goamqp.Table,
		) error {
			return results.Error
		})
}

func (obj *MockChannel) MockQueueBind(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		QueueBind(
			gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).
		DoAndReturn(
			func(
				name, key, exchange string,
				noWait bool,
				args goamqp.Table,
			) error {
				obj.AddCallLog(
					"QueueBind",
					name, key, exchange,
					noWait,
					args,
				)

				return obj.queueBind(
					name, key, exchange,
					noWait,
					args,
				)
			})
}

// *** `QueueDeclare`:

func (obj *MockChannel) SetQueueDeclare(fn QueueDeclareFn) {
	obj.queueDeclare = fn
}

func (obj *MockChannel) BuildQueueDeclare(results QueueDeclareResults) {
	obj.SetQueueDeclare(func(
		_ string,
		_, _, _, _ bool,
		_ goamqp.Table,
	) (goamqp.Queue, error) {
		if results.Error != nil {
			return goamqp.Queue{}, results.Error
		}

		return results.Queue, nil
	})
}

func (obj *MockChannel) MockQueueDeclare(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		QueueDeclare(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).
		DoAndReturn(
			func(name string,
				durable, autoDelete, exclusive, noWait bool,
				args goamqp.Table,
			) (goamqp.Queue, error) {
				return obj.queueDeclare(
					name,
					durable,
					autoDelete,
					exclusive,
					noWait,
					args,
				)
			},
		)
}

// *** `QueueDeclarePassive`:

func (obj *MockChannel) SetQueueDeclarePassive(fn QueueDeclareFn) {
	obj.queueDeclarePassive = fn
}

func (obj *MockChannel) BuildQueueDeclarePassive(results QueueDeclareResults) {
	obj.SetQueueDeclarePassive(func(
		_ string,
		_, _, _, _ bool,
		_ goamqp.Table,
	) (goamqp.Queue, error) {
		if results.Error != nil {
			return goamqp.Queue{}, results.Error
		}

		return results.Queue, nil
	})
}

func (obj *MockChannel) MockQueueDeclarePassive(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		QueueDeclarePassive(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).
		DoAndReturn(
			func(name string,
				durable, autoDelete, exclusive, noWait bool,
				args goamqp.Table,
			) (goamqp.Queue, error) {
				return obj.queueDeclarePassive(
					name,
					durable,
					autoDelete,
					exclusive,
					noWait,
					args,
				)
			},
		)
}

// *** `QueueDelete`:

func (obj *MockChannel) SetQueueDelete(fn QueueDeleteFn) {
	obj.queueDelete = fn
}

func (obj *MockChannel) BuildQueueDelete(results QueueDeleteResults) {
	obj.SetQueueDelete(func(_ string, _, _, _ bool) (int, error) {
		return results.Purged, results.Error
	})
}

func (obj *MockChannel) MockQueueDelete(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		QueueDelete(
			gomock.Any(),
			gomock.Any(), gomock.Any(), gomock.Any(),
		).
		DoAndReturn(
			func(
				name string,
				ifUnused, ifEmpty, noWait bool,
			) (int, error) {
				obj.AddCallLog(
					"QueueDelete",
					name,
					ifUnused, ifEmpty, noWait,
				)

				return obj.queueDelete(
					name,
					ifUnused, ifEmpty, noWait,
				)
			})
}

// *** `QueuePurge`:

func (obj *MockChannel) SetQueuePurge(fn QueuePurgeFn) {
	obj.queuePurge = fn
}

func (obj *MockChannel) BuildQueuePurge(results QueueDeleteResults) {
	obj.SetQueuePurge(func(_ string, _ bool) (int, error) {
		return results.Purged, results.Error
	})
}

func (obj *MockChannel) MockQueuePurge(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		QueuePurge(gomock.Any(), gomock.Any()).
		DoAndReturn(func(name string, noWait bool) (int, error) {
			obj.AddCallLog("QueuePurge", name, noWait)

			return obj.queuePurge(name, noWait)
		})
}

// *** `QueueUnbind`:

func (obj *MockChannel) SetQueueUnbind(fn QueueUnbindFn) {
	obj.queueUnbind = fn
}

func (obj *MockChannel) BuildQueueUnbind(results ErrorResults) {
	obj.SetQueueUnbind(func(_, _, _ string, _ goamqp.Table) error {
		return results.Error
	})
}

func (obj *MockChannel) MockQueueUnbind(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		QueueUnbind(
			gomock.Any(), gomock.Any(), gomock.Any(),
			gomock.Any(),
		).
		DoAndReturn(
			func(
				name, key, exchange string,
				args goamqp.Table,
			) error {
				obj.AddCallLog(
					"QueueUnbind",
					name, key, exchange,
					args,
				)

				return obj.queueUnbind(
					name, key, exchange,
					args,
				)
			})
}

// *** `Reject`:

func (obj *MockChannel) SetReject(fn RejectFn) {
	obj.reject = fn
}

func (obj *MockChannel) BuildReject(results ErrorResults) {
	obj.SetReject(func(_ uint64, _ bool) error {
		return results.Error
	})
}

func (obj *MockChannel) MockReject(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		Reject(gomock.Any(), gomock.Any()).
		DoAndReturn(func(tag uint64, requeue bool) error {
			obj.AddCallLog("Reject", tag, requeue)

			return obj.reject(tag, requeue)
		})
}

// *** `Tx`:

func (obj *MockChannel) SetTx(fn SimpleErrorFn) {
	obj.tx = fn
}

func (obj *MockChannel) BuildTx(results ErrorResults) {
	obj.SetTx(func() error {
		return results.Error
	})
}

func (obj *MockChannel) MockTx(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		Tx().
		DoAndReturn(func() error {
			obj.AddCallLog("Tx")

			return obj.tx()
		})
}

// *** `TxCommit`:

func (obj *MockChannel) SetTxCommit(fn SimpleErrorFn) {
	obj.txCommit = fn
}

func (obj *MockChannel) BuildTxCommit(results ErrorResults) {
	obj.SetTxCommit(func() error {
		return results.Error
	})
}

func (obj *MockChannel) MockTxCommit(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		TxCommit().
		DoAndReturn(func() error {
			obj.AddCallLog("TxCommit")

			return obj.txCommit()
		})
}

// *** `TxRollback`:

func (obj *MockChannel) SetTxRollback(fn SimpleErrorFn) {
	obj.txRollback = fn
}

func (obj *MockChannel) BuildTxRollback(results ErrorResults) {
	obj.SetTxRollback(func() error {
		return results.Error
	})
}

func (obj *MockChannel) MockTxRollback(mock *amqpshim.MockChannel) *gomock.Call {
	return mock.EXPECT().
		TxRollback().
		DoAndReturn(func() error {
			obj.AddCallLog("TxRollback")

			return obj.txRollback()
		})
}

// * channel.go ends here.
