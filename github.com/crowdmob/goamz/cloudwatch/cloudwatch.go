/***** BEGIN LICENSE BLOCK *****
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this file,
# You can obtain one at http://mozilla.org/MPL/2.0/.
#
# The Initial Developer of the Original Code is the Mozilla Foundation.
# Portions created by the Initial Developer are Copyright (C) 2012
# the Initial Developer. All Rights Reserved.
#
# Contributor(s):
#   Ben Bangert (bbangert@mozilla.com)
#
# ***** END LICENSE BLOCK *****/

package cloudwatch

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/crowdmob/goamz/aws"
	"github.com/feyeleanor/sets"
	"strconv"
	"time"
)

// The CloudWatch type encapsulates all the CloudWatch operations in a region.
type CloudWatch struct {
	Service   aws.AWSService
	namespace string
}

type Dimension struct {
	Name  string
	Value string
}

type StatisticSet struct {
	Maximum     float64
	Minimum     float64
	SampleCount float64
	Sum         float64
}

type MetricDatum struct {
	Dimensions      []Dimension
	MetricName      string
	StatisticValues []StatisticSet
	Timestamp       time.Time
	Unit            string
	Value           float64
}

type Datapoint struct {
	Average     float64
	Maximum     float64
	Minimum     float64
	SampleCount float64
	Sum         float64
	Timestamp   time.Time
	Unit        string
}

type GetMetricStatisticsRequest struct {
	Dimensions []Dimension
	EndTime    time.Time
	StartTime  time.Time
	MetricName string
	Unit       string
	Period     int
	Statistics []string
}

type GetMetricStatisticsResult struct {
	Datapoints []Datapoint `xml:"Datapoints>member"`
	Label      string
}

type GetMetricStatisticsResponse struct {
	GetMetricStatisticsResult GetMetricStatisticsResult
	ResponseMetadata          aws.ResponseMetadata
}

var attempts = aws.AttemptStrategy{
	Min:   5,
	Total: 5 * time.Second,
	Delay: 200 * time.Millisecond,
}

var validUnits = sets.SSet(
	"Seconds",
	"Microseconds",
	"Milliseconds",
	"Bytes",
	"Kilobytes",
	"Megabytes",
	"Gigabytes",
	"Terabytes",
	"Bits",
	"Kilobits",
	"Megabits",
	"Gigabits",
	"Terabits",
	"Percent",
	"Count",
	"Bytes/Second",
	"Kilobytes/Second",
	"Megabytes/Second",
	"Gigabytes/Second",
	"Terabytes/Second",
	"Bits/Second",
	"Kilobits/Second",
	"Megabits/Second",
	"Gigabits/Second",
	"Terabits/Second",
	"Count/Second",
)

var validMetricStatistics = sets.SSet(
	"Average",
	"Sum",
	"SampleCount",
	"Maximum",
	"Minimum",
)

// Create a new CloudWatch object for a given namespace
func NewCloudWatch(auth aws.Auth, region aws.ServiceInfo, namespace string) (*CloudWatch, error) {
	service, err := aws.NewService(auth, region)
	if err != nil {
		return nil, err
	}
	return &CloudWatch{
		Service:   service,
		namespace: namespace,
	}, nil
}

func (c *CloudWatch) query(method, path string, params map[string]string, resp interface{}) error {
	// Add basic Cloudwatch param
	params["Version"] = "2010-08-01"
	params["Namespace"] = c.namespace

	r, err := c.Service.Query(method, path, params)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		return c.Service.BuildError(r)
	}
	err = xml.NewDecoder(r.Body).Decode(resp)
	return err
}

// Get statistics for specified metric
//
// If the arguments are invalid or the server returns an error, the error will
// be set and the other values undefined.
func (c *CloudWatch) GetMetricStatistics(req *GetMetricStatisticsRequest) (result *GetMetricStatisticsResponse, err error) {
	statisticsSet := sets.SSet(req.Statistics...)
	// Kick out argument errors
	switch {
	case req.EndTime.IsZero():
		err = errors.New("No endTime specified")
	case req.StartTime.IsZero():
		err = errors.New("No startTime specified")
	case req.MetricName == "":
		err = errors.New("No metricName specified")
	case req.Period < 60 || req.Period%60 != 0:
		err = errors.New("Period not 60 seconds or a multiple of 60 seconds")
	case len(req.Statistics) < 1:
		err = errors.New("No statistics supplied")
	case validMetricStatistics.Union(statisticsSet).Len() != validMetricStatistics.Len():
		err = errors.New("Invalid statistic values supplied")
	case req.Unit != "" && !validUnits.Member(req.Unit):
		err = errors.New("Unit is not a valid value")
	}
	if err != nil {
		return
	}

	// Serialize all the params
	params := aws.MakeParams("GetMetricStatistics")
	params["EndTime"] = req.EndTime.UTC().Format(time.RFC3339)
	params["StartTime"] = req.StartTime.UTC().Format(time.RFC3339)
	params["MetricName"] = req.MetricName
	params["Period"] = strconv.Itoa(req.Period)
	if req.Unit != "" {
		params["Unit"] = req.Unit
	}

	// Serialize the lists of data
	for i, d := range req.Dimensions {
		prefix := "Dimensions.member." + strconv.Itoa(i+1)
		params[prefix+".Name"] = d.Name
		params[prefix+".Value"] = d.Value
	}
	for i, d := range req.Statistics {
		prefix := "Statistics.member." + strconv.Itoa(i+1)
		params[prefix] = d
	}
	result = new(GetMetricStatisticsResponse)
	err = c.query("GET", "/", params, result)
	return
}

func (c *CloudWatch) PutMetricData(metrics []MetricDatum) (result *aws.BaseResponse, err error) {
	// Serialize the params
	params := aws.MakeParams("PutMetricData")
	for i, metric := range metrics {
		prefix := "MetricData.member." + strconv.Itoa(i+1)
		if metric.MetricName == "" {
			err = fmt.Errorf("No metric name supplied for metric: %d", i)
			return
		}
		params[prefix+".MetricName"] = metric.MetricName
		if metric.Unit != "" {
			params[prefix+".Unit"] = metric.Unit
		}
		if metric.Value != 0 {
			params[prefix+".Value"] = strconv.FormatFloat(metric.Value, 'E', 10, 64)
		}
		if !metric.Timestamp.IsZero() {
			params[prefix+".Timestamp"] = metric.Timestamp.UTC().Format(time.RFC3339)
		}
		for j, dim := range metric.Dimensions {
			dimprefix := prefix + "Dimensions.member." + strconv.Itoa(j+1)
			params[dimprefix+".Name"] = dim.Name
			params[dimprefix+".Value"] = dim.Value
		}
		for j, stat := range metric.StatisticValues {
			statprefix := prefix + "StatisticValues.member." + strconv.Itoa(j+1)
			params[statprefix+".Maximum"] = strconv.FormatFloat(stat.Maximum, 'E', 10, 64)
			params[statprefix+".Minimum"] = strconv.FormatFloat(stat.Minimum, 'E', 10, 64)
			params[statprefix+".SampleCount"] = strconv.FormatFloat(stat.SampleCount, 'E', 10, 64)
			params[statprefix+".Sum"] = strconv.FormatFloat(stat.Sum, 'E', 10, 64)
		}
	}
	result = new(aws.BaseResponse)
	err = c.query("POST", "/", params, result)
	return
}
