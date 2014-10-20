// DO NOT EDIT!
// Code generated by ffjson <https://github.com/pquerna/ffjson>
// source: user.go
// DO NOT EDIT!
// Ups ... it's edited ...

package picnic

import (
	"bytes"
	"github.com/SchumacherFM/wanderlust/helpers"
)

func (mj *userModelCollection) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.Grow(1024)
	err := mj.MarshalJSONBuf(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (mj *userModelCollection) MarshalJSONBuf(buf *bytes.Buffer) error {
	var err error
	var obj []byte
	var first bool = true
	_ = obj
	_ = err
	_ = first
	buf.WriteString(`{`)
	if first == true {
		first = false
	} else {
		buf.WriteString(`,`)
	}
	buf.WriteString(`"Users":`)
	if mj.Users != nil {
		buf.WriteString(`[`)
		for i, v := range mj.Users {
			if i != 0 {
				buf.WriteString(`,`)
			}
			if v != nil {
				err = v.MarshalJSONBuf(buf)
				if err != nil {
					return err
				}
			} else {
				buf.WriteString(`null`)
			}
		}
		buf.WriteString(`]`)
	} else {
		buf.WriteString(`null`)
	}
	buf.WriteString(`}`)
	return nil
}

func (mj *userModel) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.Grow(1024)
	err := mj.MarshalJSONBuf(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (mj *userModel) MarshalJSONBuf(buf *bytes.Buffer) error {
	var err error
	var obj []byte
	var first bool = true
	_ = obj
	_ = err
	_ = first
	buf.WriteString(`{`)
	if first == true {
		first = false
	} else {
		buf.WriteString(`,`)
	}
	buf.WriteString(`"CreatedAt":`)
	obj, err = mj.CreatedAt.MarshalJSON()
	if err != nil {
		return err
	}
	buf.Write(obj)
	if first == true {
		first = false
	} else {
		buf.WriteString(`,`)
	}
	buf.WriteString(`"Email":`)
	helpers.Ffjson_WriteJsonString(buf, mj.Email)
	if first == true {
		first = false
	} else {
		buf.WriteString(`,`)
	}
	buf.WriteString(`"IsActivated":`)
	if mj.IsActivated {
		buf.WriteString(`true`)
	} else {
		buf.WriteString(`false`)
	}
	if first == true {
		first = false
	} else {
		buf.WriteString(`,`)
	}
	buf.WriteString(`"IsAdmin":`)
	if mj.IsAdmin {
		buf.WriteString(`true`)
	} else {
		buf.WriteString(`false`)
	}
	if first == true {
		first = false
	} else {
		buf.WriteString(`,`)
	}
	buf.WriteString(`"IsAuthenticated":`)
	if mj.IsAuthenticated {
		buf.WriteString(`true`)
	} else {
		buf.WriteString(`false`)
	}
	if first == true {
		first = false
	} else {
		buf.WriteString(`,`)
	}
	buf.WriteString(`"Name":`)
	helpers.Ffjson_WriteJsonString(buf, mj.Name)
	if first == true {
		first = false
	} else {
		buf.WriteString(`,`)
	}
	buf.WriteString(`"Password":`)
	helpers.Ffjson_WriteJsonString(buf, mj.Password)
	if first == true {
		first = false
	} else {
		buf.WriteString(`,`)
	}
	buf.WriteString(`"RecoveryCode":`)
	helpers.Ffjson_WriteJsonString(buf, mj.RecoveryCode)
	if first == true {
		first = false
	} else {
		buf.WriteString(`,`)
	}
	buf.WriteString(`"UserName":`)
	helpers.Ffjson_WriteJsonString(buf, mj.UserName)
	buf.WriteString(`}`)
	return nil
}
