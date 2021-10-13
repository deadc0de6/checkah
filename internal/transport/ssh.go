// Copyright (c) 2021 deadc0de6

package transport

// knownhost key mismatch issue: https://github.com/golang/go/issues/28870
// 		ssh-keyscan -H <host> >> ~/.ssh/known_hosts

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/knownhosts"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const (
	protocol       = "tcp"
	scpCommand     = "scp -tr %s"
	connRetry      = 3
	connRetrySleep = 3
)

// SSH the ssh struct
type SSH struct {
	config *ssh.ClientConfig
	client *ssh.Client
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func checkKnownHosts() (ssh.HostKeyCallback, error) {
	path := filepath.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	log.Debugf("reading known_hosts file from %s", path)

	f, err := knownhosts.New(path)
	if err != nil {
		return nil, err
	}

	return func(addr string, remote net.Addr, key ssh.PublicKey) error {
		log.Debugf("checking knownhost for %s (%v)", addr, remote)
		return f(addr, remote, key)
	}, nil
}

func loadAgent() ssh.AuthMethod {
	sock, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		log.Debug(err)
		return nil
	}

	log.Debug("agent socket found")
	signers := agent.NewClient(sock).Signers
	if signers != nil {
		return ssh.PublicKeysCallback(signers)
	}

	log.Debug("no signers for SSH agent")
	return nil
}

func loadKeyFile(path string) (ssh.Signer, error) {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	return signer, nil
}

// Mkdir mkdir over ssh
func (t *SSH) Mkdir(remotePath string) error {
	// get an ssh session
	session, err := t.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	remoteDir := path.Dir(remotePath)
	remoteBase := path.Base(remotePath)

	go func() {
		w, err := session.StdinPipe()
		if err != nil {
			log.Error(err)
			return
		}
		defer w.Close()
		// mkdir
		_, err = fmt.Fprintln(w, "D0755", 0, remoteBase)
		if err != nil {
			log.Error(err)
			return
		}
	}()

	cmd := fmt.Sprintf(scpCommand, remoteDir)
	out, err := session.CombinedOutput(cmd)
	if err != nil {
		return fmt.Errorf("%s: %v", out, err)
	}

	return nil
}

// Copy scp a file to the remote
func (t *SSH) Copy(localPath string, remotePath string, rights string) error {
	// read local file
	data, err := ioutil.ReadFile(localPath)
	if err != nil {
		return err
	}

	// get an ssh session
	session, err := t.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	remoteDir := path.Dir(remotePath)
	remoteBase := path.Base(remotePath)

	// content writer
	go func() {
		// open stdin session
		w, err := session.StdinPipe()
		if err != nil {
			log.Error(err)
			return
		}
		defer w.Close()

		// provide filename
		r := fmt.Sprintf("C0%s", rights)
		_, err = fmt.Fprintln(w, r, len(data), remoteBase)
		if err != nil {
			return
		}

		// write content
		_, err = io.Copy(w, bytes.NewReader(data))
		if err != nil {
			return
		}

		// terminate transfer
		_, err = fmt.Fprint(w, "\x00") // transfer end with \x00
		if err != nil {
			return
		}
	}()

	cmd := fmt.Sprintf(scpCommand, remoteDir)
	out, err := session.CombinedOutput(cmd)
	if err != nil {
		return fmt.Errorf("%s: %v", out, err)
	}

	return nil
}

// Execute executes a command through SSH
// returns stdout, stderr, error
func (t *SSH) Execute(cmd string) (string, string, error) {
	session, err := t.client.NewSession()
	if err != nil {
		log.Debugf("SSH new session for command \"%s\" failed: %v", cmd, err)
		return "", "", err
	}
	log.Debugf("SSH new session opened for: \"%s\"", cmd)
	defer session.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	session.Stdout = &stdout
	session.Stderr = &stderr

	log.Debugf("SSH run: \"%s\"", cmd)
	err = session.Run(cmd)

	_, ok := err.(*ssh.ExitMissingError)
	if ok {
		// ssh was successful but remote command didn't return an exit code
		log.Debugf("SSH command \"%s\" failed with no exit code", cmd)
		return stdout.String(), stderr.String(), fmt.Errorf("remote command is missing an exit code")
	}

	e, ok := err.(*ssh.ExitError)
	if ok {
		// an ExitError
		retCode := e.ExitStatus()
		log.Debugf("SSH command \"%s\" failed with exit code: %d", cmd, retCode)
		log.Debugf("SSH command \"%s\" failed with stdout: %s", cmd, stdout.String())
		log.Debugf("SSH command \"%s\" failed with stderr: %s", cmd, stderr.String())
		return stdout.String(), stderr.String(), fmt.Errorf("remote command exit code: %d", retCode)
	}

	if err != nil {
		// any other type of error is an I/O error
		log.Debugf("SSH command \"%s\" failed with I/O error: %v", cmd, err)
		return stdout.String(), stderr.String(), fmt.Errorf("I/O error: %v", err)
	}

	// command ran successfully
	log.Debugf("SSH command \"%s\" succeeded with stdout: %s", cmd, stdout.String())
	log.Debugf("SSH command \"%s\" succeeded with stderr: %s", cmd, stderr.String())
	return stdout.String(), stderr.String(), nil
}

// Close closes the SSH session
func (t *SSH) Close() {
	if t.client != nil {
		t.client.Close()
	}
}

// NewSSH creates an SSH instance
func NewSSH(host string, port string, user string, password string, keyfile string, timeout int, insecure bool) (*SSH, error) {
	var auths []ssh.AuthMethod

	log.Debugf("creating a new ssh connection to %s:%s", host, port)

	if len(host) < 1 {
		return nil, fmt.Errorf("SSH no host provided")
	}
	if len(port) < 1 {
		port = "22"
	}
	if len(user) < 1 {
		return nil, fmt.Errorf("SSH no user provided")
	}

	// add password as auth method
	if len(password) > 1 {
		auths = append(auths, ssh.Password(password))
		log.Info("SSH password auth method added")
	}

	// set default keyfile if unset
	if len(keyfile) < 1 {
		keyfile = filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa")
	}

	// add keyfile as auth method
	if len(keyfile) > 1 {
		if strings.HasPrefix(keyfile, "~/") {
			keyfile = filepath.Join(os.Getenv("HOME"), keyfile[2:])
		}

		log.Debugf("keyfile: %s", keyfile)

		if fileExists(keyfile) {
			log.Debugf("loading keyfile from %s", keyfile)
			s, err := loadKeyFile(keyfile)
			if err != nil {
				log.Error(err)
			} else {
				k := ssh.PublicKeys(s)
				if k != nil {
					auths = append(auths, ssh.PublicKeys(s))
				}
			}
		}
	}

	// add agent as auth method
	a := loadAgent()
	if a != nil {
		log.Debug("adding ssh agent as auth method")
		auths = append(auths, a)
	}

	var kn ssh.HostKeyCallback
	if insecure {
		log.Debug("insecure knownhost")
		kn = ssh.InsecureIgnoreHostKey()
	} else {
		var err error
		kn, err = checkKnownHosts()
		if err != nil {
			log.Debug("knownhost failed: ", err)
			return nil, err
		}
	}

	if len(auths) < 1 {
		return nil, fmt.Errorf("no SSH auths set up")
	}
	log.Debugf("%d auths methods: %#v", len(auths), auths)

	t := &SSH{}
	t.config = &ssh.ClientConfig{
		User:            user,
		Auth:            auths,
		HostKeyCallback: kn,
		Timeout:         time.Duration(timeout) * time.Second,
	}

	var c *ssh.Client
	var err error
	remote := fmt.Sprintf("%s:%s", host, port)
	for i := 0; i < connRetry; i++ {
		log.Debugf("SSH connecting to %s@%s port %s (%d/%d)", user, host, port, i+1, connRetry)
		c, err = ssh.Dial(protocol, remote, t.config)
		if err == nil {
			break
		} else {
			err = fmt.Errorf("ssh connection error (%d/%d): %s", i+1, connRetry, err.Error())
			log.Debug(err)
			time.Sleep(connRetrySleep * time.Second)
		}
	}
	if err != nil {
		return nil, err
	}
	log.Debugf("SSH connected to %s@%s port %s", user, host, port)
	t.client = c
	return t, nil
}
