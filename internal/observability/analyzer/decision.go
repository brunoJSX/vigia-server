package analyzer

// Decision is the outcome of interpreting collected data — it drives what
// happens to an Incident. UpdateIncident was considered and dropped: no
// source defines what "updating" an Incident would mean.
type Decision string

const (
	DecisionNoAction        Decision = "no_action"
	DecisionOpenIncident    Decision = "open_incident"
	DecisionResolveIncident Decision = "resolve_incident"
)
