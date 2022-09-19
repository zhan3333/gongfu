import { Injectable } from '@angular/core';
import { HttpClient, HttpEvent, HttpEventType, HttpRequest } from '@angular/common/http';
import { User } from './models/user';
import { CheckIn, CheckInExist } from './models/check-in';
import { last, map, switchMap, tap } from 'rxjs/operators';

const meUrl = '/me'
const sendValidCodeUrl = '/bind/phone'
const validCodeUrl = '/bind/phone/valid'
const todayCheckInUrl = '/check-in/today'
const checkInUrl = '/check-in'
const getUploadTokenUrl = '/storage/upload-token'
const wechatJSConfigUrl = '/wechat/js-config'
const showCheckIn = '/check-in/'
const getCheckInTop = '/check-in/top'

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  constructor(private http: HttpClient) {
  }

  me() {
    return this.http.get<User>(meUrl)
  }

  sendValidCode(phone: string) {
    return this.http.post(sendValidCodeUrl, {phone: phone})
  }

  validCode(phone: string, code: string) {
    return this.http.post(validCodeUrl, {phone: phone, code: code})
  }

  // 获取今天的打卡信息
  getTodayCheckIn() {
    return this.http.get<CheckInExist>(todayCheckInUrl)
  }

  checkIn(fileKey: string, fileName: string) {
    return this.http.post(checkInUrl, {key: fileKey, fileName: fileName})
  }

  getUploadUrl(fileName: string) {
    return this.http.get<GetUploadTokenResponse>(getUploadTokenUrl, {
      params: {
        fileName: fileName
      }
    })
  }

  uploadFile(file: File, onProgress: (value: number) => void, onCompletely: (key: string) => void) {
    return this.getUploadUrl(file.name).pipe(
      switchMap(token => {
        const req = new HttpRequest('PUT', token.uploadUrl, file, {
          reportProgress: true
        })
        return this.http.request(req).pipe(
          map(event => {
            switch (event.type) {
              case HttpEventType.UploadProgress:
                return event.total ? Math.round(100 * event.loaded / event.total) : 0;
              case HttpEventType.Response:
                onCompletely(token.key)
                return 100;
              default:
                return 0;
            }
          }),
          tap(value => onProgress(value)),
          last(), // return last (completed) message to caller
        )
      })
    )
  }


  getWechatJSConfig(uri: string) {
    return this.http.get<WechatJSConfig>(wechatJSConfigUrl, {
      params: {
        uri: uri
      }
    })
  }

  getCheckInByKey(key: string) {
    return this.http.get<CheckIn>(showCheckIn + key)
  }

  getCheckInTop(startAt: number, endAt: number) {
    return this.http.get<CheckIn[]>(getCheckInTop, {
      params: {
        startAt: startAt,
        endAt: endAt,
      }
    })
  }
}

interface GetUploadTokenResponse {
  uploadUrl: string
  key: string
}

interface WechatJSConfig {
  appID: string
  timestamp: number
  nonceStr: string
  signature: string
}
