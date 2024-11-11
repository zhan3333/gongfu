import { ChangeDetectionStrategy, Component, Inject, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { NotificationService } from '../../core/notifications/notification.service';
import { ApiService } from '../../api/api.service';
import { Profile } from '../../api/models/profile';
import { refreshSharedProfileToWechat } from '../../core/util';
import { WechatService } from '../../services/wechat.service';
import { displayLevel } from '../../services/coach-level';
import { MatLineModule, MatRippleModule } from '@angular/material/core';
import { DatePipe, NgFor, NgIf, NgOptimizedImage } from '@angular/common';
import { MatDividerModule } from '@angular/material/divider';
import { AuthService } from '../../core/auth/auth.service';
import { ROLE_ADMIN, ROLE_COACH, UserClass } from '../../api/models/userClass';
import { MAT_DIALOG_DATA, MatDialog, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { FormControl, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatInputModule } from '@angular/material/input';
import { MemberCourse } from '../../api/admin/admin-api.service';

@Component({
  selector: 'anms-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.scss'],
  changeDetection: ChangeDetectionStrategy.Default,
  standalone: true,
  imports: [NgIf, NgOptimizedImage, NgFor, MatLineModule, DatePipe, MatDividerModule, MatRippleModule, MatDialogModule]
})
export class ProfileComponent implements OnInit {
  public profile: Profile | undefined;
  public displayLevel = displayLevel;
  public user: UserClass
  private userUUID: string = ''

  constructor(
    private activeRoute: ActivatedRoute,
    private notificationService: NotificationService,
    private api: ApiService,
    private wechat: WechatService,
    public authService: AuthService,
    private dialog: MatDialog
  ) {
    this.user = this.authService.getUser()
  }

  ngOnInit(): void {
    this.activeRoute.paramMap.subscribe((data) => {
      const uuid = data.get('uuid');
      if (uuid === '' || uuid === null) {
        this.notificationService.error('invalid uuid');
        return;
      }
      this.userUUID = uuid
      this.refresh()
    });
  }

  refresh() {
    this.api.getProfile(this.userUUID).subscribe((profile) => {
      this.profile = profile;
      refreshSharedProfileToWechat(this.wechat, this.profile);
    });
  }

  clickEditMemberCourseRemain(mc: MemberCourse) {
    if (this.user.hasAnyRole([ROLE_ADMIN, ROLE_COACH])) {
      this.dialog.open(UpdateMemberCourseRemainDialog, {
        data: {
          total: mc.total,
          remain: mc.remain,
        }
      }).afterClosed().subscribe((remain) => {
        if (remain !== undefined) {
          this.api.updateMemberCourseRemain(mc.id, remain).subscribe(() => {
            this.notificationService.success('更新成功');
            this.refresh()
          });
        }
      });
    }
  }
}

// update member course remain dialog

@Component({
  selector: 'anms-update-member-course-remain-dialog',
  template: `
    <mat-dialog-content>
      <mat-form-field>
        <mat-label>剩余课时</mat-label>
        <input matInput type="number" [formControl]="remainFormControl"/>
        <mat-error *ngIf="remainFormControl.hasError('required')">必填</mat-error>
        <mat-error *ngIf="remainFormControl.hasError('max')">最大为 {{ data.total }}</mat-error>
        <mat-error *ngIf="remainFormControl.hasError('min')">最小为 0</mat-error>
      </mat-form-field>
    </mat-dialog-content>
    <mat-dialog-actions>
      <button mat-button (click)="dialogRef.close()">取消</button>
      <button mat-button color="primary" [mat-dialog-close]="remainFormControl.value"
              [disabled]="remainFormControl.invalid">确定
      </button>
    </mat-dialog-actions>
  `,
  standalone: true,
  imports: [MatFormFieldModule, MatButtonModule, MatDialogModule, MatDialogModule, MatButtonModule, MatFormFieldModule, ReactiveFormsModule, MatInputModule, NgIf]
})
class UpdateMemberCourseRemainDialog {
  public remainFormControl = new FormControl(this.data.remain, {nonNullable: true, validators: [Validators.required]});

  constructor(
    public dialogRef: MatDialogRef<UpdateMemberCourseRemainDialog>,
    @Inject(MAT_DIALOG_DATA) public data: { total: number, remain: number },
  ) {
    this.remainFormControl.addValidators([Validators.max(data.total), Validators.min(0)])
  }
}
