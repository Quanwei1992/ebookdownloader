package ebookdownloader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//UpdateCheck 检查更新
func UpdateCheck(CurVersion string) (output string, err error) {
	resp, err := http.Get("https://api.github.com/repos/sndnvaps/ebookdownloader/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", err
	}

	var obj struct {
		URL string `json:"html_url"`
		Tag string `json:"tag_name"`
	}
	err = json.Unmarshal(buf, &obj)
	if err != nil {
		return "", err
	}

	if CurVersion != "dev" {
		if !strings.HasPrefix(CurVersion, obj.Tag) {
			return fmt.Sprintf("Running version %s. Latest version is %s: %s\n", CurVersion, obj.Tag, obj.URL), nil
		} else {
			return fmt.Sprintf("Not need to update! Running version %s. Latest version is %s: %s\n", CurVersion, obj.Tag, obj.URL), nil
		}
	}

	return "", errors.New("It should not get here!")
}
