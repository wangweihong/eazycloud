package iharbor

import "time"

type ScannerBase struct {
	Name    string `json:"name"`
	Vendor  string `json:"vendor"`
	Version string `json:"version"`
}

type Scanner struct {
	UUID             string           `json:"uuid" form:"uuid"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	URL              string           `json:"url"`
	Disabled         bool             `json:"disabled"`   //扫描器禁止使用
	IsDefault        bool             `json:"is_default"` //默认扫描器
	Auth             string           `json:"auth"`       //""
	SkipCertVerify   bool             `json:"skip_certVerify"`
	UseInternalAddr  bool             `json:"use_internal_addr"`
	CreateTime       time.Time        `json:"create_time"`
	UpdateTime       time.Time        `json:"update_time"`
	CreateTimeUnix   int64            `json:"create_time_unix"`
	UpdateTimeUnix   int64            `json:"update_time_unix"`
	AccessCredential string           `json:"access_credential"`
	Metadata         *ScannerMetadata `json:"metadata"`
	LoadingMetadata  bool             `json:"loadingMetadata"`
}

func (s *Scanner) Convert() *Scanner {
	s.CreateTimeUnix = s.CreateTime.Unix()
	s.UpdateTimeUnix = s.UpdateTime.Unix()
	return s
}

type ScannerPingReq struct {
	Name             string `json:"name"`
	URL              string `json:"url"`
	Auth             string `json:"auth"`
	AccessCredential string `json:"access_credential"`
	SkipCertVerify   bool   `json:"skip_certVerify"`
	UseInternalAddr  bool   `json:"use_internal_addr"`
}

type ScannerRegisterRequest struct {
	Name             string `json:"name"`
	Description      string `json:"description"`
	URL              string `json:"url"`
	Auth             string `json:"auth"`
	AccessCredential string `json:"access_credential"`
	SkipCertVerify   bool   `json:"skip_certVerify"`
	UseInternalAddr  bool   `json:"use_internal_addr"`
	Disabled         bool   `json:"disabled"`
}

/*
{
  "scanner": {
    "name": "Clair",
    "vendor": "CoreOS",
    "version": "2.x"
  },
  "capabilities": [
    {
      "consumes_mime_types": [
        "application/vnd.oci.image.manifest.v1+json",
        "application/vnd.docker.distribution.manifest.v2+json"
      ],
      "produces_mime_types": [
        "application/vnd.scanner.adapter.vuln.report.harbor+json; version=1.0"
      ]
    }
  ],
  "properties": {
    "harbor.scanner-adapter/registry-authorization-type": "Bearer",
    "harbor.scanner-adapter/scanner-type": "os-package-vulnerability",
    "harbor.scanner-adapter/vulnerability-database-updated-at": "2020-07-05T22:34:25Z" //数据库更新时间
  }
}
*/

type ScannerMetadata struct {
	Scanner      ScannerBase `json:"scanner"`
	Capabilities []struct {
		ConsumesMimeTypes []string `json:"consumes_mime_types"`
		ProducesMimeTypes []string `json:"produces_mime_types"`
	} `json:"capabilities"`
	Properties map[string]string `json:"properties"`
}
