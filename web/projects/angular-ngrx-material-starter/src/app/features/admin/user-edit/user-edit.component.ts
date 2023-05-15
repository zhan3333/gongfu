import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { NotificationService } from '../../../core/notifications/notification.service';
import { AdminApiService } from '../../../api/admin/admin-api.service';
import { MatChipInputEvent } from '@angular/material/chips';
import { displayRoleName } from '../../../api/models/user';
import { Levels } from '../../../services/coach-level';
import { faBan } from '@fortawesome/free-solid-svg-icons';

// 用户编辑页，可以设置用户角色，设置教练信息等
@Component({
  selector: 'anms-user-edit',
  templateUrl: './user-edit.component.html',
  styleUrls: ['./user-edit.component.scss'],
  changeDetection: ChangeDetectionStrategy.Default
})
export class UserEditComponent implements OnInit {
  public faBan = faBan;
  public id = 0;
  public form = new FormGroup({
    phone: new FormControl(''),
    nickname: new FormControl('', [Validators.required]),
    level: new FormControl(''),
    teachingSpace: new FormControl(''),
    teachingAge: new FormControl(''),
    teachingExperiences: new FormControl([]),
    roleNames: new FormControl([]) // 角色组
  });
  public loading = false;
  // 角色名选项
  public roleNames: string[] = [];
  public displayRoleName = displayRoleName;

  public Levels = Levels;

  constructor(
    public adminApi: AdminApiService,
    public route: ActivatedRoute,
    public notification: NotificationService
  ) {}

  ngOnInit(): void {
    this.id = parseInt(this.route.snapshot.params['id'], 10);
    if (this.id === 0) {
      this.notification.error('invalid user id');
      return;
    }
    this.adminApi.getUser(this.id).subscribe((user) => {
      this.form.patchValue({
        nickname: user.nickname,
        phone: user.phone,
        roleNames: user.roleNames
      });
    });
    this.adminApi.getCoach(this.id).subscribe((coach) => {
      this.form.patchValue({
        level: coach.level,
        teachingSpace: coach.teachingSpace,
        teachingAge: coach.teachingAge,
        teachingExperiences: coach.teachingExperiences
      });
    });
    this.adminApi.getRoleNames().subscribe((roleNames) => {
      this.roleNames = roleNames;
    });
  }

  onClickSave() {
    this.loading = true;
    const data = this.form.value;
    this.adminApi
      .updateUser(this.id, {
        level: data.level,
        nickname: data.nickname,
        phone: data.phone,
        teachingAge: data.teachingAge,
        teachingExperiences: data.teachingExperiences,
        teachingSpace: data.teachingSpace,
        roleNames: data.roleNames
      })
      .subscribe(
        () => {
          this.notification.success('保存成功');
        },
        (error) => {
          this.notification.error('保存失败: ' + error);
        }
      )
      .add(() => {
        this.loading = false;
      });
  }

  onRemoveExp(exp: any) {
    const exps: string[] = this.form.get('teachingExperiences')?.value;
    const newExps: string[] = [];
    for (const item of exps) {
      if (exp === item) {
        continue;
      }
      newExps.push(item);
    }
    this.form.patchValue({
      teachingExperiences: [...newExps]
    });
  }

  onAddExp($event: MatChipInputEvent) {
    const value = ($event.value || '').trim();

    // Add our fruit
    if (value) {
      const exps: string[] = this.form.get('teachingExperiences')?.value;
      exps.push(value);
      this.form.patchValue({
        teachingExperiences: [...exps]
      });
    }

    // Clear the input value
    $event.chipInput?.clear();
  }
}
