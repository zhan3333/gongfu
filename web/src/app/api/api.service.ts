import { Injectable } from '@angular/core';
import { HttpClient, HttpEventType, HttpRequest } from '@angular/common/http';
import { ICoach, SimpleUser, User, UsersPage } from './models/user';
import { CheckIn, CheckInCountList, CheckInExist, CheckInList } from './models/check-in';
import { last, map, switchMap, tap } from 'rxjs/operators';
import { Profile } from './models/profile';
import { Course } from './models/course';
import { environment } from '../../environments/environment';

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
const getProfile = '/profile'
const editMe = '/me'
const getUsers = '/admin/users'
const getUser = '/admin/user'
const updateUser = '/admin/user'

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  constructor(private http: HttpClient) {
  }

  me() {
    return this.http.get<User>(meUrl).pipe(
      map(v => new User(v))
    )
  }

  // 编辑用户信息
  editMe(data: { avatarKey?: string, nickname?: string }) {
    return this.http.post(editMe, {avatarKey: data.avatarKey, nickname: data.nickname})
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

  getUploadUrl(fileName: string, dir: string = '') {
    return this.http.get<GetUploadTokenResponse>(getUploadTokenUrl, {
      params: {
        fileName: fileName,
        dir: dir,
      }
    })
  }

  uploadFile(file: File, dir: string, onProgress: (value: number) => void, onCompletely: (key: string) => void) {
    return this.getUploadUrl(file.name, dir).pipe(
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
    return this.http.get<ICoach>(getCoach)
  }

  getProfile(uuid: string) {
    return this.http.get<Profile>(getProfile + '/' + uuid).pipe(
      map(v => new Profile(v))
    )
  }

  getUsers(data: GetUsersParams) {
    return this.http.get<UsersPage>(getUsers, {
      params: {
        page: data.page,
        limit: data.limit,
        keyword: data.keyword,
        desc: data.desc,
      }
    }).pipe(
      map(v => new UsersPage(v))
    )
  }

  // 获取用户信息
  getUser(id: number) {
    return this.http.get<User>(getUser + '/' + id).pipe(
      map(v => new User(v))
    )
  }

  // 获取课程列表
  getCourses() {
    return this.http.get<Course[]>('/courses',)
  }

  // 获取课程详情
  getCourse(id: number) {
    return this.http.get<Course>('/courses/' + id)
  }

  updateCourse(id: number, course: any) {
    return this.http.put('/courses/' + id, course)
  }

  getFilePerSignedUrl(key: string) {
    return this.http.get<{ url: string }>('/storage/per-sign/' + key).pipe(
      map(v => v.url)
    )
  }

  visitFile(key: string) {
    return environment.api + '/storage/visit/' + key
  }

  // update member course remain
  updateMemberCourseRemain(id: number, courseRemain: number) {
    return this.http.put('/member-courses/' + id + '/remain', {remain: courseRemain})
  }

  // create check in comment
  createCheckInComment(checkInId: number, content: string) {
    return this.http.post('/check-in/comment', {content}, {
      params: {
        id: checkInId,
      }
    })
  }

  // get check in comments
  getCheckInComments(checkInId: number) {
    return this.http.get<CheckInComment[]>('/check-in/comments', {
      params: {
        id: checkInId
      }
    })
  }
}

export interface CheckInComment {
  id: number
  content: string
  createdAt: string
  user: SimpleUser
}

export interface GetUsersParams {
  page: number,
  limit: number,
  keyword: string,
  desc: boolean
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
