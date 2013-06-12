// Copyright 2013 Tumblr, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package scribe contains Apache Thrift protocol definitions for Apache Scribe
package scribe

import (
	"fmt"
	"tumblr/encoding/thrift"
)

import "tumblr/net/scribe/thrift/fb303"

type IScribe interface {
	fb303.IFacebookService

	/**
	 * Parameters:
	 *  - Messages
	 */
	Log(messages thrift.TList) (retval4 ResultCode, err error)
}

type ScribeClient struct {
	*fb303.FacebookServiceClient
}

func NewScribeClientFactory(t thrift.TTransport, f thrift.TProtocolFactory) *ScribeClient {
	return &ScribeClient{FacebookServiceClient: fb303.NewFacebookServiceClientFactory(t, f)}
}

func NewScribeClientProtocol(t thrift.TTransport, iprot thrift.TProtocol, oprot thrift.TProtocol) *ScribeClient {
	return &ScribeClient{FacebookServiceClient: fb303.NewFacebookServiceClientProtocol(t, iprot, oprot)}
}

/**
 * Parameters:
 *  - Messages
 */
func (p *ScribeClient) Log(messages thrift.TList) (retval5 ResultCode, err error) {
	err = p.SendLog(messages)
	if err != nil {
		return
	}
	return p.RecvLog()
}

func (p *ScribeClient) SendLog(messages thrift.TList) (err error) {
	oprot := p.OutputProtocol
	if oprot != nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	oprot.WriteMessageBegin("Log", thrift.CALL, p.SeqId)
	args6 := NewLogArgs()
	args6.Messages = messages
	err = args6.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Transport().Flush()
	return
}

func (p *ScribeClient) RecvLog() (value ResultCode, err error) {
	iprot := p.InputProtocol
	if iprot == nil {
		iprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.InputProtocol = iprot
	}
	_, mTypeId, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return
	}
	if mTypeId == thrift.EXCEPTION {
		error8 := thrift.NewTApplicationExceptionDefault()
		var error9 error
		error9, err = error8.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error9
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "ping failed: out of sequence response")
		return
	}
	result7 := NewLogResult()
	err = result7.Read(iprot)
	iprot.ReadMessageEnd()
	value = result7.Success
	return
}

type ScribeProcessor struct {
	super *fb303.FacebookServiceProcessor
}

func (p *ScribeProcessor) Handler() IScribe {
	return p.super.Handler().(IScribe)
}

func (p *ScribeProcessor) AddToProcessorMap(key string, processor thrift.TProcessorFunction) {
	p.super.AddToProcessorMap(key, processor)
}

func (p *ScribeProcessor) GetProcessorFunction(key string) (processor thrift.TProcessorFunction, exists bool) {
	return p.super.GetProcessorFunction(key)
}

func (p *ScribeProcessor) ProcessorMap() map[string]thrift.TProcessorFunction {
	return p.super.ProcessorMap()
}

func NewScribeProcessor(handler IScribe) *ScribeProcessor {
	self10 := &ScribeProcessor{super: fb303.NewFacebookServiceProcessor(handler)}
	self10.AddToProcessorMap("Log", &scribeProcessorLog{handler: handler})
	return self10
}

func (p *ScribeProcessor) Process(iprot, oprot thrift.TProtocol) (bool, thrift.TException) {
	return p.super.Process(iprot, oprot)
}

type scribeProcessorLog struct {
	handler IScribe
}

func (p *scribeProcessorLog) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := NewLogArgs()
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("Log", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Transport().Flush()
		return
	}
	iprot.ReadMessageEnd()
	result := NewLogResult()
	if result.Success, err = p.handler.Log(args.Messages); err != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing Log: "+err.Error())
		oprot.WriteMessageBegin("Log", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Transport().Flush()
		return
	}
	if err2 := oprot.WriteMessageBegin("Log", thrift.REPLY, seqId); err2 != nil {
		err = err2
	}
	if err2 := result.Write(oprot); err == nil && err2 != nil {
		err = err2
	}
	if err2 := oprot.WriteMessageEnd(); err == nil && err2 != nil {
		err = err2
	}
	if err2 := oprot.Transport().Flush(); err == nil && err2 != nil {
		err = err2
	}
	if err != nil {
		return
	}
	return true, err
}

// HELPER FUNCTIONS AND STRUCTURES

/**
 * Attributes:
 *  - Messages
 */
type LogArgs struct {
	thrift.TStruct
	Messages thrift.TList "messages" // 1
}

func NewLogArgs() *LogArgs {
	output := &LogArgs{
		TStruct: thrift.NewTStruct("Log_args", []thrift.TField{
			thrift.NewTField("messages", thrift.LIST, 1),
		}),
	}
	{
	}
	return output
}

func (p *LogArgs) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 1 || fieldName == "messages" {
			if fieldTypeId == thrift.LIST {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField1(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *LogArgs) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_etype16, _size13, err := iprot.ReadListBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "p.Messages", "", err)
	}
	p.Messages = thrift.NewTList(_etype16, _size13)
	for _i17 := 0; _i17 < _size13; _i17++ {
		_elem18 := NewLogEntry()
		err21 := _elem18.Read(iprot)
		if err21 != nil {
			return thrift.NewTProtocolExceptionReadStruct("_elem18LogEntry", err21)
		}
		p.Messages.Push(_elem18)
	}
	err = iprot.ReadListEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "", "list", err)
	}
	return err
}

func (p *LogArgs) ReadFieldMessages(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *LogArgs) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("Log_args")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField1(oprot)
	if err != nil {
		return err
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *LogArgs) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.Messages != nil {
		err = oprot.WriteFieldBegin("messages", thrift.LIST, 1)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "messages", p.ThriftName(), err)
		}
		err = oprot.WriteListBegin(thrift.STRUCT, p.Messages.Len())
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(-1, "", "list", err)
		}
		for Iter22 := range p.Messages.Iter() {
			Iter23 := Iter22.(*LogEntry)
			err = Iter23.Write(oprot)
			if err != nil {
				return thrift.NewTProtocolExceptionWriteStruct("LogEntry", err)
			}
		}
		err = oprot.WriteListEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(-1, "", "list", err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(1, "messages", p.ThriftName(), err)
		}
	}
	return err
}

func (p *LogArgs) WriteFieldMessages(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *LogArgs) TStructName() string {
	return "LogArgs"
}

func (p *LogArgs) ThriftName() string {
	return "Log_args"
}

func (p *LogArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("LogArgs(%+v)", *p)
}

func (p *LogArgs) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*LogArgs)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *LogArgs) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.Messages
	}
	return nil
}

func (p *LogArgs) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("messages", thrift.LIST, 1),
	})
}

/**
 * Attributes:
 *  - Success
 */
type LogResult struct {
	thrift.TStruct
	Success ResultCode "success" // 0
}

func NewLogResult() *LogResult {
	output := &LogResult{
		TStruct: thrift.NewTStruct("Log_result", []thrift.TField{
			thrift.NewTField("success", thrift.I32, 0),
		}),
	}
	{
	}
	return output
}

func (p *LogResult) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_, err = iprot.ReadStructBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	for {
		fieldName, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if fieldId < 0 {
			fieldId = int16(p.FieldIdFromFieldName(fieldName))
		} else if fieldName == "" {
			fieldName = p.FieldNameFromFieldId(int(fieldId))
		}
		if fieldTypeId == thrift.GENERIC {
			fieldTypeId = p.FieldFromFieldId(int(fieldId)).TypeId()
		}
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		if fieldId == 0 || fieldName == "success" {
			if fieldTypeId == thrift.I32 {
				err = p.ReadField0(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField0(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			}
		} else {
			err = iprot.Skip(fieldTypeId)
			if err != nil {
				return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
			}
		}
		err = iprot.ReadFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
		}
	}
	err = iprot.ReadStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadStruct(p.ThriftName(), err)
	}
	return err
}

func (p *LogResult) ReadField0(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v24, err25 := iprot.ReadI32()
	if err25 != nil {
		return thrift.NewTProtocolExceptionReadField(0, "success", p.ThriftName(), err25)
	}
	p.Success = ResultCode(v24)
	return err
}

func (p *LogResult) ReadFieldSuccess(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField0(iprot)
}

func (p *LogResult) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("Log_result")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	switch {
	default:
		if err = p.WriteField0(oprot); err != nil {
			return err
		}
	}
	err = oprot.WriteFieldStop()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(-1, "STOP", p.ThriftName(), err)
	}
	err = oprot.WriteStructEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	return err
}

func (p *LogResult) WriteField0(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("success", thrift.I32, 0)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(0, "success", p.ThriftName(), err)
	}
	err = oprot.WriteI32(int32(p.Success))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(0, "success", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(0, "success", p.ThriftName(), err)
	}
	return err
}

func (p *LogResult) WriteFieldSuccess(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField0(oprot)
}

func (p *LogResult) TStructName() string {
	return "LogResult"
}

func (p *LogResult) ThriftName() string {
	return "Log_result"
}

func (p *LogResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("LogResult(%+v)", *p)
}

func (p *LogResult) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*LogResult)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *LogResult) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 0:
		return p.Success
	}
	return nil
}

func (p *LogResult) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("success", thrift.I32, 0),
	})
}
