package pitlog

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

type Pitlog_base struct {
	log_title          string
	dedicated          *logrus.Logger
	log_file_name      string
	log_directory      string
	file               *os.File
	app_name           string
	app_version        string
	app_level          string
	indentation        string
	enable_log_console bool
	object_view        bool
	use_sparate        bool
}

type Log_model struct {
	Mode   string
	Detail string
	Header string
	Body   string
}

type Request_detail struct {
	Url          string `json:"url"`
	Method       string `json:"method"`
	Request_id   string `json:"request_id"`
	Request_date string `json:"request_date"`
	Ip           string `json:"ip"`
	Api_type     string `json:"api_type"`
	Channel      string `json:"channel"`
	Timezone     string `json:"timezone"`
	Bytes_in     string `json:"bytes_in"`
	App_level    string `json:"App_level"`
	App_version  string `json:"App_version"`
}

type Response_detail struct {
	Request_detail
	Bytes_out     string        `json:"bytes_out"`
	Error         string        `json:"error"`
	Latency       time.Duration `json:"latency"`
	Latency_human string        `json:"latency_human"`
	Status        int           `json:"status"`
}

type Request_logger struct {
	Request_id string      `json:"requestID,omitempty"`
	Payload    interface{} `json:"payload,omitempty"`
}

type Pit_json_message struct {
	Request_id string `json:"request_id"`
	Status     string `json:"status"`
	Title      string `json:"title"`
	Message    string `json:"message"`
	Date       string `json:"date"`
}
