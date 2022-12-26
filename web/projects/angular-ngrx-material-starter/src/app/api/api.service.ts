import { Injectable } from '@angular/core';
import { HttpClient, HttpEventType, HttpRequest } from '@angular/common/http';
import { Coach, User } from './models/user';
import { CheckIn, CheckInCountList, CheckInExist, CheckInList } from './models/check-in';
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
const getCheckInCountTop = '/check-in/top/count'
const getCheckInContinuousTop = '/check-in/top/continuous'
const getCheckInHistories = '/check-in/histories'
const login = '/auth/login'
const getCoach = '/coach'

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

  getCheckInTop(date: string) {
    return this.http.get<CheckInList>(getCheckInTop, {
      params: {
        date,
      }
    })
  }

  // 获取打卡次数排行榜
  getCheckInCountTop() {
    return this.http.get<CheckInCountList>(getCheckInCountTop)
  }

  // 获取连续打卡排行榜
  getCheckInContinuousTop() {
    return this.http.get<CheckInCountList>(getCheckInContinuousTop)
  }

  // 获取用户打卡历史
  getCheckInHistories(userID: number, startDate: string, endDate: string) {
    return this.http.get<CheckInList>(getCheckInHistories, {params: {userID, startDate, endDate}})
  }

  // 手机号登录
  login(phone: string, code: string) {
    return this.http.post(login, {phone, code})
  }

  // 获取教练信息
  getCoach() {
    return this.http.get<Coach>(getCoach)
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
