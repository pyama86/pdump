package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	logrus_stack "github.com/Gurpartap/logrus-stack"
	"github.com/k0kubun/pp"
	"github.com/sirupsen/logrus"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var (
	version   string
	revision  string
	goversion string
	builddate string
	builduser string
)

func init() {
	callerLevels := logrus.AllLevels
	stackLevels := []logrus.Level{logrus.PanicLevel, logrus.FatalLevel}
	logrus.AddHook(logrus_stack.NewHook(callerLevels, stackLevels))
}

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK     int  = 0
	ExitCodeError  int  = 1 + iota
	snaplen        int  = 65536
	counterCapa    uint = 30
	requiredSample uint = 3
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}
type cycleParams struct {
	interval uint
	alert    uint
	buffer   uint
	sec      uint
	nic      string
	exec     string
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		version bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet("pdump", flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	param := cycleParams{}
	flags.UintVar(&param.alert, "alert", 10, "alert threshould")
	flags.UintVar(&param.alert, "a", 10, "alert threshould(Short)")
	flags.UintVar(&param.buffer, "buffer", 0, "BufflerLength")
	flags.UintVar(&param.buffer, "b", 0, "BufflerLength(Short)")
	flags.UintVar(&param.sec, "sec", 5, "monitor sec")
	flags.UintVar(&param.sec, "s", 5, "monitor sec(Short)")
	flags.UintVar(&param.interval, "interval", 30, "monitor interval")
	flags.UintVar(&param.interval, "i", 30, "monitor interval(Short)")

	flags.StringVar(&param.nic, "nic", "", "monitor nic")
	flags.StringVar(&param.nic, "n", "", "monitor nic(Short)")

	flags.StringVar(&param.exec, "exec", "", "exec command")
	flags.StringVar(&param.exec, "e", "", "exec command(Short)")

	flags.BoolVar(&version, "version", false, "Print version information and quit.")

	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	if version {
		printVersion()
		return ExitCodeOK
	}

	if err := cycle(&param); err != nil {
		logrus.Error(err)
		return ExitCodeError
	}
	return ExitCodeOK

}

func cycle(p *cycleParams) error {
	logrus.SetLevel(logrus.InfoLevel)
	if os.Getenv("DEBUG") != "" {
		logrus.SetLevel(logrus.DebugLevel)
	}
	var handle *pcap.Handle
	inactive, err := pcap.NewInactiveHandle(p.nic)
	if err != nil {
		return fmt.Errorf("could not create: %v", err)
	}
	defer inactive.CleanUp()
	if err = inactive.SetSnapLen(snaplen); err != nil {
		return fmt.Errorf("could not set snap length: %v", err)
	} else if err = inactive.SetTimeout(time.Second); err != nil {
		return fmt.Errorf("could not set timeout: %v", err)
	}

	if handle, err = inactive.Activate(); err != nil {
		return fmt.Errorf("PCAP Activate error:%s", err)
	}

	ifs, err := net.InterfaceByName(p.nic)
	if err != nil {
		return err
	}

	var filters []string
	var ips []string
	counters := map[string]*counter{}
	addrs, err := ifs.Addrs()
	if err != nil {
		return err
	}

	for _, addr := range addrs {
		var ip *net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = &v.IP
		case *net.IPAddr:
			ip = &v.IP
		}
		ipstr := ip.String()
		if ip != nil {
			filters = append(filters, fmt.Sprintf(" dst %s", ipstr))
			ips = append(ips, ipstr)
			counters[ipstr] = &counter{capa: counterCapa}
		}
	}

	defer handle.Close()
	if err = handle.SetBPFFilter(strings.Join(filters, " or ")); err != nil {
		return fmt.Errorf("BPF filter error:%s", err)
	}

	source := gopacket.NewPacketSource(handle, handle.LinkType())
	source.NoCopy = true

	packetChannel := make(chan gopacket.Packet, p.buffer)

	for {
		logrus.Debug("start packet poling")
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		wg := &sync.WaitGroup{}
		go func() {
			wg.Add(1)
			logrus.Debug("start poling goroutine")
		INL:
			for {
				select {
				case <-ctx.Done():
					break INL
				case packet := <-packetChannel:
					ip4Layer := packet.Layer(layers.LayerTypeIPv4)
					if ip4Layer != nil {
						ip4 := ip4Layer.(*layers.IPv4)
						counters[ip4.DstIP.String()].increment()
					}
				}
			}
			logrus.Debug("end poling goroutine")
			wg.Done()
		}()

		timer := time.NewTimer(time.Duration(p.sec) * time.Second)
	PAC:
		for packet := range source.Packets() {
			select {
			case <-timer.C:
				cancel()
				break PAC
			default:
				packetChannel <- packet
			}
		}

		timer.Stop()
		wg.Wait()
		for _, i := range ips {
			c := counters[i]
			c.included()
			logrus.Debugf("ip:%s value:%v", i, pp.Sprint(*c))
			if c.avg()*p.alert < c.current && c.len > requiredSample && p.exec != "" {
				logrus.Infof("avg: %d, current: %d exec command:%s", c.avg(), c.current, p.exec)
				out, err := exec.Command(p.exec, i).CombinedOutput()
				if err != nil {
					return fmt.Errorf("exec cmd error:%s %s", err, string(out))
				}
			}
			c.reset()
		}
		logrus.Debug("end packet poling")
		time.Sleep(time.Duration(p.interval) * time.Second)
	}
	return nil
}

func printVersion() {
	fmt.Printf("stns version: %s (%s)\n", version, revision)
	fmt.Printf("build at %s (with %s) by %s\n", builddate, goversion, builduser)
}
