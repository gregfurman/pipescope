package gitlab

type Pipeline struct {
	ID        int
	ProjectID int

	CommitSha string
	Status    string
	URL       string
}

func IsStatusPending(status string) bool {
	switch status {
	case "created", "waiting_for_resource", "preparing", "pending", "running":
		return true
	}

	return false
}
