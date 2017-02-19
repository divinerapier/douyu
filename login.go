package douyu

import (
	"bytes"
	"errors"
	"net"
	"time"

	log "qiniupkg.com/x/log.v7"
)

// DouyuDanmuServer 斗鱼弹幕服务器地址
const DouyuDanmuServer = "openbarrage.douyutv.com:8601"

// OpenDanmu 打开斗鱼弹幕
func OpenDanmu(rid int64) (dy *Douyu, err error) {
	dy = new(Douyu)
	dy.RoomID = rid
	dailer := &net.Dialer{
		Timeout: time.Second * 10,
	}
	if dy.Conn, err = dailer.Dial("tcp", DouyuDanmuServer); err != nil {
		dy = nil
		return
	}
	dy.login()
	dy.keepLiveChan = make(chan []byte, 100)
	dy.chatMsgChan = make(chan []byte, 100)
	return
}

func (dy *Douyu) login() error {
	var resp [1024]byte
	reqData := append([]byte("type@=loginreq/roomid@="), number2bytes(dy.RoomID)...)
	reqData = append(reqData, '/')
	reqData = PackRequest(reqData)
	dy.Write(reqData)
	cnt, err := dy.Read(resp[:])
	if err != nil {
		log.Error("login failed. err: ", err)
		return err
	}
	if cnt > 12 {
		dy.LoginResponse = parseLoginResponse(resp[12:cnt])
	} else {
		return errors.New("return nothing")
	}
	return nil
}

func parseLoginResponse(data []byte) *LoginResponse {
	log.Infof("login response: %s", data)
	resp := new(LoginResponse)
	lines := bytes.Split(data, []byte{'/'})
	for _, line := range lines {
		if i := bytes.Index(line, []byte("@=")); i > 0 {
			k := line[:i]
			v := line[i+2:]
			switch string(k) {
			case "type":
				resp.Type = v
			case "userid":
				resp.UserID = bytes2number(v)
			case "roomgroup":
				resp.RoomGroup = bytes2number(v)
			case "pg":
				resp.Pg = bytes2number(v)
			case "sessionid":
				resp.SessionID = bytes2number(v)
			case "username":
				resp.Username = v
			case "nickname":
				resp.Nickname = v
			case "live_stat":
				resp.LiveStat = bytes2number(v) != 0
			case "is_signined":
				resp.IsSigned = bytes2number(v) != 0
			case "signin_count":
				resp.SignedCount = bytes2number(v)
			case "npv":
				resp.NeedPhoneVerify = bytes2number(v) != 0
			case "best_dlev":
				resp.BestDlev = bytes2number(v)
			case "cur_lev":
				resp.CurLev = bytes2number(v)
			case "is_illegal":
				if len(v) == 0 || v[0] == 0 {
					resp.IsIllegal = false
				} else {
					resp.IsIllegal = true
				}
			case "ill_ct":
				resp.IllCt = bytes2number(v)
			case "ill_ts":
				resp.IllTs = bytes2number(v)
			case "now":
				resp.Now = bytes2number(v)
			case "ps":
				resp.Ps = bytes2number(v)
			case "es":
				resp.Es = bytes2number(v)
			case "it":
				resp.It = bytes2number(v)
			case "its":
				resp.Its = bytes2number(v)
			case "nrc":
				resp.Nrc = bytes2number(v)
			case "ih":
				resp.Ih = bytes2number(v)
			case "sid":
				resp.SID = bytes2number(v)
			case "code":
				resp.ErrCode = bytes2number(v)
			default:
				log.Error("unknow field:", string(k), "value:", string(v))
			}
		}

	}
	return resp
}
