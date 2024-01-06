package hosts

type HostDatabaseMysql struct {
	Host string `json:"host"`
}

func (h *HostDatabaseMysql) GetHosts() ([]Host, error) {
	return []Host{}, nil
}

func (h *HostDatabaseMysql) AddHost(Host) error {
	return nil
}

func (h *HostDatabaseMysql) DeleteHost(string) error {
	return nil
}
