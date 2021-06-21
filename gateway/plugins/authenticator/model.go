package main

type (
	config struct {
		services    map[string]Service
		method      string
		header      Header
		responseMap map[string]string
		cacheTTL    int
		timeout     int
		whitelist   string
	}

	Header struct {
		Identity string `json:"authorization"`
		Service  string `json:"service_id"`
	}

	Service struct {
		ValidateUrl string `json:"validate_url"`
		LoginPath   string `json:"login_path"`
		LogoutPath  string `json:"logout_path"`
	}

	data struct {
		Payload map[string]string `json:"data"`
	}
)
