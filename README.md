# Nacdlow's Development Domain Name Server

This is a DNS which resolves `local.nacdlow.com` to your local computer's IP
address. This allows testing PWAs.

## Usage

1. Build with `go build`.
2. Run with `sudo ./dev-dns-server <your ip>`

Make sure to replace `<your ip>` with your computer's IP address. You can find
out with `ip addr`.

### Setup on Android

1. Go to Wi-Fi page in your Settings application.
2. Tap on the current network, then Advanced.
3. Set IP settings from DHCP to Static.
4. Put your computer's IP in DNS 1, keep DNS 2 empty. Then Save.

### Setup on iOS

1. Go to Wi-Fi page in your Settings application.
2. Tap on the (i) next to the WiFi.
3. Scroll down and click Configure DNS.
4. Tap "Manual" and add your computer's IP address. Then Save.
