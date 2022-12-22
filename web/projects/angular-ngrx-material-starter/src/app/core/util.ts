import { WechatService } from '../services/wechat.service';
import { CheckIn } from '../api/models/check-in';

export function isWechat() {
  return /MicroMessenger/i.test(window.navigator.userAgent);
}

export function refreshSharedCheckInToWechat(wechatService: WechatService, checkIn: CheckIn) {
  if (isWechat()) {
    wechatService.refresh(location.href.split('#')[0]).subscribe(
      () => {
        let checkInAtStr = '无'
        if (checkIn?.createdAt) {
          checkInAtStr = new Date(checkIn?.createdAt * 1000 ?? 0).toLocaleString('chinese', {hour12: false})
        }
        wechatService.wx.updateAppMessageShareData({
          title: checkIn?.userName + ' 的今日打卡',
          desc: `排名: ${checkIn?.dayRank}\n时间: ${checkInAtStr}`,
          link: window.location.origin + '/web/check-in/' + checkIn?.key, // 分享链接，该链接域名或路径必须与当前页面对应的公众号 JS 安全域名一致
          imgUrl: checkIn?.headImgUrl, // 分享图标
          success: function () {
            console.log('shared success')
          },
        })
      }
    )
  }
}
