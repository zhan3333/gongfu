import { ChangeDetectionStrategy, Component, ElementRef, OnInit, ViewChild } from '@angular/core';
import { ROLE_ADMIN, User } from '../../api/models/user';
import { ApiService } from '../../api/api.service';
import { NotificationService } from '../../core/notifications/notification.service';
import { MatBottomSheet } from '@angular/material/bottom-sheet';
import { BottomSheetComponent } from '../../shared/bottom-sheet.component';
import { FormControl } from '@angular/forms';
import { AuthService } from '../../core/auth/auth.service';

@Component({
  selector: 'anms-setting',
  templateUrl: './setting.component.html',
  styleUrls: ['./setting.component.scss'],
  changeDetection: ChangeDetectionStrategy.Default
})
export class SettingComponent implements OnInit {
  @ViewChild('fileUpload') fileUpload!: ElementRef;
  public user: User | undefined
  // 当前角色
  public curRole = new FormControl(this.authService.getRole());
  public loading = false
  public uploadProgressValue = 0

  constructor(
    private api: ApiService,
    private notification: NotificationService,
    private bottomSheet: MatBottomSheet,
    private authService: AuthService,
  ) {
  }

  ngOnInit(): void {
    this.refreshUser()
    this.curRole.valueChanges.subscribe(v => {
      this.authService.setRole(v)
    })
  }

  refreshUser() {
    this.api.me().subscribe(user => {
      this.user = user
    })
  }

  isAdmin() {
    if (this.user === undefined) {
      return false
    }
    return this.user.role === ROLE_ADMIN
  }

  selectHeadImage() {
    this.bottomSheet.open(BottomSheetComponent, {
      data: new Map<string, string>([
        ['ok', '更换头像'],
        ['cancel', '取消']
      ])
    }).afterDismissed().subscribe((res) => {
      if (res === 'ok') {
        this.fileUpload.nativeElement.click()
      }
    })
  }

  public onFileSelected(event: Event) {
    if (event.target == null) {
      return
    }
    const input = event.target as HTMLInputElement
    if (input.files === null || input.files.length === 0) {
      return
    }
    const file = input.files[0]

    this.loading = true
    this.api.uploadFile(
      file,
      'avatar',
      value => this.uploadProgressValue = value,
      avatarKey => {
        // @ts-ignore
        this.api.editMe(avatarKey).subscribe(
          () => {
            this.notification.success('修改成功')
            this.refreshUser()
          },
          error => this.notification.error('修改失败，请稍后重试: ' + JSON.stringify(error)),
          () => {
            this.loading = false
            this.uploadProgressValue = 0
          }
        )
      }
    ).subscribe(
      () => {
      },
      error => {
        this.notification.error('上传头像失败，请稍后重试: ' + JSON.stringify(error))
        this.loading = false
        this.uploadProgressValue = 0
      },
    )
  }
}
