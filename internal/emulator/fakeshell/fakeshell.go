package fakeshell

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func ping(command string) {
	cmd := command
	if strings.Contains(cmd, "ping") {
		result := strings.Split(cmd, " ")
		host := result[1]
		//fmt.Println(result[0])
		//fmt.Println(result[1])
		fmt.Printf("PING %v (%[1]v) 56(84) bytes of data.\n", host)
		for i := 0; i < 4; i++ {
			time.Sleep(2 * time.Second)
			time_ping := FloatToString((rand.Float64() * 10) + 20)
			ttl := strconv.Itoa((rand.Intn(5) * 10))
			fmt.Printf("64 bytes from %v: icmp_seq=%v ttl=%v time=%.2v ms\n", host, i, ttl, time_ping)

		}
	}
}

func ls(command string) {

	etc := "acpi                           ec2_version         libpaper.d      overlayroot.conf         services\n" +
		"adduser.conf                   emacs               lighttpd        overlayroot.local.conf   sgml\n" +
		"alternatives                   environment         locale.alias    pam.conf                 shadow\n" +
		"apache2                        fonts               locale.gen      pam.d                    shadow-\n" +
		"apm                            fstab               localtime       papersize                shells\n" +
		"apparmor                       fuse.conf           logcheck        passwd                   siege\n" +
		"apparmor.d                     gai.conf            login.defs      passwd-                  skel\n" +
		"apport                         groff               logrotate.conf  perl                     sos.conf\n" +
		"apt                            group               logrotate.d     pm                       ssh\n" +
		"at.deny                        group-              lsb-release     polkit-1                 ssl\n" +
		"audisp                         grub.d              ltrace.conf     pollinate                subgid\n" +
		"audit                          gshadow             lvm             popularity-contest.conf  subgid-\n" +
		"bash.bashrc                    gshadow-            machine-id      ppp                      subuid\n" +
		"bash_completion                gss                 magic           profile                  subuid-\n" +
		"bash_completion.d              gtk-2.0             magic.mime      profile.d                subversion\n" +
		"bindresvport.blacklist         hdparm.conf         mailcap         protocols                sudoers\n" +
		"binfmt.d                       host.conf           mailcap.order   proxychains.conf         sudoers.d\n" +
		"byobu                          hostname            manpath.config  python                   supervisor\n" +
		"ca-certificates                hosts               mdadm           python2.7                sysctl.conf\n" +
		"ca-certificates.conf           hosts.allow         memcached.conf  python3                  sysctl.d\n" +
		"ca-certificates.conf.dpkg-old  hosts.deny          mime.types      python3.5                sysstat\n" +
		"calendar                       init                mke2fs.conf     rc0.d                    systemd\n" +
		"checkinstallrc                 init.d              modprobe.d      rc1.d                    terminfo\n" +
		"cloud                          initramfs-tools     modules         rc2.d                    timezone\n" +
		"colordiffrc                    inputrc             modules-load.d  rc3.d                    tmpfiles.d\n" +
		"console-setup                  insserv             mtab            rc4.d                    tor\n" +
		"cron.d                         insserv.conf        mysql           rc5.d                    ucf.conf\n" +
		"cron.daily                     insserv.conf.d      node  rc6.d                    udev\n" +
		"cron.hourly                    iproute2            nanorc          rc.local                 ufw "

	fake_fs := "bin   etc   initrd.img.old  lost+found  openvpn-ca  root  snap  tmp var  volume\n" +
		"boot  home        lib             media       opt         run   srv  vmlinuz\n" +
		"dev   initrd.img  lib64           mnt         proc        sbin  sys   usr       vmlinuz.old"
	cmd := command
	if cmd == "ls" {

		fmt.Println("data export2.csv start-up.sh")
	} else {
		result := strings.Split(cmd, " ")
		dir_file := result[1]
		//fmt.Println(result[0])
		//fmt.Println(result[1])
		if dir_file == "/root" || dir_file == "/" {
			fmt.Println(fake_fs)
		}
		if dir_file == "/etc" {
			fmt.Println(etc)
		}
	}
}

func cat(command string) {
	cmd := command

	result := strings.Split(cmd, " ")
	file := result[1]
	s := []string{"fake_files/", file}
	file_path := strings.Join(s, "")
	if file == "data" || file == "export2.csv" || file == "start-up.sh" {
		b, err := ioutil.ReadFile(file_path) // just pass the file name
		if err != nil {
			fmt.Print(err)
		}

		str := string(b) // convert content to a 'string'

		fmt.Println(str) // print the content as a 'string'
	} else {
		fmt.Println("File not found")
	}
}

func FakeShell() {
	//create your file with desired read/write permissions
	f, err := os.OpenFile("./logger/honey-ssh.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}

	//defer to close when you're done with it, not because you think it's idiomatic!
	defer f.Close()
	//set output of logs to f
	log.SetOutput(f)
	//var cmd string
	shell := "root@web-1:#"
	logged := true
	ifconfig := "ens3      Link encap:Ethernet  HWaddr fa:16:3e:ea:69:d3\n" +
		"          inet addr:192.168.0.3  Bcast:192.168.0.255  Mask:255.255.255.0\n" +
		"          inet6 addr: fe80::f816:3eff:feea:69d3/64 Scope:Link\n" +
		"          UP BROADCAST RUNNING MULTICAST  MTU:1500  Metric:1\n" +
		"          RX packets:0 errors:0 dropped:0 overruns:0 frame:0\n" +
		"          TX packets:2 errors:0 dropped:0 overruns:0 carrier:0\n" +
		"          collisions:0 txqueuelen:0\n" +
		"          RX bytes:0 (0.0 B)  TX bytes:180 (180.0 B)\n\n" +
		"lo        Link encap:Local Loopback \n" +
		"          inet addr:127.0.0.1  Mask:255.0.0.0 \n" +
		"          inet6 addr: ::1/128 Scope:Host \n" +
		"          UP LOOPBACK RUNNING  MTU:65536  Metric:1 \n" +
		"          RX packets:237 errors:0 dropped:0 overruns:0 frame:0 \n" +
		"          TX packets:237 errors:0 dropped:0 overruns:0 carrier:0 \n" +
		"          collisions:0 txqueuelen:1 \n" +
		"          RX bytes:16818 (16.8 KB)  TX bytes:16818 (16.8 KB) \n"

	ipa := "1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1\n" +
		"    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00\n" +
		"    inet 127.0.0.1/8 scope host lo\n" +
		"       valid_lft forever preferred_lft forever\n" +
		"    inet6 ::1/128 scope host\n" +
		"       valid_lft forever preferred_lft forever\n" +
		"2: ens3: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP group default qlen 1000\n" +
		"    link/ether fa:16:3e:d6:f2:dd brd ff:ff:ff:ff:ff:ff\n" +
		"    inet 192.168.0.183/24 brd 192.168.0.255 scope global ens3\n" +
		"       valid_lft forever preferred_lft forever\n" +
		"    inet6 fe80::f816:3eff:fed6:f2dd/64 scope link\n" +
		"       valid_lft forever preferred_lft forever \n"

	in := bufio.NewReader(os.Stdin)
	fmt.Println(shell)
	for logged {
		// ask for command and send fake output
		fmt.Printf(shell)
		//	fmt.Scanln("%s", &cmd)

		stringa, _ := in.ReadString('\n')
		cmd := strings.Trim(stringa, " \r\n")
		log.Println("SSH-CMD:", cmd)
		if cmd == "ifconfig" {
			fmt.Printf(ifconfig)
			cmd = " "
		}
		if cmd == "ip a" {
			fmt.Printf(ipa)
			cmd = " "
		}
		if strings.Contains(cmd, "echo") {
			cmd = strings.Trim(stringa, "echo ")
			fmt.Printf(cmd)
			cmd = " "
		}
		if strings.Contains(cmd, "ls") {
			ls(cmd)
		}
		if strings.Contains(cmd, "cat") {
			cat(cmd)
		}
		if cmd == "exit" || cmd == "quit" || cmd == "logout" {
			cmd = " "
			os.Exit(3)
		}
		if strings.Contains(cmd, "ping") {
			{
				ping(cmd)
			}

		}
	}

}
