package main

import (
	"net"
	"time"
	"io/ioutil"
	"golang.org/x/crypto/ssh"
	"github.com/BurntSushi/toml"
)

type Config struct {
	SSH SSH
	Mysql Mysql
}

type SSH struct {
	IP string
	Port string
	User string
	IdentityFile string
}

type Mysql struct {
	MysqlConf string
	Database string
	DumpDir string
	DumpFilePrefix string
}

func dump() {
	var config Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		panic(err)
	}
	
	buf, err := ioutil.ReadFile(config.SSH.IdentityFile)
	if err != nil {
		panic(err)
	}
	
	key, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		panic(err)
	}
	
	conn, err := ssh.Dial("tcp", config.SSH.IP+":"+config.SSH.Port, &ssh.ClientConfig{
		User: config.SSH.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	})
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	
	session, err := conn.NewSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()
	
  	byte, err := session.Output("sudo mysqldump --defaults-file="+config.Mysql.MysqlConf+" "+config.Mysql.Database+" "+"--quick --single-transaction")
	if err != nil {
  		panic(err)
	}
	
	ioutil.WriteFile(config.Mysql.DumpDir+config.Mysql.DumpFilePrefix+time.Now().Format("2006-01-02")+".sql", byte, 0644)
}

func main() {
	dump()
}
