port {{ .Port }}
proto {{ .Proto }}

dev tun

ca {{ .Ca }}
cert {{ .Cert }}
key {{ .Key }}

dh {{ .Dh }}

server 10.8.0.0 255.255.255.0

ifconfig-pool-persist {{ .IfconfigPoolPersist }}

push "route 10.8.0.0 255.255.255.0"

push "dhcp-option DNS 8.8.8.8"
push "dhcp-option DNS 8.8.4.4"

keepalive {{ .Keepalive }}

#tls-auth {{ .TaKey }} 0
tls-crypt {{ .TaKey }}

key-direction 0

cipher {{ .Cipher }}
auth {{ .Auth }}

#comp-lzo

max-clients {{ .MaxClients }}

user nobody
group nobody

persist-key
persist-tun

status /etc/openvpn/log/openvpn-status.log

log-append /etc/openvpn/log/openvpn.log

verb 3

mute 10

management {{ .Management }}
