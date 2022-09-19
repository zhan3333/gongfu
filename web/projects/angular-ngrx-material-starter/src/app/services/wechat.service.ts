import { Injectable } from '@angular/core';
import { ApiService } from '../api/api.service';
import { AsyncSubject } from 'rxjs';
import { map } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class WechatService {
  public wx = require('weixin-js-sdk')

  constructor(
    private api: ApiService,
  ) {
  }

  refresh(uri: string) {
    this.wx.error((res: any) => console.log('wechat config error: ', res))
    return this.api.getWechatJSConfig(uri).pipe(
      map(config => {
          this.wx.config({
            debug: false,
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
}
