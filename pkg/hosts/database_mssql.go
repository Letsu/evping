package hosts

type HostDatabaseMssql struct {
	Host string `json:"host"`
}

func (h *HostDatabaseMssql) GetHosts() ([]Host, error) {
	return []Host{}, nil
}

func (h *HostDatabaseMssql) AddHost(Host) error {
	return nil
}

func (h *HostDatabaseMssql) DeleteHost(string) error {
	return nil
}
