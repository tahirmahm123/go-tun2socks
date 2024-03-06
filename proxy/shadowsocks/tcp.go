package shadowsocks

import (
	"errors"
	"fmt"
	"io"
	"net"

	sscore "github.com/shadowsocks/go-shadowsocks2/core"
	sssocks "github.com/shadowsocks/go-shadowsocks2/socks"

	"github.com/tahirmahm123/go-tun2socks/common/dns"
	"github.com/tahirmahm123/go-tun2socks/common/log"
	"github.com/tahirmahm123/go-tun2socks/core"
)

type tcpHandler struct {
	cipher  sscore.Cipher
	server  string
	fakeDns dns.FakeDns
}

func (h *tcpHandler) handleInput(conn net.Conn, input io.ReadCloser) {
	defer func() {
		conn.Close()
		input.Close()
	}()
	io.Copy(conn, input)
}

func (h *tcpHandler) handleOutput(conn net.Conn, output io.WriteCloser) {
	defer func() {
		conn.Close()
		output.Close()
	}()
	io.Copy(output, conn)
}

func NewTCPHandler(server, cipher, password string, fakeDns dns.FakeDns) core.TCPConnHandler {
	ciph, err := sscore.PickCipher(cipher, []byte{}, password)
	if err != nil {
		log.Errorf("failed to pick a cipher: %v", err)
	}
	return &tcpHandler{
		cipher:  ciph,
		server:  server,
		fakeDns: fakeDns,
	}
}

func (h *tcpHandler) Handle(conn net.Conn, target net.Addr) error {
	if target == nil {
		log.Fatalf("unexpected nil target")
	}

	// Connect the relay server.
	rc, err := net.Dial("tcp", h.server)
	if err != nil {
		return errors.New(fmt.Sprintf("dial remote server failed: %v", err))
	}
	rc = h.cipher.StreamConn(rc)

	// Replace with a domain name if target address IP is a fake IP.
	host, port, err := net.SplitHostPort(target.String())
	if err != nil {
		log.Errorf("error when split host port %v", err)
	}
	var targetHost string = host
	if h.fakeDns != nil {
		if ip := net.ParseIP(host); ip != nil {
			if h.fakeDns.IsFakeIP(ip) {
				targetHost = h.fakeDns.QueryDomain(ip)
			}
		}
	}
	dest := fmt.Sprintf("%s:%s", targetHost, port)

	// Write target address.
	tgt := sssocks.ParseAddr(dest)
	_, err = rc.Write(tgt)
	if err != nil {
		return fmt.Errorf("send target address failed: %v", err)
	}

	go h.handleInput(conn, rc)
	go h.handleOutput(conn, rc)

	log.Infof("new proxy connection for target: %s:%s", target.Network(), dest)
	return nil
}
