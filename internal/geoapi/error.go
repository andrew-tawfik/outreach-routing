package geoapi

import (
	"fmt"
	"strings"
)

type ApiErrors struct {
	FailedGuests []FailedGuest
}

type FailedGuest struct {
	Name    string
	Address string
	Reason  string
}

func (ae *ApiErrors) HasErrors() bool {
	return len(ae.FailedGuests) > 0
}

func (ae *ApiErrors) GetSummary() string {
	if !ae.HasErrors() {
		return ""
	}

	if len(ae.FailedGuests) == 1 {
		return fmt.Sprintf("Could not find address for %s at this time, please add guest manually ", ae.FailedGuests[0].Name)
	}

	return fmt.Sprintf("Could not find addresses for %d guests, please add them manually", len(ae.FailedGuests))
}

func (ae *ApiErrors) GetDetails() string {
	if len(ae.FailedGuests) == 0 {
		return ""
	}

	var details strings.Builder
	details.WriteString("The following guests will be excluded from routing:\n\n")

	for i, fg := range ae.FailedGuests {
		if i > 0 {
			details.WriteString("\n")
		}
		details.WriteString(fmt.Sprintf("• %s\n  Address: %s\n  Reason: %s",
			fg.Name, fg.Address, fg.Reason))
	}

	return details.String()
}
