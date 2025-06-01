package sftp

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pkg/sftp"
	"go.elastic.co/apm"
	"golang.org/x/crypto/ssh"
)

type SftpOop struct {
	sftpClient *sftp.Client
	Conn       *ssh.Client
}

func Init(user string, pass string, host string, port string) (*SftpOop, error) {
	conf := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(pass)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", host+":"+port, conf)
	if err != nil {
		fmt.Println("dial tcp ssh fail : ", err)
		return nil, err
	}

	sftpClientNew, err := sftp.NewClient(conn)
	if err != nil {
		fmt.Println("creation of object sftp client fail : ", err)
		return nil, err
	}

	sftpCurrent := &SftpOop{
		sftpClient: sftpClientNew,
		Conn:       conn,
	}

	go HandleReconnect(sftpCurrent, user, pass, host, port)

	return sftpCurrent, nil
}

func HandleReconnect(sftpCurrent *SftpOop, user, pass, host, port string) (*SftpOop, error) {
	closed := make(chan string)

	go func() {
		closed <- sftpCurrent.Conn.Wait().Error()
	}()

	errMsg := <-closed
	fmt.Println("closed message:", errMsg)

	conf := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(pass)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	fmt.Println("SFTP reconnection attempt ...")

	conn, err := ssh.Dial("tcp", host+":"+port, conf)
	if err != nil {
		fmt.Println("error on dial to server:", err)
		return nil, err
	}

	sftpClientNew, err := sftp.NewClient(conn)
	if err != nil {
		fmt.Println("error on creating object sftp client:", err)
		return nil, err
	}

	sftpCurrent.Conn = conn
	sftpCurrent.sftpClient = sftpClientNew

	fmt.Println("SFTP reconnection success on:", time.Now().String())

	go HandleReconnect(sftpCurrent, user, pass, host, port)

	return sftpCurrent, nil
}

func (s *SftpOop) SendLocalFileToRemote(ctx context.Context, localpath, remotepath string) (int, error) {
	apmSpan, _ := apm.StartSpan(ctx, "Send Local File to Remote", "SFTP")
	defer apmSpan.End()

	count, err := s.sendFile(remotepath, localpath)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *SftpOop) SendLocalFileToRemoteWithDelete(ctx context.Context, localpath, remotepath string) (int, error) {
	apmSpan, _ := apm.StartSpan(ctx, "Send&Delete Local File to Remote", "SFTP")
	defer apmSpan.End()

	count, err := s.sendFile(remotepath, localpath)
	if err != nil {
		return 0, err
	}

	err = os.Remove(localpath)
	if err != nil {
		fmt.Println("error on deletion local file : ", err)
		return 0, err
	}

	return count, nil
}

func (s *SftpOop) sendFile(remotepath, localpath string) (int, error) {
	remoteFile, err := s.sftpClient.Create(remotepath)
	if err != nil {
		fmt.Println("error on creating pipeline to remote hhost : ", err)
		return 0, err
	}
	defer remoteFile.Close()

	localFile, err := os.Open(localpath)
	if err != nil {
		fmt.Println("error on open local file : ", err)
		return 0, err
	}
	defer localFile.Close()

	bytes, err := io.ReadAll(localFile)
	if err != nil {
		fmt.Println("error on read local file : ", err)
		return 0, err
	}

	count, err := remoteFile.Write(bytes)
	if err != nil {
		fmt.Println("error on write to remot file : ", err)
		return 0, err
	}

	return count, nil
}

func (s *SftpOop) GetConnection(ctx context.Context) *ssh.Client {
	return s.Conn
}
