import { ChangeDetectionStrategy, Component, OnInit } from '@angular/core';
import { FormsModule, ReactiveFormsModule, UntypedFormControl, UntypedFormGroup, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { NotificationService } from '../../../core/notifications/notification.service';
import { AdminApiService } from '../../../api/admin/admin-api.service';
import { MatChipInputEvent, MatChipsModule } from '@angular/material/chips';
import { displayRoleName, User } from '../../../api/models/user';
import { Levels } from '../../../services/coach-level';
import { faBan, faChalkboardTeacher, faEdit, faFolder, faGraduationCap } from '@fortawesome/free-solid-svg-icons';
import { MatButtonModule } from '@angular/material/button';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { MatIconModule } from '@angular/material/icon';
import { MatOptionModule, MatRippleModule } from '@angular/material/core';
import { NgFor, NgIf } from '@angular/common';
import { MatSelectModule } from '@angular/material/select';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatCardModule } from '@angular/material/card';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatListModule } from '@angular/material/list';
import { TeachingRecordDialog } from './teaching-dialog';
import { StudyRecordDialog } from './study-dialog';
import { StudyRecord } from '../../../api/models/study-record';
import { TeachingRecord } from '../../../api/models/teaching-record';

// 用户编辑页，可以设置用户角色，设置教练信息等
@Component({
  selector: 'anms-user-edit',
  templateUrl: './user-edit.component.html',
  changeDetection: ChangeDetectionStrategy.Default,
  standalone: true,
  imports: [
    MatCardModule,
    FormsModule,
    ReactiveFormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    NgFor,
    MatOptionModule,
    MatChipsModule,
    MatIconModule,
    FontAwesomeModule,
    MatButtonModule,
    MatListModule,
    NgIf,
    MatDialogModule,
    MatRippleModule
  ]
})
export class UserEditComponent implements OnInit {
  public faBan = faBan;
  public id = 0;
  public form = new UntypedFormGroup({
    phone: new UntypedFormControl(''),
    nickname: new UntypedFormControl('', [Validators.required]),
    level: new UntypedFormControl(''),
    teachingSpace: new UntypedFormControl(''),
    teachingAge: new UntypedFormControl(''),
    teachingExperiences: new UntypedFormControl([]),
    roleNames: new UntypedFormControl([]) // 角色组
  });
  public loading = false;
  // 角色名选项
  public roleNames: string[] = [];
  public displayRoleName = displayRoleName;

  public Levels = Levels;
  public user: User | null = null;
  protected readonly faEdit = faEdit;
  protected readonly faFolder = faFolder;
  protected readonly faChalkboardTeacher = faChalkboardTeacher;
  protected readonly faGraduationCap = faGraduationCap;

  constructor(
    public adminApi: AdminApiService,
    public route: ActivatedRoute,
    public notification: NotificationService,
    public router: Router,
    private dialog: MatDialog
  ) {}

  ngOnInit(): void {
    this.id = parseInt(this.route.snapshot.params['id'], 10);
    if (this.id === 0) {
      this.notification.error('invalid user id');
      return;
    }
    this.refresh();
  }

  refresh() {
    this.adminApi.getUser(this.id).subscribe((user) => {
      this.user = user;
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

  // 打开新增授课记录弹窗
  openTeachingRecordDialog(record?: TeachingRecord) {
    this.dialog
      .open(TeachingRecordDialog, { data: record })
      .afterClosed()
      .subscribe(
        (result?: {
          data?: { id: number; date: string; address: string };
          isDelete?: boolean;
        }) => {
          if (!result) {
            return;
          }
          if (result.isDelete) {
            // 删除
          }
          if (result.data) {
            this.adminApi
              .editTeachingRecord({
                id: result.data.id,
                date: result.data.date,
                address: result.data.address,
                userId: this.user!.id
              })
              .subscribe({
                next: () => {
                  this.refresh();
                  this.notification.success('ok');
                }
              });
          }
        }
      );
  }

  // 打开新增授课记录弹窗
  openStudyRecordDialog(record?: StudyRecord) {
    this.dialog
      .open(StudyRecordDialog, { data: record })
      .afterClosed()
      .subscribe(
        (result?: {
          data?: { id: number; date: string; content: string };
          isDelete?: boolean;
        }) => {
          if (!result) {
            return;
          }
          if (result.isDelete) {
            this.adminApi.deleteStudyRecord(record!.id).subscribe({
              next: () => {
                this.refresh();
                this.notification.success('ok');
              }
            });
          }
          if (result.data) {
            this.adminApi
              .editStudyRecord({
                id: result.data.id,
                date: result.data.date,
                content: result.data.content,
                userId: this.user!.id
              })
              .subscribe({
                next: () => {
                  this.refresh();
                  this.notification.success('ok');
                }
              });
          }
        }
      );
  }
}
