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

// Package fb303 contains common Apache Thrift protocol definitions for Facebook services like Apache Scribe
package fb303

import (
	"fmt"
	"tumblr/encoding/thrift"
)

type IFacebookService interface { /**
	 *Standard base service
	 */

	/**
	 * Returns a descriptive name of the service
	 */
	GetName() (retval0 string, err error)
	/**
	 * Returns the version of the service
	 */
	GetVersion() (retval1 string, err error)
	/**
	 * Gets the status of this service
	 */
	GetStatus() (retval2 FbStatus, err error)
	/**
	 * User friendly description of status, such as why the service is in
	 * the dead or warning state, or what is being started or stopped.
	 */
	GetStatusDetails() (retval3 string, err error)
	/**
	 * Gets the counters for this service
	 */
	GetCounters() (retval4 thrift.TMap, err error)
	/**
	 * Gets the value of a single counter
	 *
	 * Parameters:
	 *  - Key
	 */
	GetCounter(key string) (retval5 int64, err error)
	/**
	 * Sets an option
	 *
	 * Parameters:
	 *  - Key
	 *  - Value
	 */
	SetOption(key string, value string) (err error)
	/**
	 * Gets an option
	 *
	 * Parameters:
	 *  - Key
	 */
	GetOption(key string) (retval7 string, err error)
	/**
	 * Gets all options
	 */
	GetOptions() (retval8 thrift.TMap, err error)
	/**
	 * Returns a CPU profile over the given time interval (client and server
	 * must agree on the profile format).
	 *
	 * Parameters:
	 *  - ProfileDurationInSec
	 */
	GetCpuProfile(profileDurationInSec int32) (retval9 string, err error)
	/**
	 * Returns the unix time that the server has been running since
	 */
	AliveSince() (retval10 int64, err error)
	/**
	 * Tell the server to reload its configuration, reopen log files, etc
	 */
	Reinitialize() (err error)
	/**
	 * Suggest a shutdown to the server
	 */
	Shutdown() (err error)
}

/**
 *Standard base service
 */
type FacebookServiceClient struct {
	Transport       thrift.TTransport
	ProtocolFactory thrift.TProtocolFactory
	InputProtocol   thrift.TProtocol
	OutputProtocol  thrift.TProtocol
	SeqId           int32
}

func NewFacebookServiceClientFactory(t thrift.TTransport, f thrift.TProtocolFactory) *FacebookServiceClient {
	return &FacebookServiceClient{Transport: t,
		ProtocolFactory: f,
		InputProtocol:   f.GetProtocol(t),
		OutputProtocol:  f.GetProtocol(t),
		SeqId:           0,
	}
}

func NewFacebookServiceClientProtocol(t thrift.TTransport, iprot thrift.TProtocol, oprot thrift.TProtocol) *FacebookServiceClient {
	return &FacebookServiceClient{Transport: t,
		ProtocolFactory: nil,
		InputProtocol:   iprot,
		OutputProtocol:  oprot,
		SeqId:           0,
	}
}

/**
 * Returns a descriptive name of the service
 */
func (p *FacebookServiceClient) GetName() (retval13 string, err error) {
	err = p.SendGetName()
	if err != nil {
		return
	}
	return p.RecvGetName()
}

func (p *FacebookServiceClient) SendGetName() (err error) {
	oprot := p.OutputProtocol
	if oprot != nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	oprot.WriteMessageBegin("getName", thrift.CALL, p.SeqId)
	args14 := NewGetNameArgs()
	err = args14.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Transport().Flush()
	return
}

func (p *FacebookServiceClient) RecvGetName() (value string, err error) {
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
		error16 := thrift.NewTApplicationExceptionDefault()
		var error17 error
		error17, err = error16.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error17
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "ping failed: out of sequence response")
		return
	}
	result15 := NewGetNameResult()
	err = result15.Read(iprot)
	iprot.ReadMessageEnd()
	value = result15.Success
	return
}

/**
 * Returns the version of the service
 */
func (p *FacebookServiceClient) GetVersion() (retval18 string, err error) {
	err = p.SendGetVersion()
	if err != nil {
		return
	}
	return p.RecvGetVersion()
}

func (p *FacebookServiceClient) SendGetVersion() (err error) {
	oprot := p.OutputProtocol
	if oprot != nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	oprot.WriteMessageBegin("getVersion", thrift.CALL, p.SeqId)
	args19 := NewGetVersionArgs()
	err = args19.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Transport().Flush()
	return
}

func (p *FacebookServiceClient) RecvGetVersion() (value string, err error) {
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
		error21 := thrift.NewTApplicationExceptionDefault()
		var error22 error
		error22, err = error21.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error22
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "ping failed: out of sequence response")
		return
	}
	result20 := NewGetVersionResult()
	err = result20.Read(iprot)
	iprot.ReadMessageEnd()
	value = result20.Success
	return
}

/**
 * Gets the status of this service
 */
func (p *FacebookServiceClient) GetStatus() (retval23 FbStatus, err error) {
	err = p.SendGetStatus()
	if err != nil {
		return
	}
	return p.RecvGetStatus()
}

func (p *FacebookServiceClient) SendGetStatus() (err error) {
	oprot := p.OutputProtocol
	if oprot != nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	oprot.WriteMessageBegin("getStatus", thrift.CALL, p.SeqId)
	args24 := NewGetStatusArgs()
	err = args24.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Transport().Flush()
	return
}

func (p *FacebookServiceClient) RecvGetStatus() (value FbStatus, err error) {
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
		error26 := thrift.NewTApplicationExceptionDefault()
		var error27 error
		error27, err = error26.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error27
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "ping failed: out of sequence response")
		return
	}
	result25 := NewGetStatusResult()
	err = result25.Read(iprot)
	iprot.ReadMessageEnd()
	value = result25.Success
	return
}

/**
 * User friendly description of status, such as why the service is in
 * the dead or warning state, or what is being started or stopped.
 */
func (p *FacebookServiceClient) GetStatusDetails() (retval28 string, err error) {
	err = p.SendGetStatusDetails()
	if err != nil {
		return
	}
	return p.RecvGetStatusDetails()
}

func (p *FacebookServiceClient) SendGetStatusDetails() (err error) {
	oprot := p.OutputProtocol
	if oprot != nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	oprot.WriteMessageBegin("getStatusDetails", thrift.CALL, p.SeqId)
	args29 := NewGetStatusDetailsArgs()
	err = args29.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Transport().Flush()
	return
}

func (p *FacebookServiceClient) RecvGetStatusDetails() (value string, err error) {
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
		error31 := thrift.NewTApplicationExceptionDefault()
		var error32 error
		error32, err = error31.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error32
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "ping failed: out of sequence response")
		return
	}
	result30 := NewGetStatusDetailsResult()
	err = result30.Read(iprot)
	iprot.ReadMessageEnd()
	value = result30.Success
	return
}

/**
 * Gets the counters for this service
 */
func (p *FacebookServiceClient) GetCounters() (retval33 thrift.TMap, err error) {
	err = p.SendGetCounters()
	if err != nil {
		return
	}
	return p.RecvGetCounters()
}

func (p *FacebookServiceClient) SendGetCounters() (err error) {
	oprot := p.OutputProtocol
	if oprot != nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	oprot.WriteMessageBegin("getCounters", thrift.CALL, p.SeqId)
	args34 := NewGetCountersArgs()
	err = args34.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Transport().Flush()
	return
}

func (p *FacebookServiceClient) RecvGetCounters() (value thrift.TMap, err error) {
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
		error36 := thrift.NewTApplicationExceptionDefault()
		var error37 error
		error37, err = error36.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error37
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "ping failed: out of sequence response")
		return
	}
	result35 := NewGetCountersResult()
	err = result35.Read(iprot)
	iprot.ReadMessageEnd()
	value = result35.Success
	return
}

/**
 * Gets the value of a single counter
 *
 * Parameters:
 *  - Key
 */
func (p *FacebookServiceClient) GetCounter(key string) (retval38 int64, err error) {
	err = p.SendGetCounter(key)
	if err != nil {
		return
	}
	return p.RecvGetCounter()
}

func (p *FacebookServiceClient) SendGetCounter(key string) (err error) {
	oprot := p.OutputProtocol
	if oprot != nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	oprot.WriteMessageBegin("getCounter", thrift.CALL, p.SeqId)
	args39 := NewGetCounterArgs()
	args39.Key = key
	err = args39.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Transport().Flush()
	return
}

func (p *FacebookServiceClient) RecvGetCounter() (value int64, err error) {
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
		error41 := thrift.NewTApplicationExceptionDefault()
		var error42 error
		error42, err = error41.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error42
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "ping failed: out of sequence response")
		return
	}
	result40 := NewGetCounterResult()
	err = result40.Read(iprot)
	iprot.ReadMessageEnd()
	value = result40.Success
	return
}

/**
 * Sets an option
 *
 * Parameters:
 *  - Key
 *  - Value
 */
func (p *FacebookServiceClient) SetOption(key string, value string) (err error) {
	err = p.SendSetOption(key, value)
	if err != nil {
		return
	}
	return p.RecvSetOption()
}

func (p *FacebookServiceClient) SendSetOption(key string, value string) (err error) {
	oprot := p.OutputProtocol
	if oprot != nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	oprot.WriteMessageBegin("setOption", thrift.CALL, p.SeqId)
	args44 := NewSetOptionArgs()
	args44.Key = key
	args44.Value = value
	err = args44.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Transport().Flush()
	return
}

func (p *FacebookServiceClient) RecvSetOption() (err error) {
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
		error46 := thrift.NewTApplicationExceptionDefault()
		var error47 error
		error47, err = error46.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error47
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "ping failed: out of sequence response")
		return
	}
	result45 := NewSetOptionResult()
	err = result45.Read(iprot)
	iprot.ReadMessageEnd()
	return
}

/**
 * Gets an option
 *
 * Parameters:
 *  - Key
 */
func (p *FacebookServiceClient) GetOption(key string) (retval48 string, err error) {
	err = p.SendGetOption(key)
	if err != nil {
		return
	}
	return p.RecvGetOption()
}

func (p *FacebookServiceClient) SendGetOption(key string) (err error) {
	oprot := p.OutputProtocol
	if oprot != nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	oprot.WriteMessageBegin("getOption", thrift.CALL, p.SeqId)
	args49 := NewGetOptionArgs()
	args49.Key = key
	err = args49.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Transport().Flush()
	return
}

func (p *FacebookServiceClient) RecvGetOption() (value string, err error) {
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
		error51 := thrift.NewTApplicationExceptionDefault()
		var error52 error
		error52, err = error51.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error52
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "ping failed: out of sequence response")
		return
	}
	result50 := NewGetOptionResult()
	err = result50.Read(iprot)
	iprot.ReadMessageEnd()
	value = result50.Success
	return
}

/**
 * Gets all options
 */
func (p *FacebookServiceClient) GetOptions() (retval53 thrift.TMap, err error) {
	err = p.SendGetOptions()
	if err != nil {
		return
	}
	return p.RecvGetOptions()
}

func (p *FacebookServiceClient) SendGetOptions() (err error) {
	oprot := p.OutputProtocol
	if oprot != nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	oprot.WriteMessageBegin("getOptions", thrift.CALL, p.SeqId)
	args54 := NewGetOptionsArgs()
	err = args54.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Transport().Flush()
	return
}

func (p *FacebookServiceClient) RecvGetOptions() (value thrift.TMap, err error) {
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
		error56 := thrift.NewTApplicationExceptionDefault()
		var error57 error
		error57, err = error56.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error57
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "ping failed: out of sequence response")
		return
	}
	result55 := NewGetOptionsResult()
	err = result55.Read(iprot)
	iprot.ReadMessageEnd()
	value = result55.Success
	return
}

/**
 * Returns a CPU profile over the given time interval (client and server
 * must agree on the profile format).
 *
 * Parameters:
 *  - ProfileDurationInSec
 */
func (p *FacebookServiceClient) GetCpuProfile(profileDurationInSec int32) (retval58 string, err error) {
	err = p.SendGetCpuProfile(profileDurationInSec)
	if err != nil {
		return
	}
	return p.RecvGetCpuProfile()
}

func (p *FacebookServiceClient) SendGetCpuProfile(profileDurationInSec int32) (err error) {
	oprot := p.OutputProtocol
	if oprot != nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	oprot.WriteMessageBegin("getCpuProfile", thrift.CALL, p.SeqId)
	args59 := NewGetCpuProfileArgs()
	args59.ProfileDurationInSec = profileDurationInSec
	err = args59.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Transport().Flush()
	return
}

func (p *FacebookServiceClient) RecvGetCpuProfile() (value string, err error) {
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
		error61 := thrift.NewTApplicationExceptionDefault()
		var error62 error
		error62, err = error61.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error62
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "ping failed: out of sequence response")
		return
	}
	result60 := NewGetCpuProfileResult()
	err = result60.Read(iprot)
	iprot.ReadMessageEnd()
	value = result60.Success
	return
}

/**
 * Returns the unix time that the server has been running since
 */
func (p *FacebookServiceClient) AliveSince() (retval63 int64, err error) {
	err = p.SendAliveSince()
	if err != nil {
		return
	}
	return p.RecvAliveSince()
}

func (p *FacebookServiceClient) SendAliveSince() (err error) {
	oprot := p.OutputProtocol
	if oprot != nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	oprot.WriteMessageBegin("aliveSince", thrift.CALL, p.SeqId)
	args64 := NewAliveSinceArgs()
	err = args64.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Transport().Flush()
	return
}

func (p *FacebookServiceClient) RecvAliveSince() (value int64, err error) {
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
		error66 := thrift.NewTApplicationExceptionDefault()
		var error67 error
		error67, err = error66.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error67
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "ping failed: out of sequence response")
		return
	}
	result65 := NewAliveSinceResult()
	err = result65.Read(iprot)
	iprot.ReadMessageEnd()
	value = result65.Success
	return
}

/**
 * Tell the server to reload its configuration, reopen log files, etc
 */
func (p *FacebookServiceClient) Reinitialize() (err error) {
	err = p.SendReinitialize()
	if err != nil {
		return
	}
	return
}

func (p *FacebookServiceClient) SendReinitialize() (err error) {
	oprot := p.OutputProtocol
	if oprot != nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	oprot.WriteMessageBegin("reinitialize", thrift.CALL, p.SeqId)
	args69 := NewReinitializeArgs()
	err = args69.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Transport().Flush()
	return
}

func (p *FacebookServiceClient) RecvReinitialize() (err error) {
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
		error71 := thrift.NewTApplicationExceptionDefault()
		var error72 error
		error72, err = error71.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error72
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "ping failed: out of sequence response")
		return
	}
	result70 := NewReinitializeResult()
	err = result70.Read(iprot)
	iprot.ReadMessageEnd()
	return
}

/**
 * Suggest a shutdown to the server
 */
func (p *FacebookServiceClient) Shutdown() (err error) {
	err = p.SendShutdown()
	if err != nil {
		return
	}
	return
}

func (p *FacebookServiceClient) SendShutdown() (err error) {
	oprot := p.OutputProtocol
	if oprot != nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	oprot.WriteMessageBegin("shutdown", thrift.CALL, p.SeqId)
	args74 := NewShutdownArgs()
	err = args74.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Transport().Flush()
	return
}

func (p *FacebookServiceClient) RecvShutdown() (err error) {
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
		error76 := thrift.NewTApplicationExceptionDefault()
		var error77 error
		error77, err = error76.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error77
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "ping failed: out of sequence response")
		return
	}
	result75 := NewShutdownResult()
	err = result75.Read(iprot)
	iprot.ReadMessageEnd()
	return
}

type FacebookServiceProcessor struct {
	handler      IFacebookService
	processorMap map[string]thrift.TProcessorFunction
}

func (p *FacebookServiceProcessor) Handler() IFacebookService {
	return p.handler
}

func (p *FacebookServiceProcessor) AddToProcessorMap(key string, processor thrift.TProcessorFunction) {
	p.processorMap[key] = processor
}

func (p *FacebookServiceProcessor) GetProcessorFunction(key string) (processor thrift.TProcessorFunction, exists bool) {
	processor, exists = p.processorMap[key]
	return processor, exists
}

func (p *FacebookServiceProcessor) ProcessorMap() map[string]thrift.TProcessorFunction {
	return p.processorMap
}

func NewFacebookServiceProcessor(handler IFacebookService) *FacebookServiceProcessor {

	self78 := &FacebookServiceProcessor{handler: handler, processorMap: make(map[string]thrift.TProcessorFunction)}
	self78.processorMap["getName"] = &facebookServiceProcessorGetName{handler: handler}
	self78.processorMap["getVersion"] = &facebookServiceProcessorGetVersion{handler: handler}
	self78.processorMap["getStatus"] = &facebookServiceProcessorGetStatus{handler: handler}
	self78.processorMap["getStatusDetails"] = &facebookServiceProcessorGetStatusDetails{handler: handler}
	self78.processorMap["getCounters"] = &facebookServiceProcessorGetCounters{handler: handler}
	self78.processorMap["getCounter"] = &facebookServiceProcessorGetCounter{handler: handler}
	self78.processorMap["setOption"] = &facebookServiceProcessorSetOption{handler: handler}
	self78.processorMap["getOption"] = &facebookServiceProcessorGetOption{handler: handler}
	self78.processorMap["getOptions"] = &facebookServiceProcessorGetOptions{handler: handler}
	self78.processorMap["getCpuProfile"] = &facebookServiceProcessorGetCpuProfile{handler: handler}
	self78.processorMap["aliveSince"] = &facebookServiceProcessorAliveSince{handler: handler}
	self78.processorMap["reinitialize"] = &facebookServiceProcessorReinitialize{handler: handler}
	self78.processorMap["shutdown"] = &facebookServiceProcessorShutdown{handler: handler}
	return self78
}

func (p *FacebookServiceProcessor) Process(iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	name, _, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return
	}
	process, nameFound := p.GetProcessorFunction(name)
	if !nameFound || process == nil {
		iprot.Skip(thrift.STRUCT)
		iprot.ReadMessageEnd()
		x79 := thrift.NewTApplicationException(thrift.UNKNOWN_METHOD, "Unknown function "+name)
		oprot.WriteMessageBegin(name, thrift.EXCEPTION, seqId)
		x79.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Transport().Flush()
		return false, x79
	}
	return process.Process(seqId, iprot, oprot)
}

type facebookServiceProcessorGetName struct {
	handler IFacebookService
}

func (p *facebookServiceProcessorGetName) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := NewGetNameArgs()
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("getName", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Transport().Flush()
		return
	}
	iprot.ReadMessageEnd()
	result := NewGetNameResult()
	if result.Success, err = p.handler.GetName(); err != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing getName: "+err.Error())
		oprot.WriteMessageBegin("getName", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Transport().Flush()
		return
	}
	if err2 := oprot.WriteMessageBegin("getName", thrift.REPLY, seqId); err2 != nil {
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

type facebookServiceProcessorGetVersion struct {
	handler IFacebookService
}

func (p *facebookServiceProcessorGetVersion) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := NewGetVersionArgs()
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("getVersion", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Transport().Flush()
		return
	}
	iprot.ReadMessageEnd()
	result := NewGetVersionResult()
	if result.Success, err = p.handler.GetVersion(); err != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing getVersion: "+err.Error())
		oprot.WriteMessageBegin("getVersion", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Transport().Flush()
		return
	}
	if err2 := oprot.WriteMessageBegin("getVersion", thrift.REPLY, seqId); err2 != nil {
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

type facebookServiceProcessorGetStatus struct {
	handler IFacebookService
}

func (p *facebookServiceProcessorGetStatus) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := NewGetStatusArgs()
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("getStatus", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Transport().Flush()
		return
	}
	iprot.ReadMessageEnd()
	result := NewGetStatusResult()
	if result.Success, err = p.handler.GetStatus(); err != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing getStatus: "+err.Error())
		oprot.WriteMessageBegin("getStatus", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Transport().Flush()
		return
	}
	if err2 := oprot.WriteMessageBegin("getStatus", thrift.REPLY, seqId); err2 != nil {
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

type facebookServiceProcessorGetStatusDetails struct {
	handler IFacebookService
}

func (p *facebookServiceProcessorGetStatusDetails) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := NewGetStatusDetailsArgs()
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("getStatusDetails", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Transport().Flush()
		return
	}
	iprot.ReadMessageEnd()
	result := NewGetStatusDetailsResult()
	if result.Success, err = p.handler.GetStatusDetails(); err != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing getStatusDetails: "+err.Error())
		oprot.WriteMessageBegin("getStatusDetails", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Transport().Flush()
		return
	}
	if err2 := oprot.WriteMessageBegin("getStatusDetails", thrift.REPLY, seqId); err2 != nil {
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

type facebookServiceProcessorGetCounters struct {
	handler IFacebookService
}

func (p *facebookServiceProcessorGetCounters) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := NewGetCountersArgs()
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("getCounters", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Transport().Flush()
		return
	}
	iprot.ReadMessageEnd()
	result := NewGetCountersResult()
	if result.Success, err = p.handler.GetCounters(); err != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing getCounters: "+err.Error())
		oprot.WriteMessageBegin("getCounters", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Transport().Flush()
		return
	}
	if err2 := oprot.WriteMessageBegin("getCounters", thrift.REPLY, seqId); err2 != nil {
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

type facebookServiceProcessorGetCounter struct {
	handler IFacebookService
}

func (p *facebookServiceProcessorGetCounter) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := NewGetCounterArgs()
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("getCounter", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Transport().Flush()
		return
	}
	iprot.ReadMessageEnd()
	result := NewGetCounterResult()
	if result.Success, err = p.handler.GetCounter(args.Key); err != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing getCounter: "+err.Error())
		oprot.WriteMessageBegin("getCounter", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Transport().Flush()
		return
	}
	if err2 := oprot.WriteMessageBegin("getCounter", thrift.REPLY, seqId); err2 != nil {
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

type facebookServiceProcessorSetOption struct {
	handler IFacebookService
}

func (p *facebookServiceProcessorSetOption) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := NewSetOptionArgs()
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("setOption", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Transport().Flush()
		return
	}
	iprot.ReadMessageEnd()
	result := NewSetOptionResult()
	if err = p.handler.SetOption(args.Key, args.Value); err != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing setOption: "+err.Error())
		oprot.WriteMessageBegin("setOption", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Transport().Flush()
		return
	}
	if err2 := oprot.WriteMessageBegin("setOption", thrift.REPLY, seqId); err2 != nil {
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

type facebookServiceProcessorGetOption struct {
	handler IFacebookService
}

func (p *facebookServiceProcessorGetOption) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := NewGetOptionArgs()
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("getOption", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Transport().Flush()
		return
	}
	iprot.ReadMessageEnd()
	result := NewGetOptionResult()
	if result.Success, err = p.handler.GetOption(args.Key); err != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing getOption: "+err.Error())
		oprot.WriteMessageBegin("getOption", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Transport().Flush()
		return
	}
	if err2 := oprot.WriteMessageBegin("getOption", thrift.REPLY, seqId); err2 != nil {
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

type facebookServiceProcessorGetOptions struct {
	handler IFacebookService
}

func (p *facebookServiceProcessorGetOptions) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := NewGetOptionsArgs()
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("getOptions", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Transport().Flush()
		return
	}
	iprot.ReadMessageEnd()
	result := NewGetOptionsResult()
	if result.Success, err = p.handler.GetOptions(); err != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing getOptions: "+err.Error())
		oprot.WriteMessageBegin("getOptions", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Transport().Flush()
		return
	}
	if err2 := oprot.WriteMessageBegin("getOptions", thrift.REPLY, seqId); err2 != nil {
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

type facebookServiceProcessorGetCpuProfile struct {
	handler IFacebookService
}

func (p *facebookServiceProcessorGetCpuProfile) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := NewGetCpuProfileArgs()
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("getCpuProfile", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Transport().Flush()
		return
	}
	iprot.ReadMessageEnd()
	result := NewGetCpuProfileResult()
	if result.Success, err = p.handler.GetCpuProfile(args.ProfileDurationInSec); err != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing getCpuProfile: "+err.Error())
		oprot.WriteMessageBegin("getCpuProfile", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Transport().Flush()
		return
	}
	if err2 := oprot.WriteMessageBegin("getCpuProfile", thrift.REPLY, seqId); err2 != nil {
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

type facebookServiceProcessorAliveSince struct {
	handler IFacebookService
}

func (p *facebookServiceProcessorAliveSince) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := NewAliveSinceArgs()
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("aliveSince", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Transport().Flush()
		return
	}
	iprot.ReadMessageEnd()
	result := NewAliveSinceResult()
	if result.Success, err = p.handler.AliveSince(); err != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing aliveSince: "+err.Error())
		oprot.WriteMessageBegin("aliveSince", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Transport().Flush()
		return
	}
	if err2 := oprot.WriteMessageBegin("aliveSince", thrift.REPLY, seqId); err2 != nil {
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

type facebookServiceProcessorReinitialize struct {
	handler IFacebookService
}

func (p *facebookServiceProcessorReinitialize) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := NewReinitializeArgs()
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("reinitialize", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Transport().Flush()
		return
	}
	iprot.ReadMessageEnd()
	result := NewReinitializeResult()
	if err = p.handler.Reinitialize(); err != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing reinitialize: "+err.Error())
		oprot.WriteMessageBegin("reinitialize", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Transport().Flush()
		return
	}
	if err2 := oprot.WriteMessageBegin("reinitialize", thrift.REPLY, seqId); err2 != nil {
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

type facebookServiceProcessorShutdown struct {
	handler IFacebookService
}

func (p *facebookServiceProcessorShutdown) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := NewShutdownArgs()
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("shutdown", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Transport().Flush()
		return
	}
	iprot.ReadMessageEnd()
	result := NewShutdownResult()
	if err = p.handler.Shutdown(); err != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing shutdown: "+err.Error())
		oprot.WriteMessageBegin("shutdown", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Transport().Flush()
		return
	}
	if err2 := oprot.WriteMessageBegin("shutdown", thrift.REPLY, seqId); err2 != nil {
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

type GetNameArgs struct {
	thrift.TStruct
}

func NewGetNameArgs() *GetNameArgs {
	output := &GetNameArgs{
		TStruct: thrift.NewTStruct("getName_args", []thrift.TField{}),
	}
	{
	}
	return output
}

func (p *GetNameArgs) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
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
		err = iprot.Skip(fieldTypeId)
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
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

func (p *GetNameArgs) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("getName_args")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
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

func (p *GetNameArgs) TStructName() string {
	return "GetNameArgs"
}

func (p *GetNameArgs) ThriftName() string {
	return "getName_args"
}

func (p *GetNameArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("GetNameArgs(%+v)", *p)
}

func (p *GetNameArgs) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*GetNameArgs)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *GetNameArgs) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	}
	return nil
}

func (p *GetNameArgs) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{})
}

/**
 * Attributes:
 *  - Success
 */
type GetNameResult struct {
	thrift.TStruct
	Success string "success" // 0
}

func NewGetNameResult() *GetNameResult {
	output := &GetNameResult{
		TStruct: thrift.NewTStruct("getName_result", []thrift.TField{
			thrift.NewTField("success", thrift.STRING, 0),
		}),
	}
	{
	}
	return output
}

func (p *GetNameResult) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
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
			if fieldTypeId == thrift.STRING {
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

func (p *GetNameResult) ReadField0(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v80, err81 := iprot.ReadString()
	if err81 != nil {
		return thrift.NewTProtocolExceptionReadField(0, "success", p.ThriftName(), err81)
	}
	p.Success = v80
	return err
}

func (p *GetNameResult) ReadFieldSuccess(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField0(iprot)
}

func (p *GetNameResult) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("getName_result")
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

func (p *GetNameResult) WriteField0(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("success", thrift.STRING, 0)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(0, "success", p.ThriftName(), err)
	}
	err = oprot.WriteString(string(p.Success))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(0, "success", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(0, "success", p.ThriftName(), err)
	}
	return err
}

func (p *GetNameResult) WriteFieldSuccess(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField0(oprot)
}

func (p *GetNameResult) TStructName() string {
	return "GetNameResult"
}

func (p *GetNameResult) ThriftName() string {
	return "getName_result"
}

func (p *GetNameResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("GetNameResult(%+v)", *p)
}

func (p *GetNameResult) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*GetNameResult)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *GetNameResult) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 0:
		return p.Success
	}
	return nil
}

func (p *GetNameResult) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("success", thrift.STRING, 0),
	})
}

type GetVersionArgs struct {
	thrift.TStruct
}

func NewGetVersionArgs() *GetVersionArgs {
	output := &GetVersionArgs{
		TStruct: thrift.NewTStruct("getVersion_args", []thrift.TField{}),
	}
	{
	}
	return output
}

func (p *GetVersionArgs) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
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
		err = iprot.Skip(fieldTypeId)
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
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

func (p *GetVersionArgs) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("getVersion_args")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
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

func (p *GetVersionArgs) TStructName() string {
	return "GetVersionArgs"
}

func (p *GetVersionArgs) ThriftName() string {
	return "getVersion_args"
}

func (p *GetVersionArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("GetVersionArgs(%+v)", *p)
}

func (p *GetVersionArgs) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*GetVersionArgs)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *GetVersionArgs) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	}
	return nil
}

func (p *GetVersionArgs) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{})
}

/**
 * Attributes:
 *  - Success
 */
type GetVersionResult struct {
	thrift.TStruct
	Success string "success" // 0
}

func NewGetVersionResult() *GetVersionResult {
	output := &GetVersionResult{
		TStruct: thrift.NewTStruct("getVersion_result", []thrift.TField{
			thrift.NewTField("success", thrift.STRING, 0),
		}),
	}
	{
	}
	return output
}

func (p *GetVersionResult) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
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
			if fieldTypeId == thrift.STRING {
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

func (p *GetVersionResult) ReadField0(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v82, err83 := iprot.ReadString()
	if err83 != nil {
		return thrift.NewTProtocolExceptionReadField(0, "success", p.ThriftName(), err83)
	}
	p.Success = v82
	return err
}

func (p *GetVersionResult) ReadFieldSuccess(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField0(iprot)
}

func (p *GetVersionResult) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("getVersion_result")
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

func (p *GetVersionResult) WriteField0(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("success", thrift.STRING, 0)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(0, "success", p.ThriftName(), err)
	}
	err = oprot.WriteString(string(p.Success))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(0, "success", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(0, "success", p.ThriftName(), err)
	}
	return err
}

func (p *GetVersionResult) WriteFieldSuccess(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField0(oprot)
}

func (p *GetVersionResult) TStructName() string {
	return "GetVersionResult"
}

func (p *GetVersionResult) ThriftName() string {
	return "getVersion_result"
}

func (p *GetVersionResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("GetVersionResult(%+v)", *p)
}

func (p *GetVersionResult) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*GetVersionResult)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *GetVersionResult) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 0:
		return p.Success
	}
	return nil
}

func (p *GetVersionResult) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("success", thrift.STRING, 0),
	})
}

type GetStatusArgs struct {
	thrift.TStruct
}

func NewGetStatusArgs() *GetStatusArgs {
	output := &GetStatusArgs{
		TStruct: thrift.NewTStruct("getStatus_args", []thrift.TField{}),
	}
	{
	}
	return output
}

func (p *GetStatusArgs) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
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
		err = iprot.Skip(fieldTypeId)
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
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

func (p *GetStatusArgs) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("getStatus_args")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
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

func (p *GetStatusArgs) TStructName() string {
	return "GetStatusArgs"
}

func (p *GetStatusArgs) ThriftName() string {
	return "getStatus_args"
}

func (p *GetStatusArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("GetStatusArgs(%+v)", *p)
}

func (p *GetStatusArgs) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*GetStatusArgs)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *GetStatusArgs) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	}
	return nil
}

func (p *GetStatusArgs) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{})
}

/**
 * Attributes:
 *  - Success
 */
type GetStatusResult struct {
	thrift.TStruct
	Success FbStatus "success" // 0
}

func NewGetStatusResult() *GetStatusResult {
	output := &GetStatusResult{
		TStruct: thrift.NewTStruct("getStatus_result", []thrift.TField{
			thrift.NewTField("success", thrift.I32, 0),
		}),
	}
	{
	}
	return output
}

func (p *GetStatusResult) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
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

func (p *GetStatusResult) ReadField0(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v84, err85 := iprot.ReadI32()
	if err85 != nil {
		return thrift.NewTProtocolExceptionReadField(0, "success", p.ThriftName(), err85)
	}
	p.Success = FbStatus(v84)
	return err
}

func (p *GetStatusResult) ReadFieldSuccess(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField0(iprot)
}

func (p *GetStatusResult) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("getStatus_result")
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

func (p *GetStatusResult) WriteField0(oprot thrift.TProtocol) (err thrift.TProtocolException) {
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

func (p *GetStatusResult) WriteFieldSuccess(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField0(oprot)
}

func (p *GetStatusResult) TStructName() string {
	return "GetStatusResult"
}

func (p *GetStatusResult) ThriftName() string {
	return "getStatus_result"
}

func (p *GetStatusResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("GetStatusResult(%+v)", *p)
}

func (p *GetStatusResult) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*GetStatusResult)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *GetStatusResult) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 0:
		return p.Success
	}
	return nil
}

func (p *GetStatusResult) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("success", thrift.I32, 0),
	})
}

type GetStatusDetailsArgs struct {
	thrift.TStruct
}

func NewGetStatusDetailsArgs() *GetStatusDetailsArgs {
	output := &GetStatusDetailsArgs{
		TStruct: thrift.NewTStruct("getStatusDetails_args", []thrift.TField{}),
	}
	{
	}
	return output
}

func (p *GetStatusDetailsArgs) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
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
		err = iprot.Skip(fieldTypeId)
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
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

func (p *GetStatusDetailsArgs) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("getStatusDetails_args")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
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

func (p *GetStatusDetailsArgs) TStructName() string {
	return "GetStatusDetailsArgs"
}

func (p *GetStatusDetailsArgs) ThriftName() string {
	return "getStatusDetails_args"
}

func (p *GetStatusDetailsArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("GetStatusDetailsArgs(%+v)", *p)
}

func (p *GetStatusDetailsArgs) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*GetStatusDetailsArgs)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *GetStatusDetailsArgs) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	}
	return nil
}

func (p *GetStatusDetailsArgs) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{})
}

/**
 * Attributes:
 *  - Success
 */
type GetStatusDetailsResult struct {
	thrift.TStruct
	Success string "success" // 0
}

func NewGetStatusDetailsResult() *GetStatusDetailsResult {
	output := &GetStatusDetailsResult{
		TStruct: thrift.NewTStruct("getStatusDetails_result", []thrift.TField{
			thrift.NewTField("success", thrift.STRING, 0),
		}),
	}
	{
	}
	return output
}

func (p *GetStatusDetailsResult) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
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
			if fieldTypeId == thrift.STRING {
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

func (p *GetStatusDetailsResult) ReadField0(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v86, err87 := iprot.ReadString()
	if err87 != nil {
		return thrift.NewTProtocolExceptionReadField(0, "success", p.ThriftName(), err87)
	}
	p.Success = v86
	return err
}

func (p *GetStatusDetailsResult) ReadFieldSuccess(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField0(iprot)
}

func (p *GetStatusDetailsResult) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("getStatusDetails_result")
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

func (p *GetStatusDetailsResult) WriteField0(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("success", thrift.STRING, 0)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(0, "success", p.ThriftName(), err)
	}
	err = oprot.WriteString(string(p.Success))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(0, "success", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(0, "success", p.ThriftName(), err)
	}
	return err
}

func (p *GetStatusDetailsResult) WriteFieldSuccess(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField0(oprot)
}

func (p *GetStatusDetailsResult) TStructName() string {
	return "GetStatusDetailsResult"
}

func (p *GetStatusDetailsResult) ThriftName() string {
	return "getStatusDetails_result"
}

func (p *GetStatusDetailsResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("GetStatusDetailsResult(%+v)", *p)
}

func (p *GetStatusDetailsResult) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*GetStatusDetailsResult)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *GetStatusDetailsResult) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 0:
		return p.Success
	}
	return nil
}

func (p *GetStatusDetailsResult) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("success", thrift.STRING, 0),
	})
}

type GetCountersArgs struct {
	thrift.TStruct
}

func NewGetCountersArgs() *GetCountersArgs {
	output := &GetCountersArgs{
		TStruct: thrift.NewTStruct("getCounters_args", []thrift.TField{}),
	}
	{
	}
	return output
}

func (p *GetCountersArgs) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
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
		err = iprot.Skip(fieldTypeId)
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
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

func (p *GetCountersArgs) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("getCounters_args")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
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

func (p *GetCountersArgs) TStructName() string {
	return "GetCountersArgs"
}

func (p *GetCountersArgs) ThriftName() string {
	return "getCounters_args"
}

func (p *GetCountersArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("GetCountersArgs(%+v)", *p)
}

func (p *GetCountersArgs) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*GetCountersArgs)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *GetCountersArgs) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	}
	return nil
}

func (p *GetCountersArgs) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{})
}

/**
 * Attributes:
 *  - Success
 */
type GetCountersResult struct {
	thrift.TStruct
	Success thrift.TMap "success" // 0
}

func NewGetCountersResult() *GetCountersResult {
	output := &GetCountersResult{
		TStruct: thrift.NewTStruct("getCounters_result", []thrift.TField{
			thrift.NewTField("success", thrift.MAP, 0),
		}),
	}
	{
	}
	return output
}

func (p *GetCountersResult) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
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
			if fieldTypeId == thrift.MAP {
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

func (p *GetCountersResult) ReadField0(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_ktype91, _vtype92, _size90, err := iprot.ReadMapBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "p.Success", "", err)
	}
	p.Success = thrift.NewTMap(_ktype91, _vtype92, _size90)
	for _i94 := 0; _i94 < _size90; _i94++ {
		v97, err98 := iprot.ReadString()
		if err98 != nil {
			return thrift.NewTProtocolExceptionReadField(0, "_key95", "", err98)
		}
		_key95 := v97
		v99, err100 := iprot.ReadI64()
		if err100 != nil {
			return thrift.NewTProtocolExceptionReadField(0, "_val96", "", err100)
		}
		_val96 := v99
		p.Success.Set(_key95, _val96)
	}
	err = iprot.ReadMapEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "", "map", err)
	}
	return err
}

func (p *GetCountersResult) ReadFieldSuccess(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField0(iprot)
}

func (p *GetCountersResult) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("getCounters_result")
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

func (p *GetCountersResult) WriteField0(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.Success != nil {
		err = oprot.WriteFieldBegin("success", thrift.MAP, 0)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(0, "success", p.ThriftName(), err)
		}
		err = oprot.WriteMapBegin(thrift.STRING, thrift.I64, p.Success.Len())
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(-1, "", "map", err)
		}
		for Miter101 := range p.Success.Iter() {
			Kiter102, Viter103 := Miter101.Key().(string), Miter101.Value().(int64)
			err = oprot.WriteString(string(Kiter102))
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(0, "Kiter102", "", err)
			}
			err = oprot.WriteI64(int64(Viter103))
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(0, "Viter103", "", err)
			}
		}
		err = oprot.WriteMapEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(-1, "", "map", err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(0, "success", p.ThriftName(), err)
		}
	}
	return err
}

func (p *GetCountersResult) WriteFieldSuccess(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField0(oprot)
}

func (p *GetCountersResult) TStructName() string {
	return "GetCountersResult"
}

func (p *GetCountersResult) ThriftName() string {
	return "getCounters_result"
}

func (p *GetCountersResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("GetCountersResult(%+v)", *p)
}

func (p *GetCountersResult) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*GetCountersResult)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *GetCountersResult) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 0:
		return p.Success
	}
	return nil
}

func (p *GetCountersResult) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("success", thrift.MAP, 0),
	})
}

/**
 * Attributes:
 *  - Key
 */
type GetCounterArgs struct {
	thrift.TStruct
	Key string "key" // 1
}

func NewGetCounterArgs() *GetCounterArgs {
	output := &GetCounterArgs{
		TStruct: thrift.NewTStruct("getCounter_args", []thrift.TField{
			thrift.NewTField("key", thrift.STRING, 1),
		}),
	}
	{
	}
	return output
}

func (p *GetCounterArgs) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
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
		if fieldId == 1 || fieldName == "key" {
			if fieldTypeId == thrift.STRING {
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

func (p *GetCounterArgs) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v104, err105 := iprot.ReadString()
	if err105 != nil {
		return thrift.NewTProtocolExceptionReadField(1, "key", p.ThriftName(), err105)
	}
	p.Key = v104
	return err
}

func (p *GetCounterArgs) ReadFieldKey(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *GetCounterArgs) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("getCounter_args")
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

func (p *GetCounterArgs) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("key", thrift.STRING, 1)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "key", p.ThriftName(), err)
	}
	err = oprot.WriteString(string(p.Key))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "key", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "key", p.ThriftName(), err)
	}
	return err
}

func (p *GetCounterArgs) WriteFieldKey(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *GetCounterArgs) TStructName() string {
	return "GetCounterArgs"
}

func (p *GetCounterArgs) ThriftName() string {
	return "getCounter_args"
}

func (p *GetCounterArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("GetCounterArgs(%+v)", *p)
}

func (p *GetCounterArgs) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*GetCounterArgs)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *GetCounterArgs) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.Key
	}
	return nil
}

func (p *GetCounterArgs) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("key", thrift.STRING, 1),
	})
}

/**
 * Attributes:
 *  - Success
 */
type GetCounterResult struct {
	thrift.TStruct
	Success int64 "success" // 0
}

func NewGetCounterResult() *GetCounterResult {
	output := &GetCounterResult{
		TStruct: thrift.NewTStruct("getCounter_result", []thrift.TField{
			thrift.NewTField("success", thrift.I64, 0),
		}),
	}
	{
	}
	return output
}

func (p *GetCounterResult) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
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
			if fieldTypeId == thrift.I64 {
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

func (p *GetCounterResult) ReadField0(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v106, err107 := iprot.ReadI64()
	if err107 != nil {
		return thrift.NewTProtocolExceptionReadField(0, "success", p.ThriftName(), err107)
	}
	p.Success = v106
	return err
}

func (p *GetCounterResult) ReadFieldSuccess(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField0(iprot)
}

func (p *GetCounterResult) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("getCounter_result")
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

func (p *GetCounterResult) WriteField0(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("success", thrift.I64, 0)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(0, "success", p.ThriftName(), err)
	}
	err = oprot.WriteI64(int64(p.Success))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(0, "success", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(0, "success", p.ThriftName(), err)
	}
	return err
}

func (p *GetCounterResult) WriteFieldSuccess(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField0(oprot)
}

func (p *GetCounterResult) TStructName() string {
	return "GetCounterResult"
}

func (p *GetCounterResult) ThriftName() string {
	return "getCounter_result"
}

func (p *GetCounterResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("GetCounterResult(%+v)", *p)
}

func (p *GetCounterResult) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*GetCounterResult)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *GetCounterResult) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 0:
		return p.Success
	}
	return nil
}

func (p *GetCounterResult) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("success", thrift.I64, 0),
	})
}

/**
 * Attributes:
 *  - Key
 *  - Value
 */
type SetOptionArgs struct {
	thrift.TStruct
	Key   string "key"   // 1
	Value string "value" // 2
}

func NewSetOptionArgs() *SetOptionArgs {
	output := &SetOptionArgs{
		TStruct: thrift.NewTStruct("setOption_args", []thrift.TField{
			thrift.NewTField("key", thrift.STRING, 1),
			thrift.NewTField("value", thrift.STRING, 2),
		}),
	}
	{
	}
	return output
}

func (p *SetOptionArgs) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
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
		if fieldId == 1 || fieldName == "key" {
			if fieldTypeId == thrift.STRING {
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
		} else if fieldId == 2 || fieldName == "value" {
			if fieldTypeId == thrift.STRING {
				err = p.ReadField2(iprot)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else if fieldTypeId == thrift.VOID {
				err = iprot.Skip(fieldTypeId)
				if err != nil {
					return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
				}
			} else {
				err = p.ReadField2(iprot)
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

func (p *SetOptionArgs) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v108, err109 := iprot.ReadString()
	if err109 != nil {
		return thrift.NewTProtocolExceptionReadField(1, "key", p.ThriftName(), err109)
	}
	p.Key = v108
	return err
}

func (p *SetOptionArgs) ReadFieldKey(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *SetOptionArgs) ReadField2(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v110, err111 := iprot.ReadString()
	if err111 != nil {
		return thrift.NewTProtocolExceptionReadField(2, "value", p.ThriftName(), err111)
	}
	p.Value = v110
	return err
}

func (p *SetOptionArgs) ReadFieldValue(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField2(iprot)
}

func (p *SetOptionArgs) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("setOption_args")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
	}
	err = p.WriteField1(oprot)
	if err != nil {
		return err
	}
	err = p.WriteField2(oprot)
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

func (p *SetOptionArgs) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("key", thrift.STRING, 1)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "key", p.ThriftName(), err)
	}
	err = oprot.WriteString(string(p.Key))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "key", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "key", p.ThriftName(), err)
	}
	return err
}

func (p *SetOptionArgs) WriteFieldKey(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *SetOptionArgs) WriteField2(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("value", thrift.STRING, 2)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(2, "value", p.ThriftName(), err)
	}
	err = oprot.WriteString(string(p.Value))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(2, "value", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(2, "value", p.ThriftName(), err)
	}
	return err
}

func (p *SetOptionArgs) WriteFieldValue(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField2(oprot)
}

func (p *SetOptionArgs) TStructName() string {
	return "SetOptionArgs"
}

func (p *SetOptionArgs) ThriftName() string {
	return "setOption_args"
}

func (p *SetOptionArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("SetOptionArgs(%+v)", *p)
}

func (p *SetOptionArgs) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*SetOptionArgs)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *SetOptionArgs) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.Key
	case 2:
		return p.Value
	}
	return nil
}

func (p *SetOptionArgs) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("key", thrift.STRING, 1),
		thrift.NewTField("value", thrift.STRING, 2),
	})
}

type SetOptionResult struct {
	thrift.TStruct
}

func NewSetOptionResult() *SetOptionResult {
	output := &SetOptionResult{
		TStruct: thrift.NewTStruct("setOption_result", []thrift.TField{}),
	}
	{
	}
	return output
}

func (p *SetOptionResult) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
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
		err = iprot.Skip(fieldTypeId)
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
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

func (p *SetOptionResult) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("setOption_result")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
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

func (p *SetOptionResult) TStructName() string {
	return "SetOptionResult"
}

func (p *SetOptionResult) ThriftName() string {
	return "setOption_result"
}

func (p *SetOptionResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("SetOptionResult(%+v)", *p)
}

func (p *SetOptionResult) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*SetOptionResult)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *SetOptionResult) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	}
	return nil
}

func (p *SetOptionResult) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{})
}

/**
 * Attributes:
 *  - Key
 */
type GetOptionArgs struct {
	thrift.TStruct
	Key string "key" // 1
}

func NewGetOptionArgs() *GetOptionArgs {
	output := &GetOptionArgs{
		TStruct: thrift.NewTStruct("getOption_args", []thrift.TField{
			thrift.NewTField("key", thrift.STRING, 1),
		}),
	}
	{
	}
	return output
}

func (p *GetOptionArgs) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
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
		if fieldId == 1 || fieldName == "key" {
			if fieldTypeId == thrift.STRING {
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

func (p *GetOptionArgs) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v112, err113 := iprot.ReadString()
	if err113 != nil {
		return thrift.NewTProtocolExceptionReadField(1, "key", p.ThriftName(), err113)
	}
	p.Key = v112
	return err
}

func (p *GetOptionArgs) ReadFieldKey(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *GetOptionArgs) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("getOption_args")
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

func (p *GetOptionArgs) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("key", thrift.STRING, 1)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "key", p.ThriftName(), err)
	}
	err = oprot.WriteString(string(p.Key))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "key", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "key", p.ThriftName(), err)
	}
	return err
}

func (p *GetOptionArgs) WriteFieldKey(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *GetOptionArgs) TStructName() string {
	return "GetOptionArgs"
}

func (p *GetOptionArgs) ThriftName() string {
	return "getOption_args"
}

func (p *GetOptionArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("GetOptionArgs(%+v)", *p)
}

func (p *GetOptionArgs) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*GetOptionArgs)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *GetOptionArgs) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.Key
	}
	return nil
}

func (p *GetOptionArgs) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("key", thrift.STRING, 1),
	})
}

/**
 * Attributes:
 *  - Success
 */
type GetOptionResult struct {
	thrift.TStruct
	Success string "success" // 0
}

func NewGetOptionResult() *GetOptionResult {
	output := &GetOptionResult{
		TStruct: thrift.NewTStruct("getOption_result", []thrift.TField{
			thrift.NewTField("success", thrift.STRING, 0),
		}),
	}
	{
	}
	return output
}

func (p *GetOptionResult) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
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
			if fieldTypeId == thrift.STRING {
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

func (p *GetOptionResult) ReadField0(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v114, err115 := iprot.ReadString()
	if err115 != nil {
		return thrift.NewTProtocolExceptionReadField(0, "success", p.ThriftName(), err115)
	}
	p.Success = v114
	return err
}

func (p *GetOptionResult) ReadFieldSuccess(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField0(iprot)
}

func (p *GetOptionResult) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("getOption_result")
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

func (p *GetOptionResult) WriteField0(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("success", thrift.STRING, 0)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(0, "success", p.ThriftName(), err)
	}
	err = oprot.WriteString(string(p.Success))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(0, "success", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(0, "success", p.ThriftName(), err)
	}
	return err
}

func (p *GetOptionResult) WriteFieldSuccess(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField0(oprot)
}

func (p *GetOptionResult) TStructName() string {
	return "GetOptionResult"
}

func (p *GetOptionResult) ThriftName() string {
	return "getOption_result"
}

func (p *GetOptionResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("GetOptionResult(%+v)", *p)
}

func (p *GetOptionResult) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*GetOptionResult)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *GetOptionResult) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 0:
		return p.Success
	}
	return nil
}

func (p *GetOptionResult) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("success", thrift.STRING, 0),
	})
}

type GetOptionsArgs struct {
	thrift.TStruct
}

func NewGetOptionsArgs() *GetOptionsArgs {
	output := &GetOptionsArgs{
		TStruct: thrift.NewTStruct("getOptions_args", []thrift.TField{}),
	}
	{
	}
	return output
}

func (p *GetOptionsArgs) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
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
		err = iprot.Skip(fieldTypeId)
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
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

func (p *GetOptionsArgs) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("getOptions_args")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
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

func (p *GetOptionsArgs) TStructName() string {
	return "GetOptionsArgs"
}

func (p *GetOptionsArgs) ThriftName() string {
	return "getOptions_args"
}

func (p *GetOptionsArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("GetOptionsArgs(%+v)", *p)
}

func (p *GetOptionsArgs) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*GetOptionsArgs)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *GetOptionsArgs) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	}
	return nil
}

func (p *GetOptionsArgs) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{})
}

/**
 * Attributes:
 *  - Success
 */
type GetOptionsResult struct {
	thrift.TStruct
	Success thrift.TMap "success" // 0
}

func NewGetOptionsResult() *GetOptionsResult {
	output := &GetOptionsResult{
		TStruct: thrift.NewTStruct("getOptions_result", []thrift.TField{
			thrift.NewTField("success", thrift.MAP, 0),
		}),
	}
	{
	}
	return output
}

func (p *GetOptionsResult) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
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
			if fieldTypeId == thrift.MAP {
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

func (p *GetOptionsResult) ReadField0(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	_ktype119, _vtype120, _size118, err := iprot.ReadMapBegin()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "p.Success", "", err)
	}
	p.Success = thrift.NewTMap(_ktype119, _vtype120, _size118)
	for _i122 := 0; _i122 < _size118; _i122++ {
		v125, err126 := iprot.ReadString()
		if err126 != nil {
			return thrift.NewTProtocolExceptionReadField(0, "_key123", "", err126)
		}
		_key123 := v125
		v127, err128 := iprot.ReadString()
		if err128 != nil {
			return thrift.NewTProtocolExceptionReadField(0, "_val124", "", err128)
		}
		_val124 := v127
		p.Success.Set(_key123, _val124)
	}
	err = iprot.ReadMapEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionReadField(-1, "", "map", err)
	}
	return err
}

func (p *GetOptionsResult) ReadFieldSuccess(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField0(iprot)
}

func (p *GetOptionsResult) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("getOptions_result")
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

func (p *GetOptionsResult) WriteField0(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	if p.Success != nil {
		err = oprot.WriteFieldBegin("success", thrift.MAP, 0)
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(0, "success", p.ThriftName(), err)
		}
		err = oprot.WriteMapBegin(thrift.STRING, thrift.STRING, p.Success.Len())
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(-1, "", "map", err)
		}
		for Miter129 := range p.Success.Iter() {
			Kiter130, Viter131 := Miter129.Key().(string), Miter129.Value().(string)
			err = oprot.WriteString(string(Kiter130))
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(0, "Kiter130", "", err)
			}
			err = oprot.WriteString(string(Viter131))
			if err != nil {
				return thrift.NewTProtocolExceptionWriteField(0, "Viter131", "", err)
			}
		}
		err = oprot.WriteMapEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(-1, "", "map", err)
		}
		err = oprot.WriteFieldEnd()
		if err != nil {
			return thrift.NewTProtocolExceptionWriteField(0, "success", p.ThriftName(), err)
		}
	}
	return err
}

func (p *GetOptionsResult) WriteFieldSuccess(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField0(oprot)
}

func (p *GetOptionsResult) TStructName() string {
	return "GetOptionsResult"
}

func (p *GetOptionsResult) ThriftName() string {
	return "getOptions_result"
}

func (p *GetOptionsResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("GetOptionsResult(%+v)", *p)
}

func (p *GetOptionsResult) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*GetOptionsResult)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *GetOptionsResult) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 0:
		return p.Success
	}
	return nil
}

func (p *GetOptionsResult) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("success", thrift.MAP, 0),
	})
}

/**
 * Attributes:
 *  - ProfileDurationInSec
 */
type GetCpuProfileArgs struct {
	thrift.TStruct
	ProfileDurationInSec int32 "profileDurationInSec" // 1
}

func NewGetCpuProfileArgs() *GetCpuProfileArgs {
	output := &GetCpuProfileArgs{
		TStruct: thrift.NewTStruct("getCpuProfile_args", []thrift.TField{
			thrift.NewTField("profileDurationInSec", thrift.I32, 1),
		}),
	}
	{
	}
	return output
}

func (p *GetCpuProfileArgs) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
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
		if fieldId == 1 || fieldName == "profileDurationInSec" {
			if fieldTypeId == thrift.I32 {
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

func (p *GetCpuProfileArgs) ReadField1(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v132, err133 := iprot.ReadI32()
	if err133 != nil {
		return thrift.NewTProtocolExceptionReadField(1, "profileDurationInSec", p.ThriftName(), err133)
	}
	p.ProfileDurationInSec = v132
	return err
}

func (p *GetCpuProfileArgs) ReadFieldProfileDurationInSec(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField1(iprot)
}

func (p *GetCpuProfileArgs) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("getCpuProfile_args")
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

func (p *GetCpuProfileArgs) WriteField1(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("profileDurationInSec", thrift.I32, 1)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "profileDurationInSec", p.ThriftName(), err)
	}
	err = oprot.WriteI32(int32(p.ProfileDurationInSec))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "profileDurationInSec", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(1, "profileDurationInSec", p.ThriftName(), err)
	}
	return err
}

func (p *GetCpuProfileArgs) WriteFieldProfileDurationInSec(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField1(oprot)
}

func (p *GetCpuProfileArgs) TStructName() string {
	return "GetCpuProfileArgs"
}

func (p *GetCpuProfileArgs) ThriftName() string {
	return "getCpuProfile_args"
}

func (p *GetCpuProfileArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("GetCpuProfileArgs(%+v)", *p)
}

func (p *GetCpuProfileArgs) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*GetCpuProfileArgs)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *GetCpuProfileArgs) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 1:
		return p.ProfileDurationInSec
	}
	return nil
}

func (p *GetCpuProfileArgs) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("profileDurationInSec", thrift.I32, 1),
	})
}

/**
 * Attributes:
 *  - Success
 */
type GetCpuProfileResult struct {
	thrift.TStruct
	Success string "success" // 0
}

func NewGetCpuProfileResult() *GetCpuProfileResult {
	output := &GetCpuProfileResult{
		TStruct: thrift.NewTStruct("getCpuProfile_result", []thrift.TField{
			thrift.NewTField("success", thrift.STRING, 0),
		}),
	}
	{
	}
	return output
}

func (p *GetCpuProfileResult) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
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
			if fieldTypeId == thrift.STRING {
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

func (p *GetCpuProfileResult) ReadField0(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v134, err135 := iprot.ReadString()
	if err135 != nil {
		return thrift.NewTProtocolExceptionReadField(0, "success", p.ThriftName(), err135)
	}
	p.Success = v134
	return err
}

func (p *GetCpuProfileResult) ReadFieldSuccess(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField0(iprot)
}

func (p *GetCpuProfileResult) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("getCpuProfile_result")
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

func (p *GetCpuProfileResult) WriteField0(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("success", thrift.STRING, 0)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(0, "success", p.ThriftName(), err)
	}
	err = oprot.WriteString(string(p.Success))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(0, "success", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(0, "success", p.ThriftName(), err)
	}
	return err
}

func (p *GetCpuProfileResult) WriteFieldSuccess(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField0(oprot)
}

func (p *GetCpuProfileResult) TStructName() string {
	return "GetCpuProfileResult"
}

func (p *GetCpuProfileResult) ThriftName() string {
	return "getCpuProfile_result"
}

func (p *GetCpuProfileResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("GetCpuProfileResult(%+v)", *p)
}

func (p *GetCpuProfileResult) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*GetCpuProfileResult)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *GetCpuProfileResult) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 0:
		return p.Success
	}
	return nil
}

func (p *GetCpuProfileResult) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("success", thrift.STRING, 0),
	})
}

type AliveSinceArgs struct {
	thrift.TStruct
}

func NewAliveSinceArgs() *AliveSinceArgs {
	output := &AliveSinceArgs{
		TStruct: thrift.NewTStruct("aliveSince_args", []thrift.TField{}),
	}
	{
	}
	return output
}

func (p *AliveSinceArgs) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
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
		err = iprot.Skip(fieldTypeId)
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
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

func (p *AliveSinceArgs) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("aliveSince_args")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
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

func (p *AliveSinceArgs) TStructName() string {
	return "AliveSinceArgs"
}

func (p *AliveSinceArgs) ThriftName() string {
	return "aliveSince_args"
}

func (p *AliveSinceArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("AliveSinceArgs(%+v)", *p)
}

func (p *AliveSinceArgs) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*AliveSinceArgs)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *AliveSinceArgs) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	}
	return nil
}

func (p *AliveSinceArgs) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{})
}

/**
 * Attributes:
 *  - Success
 */
type AliveSinceResult struct {
	thrift.TStruct
	Success int64 "success" // 0
}

func NewAliveSinceResult() *AliveSinceResult {
	output := &AliveSinceResult{
		TStruct: thrift.NewTStruct("aliveSince_result", []thrift.TField{
			thrift.NewTField("success", thrift.I64, 0),
		}),
	}
	{
	}
	return output
}

func (p *AliveSinceResult) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
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
			if fieldTypeId == thrift.I64 {
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

func (p *AliveSinceResult) ReadField0(iprot thrift.TProtocol) (err thrift.TProtocolException) {
	v136, err137 := iprot.ReadI64()
	if err137 != nil {
		return thrift.NewTProtocolExceptionReadField(0, "success", p.ThriftName(), err137)
	}
	p.Success = v136
	return err
}

func (p *AliveSinceResult) ReadFieldSuccess(iprot thrift.TProtocol) thrift.TProtocolException {
	return p.ReadField0(iprot)
}

func (p *AliveSinceResult) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("aliveSince_result")
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

func (p *AliveSinceResult) WriteField0(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteFieldBegin("success", thrift.I64, 0)
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(0, "success", p.ThriftName(), err)
	}
	err = oprot.WriteI64(int64(p.Success))
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(0, "success", p.ThriftName(), err)
	}
	err = oprot.WriteFieldEnd()
	if err != nil {
		return thrift.NewTProtocolExceptionWriteField(0, "success", p.ThriftName(), err)
	}
	return err
}

func (p *AliveSinceResult) WriteFieldSuccess(oprot thrift.TProtocol) thrift.TProtocolException {
	return p.WriteField0(oprot)
}

func (p *AliveSinceResult) TStructName() string {
	return "AliveSinceResult"
}

func (p *AliveSinceResult) ThriftName() string {
	return "aliveSince_result"
}

func (p *AliveSinceResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("AliveSinceResult(%+v)", *p)
}

func (p *AliveSinceResult) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*AliveSinceResult)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *AliveSinceResult) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	case 0:
		return p.Success
	}
	return nil
}

func (p *AliveSinceResult) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{
		thrift.NewTField("success", thrift.I64, 0),
	})
}

type ReinitializeArgs struct {
	thrift.TStruct
}

func NewReinitializeArgs() *ReinitializeArgs {
	output := &ReinitializeArgs{
		TStruct: thrift.NewTStruct("reinitialize_args", []thrift.TField{}),
	}
	{
	}
	return output
}

func (p *ReinitializeArgs) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
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
		err = iprot.Skip(fieldTypeId)
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
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

func (p *ReinitializeArgs) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("reinitialize_args")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
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

func (p *ReinitializeArgs) TStructName() string {
	return "ReinitializeArgs"
}

func (p *ReinitializeArgs) ThriftName() string {
	return "reinitialize_args"
}

func (p *ReinitializeArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ReinitializeArgs(%+v)", *p)
}

func (p *ReinitializeArgs) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*ReinitializeArgs)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *ReinitializeArgs) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	}
	return nil
}

func (p *ReinitializeArgs) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{})
}

type ReinitializeResult struct {
	thrift.TStruct
}

func NewReinitializeResult() *ReinitializeResult {
	output := &ReinitializeResult{
		TStruct: thrift.NewTStruct("reinitialize_result", []thrift.TField{}),
	}
	{
	}
	return output
}

func (p *ReinitializeResult) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
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
		err = iprot.Skip(fieldTypeId)
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
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

func (p *ReinitializeResult) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("reinitialize_result")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
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

func (p *ReinitializeResult) TStructName() string {
	return "ReinitializeResult"
}

func (p *ReinitializeResult) ThriftName() string {
	return "reinitialize_result"
}

func (p *ReinitializeResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ReinitializeResult(%+v)", *p)
}

func (p *ReinitializeResult) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*ReinitializeResult)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *ReinitializeResult) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	}
	return nil
}

func (p *ReinitializeResult) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{})
}

type ShutdownArgs struct {
	thrift.TStruct
}

func NewShutdownArgs() *ShutdownArgs {
	output := &ShutdownArgs{
		TStruct: thrift.NewTStruct("shutdown_args", []thrift.TField{}),
	}
	{
	}
	return output
}

func (p *ShutdownArgs) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
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
		err = iprot.Skip(fieldTypeId)
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
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

func (p *ShutdownArgs) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("shutdown_args")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
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

func (p *ShutdownArgs) TStructName() string {
	return "ShutdownArgs"
}

func (p *ShutdownArgs) ThriftName() string {
	return "shutdown_args"
}

func (p *ShutdownArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ShutdownArgs(%+v)", *p)
}

func (p *ShutdownArgs) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*ShutdownArgs)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *ShutdownArgs) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	}
	return nil
}

func (p *ShutdownArgs) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{})
}

type ShutdownResult struct {
	thrift.TStruct
}

func NewShutdownResult() *ShutdownResult {
	output := &ShutdownResult{
		TStruct: thrift.NewTStruct("shutdown_result", []thrift.TField{}),
	}
	{
	}
	return output
}

func (p *ShutdownResult) Read(iprot thrift.TProtocol) (err thrift.TProtocolException) {
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
		err = iprot.Skip(fieldTypeId)
		if err != nil {
			return thrift.NewTProtocolExceptionReadField(int(fieldId), fieldName, p.ThriftName(), err)
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

func (p *ShutdownResult) Write(oprot thrift.TProtocol) (err thrift.TProtocolException) {
	err = oprot.WriteStructBegin("shutdown_result")
	if err != nil {
		return thrift.NewTProtocolExceptionWriteStruct(p.ThriftName(), err)
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

func (p *ShutdownResult) TStructName() string {
	return "ShutdownResult"
}

func (p *ShutdownResult) ThriftName() string {
	return "shutdown_result"
}

func (p *ShutdownResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ShutdownResult(%+v)", *p)
}

func (p *ShutdownResult) CompareTo(other interface{}) (int, bool) {
	if other == nil {
		return 1, true
	}
	data, ok := other.(*ShutdownResult)
	if !ok {
		return 0, false
	}
	return thrift.TType(thrift.STRUCT).Compare(p, data)
}

func (p *ShutdownResult) AttributeByFieldId(id int) interface{} {
	switch id {
	default:
		return nil
	}
	return nil
}

func (p *ShutdownResult) TStructFields() thrift.TFieldContainer {
	return thrift.NewTFieldContainer([]thrift.TField{})
}
