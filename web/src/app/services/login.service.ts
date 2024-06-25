import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class LoginService {
  constructor() {
  }

  public setLoginRedirectUrl(url: string) {
    localStorage.setItem('loginRedirectUrl', url)
  }

  public getLoginRedirectUrl() {
    return localStorage.getItem('loginRedirectUrl')
  }
}
