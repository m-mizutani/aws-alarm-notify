package main

type cloudWatchAlarm struct {
	AWSAccountID     *string            `json:"AWSAccountId,omitempty"`
	AlarmDescription *string            `json:"AlarmDescription,omitempty"`
	AlarmName        *string            `json:"AlarmName,omitempty"`
	NewStateReason   *string            `json:"NewStateReason,omitempty"`
	NewStateValue    *string            `json:"NewStateValue,omitempty"`
	OldStateValue    *string            `json:"OldStateValue,omitempty"`
	Region           *string            `json:"Region,omitempty"`
	StateChangeTime  *string            `json:"StateChangeTime,omitempty"`
	Trigger          *cloudWatchTrigger `json:"Trigger,omitempty"`
}

type cloudWatchTrigger struct {
	ComparisonOperator               *string                `json:"ComparisonOperator,omitempty"`
	Dimensions                       []*cloudWatchDemension `json:"Dimensions,omitempty"`
	EvaluateLowSampleCountPercentile *string                `json:"EvaluateLowSampleCountPercentile,omitempty"`
	EvaluationPeriods                *float64               `json:"EvaluationPeriods,omitempty"`
	MetricName                       *string                `json:"MetricName,omitempty"`
	Namespace                        *string                `json:"Namespace,omitempty"`
	Period                           *float64               `json:"Period,omitempty"`
	Statistic                        *string                `json:"Statistic,omitempty"`
	StatisticType                    *string                `json:"StatisticType,omitempty"`
	Threshold                        *float64               `json:"Threshold,omitempty"`
	TreatMissingData                 *string                `json:"TreatMissingData,omitempty"`
	// Unit                             interface{} `json:"Unit,omitempty"`
}

type cloudWatchDemension struct {
	Name  *string     `json:"name,omitempty"`
	Value interface{} `json:"value,omitempty"`
}
