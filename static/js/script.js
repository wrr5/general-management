document.addEventListener('DOMContentLoaded', function() {
    const physicalStoreOption = document.getElementById('physical-store');
    const cloudStoreOption = document.getElementById('cloud-store');
    const saveImageBtn = document.getElementById('saveImageBtn');
    
    // 门店类型卡片点击事件
    physicalStoreOption.addEventListener('click', function() {
        confirmAndRedirect('门店', 'https://hjrtu.flctwa.cn/#/pages-store/registerStore/registerStore?qrcode_sn=LAd9w1HStzlN5XihRBiZaxL03DnxZ2xf&share_aid=75&invite_type=1&group_id=&manager_id=');
    });
    
    cloudStoreOption.addEventListener('click', function() {
        confirmAndRedirect('云店', 'https://hjrtu.flctwa.cn/#/pages-store/registerStore/registerStore?qrcode_sn=cd3ccxSJFLXW4e7n9aCGhDB0ptE1OgNN&share_aid=75&invite_type=2&group_id=&manager_id=');
    });
    
    // 保存图片按钮点击事件
    if (saveImageBtn) {
        saveImageBtn.addEventListener('click', function() {
            saveImage('../static/images/mentou.jpg', '门头照片');
        });
    }
    
    // 保存图片函数
    function saveImage(url, filename) {
        // 对于本地图片，使用fetch和blob方法
        fetch(url)
            .then(response => response.blob())
            .then(blob => {
                const link = document.createElement('a');
                link.href = URL.createObjectURL(blob);
                link.download = filename + '.jpg';
                document.body.appendChild(link);
                link.click();
                document.body.removeChild(link);
                showNotification('图片已保存到本地');
            })
            .catch(error => {
                console.error('保存图片失败:', error);
                // 如果fetch失败，尝试直接打开图片
                const newWindow = window.open(url, '_blank');
                if (!newWindow) {
                    showNotification('请右键点击图片选择"图片另存为"', 'error');
                } else {
                    showNotification('请在新窗口中右键保存图片');
                }
            });
    }
    
    // 显示通知函数
    function showNotification(message, type = 'success') {
        const notification = document.createElement('div');
        notification.className = `notification ${type}`;
        notification.textContent = message;
        
        document.body.appendChild(notification);
        
        // 触发动画
        setTimeout(() => {
            notification.classList.add('show');
        }, 10);
        
        // 3秒后自动移除通知
        setTimeout(() => {
            if (document.body.contains(notification)) {
                notification.classList.remove('show');
                setTimeout(() => {
                    if (document.body.contains(notification)) {
                        document.body.removeChild(notification);
                    }
                }, 300);
            }
        }, 3000);
    }
    
    // 确认并跳转函数
    function confirmAndRedirect(storeType, link) {
        // 创建自定义确认对话框
        const modal = document.createElement('div');
        modal.className = 'confirm-modal';
        modal.innerHTML = `
            <div class="confirm-dialog">
                <div class="confirm-header">
                    <h3>确认选择</h3>
                </div>
                <div class="confirm-body">
                    <p>您选择的门店类型是：<strong>${storeType}</strong></p>
                    <p>确认后将跳转到${storeType}注册页面</p>
                </div>
                <div class="confirm-footer">
                    <button class="btn-cancel">取消</button>
                    <button class="btn-confirm">确认</button>
                </div>
            </div>
        `;
        
        // 添加模态框到页面
        document.body.appendChild(modal);
        
        // 添加事件监听器
        const cancelBtn = modal.querySelector('.btn-cancel');
        const confirmBtn = modal.querySelector('.btn-confirm');
        
        // 取消按钮点击事件
        cancelBtn.addEventListener('click', function() {
            document.body.removeChild(modal);
        });
        
        // 确认按钮点击事件
        confirmBtn.addEventListener('click', function() {
            document.body.removeChild(modal);
            window.open(link, '_blank');
        });
        
        // 点击模态框背景关闭
        modal.addEventListener('click', function(e) {
            if (e.target === modal) {
                document.body.removeChild(modal);
            }
        });
    }
});