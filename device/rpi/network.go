// +build linux

package rpi

import (
	"errors"
	"io/ioutil"
	"net"
	"os"

	"github.com/docker/libcontainer/netlink"
	"github.com/jdevelop/passkeeper/dhcpsrv"
)

const (
	usbEthernet   = usbGadget + "/functions/ecm.usb0"
	localIPStr    = "10.101.1.1"
	leaseStartStr = "10.101.1.2"
)

var localIP = net.ParseIP("10.101.1.1")

func ethernetUp() (err error) {
	err = os.MkdirAll(usbEthernet, os.ModeDir)
	if err != nil {
		return
	}
	ioutil.WriteFile(usbEthernet+"/host_addr", []byte("48:6f:73:74:50:43"), 0600)
	ioutil.WriteFile(usbEthernet+"/self_addr", []byte("42:61:64:55:53:42"), 0600)
	os.Symlink(usbEthernet, usbConfig+"/ecm.usb0")
	return
}

func networkUp(name string) (err error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return
	}

	var iface net.Interface
	found := false

	for _, _iface := range ifaces {
		log("Checking net %d : %v", _iface.Index, _iface.Name)
		if _iface.Name == name {
			iface = _iface
			found = true
			break
		}
	}

	if !found {
		err = errors.New("Interface " + name + " not found")
		return
	}

	err = netlink.NetworkLinkAddIp(&iface, localIP, &net.IPNet{
		IP:   localIP,
		Mask: net.IPv4Mask(255, 255, 255, 0),
	})

	if err != nil {
		return
	}

	err = netlink.NetworkLinkUp(&iface)

	return

}

func dhcpUp(iface, localIP, leaseStart string) error {
	dhcpsrv.StartDHCP(dhcpsrv.IFace(iface), dhcpsrv.IP(localIP), dhcpsrv.LeaseStart(leaseStart))
	return nil
}
