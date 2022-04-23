package ebookdownloader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//LatestReleasesInfo 获取最新的Releases信息
type LatestReleasesInfo struct {
	URL string `json:"html_url"`
	Tag string `json:"tag_name"`
}

//String ReleaseInfo String()
func (ri LatestReleasesInfo) String() string {
	return fmt.Sprintf("Latest version is %s: %s", ri.Tag, ri.URL)
}

//Compare 对版本进行对比
func (ri LatestReleasesInfo) Compare(CurVersion string) string {
	if CurVersion != "dev" && !strings.HasPrefix(CurVersion, ri.Tag) {
		return fmt.Sprintf("Running version %s. Latest version is %s: %s\n", CurVersion, ri.Tag, ri.URL)
	}
	return fmt.Sprintf("Not need to update! Running version %s. Latest version is %s: %s\n", CurVersion, ri.Tag, ri.URL)
}

//UpdateCheck 检查更新
func UpdateCheck() (obj LatestReleasesInfo, err error) {
	resp, err := http.Get("https://api.github.com/repos/sndnvaps/ebookdownloader/releases/latest")
	if err != nil {
		return LatestReleasesInfo{}, err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return LatestReleasesInfo{}, err
	}

	if resp.StatusCode != 200 {
		return LatestReleasesInfo{}, err
	}

	err = json.Unmarshal(buf, &obj)
	if err != nil {
		return LatestReleasesInfo{}, err
	}

	return obj, nil
}
