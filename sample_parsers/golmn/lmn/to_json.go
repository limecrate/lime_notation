package lmn

import "encoding/json"

func ToJson(lmn string) (string, error) {
	res, err := LmnParse(lmn)

	if err != nil {
		return "", err
	}

	json, err := json.Marshal(res)

	if err != nil {
		return "", err
	}

	return string(json), nil
}

func ToJsonIndent(lmn string) (string, error) {
	res, err := LmnParse(lmn)

	if err != nil {
		return "", err
	}

	json, err := json.MarshalIndent(res, "", "\t")

	if err != nil {
		return "", err
	}

	return string(json), nil
}
