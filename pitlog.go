package pitlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"mime/multipart"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	name_app_level_development   = "development"
	name_app_level_production    = "production"
	name_type_status_info        = "INFO"
	name_type_status_debug       = "DEBUG"
	name_type_status_error       = "ERROR"
	name_type_status_log         = "LOG"
	base_format_date_time_logger = "2006-01-02 15:04:05.999 -07:00"
	name_type_request            = "Request"
	name_type_respon             = "Response"
	name_content_type            = "Content-Type"
)
var (
	sensitive_keys map[string]bool
)

type custom_string_formatter struct{}

func (f *custom_string_formatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message), nil
}

type IN_pitlog_base interface {
	Api_log_middleware(ech *echo.Echo, sensitive []string)
	Make_log_string(request_id string, status int, title string, message string)
	Make_log_object(request_id string, status int, title string, object interface{})
}

func New_pitlog(app_name, app_version, app_level, log_dir, enable_log_console, use_separate, object_view string) (IN_pitlog_base, error) {
	var err error
	dedicated_log := logrus.New()
	dedicated_log.SetLevel(logrus.InfoLevel)
	dedicated_log.SetFormatter(&custom_string_formatter{})

	log_title_default := "pitlog.go"
	log_file_name_default, _ := get_log_file_name(app_name, log_dir)

	var file *os.File
	if file, err = os.OpenFile(log_file_name_default, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err != nil {
		return nil, err
	}

	indentation := "\n"
	if app_level == name_app_level_development {
		indentation = "\n\n"
	}

	enable_log_console_bool, _ := strconv.ParseBool(enable_log_console)
	use_separate_bool, _ := strconv.ParseBool(use_separate)
	object_view_bool, _ := strconv.ParseBool(object_view)
	return &Pitlog_base{
		log_title_default,
		dedicated_log,
		log_file_name_default,
		log_dir,
		file,
		app_name,
		app_version,
		app_level,
		indentation,
		enable_log_console_bool,
		object_view_bool,
		use_separate_bool,
	}, err
}

func get_log_file_name(app_name, log_dir string) (string, error) {
	var err error
	var log_file_name string
	log_file_name = fmt.Sprintf("%s/%s_LOG.log", log_dir, strings.ToUpper(app_name))
	_ = os.Mkdir(log_dir, os.ModePerm) // Buat direktori logs jika belum ada
	return log_file_name, err
}

func (DI *Pitlog_base) Api_log_middleware(ech *echo.Echo, sensitive []string) {

	sensitive_keys = make(map[string]bool)
	for _, key := range sensitive {
		sensitive_keys[key] = true
	}

	ech.Use(middleware.RequestID())
	ech.Use(DI.logger_middleware())
	ech.Use(middleware.BodyDumpWithConfig(DI.custom_dump_body()))
}

func (DI *Pitlog_base) logger_middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request_data := DI.get_log_request_data(c)
			request_log := DI.create_api_loger(request_data, false)
			DI.dedicated.Info(request_log)
			return next(c)
		}
	}
}

func (DI *Pitlog_base) custom_dump_body() middleware.BodyDumpConfig {
	return middleware.BodyDumpConfig{
		Handler: func(c echo.Context, reqBody, resBody []byte) {
			response_data := DI.get_log_response(c, resBody)
			response_log := DI.create_api_loger(response_data, true)
			DI.dedicated.Info(response_log)
			// Simpan constant body dalam variabel di context
			c.Set("responseBody", resBody)
		},
	}
}

func (DI *Pitlog_base) get_log_response(c echo.Context, resBody []byte) Log_model {
	var err error
	var result Log_model
	start := time.Now()
	bytesOut := fmt.Sprint(c.Response().Size)
	latency := time.Since(start)
	latencyHuman := latency.String()
	status := c.Response().Status
	request_id := c.Response().Header().Get(echo.HeaderXRequestID)
	uri := c.Request().RequestURI
	bytesIn := c.Request().Header.Get(echo.HeaderContentLength)
	timezone, _ := time.Now().Local().Zone()
	method := c.Request().Method

	var response_detail_data Response_detail
	response_detail_data.Request_date = time.Now().Format(base_format_date_time_logger)
	response_detail_data.Request_id = request_id
	response_detail_data.Method = method
	response_detail_data.Url = uri
	response_detail_data.Ip = c.RealIP()
	response_detail_data.Api_type = name_type_respon
	response_detail_data.Channel = DI.app_name
	response_detail_data.Timezone = timezone
	response_detail_data.Bytes_in = bytesIn
	response_detail_data.Bytes_out = bytesOut
	//response_detail_data.Error = err.Error()
	response_detail_data.Latency = latency
	response_detail_data.Latency_human = latencyHuman
	response_detail_data.Status = status
	response_detail_data.App_level = DI.app_level
	response_detail_data.App_version = DI.app_version

	response_detail_json, _ := json.Marshal(response_detail_data)

	response_header_data := c.Request().Header
	response_header_json, _ := json.Marshal(response_header_data)

	compactResBody := new(bytes.Buffer)
	if err = json.Compact(compactResBody, resBody); err != nil {
		DI.Make_log_string(request_id, 3, DI.log_title, err.Error())
		return result
	}

	var body_json map[string]interface{}
	if err = json.Unmarshal(resBody, &body_json); err != nil {
		DI.Make_log_string(request_id, 3, DI.log_title, err.Error())
		return result
	}

	DI.modify_sensitive_fields(body_json)
	modified_body_json, _ := json.Marshal(body_json)
	result.Mode = name_type_respon
	result.Detail = string(response_detail_json)
	result.Header = string(response_header_json)
	//result.Body = compactResBody.String()
	result.Body = string(modified_body_json)
	return result
}

func (DI *Pitlog_base) get_log_request_data(c echo.Context) Log_model {
	var result Log_model
	request_id := c.Response().Header().Get(echo.HeaderXRequestID)
	uri := c.Request().RequestURI
	bytesIn := c.Request().Header.Get(echo.HeaderContentLength)
	timezone, _ := time.Now().Local().Zone()
	method := c.Request().Method
	request_headers := make(map[string]string)
	for key, values := range c.Request().Header {
		request_headers[key] = strings.Join(values, ", ")
	}

	request_header_data := c.Request().Header
	request_header_json, _ := json.Marshal(request_header_data)
	contentType := c.Request().Header.Get(name_content_type)
	request_body_string := "{}"
	switch true {
	case strings.Contains(contentType, echo.MIMEApplicationJSON):
		request_body_string = DI.get_request_body_json_string(c)
		break
	case strings.Contains(contentType, echo.MIMEApplicationForm):
		request_body_string = DI.get_request_body_form_encoded_string(c)
		break
	case strings.Contains(contentType, echo.MIMEMultipartForm):
		request_body_string = DI.get_request_body_form_data(c)
		break
	}

	var request_detail_model Request_detail
	request_detail_model.Request_date = time.Now().Format(base_format_date_time_logger)
	request_detail_model.Request_id = request_id
	request_detail_model.Method = method
	request_detail_model.Url = uri
	request_detail_model.Ip = c.RealIP()
	request_detail_model.Api_type = name_type_request
	request_detail_model.Channel = DI.app_name
	request_detail_model.Timezone = timezone
	request_detail_model.Bytes_in = bytesIn
	request_detail_model.App_level = DI.app_level
	request_detail_model.App_version = DI.app_version

	request_detail_json, _ := json.Marshal(request_detail_model)

	result.Mode = name_type_request
	result.Detail = string(request_detail_json)
	result.Header = string(request_header_json)
	result.Body = request_body_string
	return result
}

func (DI *Pitlog_base) create_api_loger(logModel Log_model, is_response bool) string {
	sparate_start := ""
	sparate_end := ""
	border := "=========="
	if is_response {
		border = "___________"
	}

	if DI.use_sparate {
		sparate_start = fmt.Sprintf("%s%s %s Information Start %s", DI.indentation, border, logModel.Mode, border)
		sparate_end = fmt.Sprintf("%s%s %s Information End %s", DI.indentation, border, logModel.Mode, border)
	}

	indentation := fmt.Sprintf("\n%s", DI.indentation)
	detail := fmt.Sprintf("%s%s Detail -=> %s", indentation, logModel.Mode, string(logModel.Detail))
	header := fmt.Sprintf("%s%s Header -=> %s", indentation, logModel.Mode, string(logModel.Header))
	body := fmt.Sprintf("%s%s Body -=> %s", indentation, logModel.Mode, string(logModel.Body))
	message_string := fmt.Sprintf("%s %s %s %s %s", sparate_start, detail, header, body, sparate_end)
	return message_string
}

func (DI *Pitlog_base) get_request_body_json_string(c echo.Context) string {
	var err error
	var body []byte
	var body_string_modified []byte
	request_id := c.Response().Header().Get(echo.HeaderXRequestID)
	if body, err = ioutil.ReadAll(c.Request().Body); err != nil {
		DI.Make_log_string(request_id, 3, DI.log_title, err.Error())
		return err.Error()
	}

	var body_json map[string]interface{}
	if err = json.Unmarshal(body, &body_json); err != nil {
		DI.Make_log_string(request_id, 3, DI.log_title, err.Error())
		return err.Error()
	}
	DI.modify_sensitive_fields(body_json)
	if body_string_modified, err = json.Marshal(body_json); err != nil {
		DI.Make_log_string(request_id, 3, DI.log_title, err.Error())
		return err.Error()
	}

	// Mengembalikan body ke body aslinya
	c.Request().Body = ioutil.NopCloser(strings.NewReader(string(body)))
	return string(body_string_modified)
}

func (DI *Pitlog_base) get_request_body_form_encoded_string(c echo.Context) string {
	var err error
	var bodyBytes []byte
	request_id := c.Response().Header().Get(echo.HeaderXRequestID)
	if bodyBytes, err = ioutil.ReadAll(c.Request().Body); err != nil {
		DI.Make_log_string(request_id, 3, DI.log_title, err.Error())
		return err.Error()
	}

	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	var formData url.Values
	if formData, err = url.ParseQuery(string(bodyBytes)); err != nil {
		DI.Make_log_string(request_id, 3, DI.log_title, err.Error())
		return err.Error()
	}

	formDataMap := make(map[string]interface{})
	for key, values := range formData {
		if len(values) > 0 {
			formDataMap[key] = values[0]
		}
	}

	DI.modify_sensitive_fields(formDataMap)
	var formDataJSON []byte
	if formDataJSON, err = json.Marshal(formDataMap); err != nil {
		DI.Make_log_string(request_id, 3, DI.log_title, err.Error())
		return err.Error()
	}

	c.Request().Body = ioutil.NopCloser(strings.NewReader(string(bodyBytes)))
	return string(formDataJSON)
}

func (DI *Pitlog_base) get_request_body_form_data(c echo.Context) string {
	var err error
	var form *multipart.Form
	request_id := c.Response().Header().Get(echo.HeaderXRequestID)
	if form, err = c.MultipartForm(); err != nil {
		DI.Make_log_string(request_id, 3, DI.log_title, err.Error())
		return err.Error()
	}

	// Menampilkan informasi files
	files := make([]map[string]interface{}, 0)
	for _, headers := range form.File {
		for _, header := range headers {
			fileInfo := map[string]interface{}{
				"file_name":    header.Filename,
				"content_type": header.Header.Get(name_content_type),
				"size":         header.Size,
			}
			files = append(files, fileInfo)
		}
	}

	// Menampilkan informasi form fields
	form_field := make(map[string]interface{})
	for key, values := range form.Value {
		// Mengambil nilai dari form fields dan menyimpannya ke dalam map
		if len(values) > 0 {
			form_field[key] = values[0]
		}
	}

	DI.modify_sensitive_fields(form_field)
	// Menampilkan hasil akhir dalam format yang diinginkan
	result := map[string]interface{}{
		"files":      files,
		"form_field": form_field,
	}

	var resultString []byte
	if resultString, err = json.Marshal(result); err != nil {
		DI.Make_log_string(request_id, 3, DI.log_title, err.Error())
		return err.Error()
	}

	return string(resultString)
}

func (DI *Pitlog_base) modify_sensitive_fields(data map[string]interface{}) {
	for key, value := range data {
		switch v := value.(type) {
		case map[string]interface{}:
			DI.modify_sensitive_fields(v)
		case string:
			if _, ok := sensitive_keys[key]; ok {
				data[key] = "*** " + key + " ***"
			}
		}
	}
}

func (DI *Pitlog_base) func_create_message_string(request_id string, status int, title string, message string) {
	var message_string string
	log_time := time.Now().Format(base_format_date_time_logger)

	switch status {
	case 1:
		message_string = fmt.Sprintf("%s%s request_id - %s - %s -=> %s x=> %s ", DI.indentation, name_type_status_info, request_id, title, message, log_time)
		DI.dedicated.Info(message_string)
		break
	case 2:
		message_string = fmt.Sprintf("%s%s request_id - %s - %s -=> %s", DI.indentation, name_type_status_debug, request_id, title, message)
		DI.dedicated.Debug(message_string)
		break
	case 3:
		message_string = fmt.Sprintf("%s%s request_id - %s - %s -=> %s", DI.indentation, name_type_status_error, request_id, title, message)
		DI.dedicated.Error(message_string)
		break
	default:
		message_string = fmt.Sprintf("%s%s request_id - %s - %s -=> %s", DI.indentation, name_type_status_error, request_id, title, message)
		DI.dedicated.Info(message_string)
	}
}

func (DI *Pitlog_base) func_create_message_object(request_id string, status int, title string, message string, is_object bool) {
	var obj Pit_json_message
	var message_string string
	log_time := time.Now().Format(base_format_date_time_logger)
	obj.Request_id = request_id
	obj.Status = name_type_status_info
	obj.Title = title
	obj.Message = message
	obj.Date = log_time
	jsonData, _ := json.Marshal(obj)
	jsonString := string(jsonData)
	if is_object {
		jsonString = strings.ReplaceAll(jsonString, "\\", "")
	}

	message_string = fmt.Sprintf("%s%s", DI.indentation, jsonString)
	switch status {
	case 1:
		DI.dedicated.Info(message_string)
		break
	case 2:
		DI.dedicated.Debug(message_string)
		break
	case 3:
		DI.dedicated.Error(message_string)
		break
	default:
		DI.dedicated.Info(message_string)
	}
}

func (DI *Pitlog_base) Make_log_string(request_id string, status int, title string, message string) {
	switch DI.object_view {
	case true:
		DI.func_create_message_object(request_id, status, title, message, false)
		break
	default:
		DI.func_create_message_string(request_id, status, title, message)
		break
	}
}

func (DI *Pitlog_base) Make_log_object(request_id string, status int, title string, object interface{}) {
	var data_string string
	if object != nil {
		json_data, _ := json.Marshal(object)
		data_string = string(json_data)
	}
	switch DI.object_view {
	case true:
		DI.func_create_message_object(request_id, status, title, data_string, true)
		break
	default:
		DI.func_create_message_string(request_id, status, title, data_string)
		break
	}

}
