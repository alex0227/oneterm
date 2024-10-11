package util

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/spf13/cast"
	mysql "github.com/veops/oneterm/db"
	ggateway "github.com/veops/oneterm/gateway"
	"github.com/veops/oneterm/model"
)

func GetAAG(assetId int, accountId int) (asset *model.Asset, account *model.Account, gateway *model.Gateway, err error) {
	asset, account, gateway = &model.Asset{}, &model.Account{}, &model.Gateway{}
	if err = mysql.DB.Model(asset).Where("id = ?", assetId).First(asset).Error; err != nil {
		return
	}
	if err = mysql.DB.Model(account).Where("id = ?", accountId).First(account).Error; err != nil {
		return
	}
	account.Password = DecryptAES(account.Password)
	account.Pk = DecryptAES(account.Pk)
	account.Phrase = DecryptAES(account.Phrase)
	if asset.GatewayId != 0 {
		if err = mysql.DB.Model(gateway).Where("id = ?", asset.GatewayId).First(gateway).Error; err != nil {
			return
		}
		gateway.Password = DecryptAES(gateway.Password)
		gateway.Pk = DecryptAES(gateway.Pk)
		gateway.Phrase = DecryptAES(gateway.Phrase)
	}

	return
}

func GetAuth(account *model.Account) (ssh.AuthMethod, error) {
	switch account.AccountType {
	case model.AUTHMETHOD_PASSWORD:
		return ssh.Password(account.Password), nil
	case model.AUTHMETHOD_PUBLICKEY:
		if account.Phrase == "" {
			pk, err := ssh.ParsePrivateKey([]byte(account.Pk))
			if err != nil {
				return nil, err
			}
			return ssh.PublicKeys(pk), nil
		} else {
			pk, err := ssh.ParsePrivateKeyWithPassphrase([]byte(account.Pk), []byte(account.Phrase))
			if err != nil {
				return nil, err
			}
			return ssh.PublicKeys(pk), nil
		}
	default:
		return nil, fmt.Errorf("invalid authmethod %d", account.AccountType)
	}
}

func Proxy(isConnectable bool, sessionId string, protocol string, asset *model.Asset, gateway *model.Gateway) (ip string, port int, err error) {
	ip, port = asset.Ip, 0
	for _, tp := range strings.Split(protocol, ",") {
		for _, p := range asset.Protocols {
			if strings.HasPrefix(strings.ToLower(p), tp) {
				if port = cast.ToInt(strings.Split(p, ":")[1]); port != 0 {
					break
				}
			}
		}
	}

	if asset.GatewayId == 0 || gateway == nil {
		return
	}

	g, err := ggateway.GetGatewayManager().Open(isConnectable, sessionId, ip, port, gateway)
	if err != nil {
		return
	}
	ip, port = g.LocalIp, g.LocalPort
	return
}
