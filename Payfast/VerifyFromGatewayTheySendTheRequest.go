package Payfast

import (
	"fmt"
	"github.com/astaxie/beego/httplib"
	"strings"
)

func (this *payfastProvider) VerifyFromGatewayTheySentTheRequest(remoteUserAgent string) {
	host := this.payfastContext.GetRemoteHost()
	url := fmt.Sprintf("https://%s/eng/query/validate", host)

	postBodyBytes := []byte(this.extractedData.ParamStringForRemoteVerify)

	request := httplib.Post(url).
		Header("Host", host).
		SetUserAgent(remoteUserAgent).
		Header("Content-Type", "application/x-www-form-urlencoded").
		Header("Content-Length", fmt.Sprintf("%d", len(postBodyBytes))).
		Body(postBodyBytes)
	// request = request.SetTimeout(connectTimeout, readWriteTimeout)
	//TODO: Proxies

	responseString, err := request.String()
	this.checkError(err)

	fmt.Println("A: ", responseString)
	if strings.ToUpper(responseString) != "VALID" {
		panic("Data is invalid")
	}
}
