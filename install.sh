curl -Lo /usr/bin/static-short https://github.com/chrissxYT/zerm.link/releases/download/1/static-short-arm
chmod +x /usr/bin/static-short

echo "[Unit]
Description=A simple link shortener

[Service]
Type=simple
Restart=always
ExecStart=/usr/bin/static-short

[Install]
WantedBy=multi-user.target" > /etc/systemd/system/static-short.service

systemctl enable --now static-short.service
