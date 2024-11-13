import { Injectable } from '@angular/core';
import { ApiService } from '../api/api.service';
import { map } from 'rxjs/operators';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class WechatService {
  public wx = require('weixin-js-sdk')

  constructor(
    private api: ApiService,
    private httpClient: HttpClient,
  ) {
  }

  setOnError(fn: (res: any) => void) {
    this.wx.error(fn)
  }

  setOnSuccess(fn: () => void) {
    this.wx.ready(fn)
  }

  refresh(uri: string) {
    this.wx.error((res: any) => console.error('wechat config error: ', res))
    return this.api.getWechatJSConfig(uri).pipe(
      map(config => {
          this.wx.config({
            debug: true,
            appId: config.appID,
            timestamp: config.timestamp,
            nonceStr: config.nonceStr,
            signature: config.signature,
            jsApiList: ['updateAppMessageShareData', 'onMenuShareAppMessage',] // 必填，需要使用的 JS 接口列表
          })
        }
      ),
    )
  }

  pay(input: PayInput) {
    this.refresh(location.href.split('#')[0]).subscribe(() => {
      this.httpClient.get<WechatPayResponse>('/wechat/pay', {
        params: {
          amount: input.amount,
          description: input.description,
          attach: input.attach,
          activity_id: input.activity_id,
          username: input.username,
          phone: input.phone,
          sex: input.sex,
        }
      }).subscribe((v: WechatPayResponse) => {
        this.wx.chooseWXPay({
          timestamp: v.timeStamp,
          nonceStr: v.nonceStr,
          package: v.package,
          signType: v.signType,
          paySign: v.paySign,
          success(res: any) {
            console.log('call wechat pay res', res)
            if (res.err_msg == "get_brand_wcpay_request:ok") {
              // 使用以上方式判断前端返回,微信团队郑重提示：
              //res.err_msg将在用户支付成功后返回ok，但并不保证它绝对可靠。
              input.onSuccess()
            } else {
              input.onFailed(res.err_msg)
            }
          }
        })
      })
    })
  }

  // 获取当前活动的报名信息
  getEnroll(activityId: string) {
    return this.httpClient.get<EnrollResponse>('/wechat/enroll', {params: {activity_id: activityId}})
  }
}

export interface EnrollResponse {
  id: number
  article_id: string
  status: string
  out_trade_no: string
  created_at: string
  transaction_id: string
  amount: number // 金额，单位分
}

interface WechatPayResponse {
  prepay_id: string
  appId: string
  timeStamp: number
  nonceStr: string
  package: string
  signType: string
  paySign: string
}

export interface PayInput {
  amount: number
  description: string
  attach: string
  activity_id: string
  onSuccess: () => void,
  onFailed: (err: string) => void

  username: string
  phone: string
  sex: string
}
