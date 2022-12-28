import { WechatService } from '../services/wechat.service';
import { CheckIn } from '../api/models/check-in';
import { Profile } from '../api/models/profile';
import { displayLevel } from '../api/models/user';

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


export function refreshSharedProfileToWechat(wechatService: WechatService, profile: Profile) {
  if (isWechat()) {
    wechatService.refresh(location.href.split('#')[0]).subscribe(
      () => {
        let title = profile.nickname + '的个人信息'
        let desc = ''
        if (profile.role === 'coach') {
          title = profile.nickname + '教练的个人信息'
          desc = `等级: ${displayLevel(profile.coach?.level)}\n任教: ${profile?.coach?.teachingSpace}`
        }
        wechatService.wx.updateAppMessageShareData({
          title: title,
          desc: desc,
          link: window.location.origin + '/web/profile/' + profile?.uuid,
          imgUrl: profile?.headimgurl, // 分享图标
          success: function () {
            console.log('shared success')
          },
        })
      }
    )
  }
}
