package hosts

type HostDatabase struct {
	Host string `json:"host"`
}

func (h HostDatabase) GetHosts() ([]Host, error) {
	return []Host{}, nil
}

func (h HostDatabase) AddHost(Host) error {
	return nil
}

func (h HostDatabase) DeleteHost(string) error {
	return nil
}
