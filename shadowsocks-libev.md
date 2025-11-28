/etc/shadowsocks-libev/config.json

shadowsocks-libev
{
    "server": "0.0.0.0",
    "server_port": 8388,
    "password": "zhouyu123",
    "method": "chacha20-ietf-poly1305",
    "timeout": 300,
    "fast_open": true
}

# 启动服务
sudo systemctl start shadowsocks-libev

# 停止服务
sudo systemctl stop shadowsocks-libev

# 重启 shadowsocks-libev 服务，重新加载配置文件
sudo systemctl restart shadowsocks-libev

# 设置服务开机自启
sudo systemctl enable shadowsocks-libev

# 检查服务运行状态，确认状态为 active (running)
sudo systemctl status shadowsocks-libev-server@config.service
✅
sudo systemctl status shadowsocks-libev