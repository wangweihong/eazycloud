package iharbor

import (
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
)

type TagComponent struct {
	Total   int                   `json:"total"`
	Summary []TagComponentSummary `json:"summary"`
}

type TagComponentSummary struct {
	Severity int `json:"severity"`
	Count    int `json:"count"`
}

//type Tag struct {
//	Digest        string       `json:"digest"`
//	Name          string       `json:"name"`
//	Size          int          `json:"size"`
//	Architecture  string       `json:"architecture"`
//	Os            string       `json:"os"`
//	OsVersion     string       `json:"os.version"`
//	DockerVersion string       `json:"docker_version"`
//	Author        string       `json:"author"`
//	Created       time.Time    `json:"created"`
//	Signature     TagSignature `json:"signature"`
//	Immutable     bool         `json:"immutable"`
//	Labels        []Label      `json:"labels"`
//	PushTime      time.Time    `json:"push_time"`
//	PullTime      time.Time    `json:"pull_time"`
//}

type Tag struct {
	Digest        string                     `json:"digest"`
	Name          string                     `json:"name"`
	Size          int                        `json:"size"`
	Architecture  string                     `json:"architecture"`
	Os            string                     `json:"os"`
	OsVersion     string                     `json:"os.version"`
	DockerVersion string                     `json:"docker_version"`
	Author        string                     `json:"author"`
	Created       *time.Time                 `json:"created"`
	Config        TagConfig                  `json:"config"`
	Immutable     bool                       `json:"immutable"`
	Signature     *TagSignature              `json:"signature,omitempty"`
	ScanOverview  map[string]TagScanOverview `json:"scan_overview,omitempty"` // 为NIL表示未扫描
	Labels        []Label                    `json:"labels"`
	PushTime      *time.Time                 `json:"push_time"`
	PullTime      *time.Time                 `json:"pull_time"`

	//v2
	ArtifactDigest string `json:"artifact_digest"`
}

type Label struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Color        string `json:"color"`
	Scope        string `json:"scope"`
	ProjectID    int    `json:"project_id"`
	CreationTime string `json:"creation_time"`
	UpdateTime   string `json:"update_time"`
	Deleted      bool   `json:"deleted"`
}

/*
	"scan_overview": {
	   "application/vnd.scanner.adapter.vuln.report.harbor+json; version=1.0": {
	     "report_id": "05a4cb7b-bcd0-11ea-9353-8ab484cdae1b",
	     "scan_status": "Success",
	     "severity": "None",
	     "duration": 63,
	     "summary": {
	       "total": 0,
	       "fixable": 0,
	       "summary": {}
	     },
	     "start_time": "2020-07-03T01:53:42.212698Z",
	     "end_time": "2020-07-03T01:54:45.286306Z",
	     "scanner": {
	       "name": "Clair",
	       "vendor": "CoreOS",
	       "version": "2.x"
	     }
	   }
	 },
*/
type TagScanOverview struct {
	ReportID   string         `json:"report_id"`
	ScanStatus string         `json:"scan_status"` //Success/Running/Pending/Error 扫描状态（结束/运行）
	Severity   string         `json:"severity"`    //None
	Duration   int            `json:"duration"`    //扫描时间
	Summary    TagScanSummary `json:"summary"`     //扫描结束后才有数据
	StartTime  *time.Time     `json:"start_time"`
	EndTime    *time.Time     `json:"end_time"`
	Scanner    *TagScanner    `json:"scanner"` //扫描结束后才有数据
}

type TagScanner struct {
	Name    string `json:"name"`
	Vendor  string `json:"vendor"`
	Version string `json:"version"`
}

type TagScanSummary struct {
	Total   int         `json:"total"`   //总的漏洞数？
	Fixable int         `json:"fixable"` //可修复漏洞数
	Summary ScanSummary `json:"summary"`
}

type TagSignature struct {
}

// 扫描出来的漏洞总结
type ScanSummary struct {
	Low        int `json:"Low"`        //严重级别为低
	Medium     int `json:"Medium"`     //严重级别为中等
	Negligible int `json:"Negligible"` //严重界别为忽略不计
	High       int `json:"High"`       //严重界别为高
	Critical   int `json:"Critical"`   //严重级别为严重
	Unknown    int `json:"Unknown"`    //严重级别未知
}

// repository用的小写, 扫描器用的首字母大写
const (
	SeverityLow        = "low"
	SeverityMedium     = "medium"
	SeverityNegligible = "negligible"
	SeverityHigh       = "high"
	SeverityCritical   = "critical"
	SeverityUnknown    = "unknown"
)

var SeverityMap = map[string]int{
	SeverityUnknown:    0,
	SeverityNegligible: 1,
	SeverityLow:        2,
	SeverityMedium:     3,
	SeverityHigh:       4,
	SeverityCritical:   5,
}

func SeverityCompare(a, b string) bool {
	levela, _ := SeverityMap[strings.ToLower(a)]
	levelb, _ := SeverityMap[strings.ToLower(b)]
	return levela >= levelb
}

type TagConfig struct {
	Labels map[string]string `json:"labels"`
}

type TagListResponse struct {
	List       []*Tag `json:"list"`
	TotalCount int64 `json:"total_count"`
}

type TagRequest struct {
	ProjectName    string `json:"project_name" form:"project_name"` //harbor 2.1.0 use this as repo  project
	RepoName       string `json:"repo_name" form:"repo_name"`       //必须是完整的repo name
	TagName        string `json:"tag_name" form:"tag_name"`
	ArtifactDigest string `json:"artifact_digest" form:"artifact_digest"`
}

type TagIdentifier struct {
	ProjectIdentifier
	RepoName string `json:"repo_name" form:"repo_name"` //必须是完整的repo name
	TagName  string `json:"tag_name" form:"tag_name"`
	//for v2
	ArtifactDigest string `json:"artifact_digest" form:"artifact_digest"`
	//PagingParam
}

type TagScanResults map[string]TagScanResult

type TagScanResult struct {
	GeneratedAt          time.Time          `json:"generated_at"`
	GeneratedAtUnix      int64              `json:"generated_at_unix"`
	Scanner              ScannerBase        `json:"scanner"`
	Severity             string             `json:"severity"` // "Medium"
	Vulnerabilities      []TagVulnerability `json:"vulnerabilities"`
	VulnerabilitiesTotal *int               `json:"vulnerabilities_total,omitempty"`
}

func (p *TagScanResult) Convert() *TagScanResult {
	p.GeneratedAtUnix = p.GeneratedAt.Unix()
	return p
}

type TagVulnerability struct {
	ID           string `json:"id"`
	Severity     string `json:"severity"`
	Package      string `json:"package"`
	Version      string `json:"version"`
	Description  string `json:"description"`
	FixedVersion string `json:"fix_version"`
}

/* 扫描结果
{
  "application/vnd.scanner.adapter.vuln.report.harbor+json; version=1.0": {
    "generated_at": "2020-07-06T08:22:49.026676067Z",
    "scanner": {
      "name": "Clair",
      "vendor": "CoreOS",
      "version": "2.x"
    },
    "severity": "Medium",
    "vulnerabilities": [
      {
        "id": "CVE-2020-13844",
        "package": "gcc-5",
        "version": "5.4.0-6ubuntu1~16.04.12",
        "fix_version": "",
        "severity": "Medium",
        "description": "Arm Armv8-A core implementations utilizing speculative execution past unconditional changes in control flow may allow unauthorized disclosure of information to an attacker with local user access via a side-channel analysis, aka \"straight-line speculation.\"",
        "links": [
          "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2020-13844"
        ]
      },
      {
        "id": "CVE-2020-14155",
        "package": "pcre3",
        "version": "2:8.38-3.1",
        "fix_version": "",
        "severity": "Medium",
        "description": "libpcre in PCRE before 8.44 allows an integer overflow via a large number after a (?C substring.",
        "links": [
          "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2020-14155"
        ]
      },
      {
        "id": "CVE-2017-7186",
        "package": "pcre3",
        "version": "2:8.38-3.1",
        "fix_version": "",
        "severity": "Low",
        "description": "libpcre1 in PCRE 8.40 and libpcre2 in PCRE2 10.23 allow remote attackers to cause a denial of service (segmentation violation for read access, and application crash) by triggering an invalid Unicode property lookup.",
        "links": [
          "http://people.ubuntu.com/~ubuntu-security/cve/CVE-2017-7186"
        ]
      },

*/

type TagListRequest struct {
	ProjectIdentifier
	RepoName string `json:"repo_name" form:"repo_name"` //必须是完整的repo name
	//PagingParam
}

type TagManifestsResponse struct {
}

type TagManifests struct {
	Manifest TagManifestsDetail `json:"manifest"`
	Config   string             `json:"config"`
	Conf     LayerDetail        `json:"conf"`
}

type LayerDetail struct {
	Architecture    string     `json:"architecture"`
	Config          BaseConfig `json:"config"`
	Container       string     `json:"container"`
	ContainerConfig BaseConfig `json:"container_config"`
	Created         string     `json:"created"`
	DockerVersion   string     `json:"docker_version"`
	History         []struct {
		Created    string `json:"created"`
		CreatedBy  string `json:"created_by"`
		Comment    string `json:"comment,omitempty"`
		EmptyLayer bool   `json:"empty_layer,omitempty"`
	} `json:"history"`
	Os     string `json:"os"`
	Rootfs struct {
		Type    string   `json:"type"`
		DiffIds []string `json:"diff_ids"`
	} `json:"rootfs"`
}

type BaseConfig struct {
	Hostname     string                 `json:"Hostname"`
	Domainname   string                 `json:"Domainname"`
	User         string                 `json:"User"`
	AttachStdin  bool                   `json:"AttachStdin"`
	AttachStdout bool                   `json:"AttachStdout"`
	AttachStderr bool                   `json:"AttachStderr"`
	ExposedPorts map[string]empty.Empty `json:"ExposedPorts"`
	Tty          bool                   `json:"Tty"`
	OpenStdin    bool                   `json:"OpenStdin"`
	StdinOnce    bool                   `json:"StdinOnce"`
	Env          []string               `json:"Env"`
	Cmd          []string               `json:"Cmd"`
	ArgsEscaped  bool                   `json:"ArgsEscaped"`
	Image        string                 `json:"Image"`
	Volumes      map[string]empty.Empty `json:"Volumes"`
	WorkingDir   string                 `json:"WorkingDir"`
	Entrypoint   []string               `json:"Entrypoint"`
	OnBuild      []string               `json:"OnBuild"`
	Labels       map[string]string      `json:"Labels"`
}

type TagManifestsDetail struct {
	SchemaVersion int    `json:"schemaVersion"`
	MediaType     string `json:"mediaType"`
	Config        struct {
		MediaType string `json:"mediaType"`
		Size      int    `json:"size"`
		Digest    string `json:"digest"`
	} `json:"config"`
	Layers []struct {
		MediaType string `json:"mediaType"`
		Size      int    `json:"size"`
		Digest    string `json:"digest"`
	} `json:"layers"`
}

type TagScanLog struct {
	Logs map[string]TagLogResult `json:"logs"`
}

type TagLogResult struct {
	Log   string `json:"log"`
	Error string `json:"error"`
}

func ConvertArtifactToTagList(artifact *Artifact) []*Tag {
	tags := make([]*Tag, 0)
	for _, v := range artifact.Tags {
		tag := &Tag{
			ArtifactDigest: artifact.Digest,
			Digest:         artifact.Digest,
			Name:           v.Name,
			Size:           artifact.Size,
			Architecture:   artifact.ExtraAttrs.Architecture,
			Os:             artifact.ExtraAttrs.Os,
			OsVersion:      "",
			DockerVersion:  "",
			//Author:        artifact.ExtraAttrs.Author,
			Config:    TagConfig{},
			Immutable: v.Immutable,
		}

		tag.PushTime = artifact.PushTime
		tag.PullTime = artifact.PullTime
		tag.Created = artifact.ExtraAttrs.Created
		if artifact.ScanOverview != nil {
			tag.ScanOverview = make(map[string]TagScanOverview)
			for k, v := range artifact.ScanOverview {
				tso := TagScanOverview{
					ReportID:   v.ReportID,
					ScanStatus: v.ScanStatus,
					Severity:   v.Severity,
					Duration:   v.Duration,
					Summary: TagScanSummary{
						Total:   v.Summary.Total,
						Fixable: v.Summary.Fixable,
						Summary: ScanSummary{
							Low:        v.Summary.Summary.Low,
							Medium:     v.Summary.Summary.Medium,
							Negligible: v.Summary.Summary.Negligible,
							High:       v.Summary.Summary.High,
							Critical:   v.Summary.Summary.Critical,
							Unknown:    v.Summary.Summary.Unknown,
						},
					},

					Scanner: &TagScanner{},
				}
				tso.StartTime = v.StartTime
				tso.EndTime = v.EndTime
				if v.Scanner != nil {
					tso.Scanner = &TagScanner{
						Name:    v.Scanner.Name,
						Vendor:  v.Scanner.Vendor,
						Version: v.Scanner.Version,
					}
				}

				tag.ScanOverview[k] = tso
			}
		}
		tags = append(tags, tag)
	}

	return tags
}
