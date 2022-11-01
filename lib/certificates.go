package lib

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"time"

	"github.com/vuonglequoc/openvpn-web-ui/models"
	"github.com/beego/beego"
)

//Cert
//https://groups.google.com/d/msg/mailing.openssl.users/gMRbePiuwV0/wTASgPhuPzkJ
type Cert struct {
	EntryType   string
	Expiration  string
	ExpirationT time.Time
	Revocation  string
	RevocationT time.Time
	Serial      string
	FileName    string
	Details     *Details
}

type Details struct {
	Name         string
	CN           string
	Country      string
	Organisation string
	Email        string
}

func ReadCerts(path string) ([]*Cert, error) {
	certs := make([]*Cert, 0, 0)
	text, err := ioutil.ReadFile(path)
	if err != nil {
		return certs, err
	}
	lines := strings.Split(trim(string(text)), "\n")
	for _, line := range lines {
		fields := strings.Split(trim(line), "\t")
		if len(fields) != 6 {
			return certs,
				fmt.Errorf("Incorrect number of lines in line: \n%s\n. Expected %d, found %d",
					line, 6, len(fields))
		}
		expT, _ := time.Parse("060102150405Z", fields[1])
		revT, _ := time.Parse("060102150405Z", fields[2])
		c := &Cert{
			EntryType:   fields[0],
			Expiration:  fields[1],
			ExpirationT: expT,
			Revocation:  fields[2],
			RevocationT: revT,
			Serial:      fields[3],
			FileName:    fields[4],
			Details:     parseDetails(fields[5]),
		}
		certs = append(certs, c)
	}

	return certs, nil
}

func parseDetails(d string) *Details {
	details := &Details{}
	lines := strings.Split(trim(string(d)), "/")
	for _, line := range lines {
		if strings.Contains(line, "") {
			fields := strings.Split(trim(line), "=")
			switch fields[0] {
			case "name":
				details.Name = fields[1]
			case "CN":
				details.CN = fields[1]
			case "C":
				details.Country = fields[1]
			case "O":
				details.Organisation = fields[1]
			case "emailAddress":
				details.Email = fields[1]
			default:
				beego.Warn(fmt.Sprintf("Undefined entry: %s", line))
			}
		}
	}
	
	if (details.Name == "") && (details.CN != "") {
		details.Name = details.CN;
	}

	if details.Country == "" {
		details.Country = "unknown";
	}

	if details.Organisation == "" {
		details.Organisation = "unknown";
	}

	if details.Email == "" {
		details.Email = "unknown";
	}

	return details
}

func trim(s string) string {
	return strings.Trim(strings.Trim(s, "\r\n"), "\n")
}

func CreateCertificate(name string) error {
	// source /etc/openvpn/keys/vars \
	// 	&& export KEY_NAME=[name] \
	// 	&& /usr/share/easy-rsa/build-key --batch [name]

	// Easy-RSA 3
	// init-pki
	// build-ca [ cmd-opts ]
	// gen-dh
	// gen-req <filename_base> [ cmd-opts ]
	// sign-req <type> <filename_base>
	// build-client-full <filename_base> [ cmd-opts ]
	// build-server-full <filename_base> [ cmd-opts ]
	// revoke <filename_base> [cmd-opts]
	// renew <filename_base> [cmd-opts]
	// build-serverClient-full <filename_base> [ cmd-opts ]
	// gen-crl
	// update-db
	// show-req <filename_base> [ cmd-opts ]
	// show-cert <filename_base> [ cmd-opts ]
	// show-ca [ cmd-opts ]
	// import-req <request_file_path> <short_basename>
	// export-p7 <filename_base> [ cmd-opts ]
	// export-p8 <filename_base> [ cmd-opts ]
	// export-p12 <filename_base> [ cmd-opts ]
	// set-rsa-pass <filename_base> [ cmd-opts ]
	// set-ec-pass <filename_base> [ cmd-opts ]
	// upgrade <type>

	rsaPath := "/usr/share/easy-rsa/"
	ovpnPath := models.GlobalCfg.OVConfigPath
	caPath := models.GlobalCfg.CAConfigPath

	// Creating an OpenVPN Server Certificate Request and Private Key
	cmd := exec.Command(
		"/bin/sh",
		"-c",
		fmt.Sprintf(
			"openssl genrsa -out %s/pki/private/client_%s.key 2048" +
			" && openssl req -new -key %s/pki/private/client_%s.key -out %s/pki/reqs/client_%s.req" +
			" -subj /emailAddress=\"%s\"/C=\"%s\"/ST=\"%s\"/L=\"%s\"/O=\"%s\"/OU=\"%s\"/CN=\"%s\"",
			ovpnPath, name, ovpnPath, name, ovpnPath, name,
			"webmaster@example.com", "US", "New York", "New York City", "DigitalOcean", "Community", name))
	cmd.Dir = ovpnPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		beego.Debug(string(output))
		beego.Error(err)
		return err
	}

	// Signing the OpenVPN Server’s Certificate Request
	cmd = exec.Command(
		"/bin/sh",
		"-c",
		fmt.Sprintf(
			"%s/easyrsa import-req %s/pki/reqs/client_%s.req client_%s" +
			" && echo -e \"yes\" | %s/easyrsa sign-req client client_%s",
			rsaPath, ovpnPath, name, name, rsaPath, name))
	cmd.Dir = caPath
	output, err = cmd.CombinedOutput()
	if err != nil {
		beego.Debug(string(output))
		beego.Error(err)
		return err
	}

	// Copy Signed Certificate
	cmd = exec.Command(
		"/bin/sh",
		"-c",
		fmt.Sprintf(
			"cp %s/pki/issued/client_%s.crt %s/client-configs/keys/client_%s.crt",
			caPath, name, ovpnPath, name))
	cmd.Dir = caPath
	output, err = cmd.CombinedOutput()
	if err != nil {
		beego.Debug(string(output))
		beego.Error(err)
		return err
	}

	return nil
}
